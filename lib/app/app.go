package app

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

type app struct {
	Name        string
	Version     string
	GoVersion   string
	Description string
	Options     []options
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
	App.Version = "0.0.1"
	App.GoVersion = runtime.Version()
	App.Description = "Goup is Golang stable installer"

	usage := []command{
		{Description: []string{"goup.exe [OPTIONS]"}},
	}

	options := []command{
		{Description: []string{"help", "Print this information"}},
		{Description: []string{"update", "Update golang stable version"}},
		{Description: []string{"version", "Print version information"}},
	}

	App.SetOptions("USAGE:", usage)
	App.SetOptions("OPTIONS:", options)
}

func (a *app) SetOptions(name string, commands []command) {
	var options options
	options.Name = name
	options.Command = commands
	a.Options = append(a.Options, options)
}

func (a *app) Print() {
	fmt.Println(a.Name, a.Version)
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
	fmt.Println("goup", a.Version)
	fmt.Println("Go", a.GoVersion)
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
