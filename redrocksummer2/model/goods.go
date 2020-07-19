package model

import "github.com/jinzhu/gorm"

//商品信息
type Goods struct {
	gorm.Model
	Name  string
	Price int
	Num   int
}


// 添加商品
func (goods *Goods)AddGoods() error{
	return DB.Create(goods).Error
}

// 查看商品
func SelectGoodsById(id uint) (goods Goods, err error){
	err = DB.Table("goods").Where("id = ?",id).First(&goods).Error
	if err != nil {
		return Goods{}, err
	}
	return goods, nil
}

// 查看所有的商品
func SelectGoods() (goods []Goods, err error){
	err = DB.Table("goods").Find(&goods).Error//搜索表中的货物加入结构体数组
	if err != nil {
		return nil, err
	}
	return goods, nil//返回货物结构体数组和错误空值
	// 疑问？？？？这个结构体数组是空的吧？？？
	// 这里搞错了。。。MySQL表不空就不空。。。。
}



