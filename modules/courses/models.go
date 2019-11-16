package courses

import ( "github.com/jinzhu/gorm"
    "smart.com/weixin/smart/modules/users"
	"smart.com/weixin/smart/utils"
	"fmt"
	"os"
	"encoding/json"
)

// Models should only be concerned with database schema, more strict checking should be put in validator.
//
// More detail you can find here: http://jinzhu.me/gorm/models.html#model-definition
//
// HINT: If you want to split null and "", you should use *string instead of string.
type CourseModel struct {
	gorm.Model
	CourseID     uint    `gorm:"column:courseid;unique_index"`
	PID          uint    `gorm:"column:parentid"`
	Name         string  `gorm:"column:name"`
	Desc         string  `gorm:"column:description;size:1024"`
	Vedio        string  `gorm:"column:vedio"`
}

type HomeworkModel struct {
	gorm.Model
	Slug        uint                    `gorm:"column:slug;unique_index"`
	Title       string
	Author      CourseUserModel
	AuthorID    uint
	CourseID    uint                    `gorm:"column:courseid"`
	Status      string                  `gorm:"column:status"`   //modify, commit, comment
	Address     string                  `gorm:"column:address"`  //作品路径
}

type CourseUserModel struct {
	gorm.Model
	UserModel          users.UserModel
	UserModelID        uint
	CourseModels       []CourseModel   `gorm:"many2many:courseid"`   //已开通的课程
	HomeWorksModels    []HomeworkModel `gorm:"ForeignKey:AuthorID"` //提交的作业
}

// Migrate the schema of database if needed
func AutoMigrate() {
	db := utils.GetDB()

	db.AutoMigrate(&CourseModel{})
	db.AutoMigrate(&HomeworkModel{})
	db.AutoMigrate(&CourseUserModel{})
}

func SaveOne(data interface{}) error {
	db := utils.GetDB()
	err := db.Save(data).Error
	return err
}

func InitCouses() {
	var courseslist []CoursesModelValidator

	config_file, err := os.Open("./courses.json")
	defer config_file.Close()

	if err != nil {
		fmt.Printf("Failed to open config file ./courses.json : %s\n",err)
		return
	}

	fi, _ := config_file.Stat()

	if fi.Size() == 0 {
		fmt.Print("config file (./courses.json) is empty, skipping")
		return
	}

	buffer := make([]byte, fi.Size())
	_, err = config_file.Read(buffer)

	buffer = []byte(os.ExpandEnv(string(buffer))) //特殊

	err = json.Unmarshal(buffer, &courseslist) //解析json格式数据
	if err != nil {
		fmt.Printf("Failed unmarshalling json: %s\n", err)
		return
	}

	//fmt.Printf("%+v", courseslist)

	var courseModel CourseModel

	db := utils.GetDB()
	courseModel = CourseModel{}
	db.First(&courseModel, "courseid = ?", "1")

	if courseModel.CourseID == 1 {
		fmt.Print("The course models database has been init!")
		return
	}

	for _, course := range courseslist {
		fmt.Printf(" %d,%d,%s,%s\n", course.CourseID,course.PID,course.Name,course.Desc)
		courseModel = CourseModel{}
		courseModel.CourseID = course.CourseID
		courseModel.PID      = course.PID
		courseModel.Name     = course.Name
		courseModel.Desc     = course.Desc
		if err := SaveOne(&courseModel); err != nil {
			fmt.Printf("database err:%+v", err)
			continue
		}
	}

	return

}

func GetCourseModel() []CourseModel{
	var courselist []CourseModel
	db := utils.GetDB()
	tx := db.Begin()     //开启事物处理
	tx.Where(CourseModel{}).Offset(0).Limit(10000).Find(&courselist) //获取course_models表中的前10000条数据
	tx.Commit()         //结束事物处理
	return courselist
}


