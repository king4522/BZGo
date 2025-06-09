package ppocr

/*
#include "ppocr.h"
#include <stdlib.h>
#cgo arm64 LDFLAGS: -L../../resources/libs/arm64-v8a -lopencv_core -lopencv_imgproc -lppocr
#cgo amd64 LDFLAGS: -L../../resources/libs/x86_64 -lopencv_core -lopencv_imgproc -lppocr
#cgo 386 LDFLAGS: -L../../resources/libs/x86 -lopencv_core -lopencv_imgproc -lppocr
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/Dasongzi1366/AutoGo/images"
)

type PpOcr struct {
	pointer *C.Ppocr
}

type DetectResult struct {
	X       int     `json:"X"`
	Y       int     `json:"Y"`
	Width   int     `json:"宽"`
	Height  int     `json:"高"`
	Label   string  `json:"标签"`
	Score   float64 `json:"精度"`
	CenterX int     `json:"-"` //中心坐标X
	CenterY int     `json:"-"` //中心坐标Y
}

func New(cpuThreadNum int) *PpOcr {
	pointer := C.newPpocr()

	if pointer == nil {
		return nil
	}
	path := filepath.Dir(os.Args[0])

	cLabelPath := C.CString(path + "/assets/label.txt")
	defer C.free(unsafe.Pointer(cLabelPath))

	cDbParamPath := C.CString(path + "/assets/db.param")
	defer C.free(unsafe.Pointer(cDbParamPath))

	cDbBinPath := C.CString(path + "/assets/db.bin")
	defer C.free(unsafe.Pointer(cDbBinPath))

	cRecParamPath := C.CString(path + "/assets/rec.param")
	defer C.free(unsafe.Pointer(cRecParamPath))

	cRecBinPath := C.CString(path + "/assets/rec.bin")
	defer C.free(unsafe.Pointer(cRecBinPath))
	result := C.loadModelPpocr(pointer, cLabelPath, cDbParamPath, cDbBinPath, cRecParamPath, cRecBinPath, C.int(cpuThreadNum))
	goResult := C.GoString(result)
	C.free(unsafe.Pointer(result))
	if goResult != "OK" {
		fmt.Println("ppocr:" + goResult)
		return nil
	}
	return &PpOcr{pointer: pointer}
}

func (p *PpOcr) Detect(x1, y1, x2, y2 int) []DetectResult {
	img := images.CaptureScreen(x1, y1, x2, y2)
	if img == nil {
		return nil
	}

	result := C.detectPpocr(
		p.pointer,
		(*C.uchar)(unsafe.Pointer(&img.Pix[0])),
		C.int(img.Rect.Dx()),
		C.int(img.Rect.Dy()),
		C.float(0.45),
		C.float(0.5),
		C.int(640),
	)
	goResult := C.GoString(result)
	C.free(unsafe.Pointer(result))

	var results []DetectResult
	err := json.Unmarshal([]byte(goResult), &results)
	if err != nil {
		return nil
	}

	var validResults []DetectResult
	for _, result := range results {
		if strings.TrimSpace(result.Label) == "" {
			continue
		}
		result.X += x1
		result.Y += y1
		result.CenterX = result.X + result.Width/2
		result.CenterY = result.Y + result.Height/2
		validResults = append(validResults, result)
	}

	return validResults
}

func (p *PpOcr) Close() {
	C.closePpocr(p.pointer)
}
