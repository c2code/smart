package student

type StudentRe struct {
	StudentID    uint    `json:"sid"`
	UserID       uint    `json:"uid"`
	RoomID       uint    `json:"rid"`
	UserName     string  `json:"uname"`
	Email        string  `json:"email"`
	phone        string  `json:"phone"`
	RoomName     string  `json:"rname"`
	Level        string  `json:"level"`
	Ccid         uint    `json:"ccid"`
}
