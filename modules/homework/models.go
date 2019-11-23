package homework

import ( "github.com/jinzhu/gorm"
)

// Models should only be concerned with database schema, more strict checking should be put in validator.
//
// More detail you can find here: http://jinzhu.me/gorm/models.html#model-definition
//
// HINT: If you want to split null and "", you should use *string instead of string.
type HomeWorkModel struct {
	gorm.Model
	HomeWoekID   uint    `gorm:"column:homeworkid;unique_index"`
	Status       string  `gorm:"column:status"`  // on or off
	Addr         string  `gorm:"column:addr"`
	UserID       uint    `gorm:"column:userid"`
	CourseID     uint    `gorm:"column:courseid"`
	RoomID       uint    `gorm:"column:roomid"`
	Comment      string  `gorm:"column:comment;size:1024"`
}
