package classroom

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
	"smart.com/weixin/smart/logp"
	"smart.com/weixin/smart/utils"
)

func GetClassroomList(c *gin.Context) {
	mlogger  := logp.NewLogger("classromm")
	logger := mlogger.Named("getlist")

	body, err := c.GetRawData()
	if err != nil {
		logger.Errorf("get raw data fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	var inReq ClassroomRe
	err = json.Unmarshal(body, &inReq)
	if err != nil {
		logger.Errorf("unmarshal fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	var classroomList []ClassroomRe
	classroommodel := GetClassroomModel(inReq.CourseID)

	for _, tmp := range classroommodel {
		classroom := ClassroomRe{
			RoomID:tmp.RoomID,
			Name:tmp.Name,
			Desc:tmp.Desc,
			StudentNum:tmp.StudentNum,
			Status:tmp.Status,
			Start:tmp.Start,
			End:tmp.End,
			CourseID:tmp.CourseID,
			TeacherID:tmp.TeacherID,
			TeacherName:"tmp,it will get"}

		classroomList = append(classroomList, classroom)
	}

	c.JSON(http.StatusOK, gin.H{"classrooms":classroomList})
}

func AddClassroom(c *gin.Context) {
	mlogger  := logp.NewLogger("classroom")
	logger := mlogger.Named("add")

	body, err := c.GetRawData()
	if err != nil {
		logger.Errorf("get raw data fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	var inReq ClassroomRe
	err = json.Unmarshal(body, &inReq)
	if err != nil {
		logger.Errorf("unmarshal fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	logger.Infof("the classroom id is %d, the name is %d, the name is %s , the desc is %s, the course level is %s",inReq.RoomID, inReq.Name)

	var classroomModel ClassroomModel
	classroomModel.RoomID     = inReq.RoomID
	classroomModel.Name       = inReq.Name
	classroomModel.Desc       = inReq.Desc
	classroomModel.StudentNum = 0
	classroomModel.Status     = inReq.Status
	classroomModel.Start      = inReq.Start
	classroomModel.End        = ""
	classroomModel.CourseID   = inReq.CourseID
	classroomModel.TeacherID  = inReq.TeacherID

	db := utils.GetDB()
	err = db.Save(&classroomModel).Error
	if err != nil {
		logger.Errorf("save data fail %+v", err)
	}

	c.JSON(http.StatusOK, gin.H{})
}

func DeleteClassroom(c *gin.Context) {
	mlogger  := logp.NewLogger("classroom")
	logger := mlogger.Named("delete")

	body, err := c.GetRawData()
	if err != nil {
		logger.Errorf("get raw data fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	var inReq ClassroomRe
	err = json.Unmarshal(body, &inReq)
	if err != nil {
		logger.Errorf("unmarshal fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	logger.Infof("the room id is %d, the name is %s will be delete",inReq.CourseID, inReq.Name)

	db := utils.GetDB()
	db.Unscoped().Where("roomid = ?", inReq.RoomID).Delete(ClassroomModel{})

	c.JSON(http.StatusOK, gin.H{})
}

func ModifyClassroom(c *gin.Context) {
	mlogger  := logp.NewLogger("classroom")
	logger := mlogger.Named("modify")

	body, err := c.GetRawData()
	if err != nil {
		logger.Errorf("get raw data fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	var inReq ClassroomRe
	err = json.Unmarshal(body, &inReq)
	if err != nil {
		logger.Errorf("unmarshal fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	logger.Infof("the room id is %d, the name is %s , the desc is %s, the room status is %s, the teacher id is %d",inReq.RoomID, inReq.Name, inReq.Desc, inReq.Status, inReq.TeacherID)

	var classroomModel ClassroomModel
	db := utils.GetDB()
	classroomModel = ClassroomModel{}
	db.First(&classroomModel, "roomid = ?", inReq.RoomID)
	db.Model(&classroomModel).Updates(map[string]interface{}{"name":inReq.Name, "description":inReq.Desc, "status":inReq.Status, "teacherid":inReq.TeacherID})
	c.JSON(http.StatusOK, gin.H{})
}
