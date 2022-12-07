package model

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Alarm struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Desc        string `json:"desc"`
	Music       int    `json:"music"`
	AlarmTime   string `json:"alarm_time"`
	CreatedAt   string `json:"created_at"`
	DelTag      int    `json:"del_tag"`
	CreatedById string `json:"created_by_id"`
}

func (a *Alarm) GetAlarmClock(sqlStr string) (alarms []Alarm, err error) {
	dbw := NewDb()
	defer dbw.Db.Close()
	var id int
	var name string
	var desc string
	var music int
	var alarm_time string
	var created_at string
	var del_tag int
	var created_by_id string
	rows, err := dbw.Db.Query(sqlStr)
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		rows.Scan(&id, &name, &desc, &music, &alarm_time, &created_at, &del_tag, &created_by_id)
		alarms = append(alarms, Alarm{
			Id:          id,
			Name:        name,
			Desc:        desc,
			Music:       music,
			AlarmTime:   alarm_time,
			CreatedAt:   created_at,
			DelTag:      del_tag,
			CreatedById: created_by_id,
		})
	}
	defer rows.Close()
	return alarms, nil
}

func (a *Alarm) CudAlarmClock(sqlStr string) error {
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
