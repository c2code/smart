package homework

import ( "github.com/jinzhu/gorm"
	"smart.com/weixin/smart/utils"
)

// Models should only be concerned with database schema, more strict checking should be put in validator.
//
// More detail you can find here: http://jinzhu.me/gorm/models.html#model-definition
//
// HINT: If you want to split null and "", you should use *string instead of string.
type HomeWorkModel struct {
	gorm.Model
	HomeWoekID            uint    `gorm:"column:homeworkid;unique_index"`
	Status                string  `gorm:"column:status"`  // on or off
	Addr                  string  `gorm:"column:addr"`
	UserID                uint    `gorm:"column:userid"`
	RoomID                uint    `gorm:"column:roomid"`
	CourseID              uint    `gorm:"column:cid"`
	Description           string  `gorm:"column:hdesc"`
	Comment               string  `gorm:"column:comment;size:1024"`
}

// Migrate the schema of database if needed
func AutoMigrate() {
	db := utils.GetDB()
	db.AutoMigrate(&HomeWorkModel{})
}

func GetLastHomeworkModel() HomeWorkModel{
	var homework HomeWorkModel
	db := utils.GetDB()
	db.Last(&homework)
	return homework
}

func GetHomeworkModelListByUid(uid int) []HomeWorkModel{
	var homeworklist []HomeWorkModel
	db := utils.GetDB()
	db.Where("userid=?", uid).Find(&homeworklist)
	return homeworklist
}

func GetHomeworkModelListByUidRid(uid int, rid int) []HomeWorkModel{
	var homeworklist []HomeWorkModel
	db := utils.GetDB()
	db.Where("userid=? AND roomid=?", uid,rid).Find(&homeworklist)
	return homeworklist
}

func GetHomeworkModelListByUidCid(uid int, cid int) HomeWorkModel{
	var homework HomeWorkModel
	db := utils.GetDB()
	db.Where("userid=? AND cid=?", uid, cid).Find(&homework)
	return homework
}