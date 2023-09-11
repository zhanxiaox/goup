package app

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

type golang struct {
	Version     string
	Arch        string
	Os          string
	Domain      string
	DownloadUrl string
}

var Golang golang

func init() {
	Golang.Version = "go0.0.0"
	Golang.Arch = runtime.GOARCH
	Golang.Os = runtime.GOOS
	Golang.Domain = "https://golang.google.cn"
	Golang.DownloadUrl = Golang.Domain + "/dl"
	Golang.GetVersion()
}

func (g *golang) GetVersion() {
	cmd := exec.Command("go", "version")
	if output, err := cmd.Output(); err == nil {
		ver_string := string(output)
		ver_string = strings.ReplaceAll(ver_string, "\n", "")
		ver_split := (strings.Split(ver_string, " "))
		if len(ver_split) == 4 {
			g.Version = ver_split[2]
			ver_os_arch := strings.Split(ver_split[3], "/")
			if len(ver_os_arch) == 2 {
				g.Os = ver_os_arch[0]
				g.Arch = ver_os_arch[1]
			}
		}
	}
}

func getStable(body_string string) (string, error) {
	// get stable block
	reg_stable_block := `id="featured".*?id="stable"`
	re := regexp.MustCompile(reg_stable_block)
	stable_block := re.FindString(body_string)
	// get stable download url
	reg_stable := `/dl/(go\d+\.\d+\.\d+)\.` + Golang.Os + `-` + Golang.Arch + `\.{1}\w+`
	re = regexp.MustCompile(reg_stable)
	stable_download_url := re.FindStringSubmatch(stable_block)
	var err error
	var stable_url = ""
	if stable_download_url[1] == Golang.Version {
		err = errors.New("已是最新: " + Golang.Version)
	} else {
		err = nil
		stable_url = Golang.Domain + stable_download_url[0]
	}
	return stable_url, err
}

func CheckUpdate() {
	if url, err := getStable(GetHtml(Golang.DownloadUrl)); err == nil {
		if filepath, err := DownloadFile(url); err == nil {
			Installer(filepath)
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}
