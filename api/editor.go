package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/MX-Steve/my-go-app/model"
	"github.com/MX-Steve/my-go-app/tools"
	"github.com/MX-Steve/my-go-app/vip"

	"github.com/gin-gonic/gin"
)

func GetArticles(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	pageNo := c.DefaultQuery("page_no", "1")
	pageSize := c.DefaultQuery("page_size", "10")
	title := c.DefaultQuery("title", "")
	tp := c.DefaultQuery("type", "")
	content := c.DefaultQuery("content", "")
	pageNo1, _ := strconv.Atoi(pageNo)
	pageSize1, _ := strconv.Atoi(pageSize)
	article := &model.Article{}
	data, _ := article.GetArticles(title, tp, content)
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

func GetArticleById(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	id := c.DefaultQuery("id", "1")
	aid, _ := strconv.Atoi(id)
	article := &model.Article{}
	data, _ := article.GetArticleById(aid)
	dir := vip.GetIniData("files.path")
	result := strings.HasPrefix(data.Content, dir)
	if result {
		fpath := data.Content
		b, err := os.ReadFile(fpath)
		if err != nil {
			c.JSON(http.StatusExpectationFailed, gin.H{
				"msg": err,
			})
			return
		}
		data.Content = string(b)
	}
	c.JSON(http.StatusOK, data)
}

func GetArticleByName(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	name := c.DefaultQuery("name", "")
	article := &model.Article{}
	var data []model.Article
	if name != "" {
		data, _ = article.GetArticleByName(name)
	} else {
		data = []model.Article{}
	}
	c.JSON(http.StatusOK, data)
}

type ArticleReq struct {
	Cud         string `json:"cud"`
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Type        string `json:"type"`
	CreatedAt   string `json:"created_at"`
	DelTag      int    `json:"del_tag"`
	CreatedById int    `json:"created_by_id"`
}

func OldArt2New(id int, title, tp, content string) string {
	fpath := vip.GetIniData("files.path")
	if strings.HasSuffix(title, ".go") || strings.HasSuffix(title, ".py") || strings.HasSuffix(title, ".sh") || strings.HasSuffix(title, ".txt") || strings.HasSuffix(title, ".md") {
		fpath = fmt.Sprintf("%s%s/%s", fpath, tp, title)
	} else {
		fpath = fmt.Sprintf("%s%s/%s.docx", fpath, tp, title)
	}
	err := os.WriteFile(fpath, []byte(content), 0644)
	if err != nil {
		return ""
	}
	return fpath
}

func CudArticles(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	var sid uint64 = 1
	k, _ := c.Get("userId")
	val, ok := k.(uint64)
	if !ok {
		sid = uint64(val)
	} else {
		sid = uint64(val)
	}
	users := &model.User{}
	us, _ := users.GetUser(sid)
	json := ArticleReq{}
	c.BindJSON(&json)
	log.Println(json)
	cud := json.Cud
	if cud == "delete" {
		tools.AddAudit(c, "删除文本", fmt.Sprintf("编号ID为 %d 的数据被删除", json.Id))
		id := json.Id
		if id == 0 {
			c.JSON(http.StatusOK, gin.H{
				"msg": "删除必须携带 id 参数",
			})
			return
		}
		article := &model.Article{}
		data, _ := article.GetArticleById(id)
		sqlStr := ""
		dir := vip.GetIniData("files.path")
		if strings.HasPrefix(data.Content, dir) {
			fpath := data.Content
			_, filename := filepath.Split(fpath)
			delPath := fmt.Sprintf("%sdeleted/%s", dir, filename)
			err := os.Rename(fpath, delPath)
			sqlStr = fmt.Sprintf("update tools_articles set del_tag=1, content='%s' where id=%d", delPath, id)
			if err != nil {
				c.JSON(http.StatusExpectationFailed, gin.H{
					"msg": err,
				})
				return
			}
		} else {
			content := OldArt2New(id, data.Title, data.Type, data.Content)
			sqlStr = fmt.Sprintf("update tools_articles set del_tag=1, content='%s' where id=%d", content, id)
		}
		err := article.CudArticle(sqlStr)
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
		tools.AddAudit(c, "新增文本", fmt.Sprintf("新增了名称为 %s 的数据文本", json.Title))
		article := &model.Article{}
		title := json.Title
		tp := json.Type
		content := json.Content
		timestamp := time.Now().Unix()
		tm := time.Unix(timestamp, 0)
		created_at := tm.Format("2006-01-02 15:04:05")
		created_by_id := us.Id
		fpath := vip.GetIniData("files.path")
		if strings.HasSuffix(title, ".go") || strings.HasSuffix(title, ".py") || strings.HasSuffix(title, ".sh") || strings.HasSuffix(title, ".txt") || strings.HasSuffix(title, ".md") {
			fpath = fmt.Sprintf("%s%s/%s", fpath, tp, title)
		} else {
			fpath = fmt.Sprintf("%s%s/%s.docx", fpath, tp, title)
		}
		err := os.WriteFile(fpath, []byte(content), 0644)
		if err != nil {
			c.JSON(http.StatusExpectationFailed, gin.H{
				"msg": err,
			})
			return
		}
		sqlStr := `insert into tools_articles(title,type,content,created_at,created_by_id,del_tag)
		 values("%s","%s",'%s',"%s",%d,0);`
		sqlStr = fmt.Sprintf(sqlStr, title, tp, fpath, created_at, created_by_id)
		err = article.CudArticle(sqlStr)
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
		tools.AddAudit(c, "更新文本", fmt.Sprintf("更新了编号为 %d 的数据", json.Id))
		article := &model.Article{}
		timestamp := time.Now().Unix()
		tm := time.Unix(timestamp, 0)
		created_at := tm.Format("2006-01-02 15:04:05")
		fpath := vip.GetIniData("files.path")
		title := json.Title
		if strings.HasSuffix(title, ".go") || strings.HasSuffix(title, ".py") || strings.HasSuffix(title, ".sh") || strings.HasSuffix(title, ".txt") || strings.HasSuffix(title, ".md") {
			fpath = fmt.Sprintf("%s%s/%s", fpath, json.Type, title)
		} else {
			fpath = fmt.Sprintf("%s%s/%s.docx", fpath, json.Type, title)
		}
		err := os.WriteFile(fpath, []byte(json.Content), 0644)
		if err != nil {
			c.JSON(http.StatusExpectationFailed, gin.H{
				"msg": err,
			})
			return
		}
		sqlStr := `UPDATE tools_articles SET 
		title="%s",type="%s",content='%s',created_at="%s",created_by_id=%d
		WHERE id=%d`
		sqlStr = fmt.Sprintf(sqlStr, json.Title, json.Type, fpath, created_at, us.Id, json.Id)
		err = article.CudArticle(sqlStr)
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
