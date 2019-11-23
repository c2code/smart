package homework

type HomeWorkRe struct {
	HomeWoekID   uint    `json:"hid"`
	Status       string  `json:"hstatus"`  // on or off
	Addr         string  `json:"haddr"`
	UserID       uint    `json:"uid"`
	CourseID     uint    `json:"cid"`
	RoomID       uint    `json:"rid"`
	Comment      string  `json:"comment"`
}