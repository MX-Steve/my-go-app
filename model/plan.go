package model

import (
	"fmt"
	"log"
	"time"
)

type Plan struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PlanStatus  int    `json:"plan_status"`
	PlanDate    string `json:"plan_date"`
	DelTag      int    `json:"del_tag"`
}

func (p *Plan) GetPlans() ([]Plan, error) {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	plan_date0 := tm.Format("2006-01-02")
	plans := []Plan{}
	dbw := NewDb()
	var id int
	var name string
	var description string
	var plan_status int
	var plan_date string
	var del_tag int

	defer dbw.Db.Close()
	sqlStr := "select id,name,description,plan_status,plan_date,del_tag from tools_plan where del_tag=0 and plan_date='%s';"
	sqlStr = fmt.Sprintf(sqlStr, plan_date0)
	log.Println(sqlStr)
	rows, err := dbw.Db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		rows.Scan(&id, &name, &description, &plan_status, &plan_date, &del_tag)
		plans = append(plans, Plan{
			Id:          id,
			Name:        name,
			Description: description,
			PlanStatus:  plan_status,
			PlanDate:    plan_date,
			DelTag:      del_tag,
		})
	}
	return plans, nil
}

func (p *Plan) GetPlanById(pid int) (Plan, error) {
	plan := Plan{}
	dbw := NewDb()
	var id int
	var name string
	var description string
	var plan_status int
	var plan_date string
	var del_tag int

	defer dbw.Db.Close()
	sqlStr := "select id,name,description,plan_status,plan_date,del_tag from tools_plan where id=%d;"
	sqlStr = fmt.Sprintf(sqlStr, pid)
	err := dbw.Db.QueryRow(sqlStr).Scan(&id, &name, &description, &plan_status, &plan_date, &del_tag)
	if err != nil {
		return plan, err
	}
	plan = Plan{
		Id:          id,
		Name:        name,
		Description: description,
		PlanStatus:  plan_status,
		PlanDate:    plan_date,
		DelTag:      del_tag,
	}
	return plan, nil
}

func (p *Plan) GetPlanByName(pname string) (Plan, error) {
	plan := Plan{}
	dbw := NewDb()
	var id int
	var name string
	var description string
	var plan_status int
	var plan_date string
	var del_tag int

	defer dbw.Db.Close()
	sqlStr := "id,name,description,plan_status,plan_date,del_tag from tools_plan  where name='%s';"
	sqlStr = fmt.Sprintf(sqlStr, pname)
	err := dbw.Db.QueryRow(sqlStr).Scan(&id, &name, &description, &plan_status, &plan_date, &del_tag)
	if err != nil {
		return plan, err
	}
	plan = Plan{
		Id:          id,
		Name:        name,
		Description: description,
		PlanStatus:  plan_status,
		PlanDate:    plan_date,
		DelTag:      del_tag,
	}
	return plan, nil
}

func (p *Plan) CudPlan(sqlStr string) error {
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
