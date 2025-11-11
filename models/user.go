package models

import (
	"avileads-web/utils"
	"database/sql"
)

type Users struct {
	Id       int           `form:"-"`
	ClientId int           `orm:"column(clientid)"`
	Password string        `orm:"column(pass)"`
	Login    string        `orm:"column(login)"`
	Name     string        `orm:"column(name)"`
	RoleId   sql.NullInt64 `orm:"column(roleid)"`
}

func (user *Users) Save() (int64, error) {
	user.Password = utils.HashUsingMD5(user.Password)
	id, err := utils.CreateDefaultDbContext().Insert(user)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (user *Users) Update() error {
	_, err := utils.CreateDefaultDbContext().Update(user, "ClientId", "Login", "Name", "RoleId")
	return err
}

func ResetPassword(userId int, newPassword string) error {
	md5HashedPassword := utils.HashUsingMD5(newPassword)

	user := &Users{}
	user.Id = userId
	user.Password = md5HashedPassword

	_, err := utils.CreateDefaultDbContext().Update(user, "Password")

	return err
}

func GetUsers() ([]Users, error) {
	o := utils.CreateDefaultDbContext()

	var users []Users
	_, err := o.QueryTable("users").OrderBy("id").All(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetUserById(id int64) (Users, error) {
	o := utils.CreateDefaultDbContext()

	var user Users
	err := o.QueryTable("users").Filter("Id", id).One(&user)
	if err != nil {
		return Users{}, err
	}

	return user, nil
}

func GetUserByLogin(login string) (Users, error) {
	o := utils.CreateDefaultDbContext()

	var user Users
	err := o.QueryTable("users").Filter("login", login).One(&user)
	if err != nil {
		return Users{}, err
	}

	return user, nil
}

func GetUserBy(user *Users, cols ...string) error {
	o := utils.CreateDefaultDbContext()

	err := o.Read(user, cols...)
	if err != nil {
		return err
	}
	return nil
}
