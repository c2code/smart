package courses

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"smart.com/weixin/smart/logp"
	"os"
	"io"
	"bufio"
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
	file_path := "/tmp/";
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

	file_path = file_path + cid + "/";

	err=os.MkdirAll(file_path ,0755)
	if err!=nil{
		logger.Errorf("create file path %s if fail : %+v", cid, err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	for i, _ := range files {
		fileName := files[i].Filename
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
				c.JSON(http.StatusInternalServerError, gin.H{"courses": ""})
				return
			}

		}

		bufferedWriter.Flush()

	}
	c.JSON(http.StatusOK, gin.H{"courses": ""})
}