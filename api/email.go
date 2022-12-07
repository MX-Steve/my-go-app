package api

import (
	"log"
	"net/http"

	"github.com/MX-Steve/my-go-app/tools"

	"github.com/gin-gonic/gin"
)

type ShuJu struct {
	Name    string `json:"ename"`
	Pwd     string `json:"epwd"`
	Host    string `json:"smtp"`
	Subject string `json:"etitle"`
	Body    string `json:"ebody"`
	Send    string `json:"elists"`
}

func SendEmail(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	json := ShuJu{}
	c.BindJSON(&json)
	log.Println(json)
	err := tools.SendMail(json.Name, json.Pwd, json.Host, json.Send, json.Subject, json.Body, "text")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "方法未知",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "发送邮件成功",
	})

}
