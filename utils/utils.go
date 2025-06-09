package utils

/*
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include <android/log.h>
#cgo arm64 LDFLAGS: -L../libs/arm64-v8a -lAutoGo
#cgo amd64 LDFLAGS: -L../libs/x86_64 -lAutoGo
#cgo 386 LDFLAGS: -L../libs/x86 -lAutoGo

void logI(const char* label,const char* message) {
    __android_log_print(ANDROID_LOG_INFO, label, "%s", message);
}

void logE(const char* label, const char* message) {
    __android_log_print(ANDROID_LOG_ERROR, label, "%s", message);
}

char* shell(const char* cmd);
*/
import "C"
import (
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func LogI(label, message string) {
	cLabel := C.CString(label)
	defer C.free(unsafe.Pointer(cLabel))
	cMessage := C.CString(message)
	defer C.free(unsafe.Pointer(cMessage))
	C.logI(cLabel, cMessage)
}

func LogE(label, message string) {
	cLabel := C.CString(label)
	defer C.free(unsafe.Pointer(cLabel))
	cMessage := C.CString(message)
	defer C.free(unsafe.Pointer(cMessage))
	C.logE(cLabel, cMessage)
}

func Shell(cmd string) string {
	cCmd := C.CString(cmd)
	defer C.free(unsafe.Pointer(cCmd))

	cResult := C.shell(cCmd)
	defer C.free(unsafe.Pointer(cResult))

	return C.GoString(cResult)
}

func Random(min, max int) int {
	if min > max {
		min, max = max, min // 如果 min 大于 max，则交换两者
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func Sleep(i int) {
	time.Sleep(time.Duration(i) * time.Millisecond)
}

func I2s(i int) string {
	return strconv.Itoa(i)
}

func S2i(s string) int {
	i, _ := strconv.Atoi(strings.TrimSpace(s))
	return i
}

// F2s 将浮点数转换为字符串。
func F2s(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// S2f 将字符串转换为浮点数。如果转换失败返回0.0。
func S2f(s string) float64 {
	f, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return f
}

// B2s 将布尔值转换为字符串 ("true" 或 "false")。
func B2s(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// S2b 将字符串转换为布尔值。如果无法转换则返回 false。
func S2b(s string) bool {
	b, _ := strconv.ParseBool(strings.TrimSpace(s))
	return b
}
