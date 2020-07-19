package service

import (
	"fmt"
	_ "fmt"
	"log"
	"redrocksummer2/model"
	"sync"
	"time"
)

type User struct {
	UserId string
	GoodsId  uint
}

var OrderChan = make(chan User, 1024)

var ItemMap = make(map[uint]*Item)

type Item struct {
	ID        uint   // 商品id
	Name      string // 名字
	Total     int    // 商品总量
	Left      int    // 商品剩余数量
	IsSoldOut bool   // 是否售罄
	leftCh    chan int       //应该是剩余管道先这么命名吧
	sellCh    chan int       //出售管道
	done      chan struct{}  //结构体管道
	Lock      sync.Mutex     //锁
}

//TODO 写一个定时任务，每天定时从数据库加载数据到Map！！！！！！！！！！！！！！！！

////加物品
//func AddShelve()  {
//	beginTime := time.Now()
//	// 获取第二天时间
//	nextTime := beginTime.Add(time.Hour * 24)
//	// 计算次日零点，即商品上架的时间
//	offShelveTime := time.Date(nextTime.Year(), nextTime.Month(), nextTime.Day(), 0, 0, 0, 0, nextTime.Location())
//	fmt.Println(offShelveTime)
//	timer := time.NewTimer(offShelveTime.Sub(beginTime))
//	<-timer.C
//
//	var good  []Goods
//	model.DB.Find(&good)
//	fmt.Println(good)
//	var S []Item
//	for i := 0; i < len(good); i++ {
//		s :=  Item{
//			ID:        good[i].ID,
//			Name:      good[i].Name,
//			Total:     good[i].Num,
//			Left:      good[i].Num,
//			IsSoldOut: false,
//		}
//
//		S := append(S, s)
//		ItemMap = S
//	}
//}

//使用的逻辑和下架商品一样
func InitMap() {
	beginTime := time.Now()//当前时间
	// 获取第二天时间
	nextTime := beginTime.Add(time.Hour * 24)
	// 计算次日零点，即商品上架的时间
	offShelveTime := time.Date(nextTime.Year(), nextTime.Month(), nextTime.Day(), 0, 0, 0, 0, nextTime.Location())
	fmt.Println(offShelveTime)
	timer := time.NewTimer(offShelveTime.Sub(beginTime))
	//	<-timer.C
	for {
		<-timer.C
		for _, i := range ItemMap {
			some := model.Goods{
				Name: i.Name,
				Num:  i.Left,
			}
			if err := some.AddGoods(); err != nil{
				log.Println(err)
				return
			}
		}
		some := SelectGoods()//找到货物，返回货物结构体
		for _, i2 := range some {
			item := &Item{
				ID:        i2.ID,
				Name:      i2.Name,
				Total:     i2.Num,
				Left:      i2.Num,
				IsSoldOut: false,
				leftCh:    make(chan int),
				sellCh:    make(chan int),
			}
			ItemMap[item.ID] = item
		}
		//timer.Reset(time.Hour * 24)
	}
}

func initMap() {
	item := &Item{
		ID:        1,
		Name:      "测试",
		Total:     100,
		Left:      100,
		IsSoldOut: false,
		leftCh:    make(chan int),  //管道
		sellCh:    make(chan int),  //管道
	}
	ItemMap[item.ID] = item  //TODO map商品ID等于这个结构体位于的位置
}

func getItem(itemId uint) *Item{
	return ItemMap[itemId]
}

//订购？应该是
func order() {
	//循环
	for {
		user := <- OrderChan //从订购管道中接受数据
		item := getItem(user.GoodsId)//获取商品信息
		item.SecKilling(user.UserId)//检测是否卖完，卖完报错，没卖完就买
	}
}

func (item *Item) SecKilling(userId string) {

	item.Lock.Lock()//锁
	defer item.Lock.Unlock()//解锁
	// 等价
	// var lock = make(chan struct{}, 1}
	// lock <- struct{}{}
	// defer func() {
	// 		<- lock
	// }
	//检查是否卖完
	if item.IsSoldOut {
		return
	}
	item.BuyGoods(1)//没卖完就买呗
	MakeOrder(userId, item.ID,1)
}


// 定时下架
func (item *Item) OffShelve() {
	beginTime := time.Now()//获取当前时间

	// 获取第二天时间
	nextTime := beginTime.Add(time.Hour * 24)
	// 计算次日零点，即商品下架的时间
	offShelveTime := time.Date(nextTime.Year(), nextTime.Month(), nextTime.Day(), 0, 0, 0, 0, nextTime.Location())

	timer := time.NewTimer(offShelveTime.Sub(beginTime))

	<-timer.C//TODO 这个有何用意？ 应该是一个时间管道

	delete(ItemMap, item.ID)//删除ID
	close(item.done)//防止时间管道死锁

}

// 出售商品
func (item *Item) SalesGoods() {
	for {
		//选择
		select {
		//num由哪个管道发进来就运行哪个
		case num := <-item.sellCh:
			if item.Left -= num; item.Left <= 0 {
				item.IsSoldOut = true
			}

		case item.leftCh <- item.Left://这个方法还没写。。。我也不知道学长想干啥

		case <- item.Done():
			log.Println("我自闭了")
			return
		}
	}
}

//TODO 这个不会死锁吗？可能会，也可能不会。。。（学长写的一定不会）
func (item *Item) Done() <-chan struct{} {
	//done就是结构体管道
	if item.done == nil {
		item.done = make(chan struct{})//如果done为空就创造不为空
	}
	d := item.done
	return d
}

//监视器？？？为啥要取这名字，可能是由于这函数是用来查询商品是否为0
func (item *Item) Monitor() {
	go item.SalesGoods()//开启协程
}

// 获取剩余库存
func (item *Item) GetLeft() int {
	var left int
	left = <-item.leftCh
	return left
}

// 购买商品
func (item *Item) BuyGoods(num int) {
	item.sellCh <- num//把购买数量传到sellCh管道
}

func InitService() {
	initMap()//封装了一个item结构体map
	//遍历item切片
	for _,item := range ItemMap{
		item.Monitor()//监管商品库存是否为0
		go item.OffShelve()//下架产品
		go InitMap()//TODO 每日上新（属于自己写的）
	}
	//TODO 为啥要循环10次？？？？这点没搞明白
	for i := 0; i < 10; i++ {
		go order()
	}
}
