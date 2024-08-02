package main

import (
	"fmt"
	"log"
	"time"

	"github.com/SWORDHealth/lib-go-core/testresource"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type MessageExchange struct {
	ID        int       `gorm:"primary_key"`
	Name      string    `gorm:"column:name"`
	MessageID string    `gorm:"column:message_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func main() {
	dbName := fmt.Sprintf("gorm_v1_poc_%d", time.Now().UnixNano())
	_, resource, err := testresource.MySQL(testresource.MySQLOptions{
		Username: "test",
		Password: "test",
		Database: dbName,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer resource.Close()

	// db, err := gorm.Open("sqlite3", "v1.db")
	db, err := gorm.Open("mysql", fmt.Sprintf("root:root@tcp(127.0.0.1:3306)/%s?parseTime=true", dbName))
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	defer db.Close()

	// Enable logging of SQL statements
	db.LogMode(true)

	// Migrate the schema
	db.AutoMigrate(&MessageExchange{})

	for i := 1; i <= 100; i++ {
		err = db.Create(&MessageExchange{
			Name: fmt.Sprintf("MessageExchange %d", i),
		}).Error
		if err != nil {
			log.Fatal(err)
		}
	}

	var messageExchange MessageExchange
	err = db.
		Table("message_exchanges").
		Where("message_id = ?", "").
		Scan(&messageExchange).Error
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(messageExchange)

	err = db.
		Table("message_exchanges").
		Where("message_id = ?", messageExchange.MessageID).
		Updates(map[string]interface{}{
			"name": fmt.Sprintf("updated %d", time.Now().Unix()),
		}).
		Error
	if err != nil {
		log.Fatal(err)
	}
}
