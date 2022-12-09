package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/MX-Steve/my-go-app/model"
	"github.com/MX-Steve/my-go-app/tools"

	"github.com/gin-gonic/gin"
)

func GetToolsmanage(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	pageNo := c.DefaultQuery("page_no", "1")
	pageSize := c.DefaultQuery("page_size", "10")
	tp := c.DefaultQuery("type", "")
	pageNo1, _ := strconv.Atoi(pageNo)
	pageSize1, _ := strconv.Atoi(pageSize)
	tm := &model.Toolsmanage{}
	data, _ := tm.GetToolsmanage(tp)
	s := 0
	e := 1
	if pageSize1 > len(data) {
		s = pageSize1 * (pageNo1 - 1)
		e = pageSize1*(pageNo1-1) + len(data)
	} else {
		s = pageSize1 * (pageNo1 - 1)
		e = pageNo1 * pageSize1
	}
	totalCount := len(data)
	pages := totalCount / pageSize1
	if totalCount%pageSize1 != 0 {
		pages += 1
	}
	if pageNo1 == pages {
		e = totalCount
	}
	data = data[s:e]
	c.JSON(http.StatusOK, gin.H{
		"data":  data,
		"pages": pages,
	})
}

func ToolsRun(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	id := c.DefaultQuery("id", "1")
	tid, _ := strconv.Atoi(id)
	tm := &model.Toolsmanage{}
	data, _ := tm.GetToolsmanageById(tid)
	tools.AddAudit(c, "执行命令", fmt.Sprintf("%s 被调用", data.Name))
	if strings.HasSuffix(data.Bin, ".py") {
		cmd := exec.Command("cmd.exe", "/C", "start python "+data.Bin)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		log.Println(cmd.Run())
	} else if strings.HasSuffix(data.Bin, ".exe") {
		cmd := exec.Command("cmd.exe", "/C", "start "+data.Bin)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		log.Println(cmd.Run())
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "程序调用",
	})
}

func GetToolsmanageById(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	id := c.DefaultQuery("id", "1")
	tid, _ := strconv.Atoi(id)
	tm := &model.Toolsmanage{}
	data, _ := tm.GetToolsmanageById(tid)
	c.JSON(http.StatusOK, data)
}

type ToolsmanageReq struct {
	Cud     string `json:"cud"`
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Bin     string `json:"bin"`
	Src     string `json:"src"`
	Type    string `json:"type"`
	Enabled int    `json:"enabled"`
	DelTag  int    `json:"del_tag"`
}

func CudToolsmanage(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	json := ToolsmanageReq{}
	c.BindJSON(&json)
	cud := json.Cud
	if cud == "delete" {
		tools.AddAudit(c, "删除工具", fmt.Sprintf("编号ID为 %d 的数据被删除", json.Id))
		id := json.Id
		if id == 0 {
			c.JSON(http.StatusOK, gin.H{
				"msg": "删除必须携带 id 参数",
			})
			return
		}
		sqlStr := fmt.Sprintf("update toolsmanage set del_tag=1 where id=%d", id)
		tm := &model.Toolsmanage{}
		err := tm.CudToolsmanage(sqlStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "删除数据失败",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": "删除数据成功",
		})
		return
	} else if cud == "post" {
		tools.AddAudit(c, "新增工具", fmt.Sprintf("新增了名称为 %s 的数据", json.Name))
		tm := &model.Toolsmanage{}
		name := json.Name
		src := json.Src
		bin := json.Bin
		tp := json.Type
		enabled := json.Enabled
		sqlStr := `insert into toolsmanage(name,bin,src, type, enabled)
		 values("%s","%s", "%s", "%s", %d);`
		sqlStr = fmt.Sprintf(sqlStr, name, bin, src, tp, enabled)
		err := tm.CudToolsmanage(sqlStr)
		if err != nil {
			c.JSON(http.StatusExpectationFailed, gin.H{
				"msg": err,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": "新增数据成功",
		})
		return
	} else if cud == "put" {
		tools.AddAudit(c, "更新工具", fmt.Sprintf("更新了编号为 %d 的数据", json.Id))
		tm := &model.Toolsmanage{}
		sqlStr := `UPDATE toolsmanage SET 
		name="%s",bin="%s", src="%s", type="%s",enabled=%d
		WHERE id=%d`
		sqlStr = fmt.Sprintf(sqlStr, json.Name, json.Bin, json.Src, json.Type, json.Enabled, json.Id)
		err := tm.CudToolsmanage(sqlStr)
		if err != nil {
			c.JSON(http.StatusExpectationFailed, gin.H{
				"msg": err,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": "更新数据成功",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "方法未知",
	})
}
