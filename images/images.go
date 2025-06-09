package images

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Dasongzi1366/AutoGo/java"
	"github.com/disintegration/imaging"
	"golang.org/x/image/bmp"
	_ "golang.org/x/image/bmp"
)

type colorInfo struct {
	x, y   int
	c2, c3 color.NRGBA
}

// opRange2D 表示遍历的范围和步进方向
type opRange2D struct {
	x1, x2, y1, y2, stepX, stepY int
}

var capimg *image.NRGBA
var mutex sync.Mutex

func init() {
	java.Send("screenShotInit")
	go readImage()
	for {
		sleep(100)
		mutex.Lock()
		if capimg != nil {
			mutex.Unlock()
			return
		}
		mutex.Unlock()
	}
}

func readImage() {
	for {
		data, err := java.GetBitMapData()
		if err != nil {
			fmt.Fprintf(os.Stderr, "[images] 读取数据出错: %v\n", err)
			os.Exit(1)
		}
		if len(data) > 1000000 {
			width := int(binary.LittleEndian.Uint32(data[0:4]))
			height := int(binary.LittleEndian.Uint32(data[4:8]))
			expectedImageDataLen := width*height*4 + 8

			if len(data) == expectedImageDataLen {
				mutex.Lock()
				capimg = image.NewNRGBA(image.Rect(0, 0, width, height))
				capimg.Pix = data
				mutex.Unlock()
			} else {
				fmt.Fprintln(os.Stderr, "[images] 解析图像数据错误:"+i2s(expectedImageDataLen)+", got "+i2s(len(data)))
				os.Exit(1)
			}
		}
		sleep(30)
	}
}

func Pixel(x, y int) string {
	mutex.Lock()
	c := capimg.At(x, y).(color.NRGBA)
	mutex.Unlock()
	return fmt.Sprintf("%02X%02X%02X", c.R, c.G, c.B)
}

func CaptureScreen(x1, y1, x2, y2 int) *image.NRGBA {
	mutex.Lock()
	defer mutex.Unlock()
	width1 := capimg.Rect.Dx()
	height1 := capimg.Rect.Dy()

	if x2 == 0 {
		x2 = width1
	}
	if y2 == 0 {
		y2 = height1
	}

	if x1 < 0 || y1 < 0 || x2 > width1 || y2 > height1 || x1 >= x2 || y1 >= y2 {
		return nil
	}

	if x1 == 0 && y1 == 0 && x2 == width1 && y2 == height1 {
		img := image.NewNRGBA(capimg.Rect)
		copy(img.Pix, capimg.Pix)
		return img
	}
	width2 := x2 - x1
	height2 := y2 - y1
	img := image.NewNRGBA(image.Rect(0, 0, width2, height2))
	for y := 0; y < height2; y++ {
		srcOffset := (y1+y)*capimg.Stride + x1*4
		destOffset := y * img.Stride
		copy(img.Pix[destOffset:], capimg.Pix[srcOffset:srcOffset+width2*4])
	}
	return img
}

func CmpColor(x, y int, colorStr string, sim float32) bool {
	colorStr = strings.ReplaceAll(colorStr, "[", "")
	colorStr = strings.ReplaceAll(colorStr, "]", "")
	colorStr = strings.ReplaceAll(colorStr, "#", "")
	colorStr = strings.ReplaceAll(colorStr, " ", "")
	colorStr = strings.ReplaceAll(colorStr, `"`, "")
	mutex.Lock()
	c1 := capimg.At(x, y).(color.NRGBA)
	mutex.Unlock()

	arr := strings.Split(colorStr, "|")
	for _, str := range arr {
		c2, c3 := str2color(str, sim)
		if isColorMatch(c1, c2, c3) {
			return true
		}
	}

	return false
}

func FindColor(x1, y1, x2, y2 int, colorStr string, sim float32, dir int) (int, int) {
	img := CaptureScreen(x1, y1, x2, y2)
	if img == nil {
		return -1, -1
	}

	width := img.Rect.Dx()
	height := img.Rect.Dy()

	var c2s, c3s []color.NRGBA
	arr := strings.Split(colorStr, "|")
	for _, str := range arr {
		c2, c3 := str2color(str, sim)
		c2s = append(c2s, c2)
		c3s = append(c3s, c3)
	}

	// 生成遍历的范围和方向
	rng := genRange(dir, 0, 0, width, height)

	for i := rng.y1; i != rng.y2; i += rng.stepY {
		for j := rng.x1; j != rng.x2; j += rng.stepX {
			c1 := getPixFromData(&img.Pix, width, j, i)
			for t := 0; t < len(c2s); t++ {
				if isColorMatch(c1, c2s[t], c3s[t]) {
					return j + x1, i + y1
				}
			}
		}
	}

	return -1, -1
}

func GetColorCountInRegion(x1, y1, x2, y2 int, colorStr string, sim float32) int {
	colorStr = strings.ReplaceAll(colorStr, "[", "")
	colorStr = strings.ReplaceAll(colorStr, "]", "")
	colorStr = strings.ReplaceAll(colorStr, "#", "")
	colorStr = strings.ReplaceAll(colorStr, " ", "")
	colorStr = strings.ReplaceAll(colorStr, `"`, "")
	img := CaptureScreen(x1, y1, x2, y2)
	if img == nil {
		return 0
	}
	width := img.Rect.Dx()
	height := img.Rect.Dy()
	var c2s, c3s []color.NRGBA
	arr := strings.Split(colorStr, "|")
	for _, str := range arr {
		c2, c3 := str2color(str, sim)
		c2s = append(c2s, c2)
		c3s = append(c3s, c3)
	}
	s := 0
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			c1 := getPixFromData(&img.Pix, width, j, i)
			for t := 0; t < len(c2s); t++ {
				if isColorMatch(c1, c2s[t], c3s[t]) {
					s++
					break
				}
			}
		}
	}
	return s
}

// 因为截图和综合工具存在x轴两个像素的偏差所以暂时这样调用多遍历几个像素
func DetectsMultiColors(colors string, sim float32) bool {
	colors = strings.ReplaceAll(colors, "[", "")
	colors = strings.ReplaceAll(colors, "]", "")
	colors = strings.ReplaceAll(colors, "#", "")
	colors = strings.ReplaceAll(colors, " ", "")
	colors = strings.ReplaceAll(colors, `"`, "")
	// 解析颜色字符串
	arr := strings.Split(colors, ",")
	if len(arr) < 6 || len(arr)%3 != 0 {
		return false
	}
	x, _ := strconv.Atoi(arr[0])
	str := strings.Join(arr[1:], ",")
	for i := 0; i < 4; i++ {
		if detectsMultiColors2(i2s(x+i)+","+str, sim) {
			return true
		}
	}
	return false
}

func detectsMultiColors2(colors string, sim float32) bool {
	colors = strings.ReplaceAll(colors, "[", "")
	colors = strings.ReplaceAll(colors, "]", "")
	colors = strings.ReplaceAll(colors, "#", "")
	colors = strings.ReplaceAll(colors, " ", "")
	colors = strings.ReplaceAll(colors, `"`, "")

	// 解析颜色字符串
	arr := strings.Split(colors, ",")
	if len(arr) < 6 || len(arr)%3 != 0 {
		return false
	}

	x, _ := strconv.Atoi(arr[0])
	y, _ := strconv.Atoi(arr[1])
	c2, c3 := str2color(arr[2], sim)

	mutex.Lock()
	c1 := capimg.At(x, y).(color.NRGBA)
	mutex.Unlock()
	if !isColorMatch(c1, c2, c3) {
		return false
	}
	for i := 3; i < len(arr); i += 3 {
		offsetX, _ := strconv.Atoi(arr[i])
		offsetY, _ := strconv.Atoi(arr[i+1])
		c2, c3 := str2color(arr[i+2], sim)
		mutex.Lock()
		c1 := capimg.At(x+offsetX, y+offsetY).(color.NRGBA)
		mutex.Unlock()
		if !isColorMatch(c1, c2, c3) {
			return false
		}
	}
	return true
}

func FindMultiColors(x1, y1, x2, y2 int, colors string, sim float32, dir int) (int, int) {
	colors = strings.ReplaceAll(colors, "[", "")
	colors = strings.ReplaceAll(colors, "]", "")
	colors = strings.ReplaceAll(colors, "#", "")
	colors = strings.ReplaceAll(colors, " ", "")
	colors = strings.ReplaceAll(colors, `"`, "")

	arr := strings.Split(colors, ",")
	if len(arr) < 4 || len(arr)%3 != 1 {
		return -1, -1
	}

	baseColor, baseTolerance := str2color(arr[0], sim)

	var infos []colorInfo

	img := CaptureScreen(x1, y1, x2, y2)
	if img == nil {
		return -1, -1
	}

	width := img.Rect.Dx()
	height := img.Rect.Dy()
	rng := genRange(dir, 0, 0, width, height)

	for i := rng.y1; i != rng.y2; i += rng.stepY {
		for j := rng.x1; j != rng.x2; j += rng.stepX {
			c1 := getPixFromData(&img.Pix, width, j, i)
			if isColorMatch(c1, baseColor, baseTolerance) {
				if infos == nil {
					infos = parseRemainingColors(arr[1:], sim)
				}
				if compareColorsInSequence(img, j, i, infos) {
					return j + x1, i + y1
				}
			}
		}
	}
	return -1, -1
}

func ReadFromPath(path string) *image.NRGBA {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil
	}
	return ImageToNRGBA(img)
}

func ImageToNRGBA(img image.Image) *image.NRGBA {
	bounds := img.Bounds()
	nrgbaImg := image.NewNRGBA(bounds)
	var wg sync.WaitGroup
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				srcColor := img.At(x, y)
				r, g, b, a := srcColor.RGBA()
				i := (y-bounds.Min.Y)*nrgbaImg.Stride + (x-bounds.Min.X)*4
				nrgbaImg.Pix[i] = uint8(r >> 8)
				nrgbaImg.Pix[i+1] = uint8(g >> 8)
				nrgbaImg.Pix[i+2] = uint8(b >> 8)
				nrgbaImg.Pix[i+3] = uint8(a >> 8)
			}
		}(y)
	}
	wg.Wait()
	return nrgbaImg
}

func ReadFromUrl(url string) *image.NRGBA {
	client := http.Client{
		Timeout: time.Duration(5000) * time.Millisecond,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil
	}
	return ImageToNRGBA(img)
}

func ReadFromBase64(base64Str string) *image.NRGBA {
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil
	}
	reader := bytes.NewReader(data)
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil
	}
	return ImageToNRGBA(img)
}

func ReadFromBytes(data []byte) *image.NRGBA {
	reader := bytes.NewReader(data)
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil
	}
	return ImageToNRGBA(img)
}

func Save(img *image.NRGBA, path string, quality int) bool {
	ext := strings.ToLower(filepath.Ext(path))
	ext = strings.ReplaceAll(ext, ".", "")
	data := EncodeToBytes(img, ext, quality)
	if data != nil {
		return os.WriteFile(path, data, 0644) == nil
	}
	return false
}

func EncodeToBase64(img *image.NRGBA, fromat string, quality int) string {
	data := EncodeToBytes(img, fromat, quality)
	if data != nil {
		return base64.StdEncoding.EncodeToString(data)
	}
	return ""
}

func EncodeToBytes(img *image.NRGBA, fromat string, quality int) []byte {
	var err error
	var buf bytes.Buffer
	fromat = strings.ToLower(fromat)
	switch fromat {
	case "png":
		err = png.Encode(&buf, img)
	case "jpg", "jpeg":
		options := &jpeg.Options{Quality: quality}
		err = jpeg.Encode(&buf, img, options)
	case "bmp":
		err = bmp.Encode(&buf, img)
	default:
		err = os.ErrInvalid
	}
	if err == nil {
		return buf.Bytes()
	}
	return nil
}

func Clip(img *image.NRGBA, x1, y1, x2, y2 int) *image.NRGBA {
	bounds := img.Bounds()
	if x1 < bounds.Min.X {
		x1 = bounds.Min.X
	}
	if y1 < bounds.Min.Y {
		y1 = bounds.Min.Y
	}
	if x2 > bounds.Max.X || x2 == 0 {
		x2 = bounds.Max.X
	}
	if y2 > bounds.Max.Y || y2 == 0 {
		y2 = bounds.Max.Y
	}
	width := x2 - x1
	height := y2 - y1
	newImg := image.NewNRGBA(image.Rect(0, 0, width, height))
	var wg sync.WaitGroup
	for y := 0; y < height; y++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			for x := 0; x < width; x++ {
				srcIdx := (y1+y)*img.Stride + (x1+x)*4
				dstIdx := y*newImg.Stride + x*4
				copy(newImg.Pix[dstIdx:dstIdx+4], img.Pix[srcIdx:srcIdx+4])
			}
		}(y)
	}
	wg.Wait()
	return newImg
}

func Resize(img *image.NRGBA, width, height int) *image.NRGBA {
	resizedImg := imaging.Resize(img, width, height, imaging.Lanczos)
	resizedNRGBA := image.NewNRGBA(resizedImg.Bounds())
	copy(resizedNRGBA.Pix, resizedImg.Pix)
	return resizedNRGBA
}

func Rotate(img *image.NRGBA, degree int) *image.NRGBA {
	rotatedImg := imaging.Rotate(img, float64(-degree), color.Transparent)
	return rotatedImg
}

func Grayscale(img *image.NRGBA) *image.Gray {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)
	var wg sync.WaitGroup
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				originalColor := img.At(x, y)
				grayColor := color.GrayModel.Convert(originalColor)
				grayImg.Set(x, y, grayColor)
			}
		}(y)
	}
	wg.Wait()
	return grayImg
}

func ApplyThreshold(img *image.NRGBA, threshold, maxVal int, typ string) *image.Gray {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			pixel := originalColor.Y
			var newPixel uint8
			switch typ {
			case "BINARY":
				if int(pixel) > threshold {
					newPixel = uint8(maxVal)
				} else {
					newPixel = 0
				}
			case "BINARY_INV":
				if int(pixel) > threshold {
					newPixel = 0
				} else {
					newPixel = uint8(maxVal)
				}
			case "TRUNC":
				if int(pixel) > threshold {
					newPixel = uint8(threshold)
				} else {
					newPixel = pixel
				}
			case "TOZERO":
				if int(pixel) > threshold {
					newPixel = pixel
				} else {
					newPixel = 0
				}
			case "TOZERO_INV":
				if int(pixel) > threshold {
					newPixel = 0
				} else {
					newPixel = pixel
				}
			default:
				newPixel = pixel
			}
			grayImg.Set(x, y, color.Gray{Y: newPixel})
		}
	}
	return grayImg
}

func ApplyAdaptiveThreshold(img *image.NRGBA, maxValue float64, adaptiveMethod string, thresholdType string, blockSize int, C float64) *image.Gray {
	bounds := img.Bounds()
	grayImg := Grayscale(img) // 先将图像灰度化
	dstImg := image.NewGray(bounds)
	var wg sync.WaitGroup
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				var sum, sumWeight, mean float64
				for j := -blockSize / 2; j <= blockSize/2; j++ {
					for i := -blockSize / 2; i <= blockSize/2; i++ {
						xx := x + i
						yy := y + j
						if xx >= bounds.Min.X && xx < bounds.Max.X && yy >= bounds.Min.Y && yy < bounds.Max.Y {
							pixel := grayImg.GrayAt(xx, yy).Y
							if adaptiveMethod == "GAUSSIAN_C" {
								weight := math.Exp(-(float64(i*i + j*j)) / (2.0 * float64(blockSize*blockSize)))
								sum += float64(pixel) * weight
								sumWeight += weight
							} else {
								sum += float64(pixel)
								sumWeight += 1.0
							}
						}
					}
				}

				if adaptiveMethod == "GAUSSIAN_C" {
					mean = sum / sumWeight
				} else {
					mean = sum / (float64(blockSize) * float64(blockSize))
				}

				threshold := mean - C
				pixel := grayImg.GrayAt(x, y).Y

				var newPixel uint8
				switch thresholdType {
				case "BINARY":
					if float64(pixel) > threshold {
						newPixel = uint8(maxValue)
					} else {
						newPixel = 0
					}
				case "BINARY_INV":
					if float64(pixel) > threshold {
						newPixel = 0
					} else {
						newPixel = uint8(maxValue)
					}
				}

				dstImg.SetGray(x, y, color.Gray{Y: newPixel})
			}
		}(y)
	}
	wg.Wait()
	return dstImg
}

func ApplyBinarization(img *image.NRGBA, threshold int) *image.Gray {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)
	var wg sync.WaitGroup
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				originalColor := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
				pixel := originalColor.Y
				var newPixel uint8
				if int(pixel) > threshold {
					newPixel = uint8(255)
				} else {
					newPixel = 0
				}
				grayImg.SetGray(x, y, color.Gray{Y: newPixel})
			}
		}(y)
	}
	wg.Wait()
	return grayImg
}

func compareColorsInSequence(img *image.NRGBA, x, y int, infos []colorInfo) bool {
	width := img.Rect.Dx()
	height := img.Rect.Dy()
	for _, info := range infos {
		offsetX := x + info.x
		offsetY := y + info.y
		if offsetX >= width || offsetY >= height {
			return false
		}
		c1 := img.At(offsetX, offsetY).(color.NRGBA)
		if !isColorMatch(c1, info.c2, info.c3) {
			return false
		}
	}
	return true
}

// 解析偏移和颜色为 colorInfo 切片
func parseRemainingColors(arr []string, sim float32) []colorInfo {
	var infos []colorInfo
	for i := 0; i < len(arr); i += 3 {
		offsetX, _ := strconv.Atoi(arr[i])
		offsetY, _ := strconv.Atoi(arr[i+1])
		c2, c3 := str2color(arr[i+2], sim)
		infos = append(infos, colorInfo{x: offsetX, y: offsetY, c2: c2, c3: c3})
	}
	return infos
}

func getPixFromData(pixData *[]byte, dataWidth, x, y int) color.NRGBA {
	idx := (y*dataWidth + x) * 4
	return color.NRGBA{
		R: (*pixData)[idx],
		G: (*pixData)[idx+1],
		B: (*pixData)[idx+2],
		A: (*pixData)[idx+3],
	}
}

// 判断两个颜色是否相似 基础颜色 要判断的颜色 偏色
func isColorMatch(c1, c2, c3 color.NRGBA) bool {
	return absDiff(c1.R, c2.R) <= c3.R && absDiff(c1.G, c2.G) <= c3.G && absDiff(c1.B, c2.B) <= c3.B
}

// 计算两个 uint8 之间的差值
func absDiff(a, b uint8) uint8 {
	if a > b {
		return a - b
	}
	return b - a
}

// 返回基础颜色和偏色,如果有相似度则返回的偏色为和相似度计算之后的偏色值
func str2color(colorStr string, sim float32) (color.NRGBA, color.NRGBA) {
	var tolerance uint8
	if sim > 0 {
		tolerance = uint8((1.0 - sim) * 255)
	}
	arr := strings.Split(colorStr, "-")
	s := arr[0]
	r, _ := strconv.ParseUint(s[0:2], 16, 8)
	g, _ := strconv.ParseUint(s[2:4], 16, 8)
	b, _ := strconv.ParseUint(s[4:6], 16, 8)
	color1 := color.NRGBA{uint8(r), uint8(g), uint8(b), 255}
	var color2 color.NRGBA
	if len(arr) > 1 {
		s := arr[1]
		r, _ := strconv.ParseUint(s[0:2], 16, 8)
		g, _ := strconv.ParseUint(s[2:4], 16, 8)
		b, _ := strconv.ParseUint(s[4:6], 16, 8)
		color2 = color.NRGBA{uint8(r) + tolerance, uint8(g) + tolerance, uint8(b) + tolerance, 255}
	} else {
		color2 = color.NRGBA{tolerance, tolerance, tolerance, 255}
	}
	return color1, color2
}

// 生成遍历的范围和步进方向
func genRange(dir, x1, y1, x2, y2 int) opRange2D {
	var rng opRange2D
	switch dir {
	case 0: // 从左到右，从上到下
		rng = opRange2D{x1, x2, y1, y2, 1, 1}
	case 1: // 从右到左，从上到下
		rng = opRange2D{x2 - 1, x1 - 1, y1, y2, -1, 1}
	case 2: // 从左到右，从下到上
		rng = opRange2D{x1, x2, y2 - 1, y1 - 1, 1, -1}
	case 3: // 从右到左，从下到上
		rng = opRange2D{x2 - 1, x1 - 1, y2 - 1, y1 - 1, -1, -1}
	}
	return rng
}

func sleep(i int) {
	time.Sleep(time.Duration(i) * time.Millisecond)
}

func i2s(i int) string {
	return strconv.Itoa(i)
}
