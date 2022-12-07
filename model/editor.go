package model

import (
	"fmt"
	"log"
)

type Article struct {
	Id          int    `json:"id"`
	Content     string `json:"content"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	CreatedAt   string `json:"created_at"`
	DelTag      int    `json:"del_tag"`
	CreatedById int    `json:"created_by_id"`
}

func (art *Article) GetArticles(ptitle, ptp, pcontent string) ([]Article, error) {
	articles := []Article{}
	dbw := NewDb()
	var id int
	var title string
	var created_by_id int
	var tp string
	var created_at string
	var del_tag int

	defer dbw.Db.Close()
	sqlStr := "select id,title,type,created_at,del_tag,created_by_id from tools_articles where del_tag=0"
	if ptitle != "" {
		sqlStr = fmt.Sprintf("%s and title like '%%%s%%'", sqlStr, ptitle)
	}
	if ptp != "" {
		sqlStr = fmt.Sprintf("%s and type='%s'", sqlStr, ptp)
	}
	if pcontent != "" {
		sqlStr = fmt.Sprintf("%s and content like '%%%s%%'", sqlStr, pcontent)
	}
	sqlStr = fmt.Sprintf("%s  order by id desc;", sqlStr)
	rows, err := dbw.Db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		rows.Scan(&id, &title, &tp, &created_at, &del_tag, &created_by_id)
		articles = append(articles, Article{
			Id:          id,
			Title:       title,
			Type:        tp,
			CreatedAt:   created_at,
			DelTag:      del_tag,
			CreatedById: created_by_id,
		})
	}
	return articles, nil
}

func (art *Article) GetArticleById(aid int) (Article, error) {
	dbw := NewDb()
	var id int
	var title string
	var created_by_id int
	var tp string
	var created_at string
	var del_tag int
	var content string

	defer dbw.Db.Close()
	sqlStr := "select id,title,type,content,created_at,del_tag,created_by_id from tools_articles where id=%d;"
	sqlStr = fmt.Sprintf(sqlStr, aid)
	err := dbw.Db.QueryRow(sqlStr).Scan(&id, &title, &tp, &content, &created_at, &del_tag, &created_by_id)
	if err != nil {
		return *art, err
	}
	art = &Article{
		Id:          id,
		Title:       title,
		Type:        tp,
		CreatedAt:   created_at,
		DelTag:      del_tag,
		CreatedById: created_by_id,
		Content:     content,
	}
	return *art, nil
}

func (art *Article) GetArticleByName(aname string) ([]Article, error) {
	articles := []Article{}
	dbw := NewDb()
	var id int
	var title string
	var created_by_id int
	var tp string
	var created_at string
	var del_tag int

	defer dbw.Db.Close()
	sqlStr := "select id,title,type,created_at,del_tag,created_by_id from tools_articles where title like '%?%';"
	err := dbw.Db.QueryRow(sqlStr, aname).Scan(&id, &title, &tp, &created_at, &del_tag, &created_by_id)
	if err != nil {
		return articles, err
	}
	articles = append(articles, Article{
		Id:          id,
		Title:       title,
		Type:        tp,
		CreatedAt:   created_at,
		DelTag:      del_tag,
		CreatedById: created_by_id,
	})
	return articles, nil
}

func (art *Article) CudArticle(sqlStr string) error {
	dbw := NewDb()
	defer dbw.Db.Close()
	log.Println(sqlStr)
	_, err := dbw.Db.Exec(sqlStr)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
