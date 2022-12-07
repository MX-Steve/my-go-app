package model

import (
	"fmt"
	"log"
)

type Task struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     int    `json:"enabled"`
	Every       int    `json:"every"`
	DelTag      int    `json:"del_tag"`
	UpdateTime  string `json:"update_time"`
	Creater     string `json:"creater"`
}

func (t *Task) GetTasks() ([]Task, error) {
	tasks := []Task{}
	dbw := NewDb()
	var id int
	var name string
	var description string
	var enabled int
	var every int
	var del_tag int
	var update_time string
	var creater string

	defer dbw.Db.Close()
	sqlStr := "select id,name,description,enabled,every,del_tag,update_time,creater from tools_interval_tasks where del_tag=0;"
	log.Println(sqlStr)
	rows, err := dbw.Db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		rows.Scan(&id, &name, &description, &enabled, &every, &del_tag, &update_time, &creater)
		tasks = append(tasks, Task{
			Id:          id,
			Name:        name,
			Description: description,
			Enabled:     enabled,
			Every:       every,
			DelTag:      del_tag,
			UpdateTime:  update_time,
			Creater:     creater})
	}
	return tasks, nil
}

func (t *Task) GetTasksForJob() ([]Task, error) {
	tasks := []Task{}
	dbw := NewDb()
	var id int
	var name string
	var description string
	var enabled int
	var every int
	var del_tag int
	var update_time string
	var creater string

	defer dbw.Db.Close()
	sqlStr := "select id,name,description,enabled,every,del_tag,update_time,creater from tools_interval_tasks where del_tag=0 and enabled=1;"
	log.Println(sqlStr)
	rows, err := dbw.Db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		rows.Scan(&id, &name, &description, &enabled, &every, &del_tag, &update_time, &creater)
		tasks = append(tasks, Task{
			Id:          id,
			Name:        name,
			Description: description,
			Enabled:     enabled,
			Every:       every,
			DelTag:      del_tag,
			UpdateTime:  update_time,
			Creater:     creater})
	}
	return tasks, nil
}

func (t *Task) GetTaskById(tid int) (Task, error) {
	task := Task{}
	dbw := NewDb()
	var id int
	var name string
	var description string
	var enabled int
	var every int
	var del_tag int
	var update_time string
	var creater string

	defer dbw.Db.Close()
	sqlStr := "select id,name,description,enabled,every,del_tag,update_time,creater from tools_interval_tasks where id=%d;"
	sqlStr = fmt.Sprintf(sqlStr, tid)
	log.Println(sqlStr)
	err := dbw.Db.QueryRow(sqlStr).Scan(&id, &name, &description, &enabled, &every, &del_tag, &update_time, &creater)
	if err != nil {
		return task, err
	}
	task = Task{
		Id:          id,
		Name:        name,
		Description: description,
		Enabled:     enabled,
		Every:       every,
		DelTag:      del_tag,
		UpdateTime:  update_time,
		Creater:     creater,
	}
	return task, nil
}

func (t *Task) GetTaskByName(tname string) (Task, error) {
	task := Task{}
	dbw := NewDb()
	var id int
	var name string
	var description string
	var enabled int
	var every int
	var del_tag int
	var update_time string
	var creater string

	defer dbw.Db.Close()
	sqlStr := "select id,name,description,enabled,every,del_tag,update_time,creater from tools_interval_tasks where name='%s';"
	sqlStr = fmt.Sprintf(sqlStr, tname)
	log.Println(sqlStr)
	err := dbw.Db.QueryRow(sqlStr).Scan(&id, &name, &description, &enabled, &every, &del_tag, &update_time, &creater)
	if err != nil {
		return task, err
	}
	task = Task{
		Id:          id,
		Name:        name,
		Description: description,
		Enabled:     enabled,
		Every:       every,
		DelTag:      del_tag,
		UpdateTime:  update_time,
		Creater:     creater,
	}
	return task, nil
}

func (t *Task) CudTask(sqlStr string) error {
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
