package model

import (
	"database/sql"
	"fmt"
	"log"
)

type User struct {
	Id       uint64 `json:"id" form:"id"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Phone    string `json:"phone" form:"phone"`
	Photo    string `json:"photo" form:"photo"`
	Email    string `json:"email" form:"email"`
}

func (u *User) GetUser(uid uint64) (usr User, err error) {
	users := User{}
	dbw := NewDb()
	defer dbw.Db.Close()
	var email sql.NullString
	var username sql.NullString
	var id uint64
	var password string
	var photo string
	var phone string
	//单行查询
	err = dbw.Db.
		QueryRow("select id,username,password,phone,email,photo from user where id = ?", uid).
		Scan(&id, &username, &password, &phone, &email, &photo)

	if err != nil {
		log.Printf("Query data error: %v\n", err)
		return users, err
	}
	if !username.Valid {
		username.String = ""
	}
	if !email.Valid {
		email.String = ""
	}
	users = User{
		Id:       id,
		Username: username.String,
		Email:    email.String,
		Password: password,
		Photo:    photo,
		Phone:    phone}
	return users, nil
}

func (u *User) Register(username, password, phone, email string) bool {
	dbw := NewDb()
	defer dbw.Db.Close()
	sqlStr := fmt.Sprintf("INSERT INTO user(username, password, phone, email) values('%s', '%s','%s','%s')", username, password, phone, email)
	log.Println(sqlStr)
	_, err := dbw.Db.Exec(sqlStr)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (u *User) Login(username, password string) (int, bool) {
	dbw := NewDb()
	defer dbw.Db.Close()
	var id int
	err := dbw.Db.
		QueryRow("select id from user where username = ? and password = ?", username, password).
		Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, false
	}
	return id, true
}

func (u *User) ModUser(id uint64, username, password, phone, photo, email string) bool {
	dbw := NewDb()
	defer dbw.Db.Close()
	sqlStr := "UPDATE user SET"
	if username != "" {
		sqlStr = fmt.Sprintf("%s username='%s',", sqlStr, username)
	}
	if password != "" {
		sqlStr = fmt.Sprintf("%s password='%s',", sqlStr, password)
	}
	if phone != "" {
		sqlStr = fmt.Sprintf("%s phone='%s',", sqlStr, phone)
	}
	if photo != "" {
		sqlStr = fmt.Sprintf("%s photo='%s',", sqlStr, photo)
	}
	if email != "" {
		sqlStr = fmt.Sprintf("%s email='%s',", sqlStr, email)
	}
	sqlStr = sqlStr[:len(sqlStr)-1]
	sqlStr = fmt.Sprintf("%s WHERE id=%d", sqlStr, id)
	log.Println(sqlStr)
	_, err := dbw.Db.Exec(sqlStr)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
