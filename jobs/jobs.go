package jobs

import (
	"fmt"
	"log"
	"os/exec"
	"syscall"
	"time"

	"github.com/MX-Steve/my-go-app/tools"
	"github.com/MX-Steve/my-go-app/vip"
)

func RunPy() {
	py := vip.GetIniData("files.autoPath")
	cmd := exec.Command("cmd.exe")
	cmdExec := fmt.Sprintf("start python %s", py)
	//核心点,直接修改执行命令方式
	cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: fmt.Sprintf(`/c %s`, cmdExec), HideWindow: true}
	output, err := cmd.Output()
	fmt.Printf("output:\n%s\n", output)
	if err != nil {
		fmt.Printf("error: %+v\n", err)
	}
}

func T1() {
	fmt.Println("T1 Hi: ", time.Now())
}

func Relax() {
	log.Printf("%s : %s", tools.RunFuncName(), time.Now())
	tools.ShowMessage2("有闹钟啦!", "可以休息一下眼睛啦！")
}

func T3() {
	log.Println("T3 Hi: ", time.Now())
	RunPy()
}

func T4() {
	log.Println("T4 Hi: ", time.Now())
}

func T5() {
	fmt.Println("T5 Hi: ", time.Now())
}
