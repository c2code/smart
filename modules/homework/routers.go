package homework

import (
	"github.com/gin-gonic/gin"
	"smart.com/weixin/smart/logp"
	"net/http"
	"strconv"
)

func GetHomework(c *gin.Context) {
	mlogger  := logp.NewLogger("homework")
	logger := mlogger.Named("get")
	var homeworkList []HomeWorkRe
	var homework     HomeWorkRe

	uid := c.Query("uid")
	cid := c.Query("cid")
	rid := c.Query("rid")

	logger.Infof("get the homework, the user id = %s, cid=%s, rid=%s ",uid ,cid, rid)

	//获取最新的homeworkID，唯一
	if (uid == ""){
		homeworkModel := GetLastHomeworkModel()
		homework.HomeWoekID  = homeworkModel.HomeWoekID
		homework.Status      = homeworkModel.Status
		homework.Addr        = homeworkModel.Addr
		homework.UserID      = homeworkModel.UserID
		homework.CourseID    = homeworkModel.CourseID
		homework.RoomID      = homeworkModel.HomeWoekID
		homework.Description = homeworkModel.Description
		homework.Comment     = homeworkModel.Comment
		c.JSON(http.StatusOK, gin.H{"homework":homework})
		return
	}

	//获取该用户在该课程下的作业，唯一
	if (cid != ""){
		userid,_ := strconv.Atoi(uid)
		courseid,_ := strconv.Atoi(cid)
		homeworkModel := GetHomeworkModelListByUidCid(userid, courseid)

		homework.HomeWoekID  = homeworkModel.HomeWoekID
		homework.Status      = homeworkModel.Status
		homework.Addr        = homeworkModel.Addr
		homework.UserID      = homeworkModel.UserID
		homework.CourseID    = homeworkModel.CourseID
		homework.RoomID      = homeworkModel.HomeWoekID
		homework.Description = homeworkModel.Description
		homework.Comment     = homeworkModel.Comment
		c.JSON(http.StatusOK, gin.H{"homework":homework})
		return
	}

	//获取该用户在该班级下的所有作业，不唯一
	if (rid != ""){
		userid,_ := strconv.Atoi(uid)
		roomid,_ := strconv.Atoi(rid)
		homeworkModelList := GetHomeworkModelListByUidRid(userid, roomid)
		for _, homeworkModel := range homeworkModelList {
			homework.HomeWoekID  = homeworkModel.HomeWoekID
			homework.Status      = homeworkModel.Status
			homework.Addr        = homeworkModel.Addr
			homework.UserID      = homeworkModel.UserID
			homework.CourseID    = homeworkModel.CourseID
			homework.RoomID      = homeworkModel.HomeWoekID
			homework.Description = homeworkModel.Description
			homework.Comment     = homeworkModel.Comment
			homeworkList = append(homeworkList, homework)
		}
		c.JSON(http.StatusOK, gin.H{"homeworks":homeworkList})
		return
	}

	//获取该用户下的全部作业，不唯一

	userid,_ := strconv.Atoi(uid)
	homeworkModelList := GetHomeworkModelListByUid(userid)

	for _, homeworkModel := range homeworkModelList {
		homework.HomeWoekID  = homeworkModel.HomeWoekID
		homework.Status      = homeworkModel.Status
		homework.Addr        = homeworkModel.Addr
		homework.UserID      = homeworkModel.UserID
		homework.CourseID    = homeworkModel.CourseID
		homework.RoomID      = homeworkModel.HomeWoekID
		homework.Description = homeworkModel.Description
		homework.Comment     = homeworkModel.Comment
		homeworkList = append(homeworkList, homework)
	}

	c.JSON(http.StatusOK, gin.H{"homeworks":homeworkList})
	return

}
