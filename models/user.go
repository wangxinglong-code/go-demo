package models

import (
	"go-demo/utils/mysql"
)

type User struct {
	Uid  int    `gorm:"column:uid" json:"uid"`
	Name string `gorm:"column:name" json:"name"`
	Age  int    `gorm:"column:age" json:"age"`
}

func (u *User) TableName() string {
	return "user"
}

//获取DB别名的规则
func (u *User) DbAliasName() string {
	return "go_demo"
}

//获取db节点的规则
func (u *User) DbNode() int {
	return 0
}

func (u *User) Get(id int) (*User, error) {
	db, err := mysql.GetConn(u.DbAliasName(), u.DbNode(), true)
	if err != nil {
		return nil, err
	}
	var table User
	if err := db.Table(u.TableName()).Where("id = ?", id).First(&table).Error; err != nil {
		return nil, err
	}
	return &table, nil
}

func (u *User) Update(table *User) error {
	db, err := mysql.GetConn(u.DbAliasName(), u.DbNode(), true)
	if err != nil {
		return err
	}
	if err := db.Table(u.TableName()).Updates(&table).Error; err != nil {
		return err
	}
	return nil
}

func (u *User) Create(table *User) (*User, error) {
	db, err := mysql.GetConn(u.DbAliasName(), u.DbNode(), true)
	if err != nil {
		return nil, err
	}
	if err := db.Table(u.TableName()).Create(&table).Error; err != nil {
		return nil, err
	}
	return table, err
}

func (u *User) List(param map[string]interface{}) ([]*User, int64, error) {
	db, err := mysql.GetConn(u.DbAliasName(), u.DbNode(), false)
	if err != nil {
		return nil, 0, err
	}
	db = db.Table(u.TableName())
	if v, ok := param["uid"]; ok {
		db = db.Where("uid = ?", v)
	}
	if v, ok := param["name"]; ok {
		db = db.Where("name = ?", v)
	}
	if v, ok := param["age"]; ok {
		db = db.Where("age = ?", v)
	}

	var count int64

	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	var list []*User

	curPage, ok1 := param["cur_page"]
	pageLimit, ok2 := param["page_limit"]
	if ok1 && ok2 {
		db = db.Offset((curPage.(int) - 1) * pageLimit.(int)).Limit(pageLimit.(int))
	}

	if err := db.Find(&list).Error; err != nil {
		return list, count, err
	}

	return list, count, nil
}
