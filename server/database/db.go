package database

import (
	"MTG/server/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"syreclabs.com/go/faker"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//connect database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	DB = db
	log.Println("Connected to database!")

	err = db.AutoMigrate(&models.Item{})
	if err != nil {
		log.Fatalf("Migrate to database failed: %v", err)
	}
	FakerItems(db, 5000)

}

func FakerItems(connDB *gorm.DB, countRow int) {
	var count int64
	connDB.Model(&models.Item{}).Count(&count)
	if count != 0 {
		fmt.Println("Table items row count:", count)
		return
	}

	for i := 0; i <= countRow; i++ {
		connDB.Create(&models.Item{
			Name:        faker.Lorem().Word(),
			Description: faker.Lorem().Sentence(5),
			Price:       faker.Number().NumberInt(3),
			Count:       faker.Number().NumberInt(5),
		})

	}
}
