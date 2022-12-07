package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/MX-Steve/my-go-app/vip"
)

// go-bindata -o asset/asset.go -pkg asset  website/...
// go build -ldflags="-H windowsgui" -o localJob.exe

// .\make_version.bat ./ ./version.h
// 		可以生成最后一位版本号和短哈希值
// rsrc.exe -manifest main.manifest -o app.syso -ico three.ico
// go build -ldflags="-H windowsgui" -o app.exe

const (
	DAEMON  = "daemon"
	FOREVER = "forever"
)

func init() {
	file := "ktzapp.log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetPrefix("[ktzgoapp] ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	go JobManage()
}

func test() {
	log.Println("初始化")
	init_http_server()
}

func StripSlice(slice []string, element string) []string {
	for i := 0; i < len(slice); {
		if slice[i] == element && i != len(slice)-1 {
			slice = append(slice[:i], slice[i+1:]...)
		} else if slice[i] == element && i == len(slice)-1 {
			slice = slice[:i]
		} else {
			i++
		}
	}
	return slice
}

func SubProcess(args []string) *exec.Cmd {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[-] Error: %s\n", err)
	}
	return cmd
}

func main() {
	daemon := flag.Bool(DAEMON, false, "run in daemon")
	forever := flag.Bool(FOREVER, false, "run forever")
	flag.Parse()
	fmt.Printf("[*] PID: %d PPID: %d ARG: %s\n", os.Getpid(), os.Getppid(), os.Args)
	if *daemon {
		SubProcess(StripSlice(os.Args, "-"+DAEMON))
		fmt.Printf("[*] Daemon running in PID: %d PPID: %d\n", os.Getpid(), os.Getppid())
		os.Exit(0)
	} else if *forever {
		for {
			cmd := SubProcess(StripSlice(os.Args, "-"+FOREVER))
			fmt.Printf("[*] Forever running in PID: %d PPID: %d\n", os.Getpid(), os.Getppid())
			cmd.Wait()
		}
		os.Exit(0)
	} else {
		fmt.Printf("[*] Service running in PID: %d PPID: %d\n", os.Getpid(), os.Getppid())
	}
	ok := vip.GetIniData("version.ok")
	if ok == "test" {
		test()
	} else if ok == "run" {
		run()
	} else {
		getWebSite()
		run()
	}
}
