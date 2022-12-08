package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/MX-Steve/my-go-app/model"
	"github.com/MX-Steve/my-go-app/tools"

	"github.com/gin-gonic/gin"
)

func GetPlans(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	pageNo := c.DefaultQuery("page_no", "1")
	pageSize := c.DefaultQuery("page_size", "10")
	pageNo1, _ := strconv.Atoi(pageNo)
	pageSize1, _ := strconv.Atoi(pageSize)
	plan := &model.Plan{}
	data, _ := plan.GetPlans()
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

func GetPlanById(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	id := c.DefaultQuery("id", "1")
	tid, _ := strconv.Atoi(id)
	plan := &model.Plan{}
	data, _ := plan.GetPlanById(tid)
	c.JSON(http.StatusOK, data)
}

type PlanReq struct {
	Cud         string `json:"cud"`
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PlanDate    string `json:"plan_date"`
	PlanStatus  int    `json:"plan_status"`
	DelTag      int    `json:"del_tag"`
}

func CudPlans(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	json := PlanReq{}
	c.BindJSON(&json)
	cud := json.Cud
	if cud == "delete" {
		tools.AddAudit(c, "删除计划", fmt.Sprintf("编号ID为 %d 的数据被删除", json.Id))
		id := json.Id
		if id == 0 {
			c.JSON(http.StatusOK, gin.H{
				"msg": "删除必须携带 id 参数",
			})
			return
		}
		sqlStr := fmt.Sprintf("update tools_plan set del_tag=1 where id=%d", id)
		plan := &model.Plan{}
		err := plan.CudPlan(sqlStr)
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
		tools.AddAudit(c, "新增计划", fmt.Sprintf("新增了名称为 %s 的数据", json.Name))
		plan := &model.Plan{}
		name := json.Name
		description := json.Description
		timestamp := time.Now().Unix()
		tm := time.Unix(timestamp, 0)
		date_changed := tm.Format("2006-01-02")
		sqlStr := `insert into tools_plan(name,description,plan_date)
		 values("%s","%s", "%s");`
		sqlStr = fmt.Sprintf(sqlStr, name, description, date_changed)
		err := plan.CudPlan(sqlStr)
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
		tools.AddAudit(c, "更新计划", fmt.Sprintf("更新了编号为 %d 的数据", json.Id))
		plan := &model.Plan{}
		timestamp := time.Now().Unix()
		tm := time.Unix(timestamp, 0)
		date_changed := tm.Format("2006-01-02")
		sqlStr := `UPDATE tools_plan SET 
		name="%s",description="%s",plan_status=%d,plan_date="%s"
		WHERE id=%d`
		sqlStr = fmt.Sprintf(sqlStr, json.Name, json.Description, json.PlanStatus, date_changed, json.Id)
		err := plan.CudPlan(sqlStr)
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
