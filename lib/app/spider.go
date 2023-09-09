package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

func GetHtml(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	body_string := string(body)
	body_string = strings.ReplaceAll(body_string, "\n", "")
	body_string = strings.ReplaceAll(body_string, "\r", "")
	return body_string
}

func DownloadFile(url string) (string, error) {
	filePath := path.Base(url)
	// 发送GET请求
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	// 将响应体内容保存到文件，并更新进度条
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Println("文件下载完成:", filePath)
	return filePath, nil
}
