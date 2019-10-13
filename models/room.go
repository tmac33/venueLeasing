package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Room struct {
	RoomId      int `orm:"pk;auto"` //pk primary key
	Name        string
	PosNum      int `orm:"size(10)"`
	Desc        string
	StartTime   string
	EndTime     string
	Status      int       `orm:"default(1)"`
	UpadateTime time.Time `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModelWithPrefix("b_", new(Room))
}
