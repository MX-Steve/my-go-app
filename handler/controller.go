package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/MX-Steve/my-go-app/model"

	"github.com/gin-gonic/gin"
)

type Label struct {
	Icon  string `json:"icon"`
	Title string `json:"title"`
	H3    string `json:"h3"`
	P     string `json:"p"`
}

type Headers struct {
	Labels map[string]Label `json:"headers"`
}

type Carousel struct {
	Title string `json:"title"`
	Img   string `json:"Img"`
}

type Carousel01 struct {
	Carousels map[string]Carousel `json:"carousel01"`
}

type Carousel02 struct {
	Carousels map[string]Carousel `json:"carousel02"`
}

func GetPs() []map[string]interface{} {
	plan := &model.Plan{}
	plans, _ := plan.GetPlans()
	result := make([]map[string]interface{}, 0)
	for _, item := range plans {
		result = append(result, map[string]interface{}{
			"id":          item.Id,
			"name":        item.Name,
			"description": item.Description,
			"plan_date":   item.PlanDate,
			"plan_status": item.PlanStatus,
			"del_tag":     item.DelTag,
		})
	}
	return result
}

func GetIndexs() (map[string]map[string]string, map[string]map[string]string, map[string]map[string]string) {
	jsonFile, err := os.Open("./website/json/index.json")
	if err != nil {
		log.Println("err open jsonfile index.json")
		return nil, nil, nil
	}
	defer jsonFile.Close()
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Println("error reading json file")
		return nil, nil, nil
	}
	var headers Headers
	json.Unmarshal(jsonData, &headers)
	data := make(map[string]map[string]string, 4)
	for label, value := range headers.Labels {
		data[label] = map[string]string{
			"icon":  value.Icon,
			"title": value.Title,
			"h3":    value.H3,
			"p":     value.P,
		}
	}
	var carousel01 Carousel01
	json.Unmarshal(jsonData, &carousel01)
	c1 := make(map[string]map[string]string, 10)
	for c, value := range carousel01.Carousels {
		c1[c] = map[string]string{
			"title": value.Title,
			"img":   value.Img,
		}
	}
	var carousel02 Carousel02
	json.Unmarshal(jsonData, &carousel02)
	c2 := make(map[string]map[string]string, 10)
	for c, value := range carousel02.Carousels {
		c2[c] = map[string]string{
			"title": value.Title,
			"img":   value.Img,
		}
	}
	return data, c1, c2
}

func IndexHandler(c *gin.Context) {
	data, c1, c2 := GetIndexs()
	plans := GetPs()
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title":      "开拓者平台",
		"Headers":    data,
		"Carousel01": c1,
		"Carousel02": c2,
		"Plans":      plans,
	})
}
func DetailHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "detail.html", gin.H{
		"Title": "开拓者",
		"Headers": map[string]map[string]string{
			"label01": map[string]string{
				"title": "人生忠告",
				"h3":    "人无完人",
				"p":     "人之大忌,好为人师",
			},
		},
	})
}

func LoginHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"Title": "开拓者平台登录",
	})
}
func HomeHandler(c *gin.Context) {
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
	data = data[s:e]
	result := make([]map[string]interface{}, 0)
	for _, item := range data {
		result = append(result, map[string]interface{}{
			"id":          item.Id,
			"name":        item.Name,
			"description": item.Description,
			"enabled":     item.Enabled,
			"every":       item.Every,
			"del_tag":     item.DelTag,
			"update_time": item.UpdateTime,
			"creater":     item.Creater,
		})
	}
	c.HTML(http.StatusOK, "home.html", gin.H{
		"Title": "开拓者平台个人中心",
		"Data":  result,
	})
}
