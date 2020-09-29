package main

import (
	"fmt"
	"time"

	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

// SampleOne has no constraint and no index during create table
type SampleOne struct {
	ID        uint
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// SampleTwo has constraint to check but no index during create table
type SampleTwo struct {
	ID    uint
	Name  string `gorm:"check:name_checker,name <> 'jinzhu'"`
	Email string `gorm:"check:email <> 'diku@gmail.com'"`
}

type Employee struct {
	Name      string `gorm:"index"`
	FirstName string `gorm:"index:idx_name,unique"`
	LastName  string `gorm:"index:idx_name,unique"`
	Username  string `gorm:"index:,sort:desc,collate:utf8,type:minmax,length:10,where:name3 != 'jinzhu'"`
	Password  string `gorm:"uniqueIndex"`
	Age       int64  `gorm:"index:,class:FULLTEXT,comment:hello \\, world,where:age > 10"`
	Age2      int64  `gorm:"index:,expression:ABS(age)"`
}

type Employer struct {
	Name      string `gorm:"index,comment:'thisisnameindex';check:name_checker,name <> 'jinzhu'"`
	FirstName string `gorm:"index:idx_name,unique; check:concat(first_name, ' ', last_name) <> 'jinzhu zhang'"`
	LastName  string `gorm:"index:idx_name,unique,comment:'thisislastnameindex'"`
	Username  string `gorm:"index:,type:minmax,length:10,where:name3 != 'jinzhu'"`
	Age       int64  `gorm:"index:,class:FULLTEXT,comment:'helloworld',where:age > 10;check:age > 20"`
	WorkYear  int64  `gorm:"index:;check:workchecker,work_year > (age + 18)"`
	Comment   string `gorm:"comment:'this is a comment field'"`
}

type DefaultAndComment struct {
	DefaultStr    string `gorm:"default:'ThisIsDefaultEyy!'"`
	DefaultNum64  int64  `gorm:"default:199000141"`
	DefaultNumm64 int64  `gorm:"default:toInt64(199000143)"`
	DefaultNum    int    `gorm:"default:1990000142"`
	CommentF      string `gorm:"comment:'this is a comment field ay'"`
}

const DSNf = "tcp://%s:%s?database=%s&username=%s&password=%s&read_timeout=10&write_timeout=20&debug=%t"

func main() {
	var (
		Host  = "localhost"
		Port  = "9000"
		Debug = true

		DBName    = "testdb"
		DBUser    = "default"
		DBPasword = ""
	)

	dsn := fmt.Sprintf(DSNf, Host, Port, DBName, DBUser, DBPasword, Debug)
	conn, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := conn.AutoMigrate(&SampleOne{}, &SampleTwo{}, &Employee{}, &Employer{}, &DefaultAndComment{}); err != nil {
		fmt.Println("errors?", err)
	}
}
