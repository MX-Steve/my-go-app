package model

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type Audit struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UpdateTime  string `json:"update_time"`
	Operator    string `json:"operator"`
}

func (a *Audit) GetAudits(sname, soperator string) ([]Audit, error) {
	audits := []Audit{}
	dbw := NewDb()
	var id int
	var name string
	var description string
	var update_time string
	var operator string

	defer dbw.Db.Close()
	sqlStr := "select id,name,description,update_time,operator from tools_audits where"
	if sname != "" {
		sqlStr = fmt.Sprintf("%s name='%s' and", sqlStr, sname)
	}
	if soperator != "" {
		sqlStr = fmt.Sprintf("%s operator='%s'", sqlStr, soperator)
	}
	if strings.HasSuffix(sqlStr, "where") {
		sqls := strings.Split(sqlStr, "where")
		sqlStr = sqls[0]
	}
	if strings.HasSuffix(sqlStr, "and") {
		sqls := strings.Split(sqlStr, "and")
		sqlStr = sqls[0]
	}
	sqlStr = fmt.Sprintf("%s  order by id desc;", sqlStr)
	log.Println(sqlStr)
	rows, err := dbw.Db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		rows.Scan(&id, &name, &description, &update_time, &operator)
		audits = append(audits, Audit{
			Id:          id,
			Name:        name,
			Description: description,
			UpdateTime:  update_time,
			Operator:    operator})
	}
	return audits, nil
}

func (a *Audit) GetAuditById(aid int) (Audit, error) {
	audit := Audit{}
	dbw := NewDb()
	var id int
	var name string
	var description string
	var update_time string
	var operator string

	defer dbw.Db.Close()
	sqlStr := "select id,name,description,update_time,operator from tools_audits where id=%d;"
	sqlStr = fmt.Sprintf(sqlStr, aid)
	log.Println(sqlStr)
	err := dbw.Db.QueryRow(sqlStr).Scan(&id, &name, &description, &update_time, &operator)
	if err != nil {
		return audit, err
	}
	audit = Audit{
		Id:          id,
		Name:        name,
		Description: description,
		UpdateTime:  update_time,
		Operator:    operator,
	}
	return audit, nil
}

func (a *Audit) GetAuditByName(aname string) (Audit, error) {
	audit := Audit{}
	dbw := NewDb()
	var id int
	var name string
	var description string
	var update_time string
	var operator string

	defer dbw.Db.Close()
	sqlStr := "select id,name,description,update_time,operator from tools_audits where name='%s';"
	sqlStr = fmt.Sprintf(sqlStr, aname)
	log.Println(sqlStr)
	err := dbw.Db.QueryRow(sqlStr).Scan(&id, &name, &description, &update_time, &operator)
	if err != nil {
		return audit, err
	}
	audit = Audit{
		Id:          id,
		Name:        name,
		Description: description,
		UpdateTime:  update_time,
		Operator:    operator,
	}
	return audit, nil
}

func (a *Audit) AddAudit(name, description, operator string) error {
	dbw := NewDb()
	tm := time.Unix(time.Now().Unix(), 0)
	update_time := tm.Format("2006-01-02 15:04:05")
	sqlStr := "INSERT INTO tools_audits(name,description,update_time,operator) VALUES('%s','%s','%s','%s');"
	sqlStr = fmt.Sprintf(sqlStr, name, description, update_time, operator)
	log.Println(sqlStr)
	_, err := dbw.Db.Exec(sqlStr)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
