package main

import (
	"fmt"
	"github.com/dnk90/chat/internal/log"
	"github.com/dnk90/chat/internal/models"
	"github.com/dnk90/chat/internal/services"
	"github.com/dnk90/chat/internal/ws"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"

	"github.com/dnk90/chat/internal/config"
	"github.com/gin-gonic/gin"
)

var ll = log.New()

type NewRoomRequest struct {
	Username string   `json:"username"`
}

func initDB(db *gorm.DB) {
	m := []interface{}{
		&models.Message{},
		&models.Room{},
	}
	if err := db.AutoMigrate(m...); err != nil {
		panic(err)
	}
}

func mustConnectMySQL(cfg *config.Config) *gorm.DB {
	if cfg.DB != nil {
		return cfg.DB
	}
	url := cfg.MySQLDSN()
	ll.S.Infow("[DB Connecting]", "url", url)
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if cfg.Environment == "D" {
		db = db.Debug()
	}

	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDb.SetConnMaxLifetime(time.Duration(cfg.MySQL.ConnMaxLifetime) * time.Hour)
	sqlDb.SetMaxIdleConns(cfg.MySQL.MaxIdleConns)
	sqlDb.SetMaxOpenConns(cfg.MySQL.MaxOpenConns)

	err = db.Raw("SELECT 1").Error
	if err != nil {
		ll.S.Fatalw("error querying SELECT 1", "err", err)
	}

	return db
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	cfg := config.Load()
	cfg.DB = mustConnectMySQL(cfg)
	initDB(cfg.DB)
	hub := ws.NewHub()
	go hub.Run()

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	router.Use(CORSMiddleware())

	service := services.NewService()
	// this api is used for checking liveness and readiness
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "OK")
	})
	router.POST("/v1/room", func(c *gin.Context){
		var newRoomRequest NewRoomRequest
		if err := c.ShouldBindJSON(&newRoomRequest); err != nil {
			ll.S.Errorw("[NewRoomRequest]binding", "err", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid params"})
			return
		}
		id, err := service.NewRoom(newRoomRequest.Username)
		if err != nil {
			ll.S.Errorw("[NewRoomRequest]NewRoom", "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create new room"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": id})
	})

	router.GET("/v1/room/:roomId/messages", func(c *gin.Context) {
		roomId := c.Param("roomId")
		fromId, err := strconv.Atoi(c.Query("fromId"))
		if err != nil {
			ll.S.Errorw("[GetMessages]Convert fromId to number:", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get fromId"})
			return
		}
		messages, err := service.GetMessages(fromId, cfg.LimitMessage, roomId)
		if err != nil {
			ll.S.Errorw("[GetMessages]Call Get Messages function", "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while getting messages"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": messages})
	})

	router.GET("/ws/:roomId", func(c *gin.Context) {
		roomId := c.Param("roomId")
		hub.ServeWs(c.Writer, c.Request, roomId)
	})

	ll.S.Infow("Live Chat is starting", "port", cfg.HttpPort)
	router.Run(fmt.Sprintf("0.0.0.0:%v", cfg.HttpPort))
}
