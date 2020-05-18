package student

type StudentRe struct {
	StudentID    uint    `json:"sid"`   //学生id
	UserID       uint    `json:"uid"`   //用户id
	RoomID       uint    `json:"rid"`   //班级id
	UserName     string  `json:"uname"` //用户名
	Email        string  `json:"email"`
	phone        string  `json:"phone"`
	RoomName     string  `json:"rname"` //班级名
	Level        string  `json:"level"` //课程级别
	Ccid         uint    `json:"ccid"`  //课程id
}
