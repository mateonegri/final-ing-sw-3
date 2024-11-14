package db

import (
	userClient "user-api/client"
	"user-api/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

var (
	db  *gorm.DB
	err error
)

func init() {
	// DB Connections Paramters
	DBName := "ing-sw-3"
	DBUser := "root"
	DBPass := "mateo222"
	//DBPass := os.Getenv("MVC_DB_PASS")
	//DBHost := "mysql"
	DBHost := "localhost"
	DBPort := "3307"
	// ------------------------

	db, err = gorm.Open("mysql", DBUser+":"+DBPass+"@tcp("+DBHost+":"+DBPort+")/"+DBName+"?charset=utf8&parseTime=True")

	if err != nil {
		log.Info("Connection Failed to Open")
		log.Fatal(err)
	} else {
		log.Info("Connection Established")
	}

	// We need to add all Clients that we build
	userClient.Db = db

}

func StartDbEngine() {
	// We need to migrate all classes model.
	db.AutoMigrate(&model.User{})

	log.Info("Finishing Migration Database Tables")
}
