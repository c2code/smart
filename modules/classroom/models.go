package classroom

import ( "github.com/jinzhu/gorm"
	"smart.com/weixin/smart/utils"
)

// Models should only be concerned with database schema, more strict checking should be put in validator.
//
// More detail you can find here: http://jinzhu.me/gorm/models.html#model-definition
//
// HINT: If you want to split null and "", you should use *string instead of string.
type ClassroomModel struct {
	gorm.Model
	RoomID       uint    `gorm:"column:roomid;unique_index"`
	Name         string  `gorm:"column:name"`
	Desc         string  `gorm:"column:description;size:1024"`
	StudentNum   uint    `gorm:"column:studentnumber"`
	Status       string  `gorm:"column:status"`  //on or off
	Start        string  `gorm:"column:start`    //start date
	End          string  `gorm:"column:end`      //end date
	CourseID     uint    `gorm:"column:courseid"`
	TeacherID    uint    `gorm:"column:teacherid"`
}

// Migrate the schema of database if needed
func AutoMigrate() {
	db := utils.GetDB()
	db.AutoMigrate(&ClassroomModel{})
}

func GetClassroomModel(cid uint) []ClassroomModel{
	var classroomlist []ClassroomModel
	db := utils.GetDB()
	db.Where("courseid = ?", cid).Find(&classroomlist)
	return classroomlist
}
