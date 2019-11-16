package courses

import (
)

type CoursesModelValidator struct {
	CourseID     uint    `json:"cid"`
	PID          uint    `json:"pid"`
	Name         string  `json:"cname"`
	Desc         string  `json:"cdes"`
	Vedio        string  `json:"cvedio"`
}


type HomeworkModelValidator struct {
	Article struct {
		Title       string   `json:"title" binding:"exists,min=4"`
		Status      string   `json:"description"`
		Address     string   `json:"body"`
	} `json:"homework"`
	homeworkModel HomeworkModel `json:"-"`
}

func NewHowmeworkModelValidator() HomeworkModelValidator {
	return HomeworkModelValidator{}
}

