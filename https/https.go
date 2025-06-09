package https

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

func Get(url string, timeout int) (code int, data []byte) {
	var client http.Client

	// 如果 timeout 不为 0，设置超时，否则使用默认的 client（没有超时）
	if timeout > 0 {
		client = http.Client{
			Timeout: time.Duration(timeout) * time.Millisecond,
		}
	} else {
		client = http.Client{}
	}

	resp, err := client.Get(url)
	if err != nil {
		return 0, nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, body
}

func PostMultipart(url string, fileName string, fileData []byte) (code int, data []byte) {
	// 创建一个缓冲区用于存储 multipart 数据
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// 创建一个 form 文件字段
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return 0, nil
	}

	// 将文件数据写入 part
	if _, err = io.Copy(part, bytes.NewReader(fileData)); err != nil {
		return 0, nil
	}

	// 关闭 multipart writer，设置结束标志
	_ = writer.Close()

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", url, &buffer)
	if err != nil {
		return 0, nil
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil
	}

	return resp.StatusCode, body
}
