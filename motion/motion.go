package motion

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	"fmt"
	"github.com/Dasongzi1366/AutoGo/java"
	"math/rand"
	"strings"
	"time"
	"unsafe"
)

func TouchDown(x, y, fingerID int) {
	fingerID = fingerID - 1
	if fingerID < 0 || fingerID > 9 {
		fingerID = 0
	}
	java.Send(fmt.Sprintf("d|%d|%d|%d", x, y, fingerID))
}

func TouchMove(x, y, fingerID int) {
	fingerID = fingerID - 1
	if fingerID < 0 || fingerID > 9 {
		fingerID = 0
	}
	java.Send(fmt.Sprintf("m|%d|%d|%d", x, y, fingerID))
}

func TouchUp(x, y, fingerID int) {
	fingerID = fingerID - 1
	if fingerID < 0 || fingerID > 9 {
		fingerID = 0
	}
	java.Send(fmt.Sprintf("u|%d|%d|%d", x, y, fingerID))
}

func Click(x, y, fingerID int) {
	TouchDown(x, y, fingerID)
	sleep(random(10, 20))
	TouchUp(x, y, fingerID)
}

func LongClick(x, y, duration int) {
	TouchDown(x, y, 1)
	sleep(duration + random(1, 20))
	TouchUp(x, y, 1)
}

func Swipe(x1, y1, x2, y2, duration int) {
	java.Send(fmt.Sprintf("s1|%d|%d|%d|%d|%d", x1, y1, x2, y2, duration))
}

func Swipe2(x1, y1, x2, y2, duration int) {
	java.Send(fmt.Sprintf("s2|%d|%d|%d|%d|%d", x1, y1, x2, y2, duration))
}

func Home() {
	KeyAction(KEYCODE_HOME)
}

func Back() {
	KeyAction(KEYCODE_BACK)
}

func Recents() {
	KeyAction(KEYCODE_APP_SWITCH)
}

func PowerDialog() {
	shell("input keyevent --longpress KEYCODE_POWER")
}

func Notifications() {
	KeyAction(KEYCODE_NOTIFICATION)
}

func QuickSettings() {
	shell("cmd statusbar expand-settings")
}

func VolumeUp() {
	KeyAction(KEYCODE_VOLUME_UP)
}

func VolumeDown() {
	KeyAction(KEYCODE_VOLUME_DOWN)
}

func Camera() {
	KeyAction(KEYCODE_CAMERA)
}

func KeyAction(code int) {
	java.Send(fmt.Sprintf("k|%d", code))
}

func random(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min+1) + min
}

func sleep(i int) {
	time.Sleep(time.Duration(i) * time.Millisecond)
}

func shell(cmd string) {
	if strings.Contains(cmd, ";") {
		cmd = "(" + cmd + ")"
	}
	cCmd := C.CString(cmd + " > /dev/null 2>&1")
	defer C.free(unsafe.Pointer(cCmd))
	C.system(cCmd)
}
