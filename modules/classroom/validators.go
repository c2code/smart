package classroom

type ClassroomRe struct {
	RoomID        uint    `json:"roomid"`
	Name          string  `json:"rname"`
	Desc          string  `json:"rdesc"`
	StudentNum    uint    `json:"stdnum"`
	Status        string  `json:"rstatus"`
	Start         string  `json:"start"`   //L1 L2 L3 L4 ... L8
	End           string  `json:"end"`
	CourseID      uint    `json:"cid"`
	TeacherID     uint    `json:"tid"`
	TeacherName   string  `json:"tname"`
}
