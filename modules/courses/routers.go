package courses

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"smart.com/weixin/smart/logp"
	"os"
	"io"
	"bufio"
	"smart.com/weixin/smart/utils"
	"net/url"
	"path"
	"encoding/json"
	"strings"
	"strconv"
)

func GetCourseList(c *gin.Context) {

	var courselist []CoursesModelValidator
	coursemodel := GetCourseModel()

	for _, tmp := range coursemodel {
		course := CoursesModelValidator {CourseID:tmp.CourseID, PID:tmp.PID, Name:tmp.Name, Desc:tmp.Desc, Vedio:tmp.Vedio, CourseLevel:tmp.CourseLevel}

		courselist = append(courselist, course)
	}

	//jsonRsp, _ := json.Marshal(courselist)
	//c.Data(http.StatusOK, "application/json", jsonRsp)
	c.JSON(http.StatusOK, gin.H{"courses": courselist})
}

func UploadCourseFile(c *gin.Context) {

	mlogger  := logp.NewLogger("courses")
	logger := mlogger.Named("upload")
	file_path := "/tmp/course/";
	//logger.Info("Begin to upload files!")

	err := c.Request.ParseMultipartForm(32 << 10)  //32M

	if err != nil {
		logger.Errorf("pares formdata error : %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	form := c.Request.MultipartForm

	cid := form.Value["cid"][0]
	files := form.File["upload"]
	fileName := ""

	file_path = file_path + cid + "/";

	err=os.MkdirAll(file_path ,0755)
	if err!=nil{
		logger.Errorf("create file path %s if fail : %+v", cid, err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	for i, _ := range files {
		fileName = files[i].Filename
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

	var courseModel CourseModel

	db := utils.GetDB()
	courseModel = CourseModel{}
	db.First(&courseModel, "courseid = ?", cid)
	db.Model(&courseModel).Update("vedio", fileName)

	c.JSON(http.StatusOK, gin.H{})
}

func ReadCourseFile (c *gin.Context)  {
	mlogger  := logp.NewLogger("courses")
	logger := mlogger.Named("readfile")
	file_path := "/tmp/course/";
	//logger.Info("Begin to read files!")

	cid := c.Query("cid")
	fileName := c.Query("file_name")

	file_path = file_path + cid + "/" + fileName

	c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(path.Base(file_path)))

	file, err := os.OpenFile(file_path, os.O_RDONLY, 0666)
	defer file.Close()

	if err != nil {
		logger.Errorf("open file fail on disk: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	fi, err := file.Stat()

	if err != nil {
		logger.Errorf("open file fail on disk: %+v", err)
		c.Data(http.StatusOK, "application/octet-stream", nil)
		return
	}

	fileSize := int(fi.Size())
	logger.Infof("mp4 file size is %d", fileSize)

	if len(c.GetHeader("Range")) == 0 {
		logger.Info("mp4 file will play on begin")
		contentLength := strconv.Itoa(fileSize)
		contentEnd := strconv.Itoa(fileSize - 1)
		c.Header("Content-Type", "video/mp4")
		c.Header("Accept-Ranges", "bytes")
		c.Header("Content-Length", contentLength)
		c.Header("Content-Range", "bytes 0-"+contentEnd+"/"+contentLength)
		c.Writer.WriteHeader(206)

		c.Stream(func(w io.Writer) bool {
			io.Copy(w, file)
			return false
		})

		return

	} else {
		rangeParam := strings.Split(c.GetHeader("Range"), "=")[1]
		splitParams := strings.Split(rangeParam, "-")

		logger.Info("mp4 file will play on middle")

		contentStartValue := 0
		contentStart := strconv.Itoa(contentStartValue)
		contentEndValue := fileSize - 1
		contentEnd := strconv.Itoa(contentEndValue)
		contentSize := strconv.Itoa(fileSize)

		if len(splitParams) > 0 {
			contentStartValue, err = strconv.Atoi(splitParams[0])
			if err != nil {
				contentStartValue = 0
			}
			contentStart = strconv.Itoa(contentStartValue)
		}

		if len(splitParams) > 1 {
			contentEndValue, err = strconv.Atoi(splitParams[1])

			if err != nil {
				contentEndValue = fileSize - 1
			}

			contentEnd = strconv.Itoa(contentEndValue)
		}

		contentLength := strconv.Itoa(contentEndValue - contentStartValue + 1)

		c.Header("Content-Type", "video/mp4")
		c.Header("Accept-Ranges", "bytes")
		c.Header("Content-Length", contentLength)
		c.Header("Content-Range", "bytes "+contentStart+"-"+contentEnd+"/"+contentSize)
		c.Writer.WriteHeader(206)

		file.Seek(int64(contentStartValue), 0)

		c.Stream(func(w io.Writer) bool {
			io.Copy(w, file)
			return false
		})
	}

	//c.Header("Content-Type", "video/mp4")

	//c.Stream(func(w io.Writer) bool {
	//	io.Copy(w, file)
	//	return false
	//})

	//return
}

func ModifyCourse (c *gin.Context)  {
	mlogger  := logp.NewLogger("courses")
	logger := mlogger.Named("modify")

	body, err := c.GetRawData()
	if err != nil {
		logger.Errorf("get raw data fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	var inReq CoursesModelValidator
	err = json.Unmarshal(body, &inReq)
	if err != nil {
		logger.Errorf("unmarshal fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	logger.Infof("the cid is %d, the name is %s , the desc is %s, the course level is %s",inReq.CourseID, inReq.Name, inReq.Desc, inReq.CourseLevel)

	var courseModel CourseModel
	db := utils.GetDB()
	courseModel = CourseModel{}
	db.First(&courseModel, "courseid = ?", inReq.CourseID)
	db.Model(&courseModel).Updates(map[string]interface{}{"name":inReq.Name, "description":inReq.Desc, "clevel":inReq.CourseLevel})
	c.JSON(http.StatusOK, gin.H{})
}

func AddCourse (c *gin.Context)  {
	mlogger  := logp.NewLogger("courses")
	logger := mlogger.Named("add")

	body, err := c.GetRawData()
	if err != nil {
		logger.Errorf("get raw data fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	var inReq CoursesModelValidator
	err = json.Unmarshal(body, &inReq)
	if err != nil {
		logger.Errorf("unmarshal fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	logger.Infof("the cid is %d, the PID is %d, the name is %s , the desc is %s, the course level is %s",inReq.CourseID, inReq.PID, inReq.Name, inReq.Desc, inReq.CourseLevel)

	var courseModel CourseModel
	courseModel.CourseID = inReq.CourseID
	courseModel.PID      = inReq.PID
	courseModel.Name     = inReq.Name
	courseModel.Desc     = inReq.Desc
	courseModel.Vedio    = inReq.Vedio
	courseModel.CourseLevel = inReq.CourseLevel
	db := utils.GetDB()
	err = db.Save(&courseModel).Error
	if err != nil {
		logger.Errorf("save data fail %+v", err)
	}

	c.JSON(http.StatusOK, gin.H{})
}

func DeleteCourse (c *gin.Context)  {
	mlogger  := logp.NewLogger("courses")
	logger := mlogger.Named("delete")

	body, err := c.GetRawData()
	if err != nil {
		logger.Errorf("get raw data fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	var inReq CoursesModelValidator
	err = json.Unmarshal(body, &inReq)
	if err != nil {
		logger.Errorf("unmarshal fail %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	logger.Infof("the cid is %d, the name is %s , the desc is %s, the course level is %s will be delete",inReq.CourseID, inReq.Name, inReq.Desc, inReq.CourseLevel)

	var courseModel CourseModel
	courseModel.CourseID = inReq.CourseID
	courseModel.PID      = inReq.PID
	courseModel.Name     = inReq.Name
	courseModel.Desc     = inReq.Desc
	courseModel.CourseLevel = inReq.CourseLevel
	db := utils.GetDB()
	db.Unscoped().Where("courseid = ?", inReq.CourseID).Delete(CourseModel{})

	c.JSON(http.StatusOK, gin.H{})
}