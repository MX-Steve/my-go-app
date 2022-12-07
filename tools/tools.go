package tools

import (
	"runtime"
	"syscall"
	"unsafe"

	"github.com/MX-Steve/my-go-app/model"

	"github.com/gin-gonic/gin"
)

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

func RunFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

func AddAudit(c *gin.Context, name string, description string) {
	var sid uint64
	k, _ := c.Get("userId")
	val, ok := k.(uint64)
	if !ok {
		sid = uint64(val)
	} else {
		sid = val
	}
	user := model.User{}
	u, err := user.GetUser(sid)
	if err != nil {
		return
	}
	uname := u.Username
	audit := model.Audit{}
	err = audit.AddAudit(name, description, uname)
	if err != nil {
		return
	}
}
