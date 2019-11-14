package courses

import (
	"github.com/gin-gonic/gin"
	"net/http"
	//"encoding/json"
//	"io"
//	"bytes"
	"smart.com/weixin/smart/logp"
)

func GetCourseList(c *gin.Context) {

	var courselist []CoursesModelValidator
	coursemodel := GetCourseModel()

	for _, tmp := range coursemodel {
		course := CoursesModelValidator {CourseID:tmp.CourseID, PID:tmp.PID, Name:tmp.Name, Desc:tmp.Desc}

		courselist = append(courselist, course)
	}

	//jsonRsp, _ := json.Marshal(courselist)
	//c.Data(http.StatusOK, "application/json", jsonRsp)
	c.JSON(http.StatusOK, gin.H{"courses": courselist})
}

func UploadCourseFile(c *gin.Context) {

	mlogger  := logp.NewLogger("courses")
	logger := mlogger.Named("upload")
	logger.Info("Begin to upload files!")

	err := c.Request.ParseMultipartForm(32 << 10)  //32M

	if err != nil {
		logger.Errorf("pares formdata error : %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	form := c.Request.MultipartForm

	cid := form.Value["cid"][0]
	files := form.File["upload"]

	for i, _ := range files {
		fileName := files[i].Filename
		logger.Infof("cid %s of %s will create", cid, fileName)
		file, err := files[i].Open()

		defer file.Close()


		if err != nil {
			logger.Errorf("Open source file : %+v", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
/*
		writeType := new(wire.WriteType)
		//*writeType = wire.WriteTypeCacheThrough
		*writeType = wire.WriteTypeThrough

		id, err := m.fs.CreateFile(object+fileName, &option.CreateFile{ WriteType: writeType})
		defer m.fs.Close(id)

		if err != nil {
			logger.Errorf("Create destination file fail on disk: %+v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"courses": ""})
			return
		}

		buf := make([]byte, 35*1024)

		for i := 0; i < 2*1024*1024*1024/35; i++ {//break in 2*1024*1024/(35*35)
			r, e := file.Read(buf)
			if r == 0 {
				if e != nil || e == io.EOF {
					break
				}
			}

			_, err = m.fs.Write(id, bytes.NewReader(buf))

			if err != nil {
				logger.Errorf("Write destination file fail on disk: %+v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"courses": ""})
				return
			}

		}
*/
	}
	c.JSON(http.StatusOK, gin.H{"courses": ""})
}