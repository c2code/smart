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
)


func GetStudent(c *gin.Context) {
	mlogger  := logp.NewLogger("student")
	logger := mlogger.Named("get")

	roomid := c.Query("rid")

	logger.Infof("get the students, the classroom id is %s",roomid)

	if (roomid == "") {
		tmp := StudentModel{}
		db := utils.GetDB()
		db.Last(&tmp)
		myuser := users.UserModel{}
		db.Where("ID=?", tmp.UserID).Find(&myuser)
		student := StudentRe{
			StudentID:tmp.StudentID,
			UserID:tmp.UserID,
			RoomID:tmp.RoomID,
			UserName:myuser.Username,
			Email:myuser.Email,
			phone:"",
			RoomName:""}
		c.JSON(http.StatusOK, gin.H{"student":student})
		return
	}

	rid,_ := strconv.Atoi(roomid)
	var studentList []StudentRe
	studentmodel := GetStudentModelList(rid)

	db := utils.GetDB()
	myclassroom := classroom.ClassroomModel{}
	db.Where("roomid=?", rid).Find(&myclassroom)

	for _, tmp := range studentmodel {
		myuser := users.UserModel{}
		db.Where("ID=?", tmp.UserID).Find(&myuser)
		student := StudentRe{
			StudentID:tmp.StudentID,
			UserID:tmp.UserID,
			RoomID:tmp.RoomID,
			UserName:myuser.Username,
			Email:myuser.Email,
			phone:"",
			RoomName:myclassroom.Name}

		studentList = append(studentList, student)
	}

	c.JSON(http.StatusOK, gin.H{"students":studentList})
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

	db.Where("userid=? AND roomid=?",inReq.UserID, inReq.RoomID).Find(&studentModel)

	if (studentModel.StudentID != 0){
		logger.Infof("the student id is %d, user id is %d, has been added into the classroom!",inReq.StudentID,inReq.UserID)
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	studentModel.StudentID    = inReq.StudentID
	studentModel.UserID       = inReq.UserID
	studentModel.RoomID       = inReq.RoomID

	err = db.Save(&studentModel).Error
	if err != nil {
		logger.Errorf("save data fail %+v", err)
	}

	classroomModel := classroom.ClassroomModel{}
	db.First(&classroomModel, "roomid = ?", inReq.RoomID)
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

	c.JSON(http.StatusOK, gin.H{})
}
