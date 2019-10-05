package controllers

import (
	"time"
	"venueleasing/models"

	"github.com/astaxie/beego/validation"

	"github.com/astaxie/beego/orm"
)

type BookingController struct {
	Base
}

func (c *BookingController) URLMapping() {
	c.Mapping("List", c.List)
	c.Mapping("Time", c.Time)
	c.Mapping("Add", c.Add)
	c.Mapping("GetList", c.GetList)
}

//new appointment
// @router /booking/new [get]
func (c *BookingController) Show() {
	o := orm.NewOrm()
	var rooms []*models.Room
	o.QueryTable("b_room").Filter("status", 1).All(&rooms)
	c.Data["RoomList"] = rooms
	var designs []*models.Design
	o.QueryTable("b_design").Filter("status", 1).All(&designs)
	c.Data["DesignList"] = designs
	c.TplName = "NewBooking.html"
}

//appointment list
// @router /booking/list [get]
func (c *BookingController) List() {
	c.TplName = "ShowBooking.html"
}

//obtain the duration can leased
// @router /booking/time [post]
func (c *BookingController) Time() {
	o := orm.NewOrm()
	id, _ := c.GetInt("room_id")
	date := c.GetString("date")
	room := models.Room{RoomId: id}
	o.Read(&room)
	if room.Status != 1 {
		c.AjaxErr("error message")
	}
	var list orm.ParamsList
	o.QueryTable("b_booking").Filter("room_id", id).Filter("booking_date", date).ValuesFlat(&list, "booking_time_type")
	startList := seralizeTime(room.StartTime)
	endList := seralizeTime(room.EndTime)
	res := make([]interface{}, 0)
	for k, v := range startList {
		if !inArr(k, list) {
			single := make(map[string]interface{})
			single["type"] = k
			single["str"] = v + "~" + endList[k]
			res = append(res, single)
		}
	}
	json := make(map[string]interface{})
	json["status"] = 1
	json["msg"] = "success"
	json["data"] = res
	c.Data["json"] = json
	c.ServeJSON()
}

func inArr(str int, arr []interface{}) bool {
	for _, v := range arr {
		if v == int64(str) {
			return true
		}
	}
	return false
}

//add
// @router /booking/add [post]
func (c *BookingController) Add() {
	//1.idetify weather the venue is exited，weather is leaseable,if can get the timetable
	o := orm.NewOrm()
	var newBooking models.Booking
	var room models.Room
	id, _ := c.GetInt("room_id")
	rooms := models.Room{RoomId: id}
	o.Read(&rooms)
	newBooking.Room = &rooms
	err := o.QueryTable("b_room").Filter("room_id", id).Filter("status", 1).One(&room)
	if err != nil || room.RoomId == 0 {
		c.AjaxErr("the room is not exits!")
	}

	//2.check the room's time
	timeArr := seralizeTime(room.StartTime)
	useType, err := c.GetInt("use_type")
	newBooking.BookingTimeType = useType
	if err != nil || useType < 0 || useType > len(timeArr)-1 {
		c.AjaxErr("Illegal lease time")
	}

	//3.check weather already leased
	valid := validation.Validation{}
	var booking models.Booking
	useDate := c.GetString("use_date")
	date, _ := time.Parse("2006-01-02", useDate)
	newBooking.BookingDate = date
	newBooking.BookingContent = c.GetString("content")
	newBooking.BookingOrg = c.GetString("org")
	newBooking.BookingUseNum, _ = c.GetInt("num")
	newBooking.BookingUserName = c.GetString("name")
	newBooking.BookingUserTel = c.GetString("tel")
	valid.Required(newBooking.BookingDate, "type")
	valid.Required(newBooking.BookingContent, "content")
	valid.Required(newBooking.BookingUseNum, "num")
	valid.Required(newBooking.BookingUserName, "name")
	valid.Required(newBooking.BookingUserTel, "tel")
	valid.Required(newBooking.BookingOrg, "org")
	if valid.HasErrors() {
		c.AjaxErr("error information, please fill again ")
	}
	err = o.QueryTable("b_booking").Filter("booking_date", useDate).Filter("booking_time_type", useType).One(&booking)
	if err == nil {
		c.AjaxErr("This duration has been rented!")
	}

	//format information
	var design []*models.Design
	o.QueryTable("b_design").Filter("status", 1).All(&design)
	var ext string
	for _, v := range design {
		ext += v.Name + ":" + c.GetString(v.Key) + "\n"
	}
	newBooking.BookingExt = ext
	//6、insert
	_, err = o.Insert(&newBooking)
	if err != nil {
		c.AjaxErr("lease failed")
	}
	//7.result
	c.AjaxOk("lease successfully, please wait auit!")
}

//obtain appointment list
func (c *BookingController) GetList() {
	var page, limit int
	c.Ctx.Input.Bind(&page, "page")
	c.Ctx.Input.Bind(&limit, "rows")
	o := orm.NewOrm()
	var bookings []*models.Booking
	o.QueryTable("b_booking").RelatedSel().OrderBy("booking_status").Limit(limit, (page-1)*10).All(&bookings)
	for _, v := range bookings {
		timeStart := seralizeTime(v.Room.StartTime)[v.BookingTimeType]
		timeEnd := seralizeTime(v.Room.EndTime)[v.BookingTimeType]
		v.Date = v.BookingDate.Format("2006-01-02") + " " + timeStart + "~" + timeEnd
		v.RoomName = v.Room.Name
	}
	json := make(map[string]interface{})
	json["rows"] = bookings
	json["code"] = 0
	json["msg"] = "success"
	cnt, _ := o.QueryTable("b_booking").Count() // SELECT COUNT(*) FROM USER
	json["total"] = cnt
	c.Data["json"] = json
	c.ServeJSON()
}
