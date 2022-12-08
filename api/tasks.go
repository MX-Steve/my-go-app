package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MX-Steve/my-go-app/model"
	"github.com/MX-Steve/my-go-app/tools"

	"github.com/gin-gonic/gin"
)

func GetTasks(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	pageNo := c.DefaultQuery("page_no", "1")
	pageSize := c.DefaultQuery("page_size", "10")
	pageNo1, _ := strconv.Atoi(pageNo)
	pageSize1, _ := strconv.Atoi(pageSize)
	task := &model.Task{}
	data, _ := task.GetTasks()
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

func GetTaskById(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	id := c.DefaultQuery("id", "1")
	tid, _ := strconv.Atoi(id)
	task := &model.Task{}
	data, _ := task.GetTaskById(tid)
	c.JSON(http.StatusOK, data)
}

type TaskReq struct {
	Cud         string `json:"cud"`
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     int    `json:"enabled"`
	Every       int    `json:"every"`
	DelTag      int    `json:"del_tag"`
	UpdateTime  string `json:"update_time"`
	Creater     string `json:"creater"`
}

func CudTasks(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	json := TaskReq{}
	c.BindJSON(&json)
	cud := json.Cud
	if cud == "delete" {
		tools.AddAudit(c, "删除任务", fmt.Sprintf("删除了编号为 %d 的数据", json.Id))
		id := json.Id
		if id == 0 {
			c.JSON(http.StatusOK, gin.H{
				"msg": "删除必须携带 id 参数",
			})
			return
		}
		sqlStr := fmt.Sprintf("update tools_interval_tasks set del_tag=1 where id=%d", id)
		log.Println(sqlStr)
		task := &model.Task{}
		err := task.CudTask(sqlStr)
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
		tools.AddAudit(c, "新增任务", fmt.Sprintf("新增了名称为 %s 的数据", json.Name))
		task := &model.Task{}
		name := json.Name
		description := json.Description
		every := json.Every
		timestamp := time.Now().Unix()
		tm := time.Unix(timestamp, 0)
		date_changed := tm.Format("2006-01-02 15:04:05")
		sqlStr := `insert into tools_interval_tasks(name,description,enabled,every,del_tag,update_time,creater)
		 values("%s","%s",0,%d,0,"%s","admin");`
		sqlStr = fmt.Sprintf(sqlStr, name, description, every, date_changed)
		log.Println(sqlStr)
		err := task.CudTask(sqlStr)
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
		tools.AddAudit(c, "更新任务", fmt.Sprintf("更新了编号为 %d 的数据", json.Id))
		task := &model.Task{}
		timestamp := time.Now().Unix()
		tm := time.Unix(timestamp, 0)
		date_changed := tm.Format("2006-01-02 15:04:05")
		sqlStr := `UPDATE tools_interval_tasks SET 
		name="%s",description="%s",enabled=%d,update_time="%s",every=%d
		WHERE id=%d`
		sqlStr = fmt.Sprintf(sqlStr, json.Name, json.Description, json.Enabled, date_changed, json.Every, json.Id)
		log.Println(sqlStr)
		err := task.CudTask(sqlStr)
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": "更新数据成功",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "方法未知",
	})
}
