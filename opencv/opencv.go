package opencv

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/Dasongzi1366/AutoGo/images"
)

var templateMap = make(map[string]Mat)
var maskMap = make(map[string]Mat)

// FindImage 在指定区域内查找匹配的图片模板，支持透明图像处理。
//
// 参数：
//   - x1, y1: 区域左上角的坐标。
//   - x2, y2: 区域右下角的坐标。当 x2 或 y2 为 0 时，表示使用图像的最大宽度或高度。
//   - template: 模板图片的字节数组指针，表示要在区域内查找的图片。
//   - isGray: 布尔值，指示是否将图像转换为灰度图进行匹配，提升匹配速度和鲁棒性。
//   - scalingFactor: 缩放比例，应用于模板图片和截取的待匹配区域，确保模板能够适应不同分辨率的图像。
//   - 0.5 表示缩小为原图的 50%。
//   - 2.0 表示放大为原图的 200%。
//   - sim: 相似度阈值，取值范围为 0.1 到 1.0，值越高表示匹配要求越精确。
//
// 返回值：
//   - (int, int): 返回找到的图片左上角坐标。如果未找到则返回 (-1, -1)。
//
// 透明图支持：
//   - 当模板图片为透明图时，函数会自动生成遮罩来忽略透明区域。
//   - 透明图的判定标准：图像的四个角颜色相同，且透明像素数量占总像素数的 30% 到 99%。
//
// 缩放说明：
//   - scalingFactor 控制模板和待匹配区域的缩放比例，通过对模板图像的缩放，实现在不同分辨率的屏幕上进行一致的匹配。
//   - 当模板图像在原分辨率（例如 540x960）下创建，但需要在更高或更低分辨率（例如 1080x1920 或 270x480）上查找时，
//     通过调整 scalingFactor 使模板与不同分辨率的目标图像比例一致。
//   - 例如：scalingFactor 为 2.0 时，模板会放大 200%，适用于将 540x960 的模板匹配到 1080x1920 的图像中。
//     scalingFactor 为 0.5 时，模板会缩小 50%，适用于将 540x960 的模板匹配到 270x480 的图像中。
//   - 若 scalingFactor 为 1.0，则模板和截取区域保持原尺寸，不进行缩放。
func FindImage(x1, y1, x2, y2 int, template *[]byte, isGray bool, scalingFactor, sim float32) (int, int) {
	if scalingFactor < 0.1 {
		scalingFactor = 1
	}
	mat2, mat3 := byte2mat(template, isGray, scalingFactor)
	if mat2.Empty() {
		return -1, -1
	}

	img := images.CaptureScreen(x1, y1, x2, y2)
	if img == nil {
		return -1, -1
	}

	bounds := img.Bounds()
	mat1, err := NewMatFromBytes(bounds.Dy(), bounds.Dx(), MatTypeCV8UC4, img.Pix)
	defer mat1.Close()
	if err != nil {
		return -1, -1
	}

	if isGray {
		mat1 = matGray(mat1)
	}

	result := NewMat()
	defer result.Close()

	MatchTemplate(mat1, mat2, &result, TmCcoeffNormed, mat3)

	_, maxVal, _, maxLoc := MinMaxLoc(result)
	if maxVal >= 0.5+sim*0.5 {
		return int(float32(maxLoc.X)/scalingFactor) + x1, int(float32(maxLoc.Y)/scalingFactor) + y1
	}

	return -1, -1
}

func byte2mat(pngData *[]byte, isGray bool, scale float32) (Mat, Mat) {
	sign := fmt.Sprintf("%p-%t-%.2f", pngData, isGray, scale)

	if cachedMat, ok := templateMap[sign]; ok {
		return cachedMat, maskMap[sign]
	}

	img, _, err := image.Decode(bytes.NewReader(*pngData))
	if err != nil {
		fmt.Println("图像解码失败")
		return NewMat(), NewMat()
	}
	imgNrgba := images.ImageToNRGBA(img)

	bounds := imgNrgba.Bounds()
	templateMat, _ := NewMatFromBytes(bounds.Dy(), bounds.Dx(), MatTypeCV8UC4, imgNrgba.Pix)

	isTransparent := checkTransparent(imgNrgba)

	if isGray {
		templateMat = matGray(templateMat)
	}
	templateMat = matScale(templateMat, scale)

	var maskMat Mat
	if isTransparent {
		maskMat = createMask(imgNrgba)
	} else {
		maskMat = NewMat()
	}

	templateMap[sign] = templateMat
	maskMap[sign] = maskMat

	return templateMap[sign], maskMap[sign]
}

func matGray(mat Mat) Mat {
	grayMat := NewMat()
	CvtColor(mat, &grayMat, ColorBGRToGray)
	_ = mat.Close()
	return grayMat
}

func matScale(mat Mat, scale float32) Mat {
	const epsilon = 1e-6
	if math.Abs(float64(scale-1)) < epsilon {
		return mat
	}
	scaledMat := NewMat()
	Resize(mat, &scaledMat, image.Point{X: int(float32(mat.Cols()) * scale), Y: int(float32(mat.Rows()) * scale)}, 0, 0, InterpolationLinear)
	_ = mat.Close()
	return scaledMat
}

// 判断是否是透明图
func checkTransparent(img image.Image) bool {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	if width < 2 || height < 2 {
		return false
	}

	c0 := getRGB(img.At(0, 0))
	c1 := getRGB(img.At(width-1, 0))
	c2 := getRGB(img.At(0, height-1))
	c3 := getRGB(img.At(width-1, height-1))

	if c0 != c1 || c0 != c2 || c0 != c3 {
		return false
	}

	transparentCount := 0
	totalPixels := width * height
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if getRGB(img.At(x, y)) == c0 {
				transparentCount++
			}
		}
	}

	if transparentCount >= int(float32(totalPixels)*0.3) && transparentCount < totalPixels {
		return true
	}

	return false
}

// 创建透明图遮罩
func createMask(img image.Image) Mat {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	c0 := getRGB(img.At(0, 0))

	mask := NewMatWithSize(height, width, MatTypeCV8U)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if getRGB(img.At(x, y)) == c0 {
				mask.SetUCharAt(y, x, 1)
			} else {
				mask.SetUCharAt(y, x, 0)
			}
		}
	}

	return mask
}

func getRGB(c color.Color) color.RGBA {
	r, g, b, _ := c.RGBA() // 忽略 Alpha 通道
	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 255}
}
