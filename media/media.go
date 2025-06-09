package media

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	"strings"
	"unsafe"
)

// ScanFile 扫描路径path的媒体文件，将它加入媒体库中
func ScanFile(path string) {
	shell("am broadcast -a android.intent.action.MEDIA_SCANNER_SCAN_FILE -d file://" + path)
}

func shell(cmd string) {
	if strings.Contains(cmd, ";") {
		cmd = "(" + cmd + ")"
	}
	cCmd := C.CString(cmd + " > /dev/null 2>&1")
	defer C.free(unsafe.Pointer(cCmd))
	C.system(cCmd)
}
