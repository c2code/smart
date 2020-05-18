package student

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
	"smart.com/weixin/smart/logp"
	"smart.com/weixin/smart/utils"
	"strconv"
	"smart.com/weixin/smart/modules/users"
	"smart.com/weixin/smart/modules/classroom"
	"smart.com/weixin/smart/modules/courses"
	"github.com/biezhi/gorm-paginator/pagination"
)


func GetStudent(c *gin.Context) {
	mlogger  := logp.NewLogger("student")
	logger := mlogger.Named("get")

	roomid := c.Query("rid")
	uid := c.Query("uid")
	clevel := c.Query("clevel")
	page := c.Query("page")

	logger.Infof("get the students, the classroom id is %s, uid is %d, clevel is %s,page is %s.",roomid, uid, clevel, page)

	if (page != "") {
		/*分页处理，每页两条，取第二页内容*/
		mypage, _ := strconv.Atoi(page)
		var studentList []StudentRe
		var studentmodel []StudentModel
		db := utils.GetDB()

		pagination.Paging(&pagination.Param{
			DB:      db.Where("id > ?", 0),
			Page:    mypage,
			Limit:   20,
			OrderBy: []string{"id"}, //降序是id desc
			ShowSQL: true,
		}, &studentmodel)

		for _, tmp := range studentmodel {
			myclassroom := classroom.ClassroomModel{}
			db.Where("roomid=?", tmp.RoomID).Find(&myclassroom)
			myuser := users.UserModel{}
			db.Where("ID=?", tmp.UserID).Find(&myuser)
			student := StudentRe{
				StudentID: tmp.StudentID,
				UserID:    tmp.UserID,
				RoomID:    tmp.RoomID,
				UserName:  myuser.Username,
				Email:     myuser.Email,
				phone:     "",
				RoomName:  myclassroom.Name,
				Level:     tmp.Level,
				Ccid:      tmp.Ccid}

			studentList = append(studentList, student)
		}

		c.JSON(http.StatusOK, gin.H{"students": studentList})

	}else {

		if (uid != "" && clevel != "") {
			tmp := StudentModel{}
			db := utils.GetDB()
			userid, _ := strconv.Atoi(uid)

			db.Where("userid=? AND level=?", userid, clevel).Find(&tmp)
			myuser := users.UserModel{}
			db.Where("ID=?", tmp.UserID).Find(&myuser)
			student := StudentRe{
				StudentID: tmp.StudentID,
				UserID:    tmp.UserID,
				RoomID:    tmp.RoomID,
				UserName:  myuser.Username,
				Email:     myuser.Email,
				phone:     "",
				RoomName:  "",
				Level:     tmp.Level,
				Ccid:      tmp.Ccid}
			c.JSON(http.StatusOK, gin.H{"student": student})
			return
		}

		if (roomid == "") {
			tmp := StudentModel{}
			db := utils.GetDB()
			db.Last(&tmp)
			myuser := users.UserModel{}
			db.Where("ID=?", tmp.UserID).Find(&myuser)
			student := StudentRe{
				StudentID: tmp.StudentID,
				UserID:    tmp.UserID,
				RoomID:    tmp.RoomID,
				UserName:  myuser.Username,
				Email:     myuser.Email,
				phone:     "",
				RoomName:  "",
				Level:     tmp.Level,
				Ccid:      tmp.Ccid}
			c.JSON(http.StatusOK, gin.H{"student": student})
			return
		}

		rid, _ := strconv.Atoi(roomid)
		var studentList []StudentRe
		studentmodel := GetStudentModelList(rid)

		db := utils.GetDB()
		myclassroom := classroom.ClassroomModel{}
		db.Where("roomid=?", rid).Find(&myclassroom)

		for _, tmp := range studentmodel {
			myuser := users.UserModel{}
			db.Where("ID=?", tmp.UserID).Find(&myuser)
			student := StudentRe{
				StudentID: tmp.StudentID,
				UserID:    tmp.UserID,
				RoomID:    tmp.RoomID,
				UserName:  myuser.Username,
				Email:     myuser.Email,
				phone:     "",
				RoomName:  myclassroom.Name,
				Level:     tmp.Level,
				Ccid:      tmp.Ccid}

			studentList = append(studentList, student)
		}

		c.JSON(http.StatusOK, gin.H{"students": studentList})
	}
}

func GetStudentCount(c *gin.Context) {
	mlogger  := logp.NewLogger("student")
	logger := mlogger.Named("get")

	roomid := c.Query("rid")
	var count uint
	var studentmodel []StudentModel

	logger.Infof("get the students, the classroom id is %s.",roomid)

	db := utils.GetDB()
	/*分页处理，每页两条，取第二页内容*/
	if (roomid != "") {
		rid, _ := strconv.Atoi(roomid)
		db.Where("roomid=?", rid).Find(&studentmodel).Count(&count)
		c.JSON(http.StatusOK, gin.H{"count": count})
		return
	}

	db.Find(&studentmodel).Count(&count)
	c.JSON(http.StatusOK, gin.H{"count": count})
	return
}

func AddStudent(c *gin.Context) {
	mlogger  := logp.NewLogger("student")
	logger := mlogger.Named("add")

	body, err := c.GetRawData()
	if err != nil {
		logger.Errorf("get raw data fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	var inReq StudentRe
	err = json.Unmarshal(body, &inReq)
	if err != nil {
		logger.Errorf("unmarshal fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	logger.Infof("the student id is %d, user id is %d, room id is %d ",inReq.StudentID,inReq.UserID,inReq.RoomID)

	db := utils.GetDB()

	var studentModel StudentModel
	var studentlist []StudentModel
	classroomModel := classroom.ClassroomModel{}
	var classroom classroom.ClassroomModel

	db.First(&classroomModel, "roomid = ?", inReq.RoomID)

	db.Where("userid=?",inReq.UserID).Find(&studentlist)

	for _, tmp := range studentlist {
		db.Where("roomid=?",tmp.RoomID).Find(&classroom)

		if (classroom.CourseID == classroomModel.CourseID){
			logger.Infof("the student id is %d, user id is %d, has been added into the classroom %d!",inReq.StudentID,inReq.UserID, tmp.RoomID)
			c.JSON(http.StatusInternalServerError, gin.H{"result": false, "error": "the student has been add other classroom of the course"})
			return
		}
	}

	var courseModel courses.CourseModel
	var courseList  []courses.CourseModel
	var userModel users.UserModel
	db.Where("courseid=?",classroomModel.CourseID).Find(&courseModel)
	db.Where("clevel=? AND depth=?",courseModel.CourseLevel, 3).Find(&courseList)

	laststudent := StudentModel{}
	db.Last(&laststudent)

	studentModel.StudentID    = laststudent.StudentID + 1
	studentModel.UserID       = inReq.UserID
	studentModel.RoomID       = inReq.RoomID
	studentModel.Level        = courseModel.CourseLevel
	studentModel.Ccid         = courseList[0].CourseID

	err = db.Save(&studentModel).Error
	if err != nil {
		logger.Errorf("save data fail %+v", err)
	}

	level := courseModel.CourseLevel //获取位图信息，通过level的定位位图信息
	tmp := level[1:len(level)]
	i, _ := strconv.Atoi(tmp)
	count := 1 << uint(i - 1)

	db.Where("ID=?",inReq.UserID).Find(&userModel)
	userModel.Rights = userModel.Rights | count
	db.Model(&userModel).Updates(map[string]interface{}{"rights":userModel.Rights})

	number := classroomModel.StudentNum + 1;
	db.Model(&classroomModel).Updates(map[string]interface{}{"studentnumber":number})

	c.JSON(http.StatusOK, gin.H{})
}

func DeleteStudent(c *gin.Context) {
	mlogger  := logp.NewLogger("student")
	logger := mlogger.Named("delete")

	body, err := c.GetRawData()
	if err != nil {
		logger.Errorf("get raw data fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	var inReq StudentRe
	err = json.Unmarshal(body, &inReq)
	if err != nil {
		logger.Errorf("unmarshal fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	logger.Infof("the student id is %d, user id is %d, room id is %d ",inReq.StudentID,inReq.UserID,inReq.RoomID)

	db := utils.GetDB()
	db.Unscoped().Where("studentid = ?", inReq.StudentID).Delete(StudentModel{})

	classroomModel := classroom.ClassroomModel{}
	db.First(&classroomModel, "roomid = ?", inReq.RoomID)
	number := classroomModel.StudentNum - 1;
	db.Model(&classroomModel).Updates(map[string]interface{}{"studentnumber":number})

	var courseModel courses.CourseModel
	var userModel users.UserModel
	db.Where("courseid=?",classroomModel.CourseID).Find(&courseModel)

	level := courseModel.CourseLevel //获取位图信息，通过level的定位位图信息
	tmp := level[1:len(level)]
	i, _ := strconv.Atoi(tmp)
	count := 1 << uint(i - 1)

	db.Where("ID=?",inReq.UserID).Find(&userModel)
	userModel.Rights = userModel.Rights & (^count)
	db.Model(&userModel).Updates(map[string]interface{}{"rights":userModel.Rights})

	c.JSON(http.StatusOK, gin.H{})
}
