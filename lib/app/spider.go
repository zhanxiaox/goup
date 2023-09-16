package app

import (
	"fmt"
	"io"
	"log"
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
	log.Println("Downloading...")
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

	// 创建缓冲区
	buffer := make([]byte, 1024)

	totalSize := resp.ContentLength
	downloadedSize := 0

	for {
		// 从响应体读取数据
		n, err := resp.Body.Read(buffer)
		if err != nil && err != io.EOF {
			log.Fatal("读取响应体失败:", err)
			break
		}

		// 写入文件
		_, err = file.Write(buffer[:n])
		if err != nil {
			log.Fatal("写入文件失败:", err)
			break
		}

		// 更新下载进度
		downloadedSize += n
		progress := float64(downloadedSize) / float64(totalSize) * 100

		// 生成进度条
		bar := strings.Repeat("█", int(progress/2))
		fmt.Printf("\r"+filePath+":[%-50s] %.2f%%", bar, progress)
		// 判断是否已经下载完成
		if downloadedSize >= int(totalSize) {
			fmt.Println()
			log.Println("Downloading complete")
			break
		}
	}

	// 将响应体内容保存到文件，并更新进度条
	// _, err = io.Copy(file, resp.Body)
	// if err != nil {
	// 	return "", err
	// }
	// log.Println("Downloading complete")
	return filePath, nil
}
