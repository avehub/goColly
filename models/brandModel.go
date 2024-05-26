package models

type Brand struct {
	ID         uint `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	Pid        uint
	Type       uint
	Show       uint
	Letter     uint
	Name       string
	Logo       string
	CreateTime int `gorm:"autoCreateTime"`
	UpdateTime int `gorm:"autoUpdateTime"`
	DeleteTime int `gorm:"autoUpdateTime"`
}

func AddBrandOne(data *Brand) {
	result := DB.Create(data)
	if result.Error != nil {
		panic(result.Error)
	}
}

func AddBrandBatch(data *[]Brand) {
	result := DB.Create(&data)
	if result.Error != nil {
		panic(result.Error)
	}
}
