package app

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"syscall"
)

type app struct {
	Name           string
	AppVersion     string
	AppBuildVerion string
	GoVersion      string
	Description    string
	Options        []options
}

type options struct {
	Name    string
	Command []command
}

type command struct {
	Description []string
	fn          func()
}

var App app

func init() {
	App.Name = "goup"
	App.AppVersion = "0.0.1"
	App.AppBuildVerion = runtime.Version()
	App.GoVersion = Golang.Version
	App.Description = "Goup is Golang stable installer"

	usage := []command{
		{Description: []string{"goup.exe [OPTIONS]"}},
	}

	options := []command{
		{Description: []string{"help", "Print this information"}},
		{Description: []string{"update", "Update golang stable version"}},
		{Description: []string{"version", "Print version information"}},
		{Description: []string{"install", "Install goup into Golang's system path"}},
		{Description: []string{"uninstall", "Remove goup from Golang's system path"}},
	}

	App.SetOptions("USAGE:", usage)
	App.SetOptions("OPTIONS:", options)
}

func (a *app) Install() {
	re := regexp.MustCompile(`[^;]+\\Go\\bin|[^;]+\\go\\bin`)
	match := re.FindString(os.Getenv("path"))
	if match == "" {
		fmt.Println("未找到 Golang 系统变量路径，确定已安装最新版本")
	} else {
	}
}

func (a *app) SetOptions(name string, commands []command) {
	var options options
	options.Name = name
	options.Command = commands
	a.Options = append(a.Options, options)
}

func (a *app) Print() {
	fmt.Println(a.Name, a.AppVersion)
	fmt.Println(a.Description)
	fmt.Println()
	for _, v := range a.Options {
		fmt.Println(v.Name)
		for _, v := range v.Command {
			fmt.Print("- ")
			for _, v := range v.Description {
				fmt.Printf("%-15s", v+" ")
			}
			fmt.Println()
		}
		fmt.Println()
	}
}

func (a *app) GetVersion() {
	fmt.Println("goup", a.AppVersion, "( build in", a.AppBuildVerion, ")")
	fmt.Println("Go", Golang.Version)
}

func Installer(filePath string) {
	if runtime.GOOS == "windows" {
		msi(filePath)
	} else {
		fmt.Println("no support")
	}
}

func msi(filePath string) {
	cmd := exec.Command("msiexec", "/i", filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Run()
	if err != nil {
		fmt.Println("MSI文件执行失败:", err)
		return
	}
	fmt.Println("MSI文件执行完成")
}

func pkg(filePath string) {}

func gz(filePath string) {}
