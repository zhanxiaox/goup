package app

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type golang struct {
	Version     string
	Arch        string
	Os          string
	Domain      string
	DownloadUrl string
	IsInstall   bool
	Path        string
}

var Go golang

func init() {
	Go.IsInstall = false
	Go.GetDomain()
	Go.GetVersion()
	Go.GetSystemPath()
}

func (g *golang) GetDomain() {
	domains := []string{"https://golang.google.cn", "https://golang.google.com"}
	complete := make(chan string, 1)
	for _, domain := range domains {
		go func(domain string) {
			if _, err := http.Get(domain); err == nil {
				complete <- domain
			}
		}(domain)
		if len(complete) == 1 {
			break
		}
	}
	select {
	case g.Domain = <-complete:
		g.DownloadUrl = g.Domain + "/dl"
		fmt.Println("域名验证成功：", g.Domain)
	case <-time.After(5 * time.Second):
		fmt.Println("域名验证失败", domains)
	}
}

func (g *golang) GetSystemPath() {
	re := regexp.MustCompile(`[^;]+\\Go\\bin|[^;]+\\go\\bin`)
	match := re.FindString(os.Getenv("path"))
	if match != "" {
		g.Path = match
	}
}

func (g *golang) GetVersion() {
	cmd := exec.Command("go", "version")
	if output, err := cmd.Output(); err == nil {
		ver_string := strings.ReplaceAll(string(output), "\n", "")
		ver_split := strings.Split(ver_string, " ")
		if len(ver_split) == 4 {
			g.Version = ver_split[2]
			ver_os_arch := strings.Split(ver_split[3], "/")
			if len(ver_os_arch) == 2 {
				g.Os = ver_os_arch[0]
				g.Arch = ver_os_arch[1]
				g.IsInstall = true
			}
		}
	}
}

func GetStable(body_string string) (string, error) {
	// get stable block
	reg_stable_block := `id="featured".*?id="stable"`
	re := regexp.MustCompile(reg_stable_block)
	stable_block := re.FindString(body_string)
	// get stable download url
	reg_stable := `/dl/(go\d+\.\d+\.\d+)\.` + Go.Os + `-` + Go.Arch + `\.{1}\w+`
	re = regexp.MustCompile(reg_stable)
	stable_download_url := re.FindStringSubmatch(stable_block)
	var err error
	var stable_url = ""
	if stable_download_url[1] == Go.Version {
		err = errors.New("已是最新: " + Go.Version)
	} else {
		err = nil
		stable_url = Go.Domain + stable_download_url[0]
	}
	return stable_url, err
}

func (g *golang) CheckUpdate() {
	if url, err := GetStable(GetHtml(g.DownloadUrl)); err == nil {
		if filepath, err := DownloadFile(url); err == nil {
			Installer(filepath)
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}
