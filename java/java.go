package java

/*
#include <stdint.h>
#include <stdlib.h>
#cgo arm64 LDFLAGS: -L../libs/arm64-v8a -lAutoGo
#cgo amd64 LDFLAGS: -L../libs/x86_64 -lAutoGo
#cgo 386 LDFLAGS: -L../libs/x86 -lAutoGo

int shmem_create(const char *name, size_t size);
int shmem_read(void *buffer, size_t buffer_size, int clear_data);
int shmem_getfb();
*/
import "C"
import (
	"bufio"
	_ "embed"
	"fmt"
	"github.com/Dasongzi1366/AutoGo/utils"
	"io"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"
)

//go:embed utils.dex
var _utils []byte

//go:embed libashmem_arm64.so
var _libashmem_arm64 []byte

//go:embed libashmem_amd64.so
var _libashmem_amd64 []byte

//go:embed libashmem_amd.so
var _libashmem_amd []byte

var conn net.Conn
var msgChanMap = make(map[string]chan string)
var mu sync.Mutex
var currentMsgId int

func init() {
	path := filepath.Dir(os.Args[0])
	width, height := getWmSize()
	size := (width + 40) * (height + 40) * 4
	handleError(shmemCreate("my_shared_memory", size))
	fd := shmemGetFD()
	listen, err := net.Listen("tcp", ":56788")
	handleError(err)
	utils.Shell("pkill -f com.autogo.Main")
	if runtime.GOARCH == "arm64" {
		_ = os.WriteFile(path+"/libashmem.so", _libashmem_arm64, 0644)
	} else if runtime.GOARCH == "amd64" {
		_ = os.WriteFile(path+"/libashmem.so", _libashmem_amd64, 0644)
	} else {
		_ = os.WriteFile(path+"/libashmem.so", _libashmem_amd, 0644)
	}
	_ = os.WriteFile(path+"/autogo.dex", _utils, 0644)
	go func() {
		cmdStr := "app_process -Djava.class.path=" + path + "/autogo.dex / com.autogo.Main " + fmt.Sprintf("%s %d %d", "56788", fd, size)
		cCmd := C.CString(cmdStr)
		defer C.free(unsafe.Pointer(cCmd))
		C.system(cCmd)
		_, _ = fmt.Fprintln(os.Stderr, "[AutoGo] 服务异常崩溃")
		os.Exit(1)
	}()
	conn, err = listen.Accept()
	_ = os.Remove(path + "/autogo.dex")
	handleError(err)
	go func() {
		reader := bufio.NewReader(conn)
		for {
			// 先读取消息ID和消息长度，共12个字节
			header := make([]byte, 12)
			_, err = reader.Read(header)
			if err != nil {
				fmt.Println("Error reading message header:", err)
				return
			}
			// 解析消息ID（前6个字节）
			msgId := string(header[:6])

			// 解析消息长度（后6个字节）
			lengthStr := string(header[6:12])
			msgLength, err := strconv.Atoi(lengthStr)
			if err != nil {
				fmt.Println("Error parsing message length:", err)
				return
			}

			// 读取消息体
			message := make([]byte, msgLength)
			_, err = io.ReadFull(reader, message)
			if err != nil {
				fmt.Println("Error reading message body:", err)
				return
			}

			// 将消息发送到对应的通道
			mu.Lock()
			if ch, ok := msgChanMap[msgId]; ok {
				ch <- string(message)
				close(ch)
				delete(msgChanMap, msgId)
			}
			mu.Unlock()
		}
	}()
}

func CallJavaMethod(model string, str string) string {
	respChan := make(chan string)
	mu.Lock()
	currentMsgId++
	if currentMsgId > 999999 {
		currentMsgId = 1
	}
	msgId := fmt.Sprintf("%06d", currentMsgId)
	_, _ = conn.Write([]byte(model + "|" + msgId + "|" + str + "\n"))
	msgChanMap[msgId] = respChan
	mu.Unlock()
	select {
	case resp := <-respChan:
		mu.Lock()
		delete(msgChanMap, msgId) // 收到响应后清理 map
		mu.Unlock()
		return resp
	case <-time.After(10 * time.Second):
		mu.Lock()
		delete(msgChanMap, msgId) // 超时后清理 map
		mu.Unlock()
		fmt.Println("Timeout waiting for response")
		return ""
	}
}

func Send(str string) {
	mu.Lock()
	_, _ = conn.Write([]byte(str + "\n"))
	mu.Unlock()
}

func GetBitMapData() ([]byte, error) {
	return shmemRead(1)
}

func shmemCreate(name string, size int) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	ret := C.shmem_create(cName, C.size_t(size))
	if ret < 0 {
		return fmt.Errorf("failed to create ashmem")
	}

	return nil
}

func shmemRead(clearData int) ([]byte, error) {
	// 第一次读取时，传入 nil 来获取数据长度
	dataLen := int(C.shmem_read(nil, 0, 0)) // 先获取数据长度
	if dataLen < 0 {
		return nil, fmt.Errorf("failed to get data length from shared memory")
	}

	if dataLen == 0 {
		return nil, nil
	}

	// 创建一个大小为 dataLen 的缓冲区
	buffer := make([]byte, dataLen)

	// 读取数据到缓冲区，传入 clearData 参数
	bufferPtr := unsafe.Pointer(&buffer[0])
	readSize := C.shmem_read(bufferPtr, C.size_t(len(buffer)), C.int(clearData))
	if int(readSize) != dataLen {
		return nil, fmt.Errorf("buffer size mismatch: expected %d, got %d", dataLen, int(readSize))
	}

	return buffer, nil
}

func shmemGetFD() int {
	return int(C.shmem_getfb())
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "[AutoGo] 出现错误:", err)
		os.Exit(1)
	}
}

func getWmSize() (int, int) {
	input := utils.Shell("wm size")
	lastLine := strings.Split(input, "\n")[len(strings.Split(input, "\n"))-1]
	re := regexp.MustCompile(`(\d+)x(\d+)`)
	matches := re.FindStringSubmatch(lastLine)
	if len(matches) == 3 {
		return s2i(matches[1]), s2i(matches[2])
	}
	fmt.Fprintln(os.Stderr, "[AutoGo] 设备分辨率获取失败")
	os.Exit(1)
	return 0, 0
}

func s2i(s string) int {
	i, _ := strconv.Atoi(strings.TrimSpace(s))
	return i
}
