package infra

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"smart.com/weixin/smart/thirdparty/github.com/gin-contrib/cors"
	"smart.com/weixin/smart/thirdparty/github.com/gin-contrib/static"
	"smart.com/weixin/smart/utils"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"smart.com/weixin/smart/modules/users"
	"smart.com/weixin/smart/modules/courses"
	"smart.com/weixin/smart/modules/classroom"
)
type PingResponse struct {
	BaseResponse
	Message string `json:"message"`
}

type Login struct {
	User     string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

//entry of web server
func (m Manager) webListen() {
	ginLogger := m.logger.Named("gin")

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	//router := gin.Default()
	router.Use(utils.Ginzap(ginLogger))
	router.Use(gin.Recovery())
	router.Use(cors.Default())
	router.Use(static.Serve("/", static.LocalFile("./dist", true)))
	router.Use(static.Serve("/api", static.LocalFile("./dist", true)))

	//provide a internal access rest api
	tuna_v1:= router.Group("/")
	{
		tuna_v1.POST("/example", m.exampleRestCall)
		tuna_v1.POST("/login", func(c *gin.Context) {
			var json Login
			fmt.Println("login begin")
			if err := c.ShouldBindWith(&json, binding.JSON); err == nil {
				if json.User == "test" && json.Password == "7c4a8d09ca3762af61e59520943dc26494f8941b" {
					fmt.Println("you are logged in")
					c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
				} else {
					fmt.Println("unauthorized")
					fmt.Printf("%+v", json)
					c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
				}
			} else {
				fmt.Printf("%+v", json)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		})
	}

	//provide a external access rest api
	tuna_v2 := router.Group("/api")
	{
		tuna_v2.GET("/log-level", m.onGetLogLevel) //get log level
		tuna_v2.POST("/log-level", m.onSetLogLevel) //set log level

		tuna_v3 := tuna_v2.Group("/users")
		tuna_v3.POST("/", users.UsersRegistration)// UsersRegistration)
		tuna_v3.POST("/login",users.UsersLogin)// UsersLogin)

		tuna_v2.Use(users.AuthMiddleware(true))

		tuna_v4 := tuna_v2.Group("/profiles")
		tuna_v4.GET("/:username", users.ProfileRetrieve)//ProfileRetrieve)
		tuna_v4.POST("/:username/follow", users.ProfileFollow)//ProfileFollow)
		tuna_v4.DELETE("/:username/follow", users.ProfileUnfollow)//ProfileUnfollow)

		tuna_v5 := router.Group("/api/ping")
		tuna_v5.GET("/", m.onPing) //to check the tuna service is accessful

		tuna_v6 := tuna_v2.Group("/user")
		tuna_v6.GET("/", users.UserRetrieve)
		tuna_v6.PUT("/", users.UserUpdate)

		//#######################################################//

		tuna_v7 := tuna_v2.Group("/course")
		tuna_v7.GET("/",courses.GetCourseList)
		tuna_v7.POST("/upload",courses.UploadCourseFile)
		tuna_v7.GET("/download", courses.ReadCourseFile)
		tuna_v7.POST("/modify",courses.ModifyCourse)
		tuna_v7.POST("/add",courses.AddCourse)
		tuna_v7.POST("/delete",courses.DeleteCourse)

		tuna_v8 := tuna_v2.Group("/classroom")
		tuna_v8.GET("/",classroom.GetLastClassroom)
		tuna_v8.POST("/",classroom.GetClassroomList)
		tuna_v8.POST("/add",classroom.AddClassroom)
		tuna_v8.POST("/modify",classroom.ModifyClassroom)
		tuna_v8.POST("/delete",classroom.DeleteClassroom)

		//tuna_v5 := tuna_v2.Group("/articles")
	}

	portSpec := fmt.Sprintf(":%d", m.config.WebPort)

	router.Run(portSpec)
}

func (m Manager) onPing(c *gin.Context) {
	/*c.JSON(200, PingResponse{
		BaseResponse: BaseResponse{
			ErrCode: ErrCodeOk,
			ErrInfo: ErrInfoOk,
		},
		Message: "pong"})*/
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

