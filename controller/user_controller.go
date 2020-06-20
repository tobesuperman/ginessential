package controller

import (
	"chao.com/ginessential/common"
	"chao.com/ginessential/dto"
	"chao.com/ginessential/model"
	"chao.com/ginessential/response"
	"chao.com/ginessential/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func Register(context *gin.Context) {
	db := common.GetDB()
	// 使用map获取请求的参数
	//var requestMap = make(map[string]string)
	//json.NewDecoder(context.Request.Body).Decode(&requestMap)
	var requestUser = model.User{}
	// 使用结构体获取请求的参数
	//json.NewDecoder(context.Request.Body).Decode(&requestUser)
	context.Bind(&requestUser)
	// 获取参数
	name := requestUser.Name
	telephone := requestUser.Telephone
	password := requestUser.Password
	// 数据验证
	if len(telephone) != 11 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位！")
		return
	}
	if len(password) < 6 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "密码不少于6位！")
		return
	}
	if len(name) == 0 {
		name = util.RandomString(10)
	}
	// 判断手机号是否存在
	if IsTelephoneExist(db, telephone) {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "用户已存在！")
		return
	}
	// 创建用户
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(context, http.StatusInternalServerError, 500, nil, "加密错误！")
		return
	}
	user := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hashedPassword),
	}
	db.Create(&user)
	// 发放token
	token, err := common.ReleaseToken(&user)
	if err != nil {
		response.Response(context, http.StatusInternalServerError, 500, nil, "系统异常！")
		log.Printf("token generate error: %v", err)
		return
	}
	// 返回结果
	response.Success(context, gin.H{"token": token}, "注册成功！")
}

func Login(context *gin.Context) {
	var requestUser = model.User{}
	context.Bind(&requestUser)
	telephone := requestUser.Telephone
	password := requestUser.Password
	db := common.GetDB()
	// 数据验证
	if len(telephone) != 11 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位！")
		return
	}
	if len(password) < 6 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "密码不少于6位！")
		return
	}
	// 判断手机号是否存在
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "用户不存在！")
		return
	}
	// 判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Fail(context, nil, "密码错误！")
		return
	}
	// 发放token
	token, err := common.ReleaseToken(&user)
	if err != nil {
		response.Response(context, http.StatusInternalServerError, 500, nil, "系统异常！")
		log.Printf("token generate error: %v", err)
		return
	}
	// 返回结果
	response.Success(context, gin.H{"token": token}, "登录成功！")
}

func IsTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	db.AutoMigrate(&model.User{})
	return false
}

func Info(context *gin.Context) {
	user, _ := context.Get("user")
	response.Success(context, gin.H{"user": dto.ToUserDto(user.(model.User))}, "")
}
