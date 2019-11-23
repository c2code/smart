package teacher

type TeacherRe struct {
	TeacherID    uint    `json:"tid"`
	Status       string  `json:"tstatus"`  // on or off
	UserID       uint    `json:"uid"`
	UserName     uint    `json:"uname"`
	Email        string  `json:"email"`
	Phone        string  `json:"phone"`
}
