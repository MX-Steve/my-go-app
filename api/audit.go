package api

import (
	"net/http"
	"strconv"

	"github.com/MX-Steve/my-go-app/model"

	"github.com/gin-gonic/gin"
)

func GetAudits(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	pageNo := c.DefaultQuery("page_no", "1")
	pageSize := c.DefaultQuery("page_size", "10")
	name := c.DefaultQuery("name", "")
	operator := c.DefaultQuery("operator", "")
	pageNo1, _ := strconv.Atoi(pageNo)
	pageSize1, _ := strconv.Atoi(pageSize)
	audit := &model.Audit{}
	data, _ := audit.GetAudits(name, operator)
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

func GetAuditById(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	id := c.DefaultQuery("id", "1")
	aid, _ := strconv.Atoi(id)
	audit := &model.Audit{}
	data, _ := audit.GetAuditById(aid)
	c.JSON(http.StatusOK, data)
}

func GetAuditByName(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	name := c.DefaultQuery("name", "")
	if name == "" {
		c.JSON(http.StatusOK, gin.H{
			"msg": "参数 name 不能为空",
		})
	}
	audit := &model.Audit{}
	data, _ := audit.GetAuditByName(name)
	c.JSON(http.StatusOK, data)
}
