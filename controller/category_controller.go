package controller

import (
	"chao.com/ginessential/model"
	"chao.com/ginessential/repository"
	"chao.com/ginessential/response"
	"chao.com/ginessential/vo"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ICategoryController interface {
	RestController
}

type CategoryController struct {
	repository repository.CategoryRepository
}

func NewCategoryController() ICategoryController {
	repo := repository.NewCategoryRepository()
	repo.DB.AutoMigrate(model.Category{})
	return CategoryController{
		repository: repo,
	}
}

func (c CategoryController) Create(ctx *gin.Context) {
	var requestCategory vo.CreateCategoryRequest
	if err := ctx.ShouldBind(&requestCategory); err != nil {
		response.Fail(ctx, nil, "数据验证错误！")
		return
	}
	category, err := c.repository.Create(requestCategory.Name)
	if err != nil {
		//response.Fail(ctx, nil, "创建分类失败！")
		panic(err)
		return
	}
	response.Success(ctx, gin.H{"category": category}, "创建分类成功！")
}

func (c CategoryController) Update(ctx *gin.Context) {
	// 绑定body中的参数
	// 新的分类
	var requestCategory vo.CreateCategoryRequest
	if err := ctx.ShouldBind(&requestCategory); err != nil {
		response.Fail(ctx, nil, "数据验证错误！")
		return
	}
	// 获取path中的参数
	// 需要更新的分类
	categoryId, _ := strconv.Atoi(ctx.Params.ByName("id"))
	updateCategory, err := c.repository.SelectById(categoryId)
	if err != nil {
		response.Fail(ctx, nil, "分类不存在！")
		return
	}
	// 更新分类
	category, err := c.repository.Update(*updateCategory, requestCategory.Name)
	if err != nil {
		response.Fail(ctx, nil, "更新分类失败！")
		return
	}
	response.Success(ctx, gin.H{"category": category}, "更新分类成功！")
}

func (c CategoryController) Show(ctx *gin.Context) {
	categoryId, _ := strconv.Atoi(ctx.Params.ByName("id"))
	category, err := c.repository.SelectById(categoryId)
	if err != nil {
		response.Fail(ctx, nil, "分类不存在！")
		return
	}
	response.Success(ctx, gin.H{"category": category}, "")
}

func (c CategoryController) Delete(ctx *gin.Context) {
	categoryId, _ := strconv.Atoi(ctx.Params.ByName("id"))
	if err := c.repository.DeleteById(categoryId); err != nil {
		response.Fail(ctx, nil, "删除分类失败！")
		return
	}
	response.Success(ctx, nil, "删除分类成功！")
}
