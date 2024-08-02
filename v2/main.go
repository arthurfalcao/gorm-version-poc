package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/SWORDHealth/lib-go-core/testresource"
)

type MessageExchange struct {
	ID        int       `gorm:"primaryKey"`
	Name      string    `gorm:"column:name"`
	MessageID string    `gorm:"column:message_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func main() {
	conn, resource, err := testresource.MySQL(testresource.MySQLOptions{
		Username: "test",
		Password: "test",
		Database: fmt.Sprintf("gorm_v2_poc_%d", time.Now().UnixNano()),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer resource.Close()

	// db, err := gorm.Open(sqlite.Open("v2.db"), &gorm.Config{})
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: conn}))
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// Migrate the schema
	db.AutoMigrate(&MessageExchange{})

	var messagesExchange []MessageExchange
	for i := 1; i <= 100; i++ {
		messagesExchange = append(messagesExchange, MessageExchange{
			Name: fmt.Sprintf("MessageExchange %d", i),
		})
	}

	err = db.Create(&messagesExchange).Error
	if err != nil {
		log.Fatal(err)
	}

	var messageExchange MessageExchange
	err = db.
		Debug().
		Table("message_exchanges").
		Where("message_id = ?", "").
		Scan(&messageExchange).Error
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(messageExchange)

	err = db.
		Debug().
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
