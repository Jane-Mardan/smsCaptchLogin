package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"server/api/v1"
	"server/storage"
)

func main() {
	storage.Init()
	router := gin.Default()
	setRouter(router)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 9992),
		Handler: router,
	}
	srv.ListenAndServe()
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listen: %s\n", zap.Error(err))
		}
	}()
}

func setRouter(router *gin.Engine) {
	f := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ping": "pong"})
	}
	router.Use(cors.Default())
	router.GET("/ping", f)

	apiRouter := router.Group("/api")
	{
		apiRouter.GET("/player/v1/SmsCaptcha/get", v1.Player.SMSCaptchaGet)
		apiRouter.POST("/player/v1/SmsCaptcha/login", v1.Player.SMSCaptchaLogin)
	}
}
