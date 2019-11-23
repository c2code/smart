package teacher

import ( "github.com/jinzhu/gorm"
)

// Models should only be concerned with database schema, more strict checking should be put in validator.
//
// More detail you can find here: http://jinzhu.me/gorm/models.html#model-definition
//
// HINT: If you want to split null and "", you should use *string instead of string.
type TeacherModel struct {
	gorm.Model
	TeacherID    uint    `gorm:"column:teacherid;unique_index"`
	Status       string  `gorm:"column:status"`  // on or off
	UserID       uint    `gorm:"column:userid"`
	RoomID       uint    `gorm:"column:roomid"`
}