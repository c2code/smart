package homework

import (
	"github.com/gin-gonic/gin"
	"smart.com/weixin/smart/logp"
	"net/http"
	"strconv"
	"encoding/json"
	"smart.com/weixin/smart/modules/classroom"
	"net/url"
	"path"
	"os"
	"io"
	"bufio"
	"smart.com/weixin/smart/utils"
	"smart.com/weixin/smart/modules/student"
	"smart.com/weixin/smart/modules/courses"
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
		homework.RoomID      = homeworkModel.RoomID
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
		homework.RoomID      = homeworkModel.RoomID
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
		logger.Infof("get the homework, the user id = %d, roomid=%d",userid ,roomid)
		for _, homeworkModel := range homeworkModelList {
			homework.HomeWoekID  = homeworkModel.HomeWoekID
			homework.Status      = homeworkModel.Status
			homework.Addr        = homeworkModel.Addr
			homework.UserID      = homeworkModel.UserID
			homework.CourseID    = homeworkModel.CourseID
			homework.RoomID      = homeworkModel.RoomID
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
		if(homeworkModel.CourseID == 0) {
			continue
		}
		homework.HomeWoekID  = homeworkModel.HomeWoekID
		homework.Status      = homeworkModel.Status
		homework.Addr        = homeworkModel.Addr
		homework.UserID      = homeworkModel.UserID
		homework.CourseID    = homeworkModel.CourseID
		homework.RoomID      = homeworkModel.RoomID
		homework.Description = homeworkModel.Description
		homework.Comment     = homeworkModel.Comment
		homeworkList = append(homeworkList, homework)
	}

	c.JSON(http.StatusOK, gin.H{"homeworks":homeworkList})
	return

}

func AddHomework(c *gin.Context) {
	mlogger  := logp.NewLogger("homework")
	logger := mlogger.Named("add")

	body, err := c.GetRawData()
	if err != nil {
		logger.Errorf("get raw data fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	var inReq HomeWorkRe
	err = json.Unmarshal(body, &inReq)
	if err != nil {
		logger.Errorf("unmarshal fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	logger.Infof("the homework id is %d, user id is %d, course id is %d",inReq.HomeWoekID,inReq.UserID,inReq.CourseID)

    var homeworkModel HomeWorkModel
    var lasthowwork HomeWorkModel
	var studentlist []student.StudentModel
	var classroom       classroom.ClassroomModel
	var course          courses.CourseModel
	var current_course  courses.CourseModel

	db := utils.GetDB()

	db.Where("userid=? AND cid=?",inReq.UserID, inReq.CourseID).Find(&homeworkModel)

	if(homeworkModel.HomeWoekID != 0){
		logger.Infof("the user id is %d, homework has been added!",inReq.UserID)
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	logger.Infof("the user id is %d, homework will be added!",inReq.UserID)

	if (inReq.CourseID != 0) {

		db.Where("userid=?", inReq.UserID).Find(&studentlist)

		if (len(studentlist) == 0) {
			logger.Errorf("the user %d doese not join any classroom", inReq.UserID)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		db.Where("courseid=?", inReq.CourseID).Find(&current_course)

		i := 0;
		for _, student := range studentlist {
			db.Where("roomid=?", student.RoomID).Find(&classroom)
			db.Where("courseid=?", classroom.CourseID).Find(&course)
			if (course.CourseLevel == current_course.CourseLevel) {
				i = i + 1
				break
			}
		}

		if (i == 0) {
			logger.Errorf("the user %d doese not join the classroom", inReq.UserID)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
	}

	waitgroup.Wait()

	waitgroup.Add(1)
	defer waitgroup.Done()

	db.Last(&lasthowwork)

	homeworkModel.HomeWoekID  = lasthowwork.HomeWoekID + 1;
	homeworkModel.Status      = "未提交";
	homeworkModel.Addr        = inReq.Addr;
	homeworkModel.UserID      = inReq.UserID;
	homeworkModel.RoomID      = classroom.RoomID;
	homeworkModel.CourseID    = inReq.CourseID;
	homeworkModel.Description = inReq.Description;
	homeworkModel.Comment     = ""

	err = db.Save(&homeworkModel).Error
	if err != nil {
		logger.Errorf("save data fail %+v", err)
	}

	c.JSON(http.StatusOK, gin.H{})
}

func ModifyHomework(c *gin.Context) {
	mlogger  := logp.NewLogger("homework")
	logger := mlogger.Named("modify")

	body, err := c.GetRawData()
	if err != nil {
		logger.Errorf("get raw data fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	var inReq HomeWorkRe
	err = json.Unmarshal(body, &inReq)
	if err != nil {
		logger.Errorf("unmarshal fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	logger.Infof("the user id is %d, course id is %d will be modified ",inReq.UserID,inReq.CourseID)

	db := utils.GetDB()

	var homeworkModel HomeWorkModel
	db.Where("userid = ? AND cid=?", inReq.UserID,inReq.CourseID).Find(&homeworkModel)
	db.Model(&homeworkModel).Updates(map[string]interface{}{"status":inReq.Status, "comment":inReq.Comment})

	c.JSON(http.StatusOK, gin.H{})
}

func UploadCourseFile(c *gin.Context) {

	mlogger  := logp.NewLogger("homework")
	logger := mlogger.Named("upload")
	file_path := "/tmp/homework/";
	//logger.Info("Begin to upload files!")

	err := c.Request.ParseMultipartForm(32 << 10)  //32M

	if err != nil {
		logger.Errorf("pares formdata error : %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	form := c.Request.MultipartForm

	cid := form.Value["cid"][0]
	uid := form.Value["uid"][0]
	files := form.File["upload"]
	fileName := cid + ".sb3"

	userid,_ := strconv.Atoi(uid)
	courseid,_ := strconv.Atoi(cid)

	var homeworkModel HomeWorkModel
	db := utils.GetDB()
	db.Where("userid = ? AND cid=?", userid,courseid).Find(&homeworkModel)

	file_path = file_path + strconv.Itoa((int)(homeworkModel.CourseID)) + "/" + strconv.Itoa((int)(homeworkModel.RoomID)) + "/" + uid + "/";

	err=os.MkdirAll(file_path ,0755)
	if err!=nil{
		logger.Errorf("create file path %s if fail : %+v", cid, err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	for i, _ := range files {
		//fileName = files[i].Filename
		logger.Infof("%s of cid=%s will create", fileName, cid)
		file, err := files[i].Open()

		defer file.Close()

		if err != nil {
			logger.Errorf("Open source file : %+v", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		// Open file for writing
		newfile, err := os.OpenFile(file_path + fileName,
			os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
			0755,
		)

		if err != nil {
			logger.Errorf("Create destination file fail on disk: %+v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"courses": ""})
			return
		}

		defer newfile.Close();

		// Create a buffered writer from the file
		bufferedWriter := bufio.NewWriter(newfile)

		buf := make([]byte, 35*1024)

		for i := 0; i < 2*1024*1024*1024/35; i++ {//break in 2*1024*1024/(35*35)
			r, e := file.Read(buf)
			if r == 0 {
				if e != nil || e == io.EOF {
					break
				}
			}

			_, err = bufferedWriter.Write(buf)

			if err != nil {
				logger.Errorf("Write destination file fail on disk: %+v", err)
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}

		}

		bufferedWriter.Flush()
	}

	c.JSON(http.StatusOK, gin.H{})
}

func ReadCourseFile (c *gin.Context)  {
	mlogger  := logp.NewLogger("courses")
	logger := mlogger.Named("readfile")
	file_path := "/tmp/homework/";
	//logger.Info("Begin to read files!")

	cid := c.Query("cid")
	uid := c.Query("uid")
	fileName := cid + ".sb3"

	userid,_ := strconv.Atoi(uid)
	courseid,_ := strconv.Atoi(cid)
	var homeworkModel HomeWorkModel
	db := utils.GetDB()
	db.Where("userid = ? AND cid=?", userid,courseid).Find(&homeworkModel)

	file_path = file_path + strconv.Itoa((int)(homeworkModel.CourseID)) + "/" + strconv.Itoa((int)(homeworkModel.RoomID)) + "/" + uid + "/" + fileName;

	c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(path.Base(file_path)))

	file, err := os.OpenFile(file_path, os.O_RDONLY, 0666)
	defer file.Close()

	if err != nil {
		logger.Errorf("open file fail on disk: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.Header("Content-Type", "application/octet-stream")

	c.Stream(func(w io.Writer) bool {
		io.Copy(w, file)
		return false
	})

	return
}