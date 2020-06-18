package controller

import (
	"chao.com/ginessential/common"
	"chao.com/ginessential/model"
	"chao.com/ginessential/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

func Register(context *gin.Context) {
	db := common.GetDB()
	// 获取参数
	name := context.PostForm("name")
	telephone := context.PostForm("telephone")
	password := context.PostForm("password")
	// 数据验证
	if len(telephone) != 11 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "手机号必须为11位！",
		})
		return
	}
	if len(password) < 6 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "密码不能少于11位！",
		})
		return
	}
	if len(name) == 0 {
		name = util.RandomString(10)
	}
	log.Println(name, telephone, password)
	// 判断手机号是否存在
	if IsTelephoneExist(db, telephone) {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "用户已存在！",
		})
		return
	}
	// 创建用户
	user := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  password,
	}
	db.Create(&user)
	// 返回结果
	context.JSON(200, gin.H{
		"msg": "注册成功！",
	})
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
