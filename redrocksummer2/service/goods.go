package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"

	"redrocksummer2/model"
)

type Goods struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Num   int    `json:"num"`
}

//TODO 添加商品（自己写的）
func AddGoods(c *gin.Context) {
	//接受Postman数据
	name  := c.PostForm("goods_name")
	price := c.PostForm("goods_price")
	num   := c.PostForm("goods_num")

	//string->int
	price_int, _ := strconv.Atoi(price)
	num_int, _   := strconv.Atoi(num)

	//创建表
	model.DB.Create(&Goods{
		Name:  name,
		Price: price_int,
		Num:   num_int,
	})
}

//找寻货物
func SelectGoods() (goods []Goods) {
	_goods, err := model.SelectGoods()//这个返回的goods是啥？ 空值？？？？（又蠢了） 看MySQL表吧。。。
	if err != nil {
		log.Printf("Error get goods info. Error: %s", err)
	}
	//遍历结构体数组
	for _, v := range _goods {
		good := Goods{
			ID:    v.ID,
			Name:  v.Name,
			Price: v.Price,
			Num:   v.Num,
		}
		goods = append(goods, good)//加入切片
	}
	return goods//返回切片
}
