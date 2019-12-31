package courses

import (
)

type CoursesModelValidator struct {
	CourseID     uint    `json:"cid"`
	PID          uint    `json:"pid"`
	Name         string  `json:"cname"`
	Desc         string  `json:"cdes"`
	Vedio        string  `json:"cvedio"`
	CourseLevel  string  `json:"clevel"`   //L1 L2 L3 L4 ... L8
	Depth        int     `json:"depth"`
}

