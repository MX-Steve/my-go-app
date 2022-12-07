package main

import (
	"fmt"
	"log"
	"syscall"
	"time"
	"unsafe"

	"github.com/MX-Steve/my-go-app/jobs"
	"github.com/MX-Steve/my-go-app/model"

	"github.com/robfig/cron/v3"
)

// go-bindata -o asset/asset.go -pkg asset  website/...

func IntPtr(n int) uintptr {
	return uintptr(n)
}
func StrPtr(s string) uintptr {
	return uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(s)))
}

//windows下的另一种DLL方法调用
func ShowMessage2(tittle, text string) {
	user32dll, _ := syscall.LoadLibrary("user32.dll")
	user32 := syscall.NewLazyDLL("user32.dll")
	MessageBoxW := user32.NewProc("MessageBoxW")
	MessageBoxW.Call(IntPtr(0), StrPtr(text), StrPtr(tittle), IntPtr(0))
	defer syscall.FreeLibrary(user32dll)
}

type EntryID int

var JobBulk = make(map[string]EntryID)
var c = cron.New()

func RTask(task model.Task) {
	var T func()
	if task.Name == "Relax" {
		T = jobs.Relax
	} else if task.Name == "T3" {
		T = jobs.T3
	} else if task.Name == "T4" {
		T = jobs.T4
	} else if task.Name == "T5" {
		T = jobs.T5
	} else {
		T = func() {
			log.Printf("当前任务 %s 未知", task.Name)
		}
	}
	sec := fmt.Sprintf("@every %ds", task.Every)
	EId, _ := c.AddFunc(sec, func() {
		T()
	})
	JobBulk[task.Name] = EntryID(EId)
	c.Start()
}

func Job() {
	log.Println("触发 JobManage 任务")
	var JobBulkNeeds = make(map[string]EntryID)
	task := model.Task{}
	tasks, err := task.GetTasksForJob()
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range tasks {
		JobBulkNeeds[t.Name] = 0
	}
	log.Println("最新任务列表：", JobBulkNeeds)
	log.Println("正在跑的任务列表：", JobBulk)
	for k := range JobBulkNeeds {
		_, ok := JobBulk[k]
		if !ok {
			log.Println("开始新增任务：", k)
			for _, task := range tasks {
				if task.Name == k {
					log.Println(task)
					RTask(task)
				}
			}
		}
	}
	for k, v := range JobBulk {
		_, ok := JobBulkNeeds[k]
		if !ok {
			var v2 = cron.EntryID(v)
			log.Printf("开始移除任务: %s , 任务ID: %d", k, v2)
			c.Remove(v2)
			delete(JobBulk, k)
		}
	}
}

func JobManage() {
	for {
		Job()
		time.Sleep(30 * time.Second)
	}
}
