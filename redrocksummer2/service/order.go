package service

import (
	"log"

	"redrocksummer2/model"
)

//order生成订单
func MakeOrder(userId string, goodsId uint, num int) {

	order := model.Order{
		UserID:  userId,
		GoodsID: goodsId,
		Num:     num,
	}
	err := order.MakeOrder()//创建了订单表。。。然后？
	// 订单表好像没啥用的哈。。。
	// 应该就是个证明，证明你买过
	if err != nil {
		log.Printf("Error make an order. Error: %s",err)
	}
	log.Println("success")
}