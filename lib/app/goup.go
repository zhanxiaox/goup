package app

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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
	Fn          func()
}

var Goup app

func (a *app) Run() {
	isCalled := false
	if len(os.Args) > 1 {
		fn := os.Args[1]
		for _, option := range Goup.Options {
			for _, desc := range option.Command {
				name := desc.Description[0]
				if name == fn {
					desc.Fn()
					isCalled = true
				}
			}
		}
		if !isCalled {
			a.NoSupport(fn)
		}
	} else {
		a.Print()
	}
}

func init() {
	Goup.Name = "goup"
	Goup.AppVersion = "0.0.1"
	Goup.AppBuildVerion = runtime.Version()
	Goup.GoVersion = Go.Version
	Goup.Description = "Goup is Golang toolchain installer"

	usage := []command{
		{Description: []string{"goup.exe [OPTIONS]"}},
	}

	options := []command{
		{Description: []string{"help", "Print this information"}, Fn: Goup.Print},
		{Description: []string{"update", "Update golang stable version"}, Fn: Go.CheckUpdate},
		{Description: []string{"version", "Print version information"}, Fn: Goup.GetVersion},
		{Description: []string{"install", "Install goup into Golang's system path"}, Fn: Goup.Install},
		{Description: []string{"uninstall", "Remove goup from Golang's system path"}, Fn: Goup.Uninstall},
	}

	Goup.SetOptions("USAGE:", usage)
	Goup.SetOptions("OPTIONS:", options)
}

func (a *app) Uninstall() {
	baseName := filepath.Base(os.Args[0])
	pathName := filepath.Join(Go.Path, baseName)
	_, err := os.Stat(pathName)
	if os.IsNotExist(err) {
		fmt.Printf("还未安装 goup")
	} else {
		if err := os.Remove(pathName); err == nil {
			fmt.Printf("uninstall success")
		} else {
			fmt.Println("uninstall fail:", err)
		}
	}
}

func (a *app) Install() {
	if Go.IsInstall {
		baseName := filepath.Base(os.Args[0])
		pathName := filepath.Join(Go.Path, baseName)

		srcFile, err := os.Open(os.Args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		defer srcFile.Close()

		destFile, err := os.Create(pathName)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Install success")

	} else {
		fmt.Println("还未安装 golang")
	}
}

func (a *app) NoSupport(fn string) {
	fmt.Println("No support this command", fn)
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
	fmt.Println("Go", Go.Version)
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
