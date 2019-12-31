package student

import ( "github.com/jinzhu/gorm"
	"smart.com/weixin/smart/utils"
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
	Level        string  `gorm:"column:level"`
	Ccid         uint    `gorm:"column:ccid"`
}

// Migrate the schema of database if needed
func AutoMigrate() {
	db := utils.GetDB()
	db.AutoMigrate(&StudentModel{})
}

func GetStudentModelList(rid int) []StudentModel{
	var studentlist []StudentModel
	db := utils.GetDB()
	db.Where("roomid = ?", rid).Find(&studentlist)
	return studentlist
}