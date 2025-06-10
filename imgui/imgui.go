package imgui

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/wait.h>

static int input_pipe[2] = {-1, -1}; // 管道用于向命令输入数据
static pid_t child_pid = -1;         // 子进程 ID

// 启动阻塞命令
int start_command(const char* cmd) {
    if (pipe(input_pipe) == -1) {
        perror("pipe creation failed");
        return -1;
    }

    child_pid = fork();
    if (child_pid == -1) {
        perror("fork failed");
        return -1;
    }

    if (child_pid == 0) { // 子进程
        close(input_pipe[1]);  // 关闭写端
        dup2(input_pipe[0], STDIN_FILENO); // 标准输入重定向
        close(input_pipe[0]);  // 关闭读端

        execl("/system/bin/sh", "sh", "-c", cmd, NULL);
        exit(1); // 如果 execl 失败，退出
    } else { // 父进程
        close(input_pipe[0]); // 父进程关闭管道读端

        // 阻塞等待子进程完成
        int status;
        waitpid(child_pid, &status, 0); // 等待子进程
        if (WIFEXITED(status)) {
            return WEXITSTATUS(status); // 返回子进程的退出状态
        } else {
            return -1; // 子进程异常
        }
    }
}

// 向命令输入数据
int write_to_command(const char* input) {
    if (child_pid == -1 || input_pipe[1] == -1) {
        fprintf(stderr, "Command is not running\n");
        return -1;
    }

    // 写入数据到管道
    ssize_t written = write(input_pipe[1], input, strlen(input));
    if (written == -1) {
        perror("write to pipe failed");
        return -1;
    }

    // 写入换行符，模拟控制台行为
    write(input_pipe[1], "\n", 1);

    return 0;
}
*/
import "C"
import (
	"encoding/base64"
	"fmt"
	"github.com/Dasongzi1366/AutoGo/device"
	"github.com/Dasongzi1366/AutoGo/utils"
	"image/color"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// TextItem 表示颜色和文本的组合
type TextItem struct {
	TextColor color.Color //文字颜色 如果为 nil，则默认使用白色 (255,255,255)
	Text      string
}

func init() {
	if s2i(device.SdkInt) > 27 {
		imguiPath := filepath.Dir(os.Args[0]) + "/imgui"
		utils.Shell("pkill -f " + imguiPath)
		os.WriteFile(imguiPath, _imgui, 0644)
		utils.Shell("chmod 755 " + imguiPath)
		h := strconv.Itoa(maxInt(getWmSize()))
		go func() {
			cCmd := C.CString(imguiPath + " " + h)
			defer C.free(unsafe.Pointer(cCmd))
			C.start_command(cCmd)
			_, _ = fmt.Fprintln(os.Stderr, "[imgui] 服务异常崩溃")
			os.Exit(1)
		}()
		sleep(150)
	}
}

// HudInit Hud初始化
// x1, y1, x2, y2 是状态条的坐标
// bgColor 是状态条的背景颜色，如果为 nil，则默认使用灰色 (100, 100, 100)
// textSize 是状态条上的文字大小，如果小于等于 0，则默认使用 45
func HudInit(x1, y1, x2, y2 int, bgColor color.Color, textSize int) {
	fmt.Fprintln(bgColor)
	if s2i(device.SdkInt) > 27 {
		var rInt, gInt, bInt int
		if bgColor == nil {
			rInt, gInt, bInt = 100, 100, 100
		} else {
			r, g, b, _ := bgColor.RGBA()
			rInt = int(r >> 8)
			gInt = int(g >> 8)
			bInt = int(b >> 8)
		}
		if textSize <= 0 {
			textSize = 45
		}
		writeToCommand(fmt.Sprintf("draw_text_rect %d %d %d %d %d %d %d %d", x1, y1, x2, y2, rInt, gInt, bInt, textSize))
	}
}

// HudSetText 设置Hud文本
func HudSetText(items []TextItem) {
	if s2i(device.SdkInt) > 27 {
		str := ""
		for _, item := range items {
			if item.Text != "" {
				var rInt, gInt, bInt string
				if item.TextColor == nil {
					rInt, gInt, bInt = "255", "255", "255"
				} else {
					r, g, b, _ := item.TextColor.RGBA()
					rInt = i2s(int(r >> 8))
					gInt = i2s(int(g >> 8))
					bInt = i2s(int(b >> 8))
				}
				str = str + rInt + "|" + gInt + "|" + bInt + "|" + item.Text + "|"
			}
		}
		writeToCommand("draw_text " + str)
	}
}

// HudClose 状态条销毁
func HudClose() {
	if s2i(device.SdkInt) > 27 {
		writeToCommand("remove_text")
	}
}

// DrawRect 绘制矩形
func DrawRect(x1, y1, x2, y2 int, c color.Color) {
	if s2i(device.SdkInt) > 27 {
		r, g, b, _ := c.RGBA()
		rInt := int(r >> 8)
		gInt := int(g >> 8)
		bInt := int(b >> 8)
		writeToCommand(fmt.Sprintf("draw_rect %d %d %d %d %d %d %d", x1, y1, x2, y2, rInt, gInt, bInt))
	}
}

// ClearRect 清除绘制的矩形
func ClearRect() {
	if s2i(device.SdkInt) > 27 {
		writeToCommand("clear")
	}
}

// Toast 显示Toast提示信息
func Toast(message string) {
	if message == "" {
		message = " "
	}
	if s2i(device.SdkInt) < 28 {
		utils.Shell("am broadcast -a com.autogo --es message toast --es data " + base64.StdEncoding.EncodeToString([]byte(message)))
	} else {
		writeToCommand(fmt.Sprintf("toast %d %d %s", device.Width, device.Height, message))
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
	_, _ = fmt.Fprintln(os.Stderr, "[AutoGo] 设备分辨率获取失败")
	os.Exit(1)
	return 0, 0
}

func writeToCommand(input string) {
	cInput := C.CString(input)
	defer C.free(unsafe.Pointer(cInput))

	// 调用 C 函数向命令输入数据
	ret := C.write_to_command(cInput)
	if ret != 0 {
		_, _ = fmt.Fprintln(os.Stderr, "[imgui] 写入数据失败")
		os.Exit(1)
	}
}

func sleep(i int) {
	time.Sleep(time.Duration(i) * time.Millisecond)
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func s2i(s string) int {
	i, _ := strconv.Atoi(strings.TrimSpace(s))
	return i
}

func i2s(i int) string {
	return strconv.Itoa(i)
}
