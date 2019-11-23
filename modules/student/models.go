package student

import ( "github.com/jinzhu/gorm"
)

// Models should only be concerned with database schema, more strict checking should be put in validator.
//
// More detail you can find here: http://jinzhu.me/gorm/models.html#model-definition
//
// HINT: If you want to split null and "", you should use *string instead of string.
type StudentModel struct {
	gorm.Model
	StudentID    uint    `gorm:"column:studentid;unique_index"`
	UserID       uint    `gorm:"column:userid"`
	RoomID       uint    `gorm:"column:roomid"`
}