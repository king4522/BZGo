package yolov5

/*
#include "yolov5.h"
#include <stdlib.h>
#cgo arm64 LDFLAGS: -L../../resources/libs/arm64-v8a -lyolov5
#cgo amd64 LDFLAGS: -L../../resources/libs/x86_64 -lyolov5
#cgo 386 LDFLAGS: -L../../resources/libs/x86 -lyolov5
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/Dasongzi1366/AutoGo/images"
	"strings"
	"unsafe"
)

type YoloV5 struct {
	pointer *C.YoloV5
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

func New(cpuThreadNum int, paramPath, binPath, labels string) *YoloV5 {
	pointer := C.newYoloV5()
	if pointer == nil {
		return nil
	}

	cParamPath := C.CString(paramPath)
	defer C.free(unsafe.Pointer(cParamPath))

	cBinPath := C.CString(binPath)
	defer C.free(unsafe.Pointer(cBinPath))

	cLabels := C.CString(labels)
	defer C.free(unsafe.Pointer(cLabels))

	result := C.loadModelYoloV5(pointer, cParamPath, cBinPath, cLabels, C.int(cpuThreadNum))
	goResult := C.GoString(result)
	C.free(unsafe.Pointer(result))
	if goResult != "OK" {
		fmt.Println("yolov5:" + goResult)
		return nil
	}
	return &YoloV5{pointer: pointer}
}

func (y *YoloV5) Detect(x1, y1, x2, y2 int) []DetectResult {
	img := images.CaptureScreen(x1, y1, x2, y2)
	if img == nil {
		return nil
	}

	result := C.detectYoloV5(
		y.pointer,
		(*C.uchar)(unsafe.Pointer(&img.Pix[0])),
		C.int(img.Rect.Dx()),
		C.int(img.Rect.Dy()),
		C.float(0.5),
		C.float(0.45),
		C.int(640),
	)
	goResult := C.GoString(result)
	C.free(unsafe.Pointer(result))
	fmt.Println(goResult)
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

func (y *YoloV5) Close() {
	C.closeYoloV5(y.pointer)
}
