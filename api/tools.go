package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/MX-Steve/my-go-app/model"

	"github.com/gin-gonic/gin"
)

func GetAlarmClock(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	alarms := &model.Alarm{}
	sqlStr := "select * from tools_alarmclock where del_tag=0 "
	id := c.DefaultQuery("id", "")
	if id != "" {
		sqlStr = fmt.Sprintf("%s and id=%s", sqlStr, id)
	}
	name := c.DefaultQuery("name", "")
	if name != "" {
		sqlStr = fmt.Sprintf("%s and name='%s'", sqlStr, name)
	}
	pageNo := c.DefaultQuery("page_no", "1")
	pageSize := c.DefaultQuery("page_size", "10")
	start := c.DefaultQuery("start", "")
	if start != "" {
		sqlStr = fmt.Sprintf("%s and created_at>='%s'", sqlStr, start)
	}
	end := c.DefaultQuery("end", "")
	if end != "" {
		sqlStr = fmt.Sprintf("%s and created_at<='%s'", sqlStr, end)
	}
	pageNo1, _ := strconv.Atoi(pageNo)
	pageSize1, _ := strconv.Atoi(pageSize)
	data, _ := alarms.GetAlarmClock(sqlStr)
	s := 0
	e := 1
	if pageSize1 > len(data) {
		s = pageSize1 * (pageNo1 - 1)
		e = pageSize1*(pageNo1-1) + len(data)
	} else {
		s = pageSize1 * (pageNo1 - 1)
		e = pageNo1 * pageSize1
	}
	data = data[s:e]
	c.JSON(http.StatusOK, data)
}

type AlarmReq struct {
	Cud       string `json:"cud"`
	Id        int    `json:"id"`
	Desc      string `json:"desc"`
	Name      string `json:"name"`
	AlarmTime string `json:"alarm_time"`
	Music     int    `json:"music"`
}

func CudAlarmClock(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	json := AlarmReq{}
	c.BindJSON(&json)
	cud := json.Cud
	fmt.Println(cud)
	if cud == "post" {
		alarm := &model.Alarm{}
		timestamp := time.Now().Unix()
		tm := time.Unix(timestamp, 0)
		created_at := tm.Format("2006-01-02 15:04:05")
		sqlStr := "INSERT INTO tools_alarmclock(name, `desc`, music, alarm_time, created_at, del_tag, created_by_id) values('%s','%s',%d,'%s','%s',0,'%s')"
		sqlStr = fmt.Sprintf(sqlStr, json.Name, json.Desc, json.Music, json.AlarmTime, created_at, "1")
		err := alarm.CudAlarmClock(sqlStr)
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
		if json.Id == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "必须提供参数id",
			})
			return
		}
		timestamp := time.Now().Unix()
		tm := time.Unix(timestamp, 0)
		created_at := tm.Format("2006-01-02 15:04:05")
		alarm := &model.Alarm{}
		sqlStr := ""
		if json.Name != "" {
			sqlStr = fmt.Sprintf("name='%s'", json.Name)
		}
		if json.Desc != "" {
			sqlStr = fmt.Sprintf("%s, desc='%s'", sqlStr, json.Desc)
		}
		if json.AlarmTime != "" {
			sqlStr = fmt.Sprintf("%s, alarm_time='%s'", sqlStr, json.AlarmTime)
		}
		sqlStr = fmt.Sprintf("%s, music=%d", sqlStr, json.Music)
		sqlStr = fmt.Sprintf("update tools_alarmclock set %s, created_time='%s' where id=%d", sqlStr, created_at, json.Id)
		err := alarm.CudAlarmClock(sqlStr)
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
	} else if cud == "delete" {
		id := json.Id
		if id == 0 {
			c.JSON(http.StatusOK, gin.H{
				"msg": "删除必须携带 id 参数",
			})
			return
		}
		sqlStr := fmt.Sprintf("update tools_alarmclock set del_tag=1 where id=%d", id)
		alarm := &model.Alarm{}
		err := alarm.CudAlarmClock(sqlStr)
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
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "方法未知",
	})
}
