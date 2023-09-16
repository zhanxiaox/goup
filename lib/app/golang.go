package app

import (
	"errors"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"syscall"
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
	case <-time.After(5 * time.Second):
		log.Fatalln("域名验证超时")
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
	g.Os = runtime.GOOS
	g.Arch = runtime.GOARCH
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
	// Get stable block
	reg_stable_block := `id="featured".*?id="stable"`
	re := regexp.MustCompile(reg_stable_block)
	stable_block := re.FindString(body_string)
	// Get stable download url
	reg_stable := `/dl/(go\d+\.\d+\.\d+)\.` + Go.Os + `-` + Go.Arch + `\.{1}\w+`
	re = regexp.MustCompile(reg_stable)
	stable_download_url := re.FindStringSubmatch(stable_block)
	var err error
	var stable_url = ""
	if stable_download_url[1] == Go.Version && false {
		err = errors.New("You have latest version: " + Go.Version)
	} else {
		err = nil
		stable_url = Go.Domain + stable_download_url[0]
	}
	return stable_url, err
}

func (g *golang) CheckUpdate() {
	log.Println("Syncing channel update for " + g.Os + "/" + g.Arch)
	if url, err := GetStable(GetHtml(g.DownloadUrl)); err == nil {
		if filepath, err := DownloadFile(url); err == nil {
			g.Installer(filepath)
		} else {
			log.Fatalln(err.Error())
		}
	} else {
		log.Fatalln(err.Error())
	}
}

func (g *golang) Installer(filePath string) {
	if runtime.GOOS == "windows" {
		msi_installer(filePath)
	} else {
		log.Println("No support:", runtime.GOOS)
	}
}

func msi_installer(filePath string) {
	log.Println("Installing...")
	cmd := exec.Command("msiexec", "/i", filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Run(); err != nil {
		log.Fatalln("MSI install fail:", err)
	}
	log.Println("MSI install success")
}

func pkg_install(filePath string) {}

func gz_install(filePath string) {}
