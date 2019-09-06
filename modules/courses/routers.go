package courses

import (
	"github.com/gin-gonic/gin"
	"net/http"
	//"encoding/json"
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
