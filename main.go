package main

import (
	"chao.com/ginessential/common"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := common.GetDB()
	defer db.Close()
	r := gin.Default()
	r = CollectRoute(r)
	panic(r.Run())
}
