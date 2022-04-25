package migrations

import (
	"github.com/akifanabil/synapsis-backend-challenge/helpers"
	"github.com/akifanabil/synapsis-backend-challenge/interfaces"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	er := godotenv.Load(".env")
	if er != nil {
		panic(er.Error())
	}

	db_host := os.Getenv("DB_HOST")
	db_driver := os.Getenv("DB_DRIVER")
	db_port := os.Getenv("DB_PORT")
	db_user := os.Getenv("DB_USER")
	db_pass := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")

	dbURL := db_driver + "://" + db_user + ":" + db_pass + "@" + db_host + ":" + db_port + "/" + db_name +"?sslmode=require" 
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	helpers.HandleError(err)

	db.AutoMigrate(&interfaces.Customer{}, &interfaces.Product{}, &interfaces.Cart{}, &interfaces.Transaction{})
	return db
}
