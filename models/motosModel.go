package models

type Motos struct {
	ID              uint `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	BrandId         float64
	Aleph           string
	BrandName       string
	BrandLogo       string
	Keywords        string
	Spelling        string
	BrandEnergyType float64
	GoodId          float64
	GoodName        string
	GoodPic         string
	OriginGoodPic   string
	SeriesId        float64
	SeriesName      string
	CreateTime      int `gorm:"autoCreateTime"`
	UpdateTime      int `gorm:"autoUpdateTime"`
	DeleteTime      int `gorm:"autoUpdateTime"`
}

func AddMotoBatch(data *[]Motos) {
	result := DB.Create(&data)
	if result.Error != nil {
		panic(result.Error)
	}
}

func AddMotoOne(data *Motos) {
	result := DB.Create(&data)
	if result.Error != nil {
		panic(result.Error)
	}
}
