package app

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
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
	called := false
	if len(os.Args) == 1 {
		a.PrintHelpMessage()
	} else {
		user_called_command := os.Args[1]
		for _, option := range Goup.Options {
			for _, desc := range option.Command {
				command := desc.Description[0]
				if command == user_called_command {
					desc.Fn()
					called = true
					break
				}
			}
		}
		if !called {
			a.NoSupport(user_called_command)
		}
	}
}

func init() {
	Goup.Name = "Goup"
	Goup.AppVersion = "1.0.0"
	Goup.AppBuildVerion = runtime.Version()
	Goup.GoVersion = Go.Version
	Goup.Description = "Goup is Golang toolchain installer"

	usage := []command{
		{Description: []string{"goup.exe [OPTIONS]"}},
	}

	options := []command{
		{Description: []string{"help", "Print this information"}, Fn: Goup.PrintHelpMessage},
		{Description: []string{"update", "Update golang stable version"}, Fn: Go.CheckUpdate},
		{Description: []string{"version", "Print version information"}, Fn: Goup.GetVersion},
		{Description: []string{"install", "Install goup into Golang's system path (need root permisson)"}, Fn: Goup.Install},
		{Description: []string{"uninstall", "Remove goup from Golang's system path (need root permisson)"}, Fn: Goup.Uninstall},
	}

	Goup.SetOptions("USAGE:", usage)
	Goup.SetOptions("OPTIONS:", options)
}

func (a *app) Uninstall() {
	baseName := filepath.Base(os.Args[0])
	pathName := filepath.Join(Go.Path, baseName)
	_, err := os.Stat(pathName)
	if os.IsNotExist(err) {
		log.Println("Goup not install")
	} else {
		if err := os.Remove(pathName); err == nil {
			log.Println("Uninstall success")
		} else {
			log.Println("Uninstall fail:", err)
		}
	}
}

func (a *app) Install() {
	if Go.IsInstall {
		baseName := filepath.Base(os.Args[0])
		pathName := filepath.Join(Go.Path, baseName)
		srcFile, err := os.Open(os.Args[0])
		if err != nil {
			log.Fatalln(err)
		}
		defer srcFile.Close()
		destFile, err := os.Create(pathName)
		if err != nil {
			log.Fatalln(err)
		}
		defer destFile.Close()
		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Install success")
	} else {
		log.Println("Golang 环境变量不存在或未配置，无法安装 Goup。")
	}
}

func (a *app) NoSupport(fn string) {
	log.Println(("No support this command: " + fn))
}

func (a *app) SetOptions(name string, commands []command) {
	var options options
	options.Name = name
	options.Command = commands
	a.Options = append(a.Options, options)
}

func (a *app) PrintHelpMessage() {
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
	fmt.Println("Goup", a.AppVersion)
	fmt.Println("Build with", a.AppBuildVerion)
}
