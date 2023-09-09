package app

import (
	"errors"
	"fmt"
	"regexp"
	"runtime"
)

func GetGoVersion() string {
	return runtime.Version()
}

const domain = "https://golang.google.cn"
const downloadUrl = domain + "/dl"

func getStable(body_string string) (string, error) {
	// get stable block
	reg_stable_block := `id="featured".*?id="stable"`
	re := regexp.MustCompile(reg_stable_block)
	stable_block := re.FindString(body_string)
	// get stable download url
	reg_stable := `/dl/(go\d+\.\d+\.\d+)\.` + runtime.GOOS + `-` + runtime.GOARCH + `\.{1}\w+`
	re = regexp.MustCompile(reg_stable)
	stable_download_url := re.FindStringSubmatch(stable_block)
	var err error
	var stable_url = ""
	if stable_download_url[1] == runtime.Version() {
		err = errors.New("已是最新")
	} else {
		err = nil
		stable_url = domain + stable_download_url[0]
	}
	return stable_url, err
}

func CheckUpdate() {
	if url, err := getStable(GetHtml(downloadUrl)); err == nil {
		if filepath, err := DownloadFile(url); err == nil {
			Installer(filepath)
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}
