package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"redrocksummer2/service"
)

func SelectGoods(ctx *gin.Context) {
	goods := service.SelectGoods()//接受返回的切片
	ctx.JSON(http.StatusOK, gin.H{
		"status": 200,
		"info": "success",
		//打印所含的货物切片
		//TODO 打印出来是啥样很好奇，struct{。。。}中间没懂打印出来是啥等会试试。。。！！！
		"data": struct {
			Goods []service.Goods `json:"goods"`
		}{goods},
	})
}



