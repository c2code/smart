package teacher

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
	"smart.com/weixin/smart/logp"
	"smart.com/weixin/smart/utils"
	"smart.com/weixin/smart/modules/users"
	"strconv"
	"smart.com/weixin/smart/modules/classroom"
)


func GetTeacher(c *gin.Context) {
	mlogger  := logp.NewLogger("teacher")
	logger := mlogger.Named("get")

	tid := c.Query("tid")

	logger.Infof("get the teachers, the teacher id is %s",tid)

	if (tid != "") {
		tmp := TeacherModel{}
		db := utils.GetDB()
		db.Last(&tmp)
		myuser := users.UserModel{}
		db.Where("ID=?", tmp.UserID).Find(&myuser)
		teacher := TeacherRe{
			TeacherID:tmp.TeacherID,
			UserID:tmp.UserID,
			Status:tmp.Status,
			UserName:myuser.Username,
			Email:myuser.Email,
			Phone:"",
			RoomCount:tmp.RoomCount}
		c.JSON(http.StatusOK, gin.H{"teacher":teacher})
		return
	}

	var teacherList []TeacherRe
	studentmodel := GetTeacherModelList()

	db := utils.GetDB()
	for _, tmp := range studentmodel {
		myuser := users.UserModel{}
		db.Where("ID=?", tmp.UserID).Find(&myuser)
		teacher := TeacherRe{
			TeacherID:tmp.TeacherID,
			UserID:tmp.UserID,
			Status:tmp.Status,
			UserName:myuser.Username,
			Email:myuser.Email,
			Phone:"",
			RoomCount:tmp.RoomCount}

		teacherList = append(teacherList, teacher)
	}

	c.JSON(http.StatusOK, gin.H{"teachers":teacherList})
}

func AddTeacher(c *gin.Context) {
	mlogger  := logp.NewLogger("teacher")
	logger := mlogger.Named("add")

	body, err := c.GetRawData()
	if err != nil {
		logger.Errorf("get raw data fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	var inReq TeacherRe
	err = json.Unmarshal(body, &inReq)
	if err != nil {
		logger.Errorf("unmarshal fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	logger.Infof("the teacher id is %d, user id is %d",inReq.TeacherID,inReq.UserID)


	var teacherModel TeacherModel

	db := utils.GetDB()
	db.Where("userid=?",inReq.UserID).Find(&teacherModel)

	if (teacherModel.TeacherID != 0){
		logger.Infof("the teacher id is %d, user id is %d, has been added!",inReq.TeacherID,inReq.UserID)
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	teacherModel.TeacherID    = inReq.TeacherID
	teacherModel.UserID       = inReq.UserID
	teacherModel.Status       = inReq.Status
	teacherModel.RoomCount    = 0

	err = db.Save(&teacherModel).Error
	if err != nil {
		logger.Errorf("save data fail %+v", err)
	}

	c.JSON(http.StatusOK, gin.H{})
}

func DeleteTeacher(c *gin.Context) {
	mlogger  := logp.NewLogger("teacher")
	logger := mlogger.Named("delete")

	body, err := c.GetRawData()
	if err != nil {
		logger.Errorf("get raw data fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	var inReq TeacherRe
	err = json.Unmarshal(body, &inReq)
	if err != nil {
		logger.Errorf("unmarshal fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	logger.Infof("the teacher id is %d, user id is %d will be deleted ",inReq.TeacherID,inReq.UserID)

	db := utils.GetDB()
	db.Unscoped().Where("teacherid = ?", inReq.TeacherID).Delete(TeacherModel{})

	c.JSON(http.StatusOK, gin.H{})
}

func ModifyTeacher(c *gin.Context) {
	mlogger  := logp.NewLogger("teacher")
	logger := mlogger.Named("modify")

	roomid := c.Query("rid")
	rid,_ := strconv.Atoi(roomid)

	body, err := c.GetRawData()
	if err != nil {
		logger.Errorf("get raw data fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	var inReq TeacherRe
	err = json.Unmarshal(body, &inReq)
	if err != nil {
		logger.Errorf("unmarshal fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	logger.Infof("the teacher id is %d, user id is %d, the room id is %d will be modified ",inReq.TeacherID,inReq.UserID,rid)

	db := utils.GetDB()

	classroomModel := classroom.ClassroomModel{}
	db.First(&classroomModel, "roomid = ?", rid)
	if (classroomModel.TeacherID == inReq.TeacherID){
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	//上次绑定的教师的班级数量减一
	oldteacher := TeacherModel{}
	db.Where("teacherid = ?", classroomModel.TeacherID).Find(&oldteacher)
	if (oldteacher.RoomCount != 0){
		count := oldteacher.RoomCount - 1
		db.Model(&oldteacher).Updates(map[string]interface{}{"rcount":count})
	}

	//当前绑定的教师的班级数量加1
	newteacher := TeacherModel{}
	db.Where("teacherid = ?", inReq.TeacherID).Find(&newteacher)
	count := newteacher.RoomCount + 1
	db.Model(&newteacher).Updates(map[string]interface{}{"rcount":count})

	//变更班级老师ID
	db.Model(&classroomModel).Updates(map[string]interface{}{"teacherid":inReq.TeacherID})

	c.JSON(http.StatusOK, gin.H{})
}