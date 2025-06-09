package ime

import (
	"bufio"
	"encoding/base64"
	"github.com/Dasongzi1366/AutoGo/java"
	"github.com/Dasongzi1366/AutoGo/motion"
	"github.com/Dasongzi1366/AutoGo/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// GetClipText 获取剪切板内容
func GetClipText() string {
	return java.CallJavaMethod("utils", "getClipText")
}

// SetClipText 设置剪切板内容
func SetClipText(text string) bool {
	return java.CallJavaMethod("utils", "setClipText|"+text) == "true"
}

// InputText 输入文本
func InputText(text string) {
	execPath, _ := filepath.Abs(os.Args[0])
	pkg := filepath.Base(filepath.Dir(execPath))

	if !strings.Contains(pkg, ".") {
		// 判断字符串是否包含中文
		if containsChinese(text) {
			SetClipText(text)
			motion.KeyAction(279)
			time.Sleep(time.Duration(30) * time.Millisecond)
			SetClipText("")
		} else {
			// 如果不包含中文，则使用 Shell 输入文本
			utils.Shell("input text " + strconv.Quote(text))
		}
	} else {
		utils.Shell("am broadcast -a com.autogo --es message inputtext --es data " + base64.StdEncoding.EncodeToString([]byte(text)))
	}
}

// KeyAction 模拟按键
func KeyAction(code int) {
	utils.Shell("am broadcast -a com.autogo --es message keyaction --es data " + strconv.Itoa(code))
}

// GetIMEList 获取输入法列表
func GetIMEList() []string {
	shellOutput := utils.Shell("ime list -a | grep mId")

	var imeList []string

	scanner := bufio.NewScanner(strings.NewReader(shellOutput))
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "mId=") {
			mId := strings.SplitN(line, " ", 2)[0]
			mId = strings.TrimPrefix(mId, "mId=")
			imeList = append(imeList, mId)
		}
	}

	if err := scanner.Err(); err != nil {
		//fmt.Println("Scanner error:", err)
		return nil
	}

	return imeList
}

// SetCurrentIME 设置当前输入法
func SetCurrentIME(imeID string) {
	if imeID == "" {
		execPath, _ := filepath.Abs(os.Args[0])
		pkg := filepath.Base(filepath.Dir(execPath))
		if !strings.Contains(pkg, ".") {
			return
		}
		imeID = pkg + "/.keyboard.IME"
	}
	utils.Shell("ime enable " + imeID)
	utils.Shell("ime set " + imeID)
}

// containsChinese 检查字符串是否包含中文
func containsChinese(s string) bool {
	for _, r := range s {
		if r >= '\u4e00' && r <= '\u9fff' { // 中文字符范围
			return true
		}
	}
	return false
}
