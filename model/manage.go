package model

import (
	"fmt"
	"log"
)

type Toolsmanage struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Bin     string `json:"bin"`
	Src     string `json:"src"`
	Type    string `json:"type"`
	Enabled int    `json:"enabled"`
	DelTag  int    `json:"del_tag"`
}

func (tm *Toolsmanage) GetToolsmanage(stp string) ([]Toolsmanage, error) {
	tms := []Toolsmanage{}
	dbw := NewDb()
	var id int
	var name string
	var bin string
	var src string
	var tp string
	var enabled int
	var del_tag int

	defer dbw.Db.Close()
	sqlStr := "select id,name,bin,src,type,enabled,del_tag from toolsmanage where "
	if stp != "" {
		sqlStr = fmt.Sprintf("%s type='%s'", sqlStr, stp)
	}
	sqlStr = fmt.Sprintf("%s del_tag=0;", sqlStr)
	log.Println(sqlStr)
	rows, err := dbw.Db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		rows.Scan(&id, &name, &bin, &src, &tp, &enabled, &del_tag)
		tms = append(tms, Toolsmanage{
			Id:      id,
			Name:    name,
			Bin:     bin,
			Src:     src,
			Type:    tp,
			Enabled: enabled,
			DelTag:  del_tag,
		})
	}
	return tms, nil
}

func (tm *Toolsmanage) GetToolsmanageById(tid int) (Toolsmanage, error) {
	tm0 := Toolsmanage{}
	dbw := NewDb()
	var id int
	var name string
	var bin string
	var src string
	var tp string
	var enabled int
	var del_tag int
	defer dbw.Db.Close()
	sqlStr := "select id,name,bin,src,type,enabled,del_tag from toolsmanage where id=%d;"
	sqlStr = fmt.Sprintf(sqlStr, tid)
	err := dbw.Db.QueryRow(sqlStr).Scan(&id, &name, &bin, &src, &tp, &enabled, &del_tag)
	if err != nil {
		return tm0, err
	}
	tm0 = Toolsmanage{
		Id:      id,
		Name:    name,
		Bin:     bin,
		Src:     src,
		Type:    tp,
		Enabled: enabled,
		DelTag:  del_tag,
	}
	return tm0, nil
}

func (tm *Toolsmanage) CudToolsmanage(sqlStr string) error {
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
