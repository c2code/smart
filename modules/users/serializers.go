package users

import (
	"github.com/gin-gonic/gin"

	"smart.com/weixin/smart/utils"
)

type ProfileSerializer struct {
	C *gin.Context
	UserModel
}

// Declare your response schema here
type ProfileResponse struct {
	ID        uint    `json:"-"`
	Username  string  `json:"username"`
	Bio       string  `json:"bio"`
	Image     *string `json:"image"`
	Following bool    `json:"following"`
}

// Put your response logic including wrap the userModel here.
func (self *ProfileSerializer) Response() ProfileResponse {
	myUserModel := self.C.MustGet("my_user_model").(UserModel)
	profile := ProfileResponse{
		ID:        self.ID,
		Username:  self.Username,
		Bio:       self.Bio,
		Image:     self.Image,
		Following: myUserModel.isFollowing(self.UserModel),
	}
	return profile
}

type UserSerializer struct {
	c *gin.Context
}

type UserResponse struct {
	ID       uint    `json:"uid"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Bio      string  `json:"bio"`
	Image    *string `json:"image"`
	Token    string  `json:"token"`
	Role     string  `json:"role"`
	Rights   int     `json:"rights"`
	Schedule int     `json:"schedule"`
	Phone    string  `json:"phone"`
	Nickname string  `json:"nickname"`
}

func (self *UserSerializer) Response() UserResponse {
	myUserModel := self.c.MustGet("my_user_model").(UserModel)
	user := UserResponse{
		ID:       myUserModel.ID,
		Username: myUserModel.Username,
		Email:    myUserModel.Email,
		Bio:      myUserModel.Bio,
		Image:    myUserModel.Image,
		Token:    utils.GenToken(myUserModel.ID),
		Role:     myUserModel.Role,
		Rights:   myUserModel.Rights,
		Schedule: myUserModel.Schedule,
		Phone:    myUserModel.Phone,
		Nickname: myUserModel.Nickname,
	}
	return user
}

