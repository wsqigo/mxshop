package model

// 实际开发过程中 尽量设置为不为null
// https://zhuanlan.zhihu.com/p/73997266
// 这些类型我们使用int32还是int,减小到proto文件的转换

type Category struct {
	BaseModel
	Name  string `gorm:"type:varchar(20);not null"`
	Level int32  `gorm:"type:int;not null;default:1"`
	IsTab bool   `gorm:"default:false;not null"`

	// 外键
	ParentCategoryID int32
	ParentCategory   *Category
}

type Brand struct {
	BaseModel

	Name string `gorm:"type:varchar(20);not null"`
	Logo string `gorm:"type:varchar(200);default:'';not null"`
}

type GoodsCategoryBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Category   Category

	BrandID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Brand   Brand
}

// TableName 可以自定义表明

type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null"`
	Url   string `gorm:"type:varchar(200);not null"`
	Index int32  `gorm:"type:int;default:1;not null"`
}

type Goods struct {
	BaseModel

	CategoryID int32 `gorm:"type:int;index:not null"`
	Category   Category
	BrandID    int32
	Brand      Brand

	OnSale          bool     `gorm:"default:false;not null"`
	ShipFree        bool     `gorm:"default:false;not null"`
	IsNew           bool     `gorm:"default:false;not null"`
	IsHot           bool     `gorm:"default:false;not null"`
	Name            string   `gorm:"type:varchar(50);not null"`
	GoodsSn         string   `gorm:"type:varchar(50);not null"`
	ClickNum        int32    `gorm:"type:int;default:0;not null"`
	SoldNum         int32    `gorm:"type:int;default:0;not null"`
	FavNum          int32    `grom:"type:int;default:0;not null"`
	MarketPrice     float32  `gorm:"not null"`
	ShopPrice       float32  `gorm:"not null"`
	GoodsBrief      string   `gorm:"type:varchar(100);not null"`
	Images          GormList `gorm:"type:varchar(1000);not null"` // 数据库没有数组类型
	DescImages      GormList `gorm:"type:varchar(1000);not null"`
	GoodsFrontImage string   `gorm:"type:varchar(200);not null"`
}