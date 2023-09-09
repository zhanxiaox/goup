package app

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

const Version = "0.0.1"

func GetVersion() {
	fmt.Println("Go", GetGoVersion())
	fmt.Println("goup", Version)
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
