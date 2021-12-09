package database

import (
	"gorm.io/gorm"
	"os"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&Product{})
	if err != nil {
		return err
	}


	seedTmp:=os.Getenv("SEED_IF_EMPTY")

	seed:=false

	if seedTmp=="true" {
		seed=true
	}

	if !seed {
		return nil
	}
	var count int64
	count=0
	db.Model(&Product{}).Count(&count)

	if count!=0 {
		return nil
	}

	var products = []Product{
		{
			Name:       "Nvidia GTX 580",
			Img:        "gtx580.jpeg",
			Price:      85,
			Discount:   0,
			CategoryID: 1,
		},
		{
			Name:       "Nvidia GTX 1050 Ti",
			Img:        "gtx1050ti.jpg",
			Price:      110,
			Discount:   15,
			CategoryID: 1,
		},
		{
			Name:       "AMD Radeon RX 470",
			Img:        "rx470.jpg",
			Price:      110,
			Discount:   0,
			CategoryID: 1,
		},
		{
			Name:       "AMD Radeon RX 570",
			Img:        "rx570.jpg",
			Price:      130,
			Discount:   0,
			CategoryID: 1,
		},
		{
			Name:       "AMD FX-6100",
			Img:        "fx6100.jpeg",
			Price:      50,
			Discount:   0,
			CategoryID: 2,
		},
		{
			Name:       "AMD FX-4170",
			Img:        "fx4170.jpg",
			Price:      80,
			Discount:   5,
			CategoryID: 2,
		},
		{
			Name:       "MSI 760GM-P21",
			Img:        "msi760gm-p21.jpg",
			Price:      30,
			Discount:   0,
			CategoryID: 3,
		},
	}
	db.Create(&products)
	
	return nil
}
