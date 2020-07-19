package controller

import (
	"github.com/gin-gonic/gin"
	"redrocksummer2/service"
	"strconv"

)


func MakeOrder(ctx *gin.Context) {
	//接受Postman参数
	userId := ctx.PostForm("userId")
	goodsId := ctx.PostForm("goodsId")

	//string->int
	itemId,_ := strconv.Atoi(goodsId)

	//写入数据（User结构体），传入order购买的管道中
	service.OrderChan <- service.User{
		UserId:  userId,
		GoodsId: uint(itemId),
	}
	//应该是购买成功，或者是订单成功
	ctx.JSON(200, gin.H{
		"status": 200,
		"info": "success",
	})
}



