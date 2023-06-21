package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const requestTimeout = 30 * time.Second

func main() {
	config := LoadConfig()

	now := time.Now()
	now.Add(3 * time.Second)

	if now.After(time.Now()) {

	}

	db, err := gorm.Open(mysql.Open(config.DatabaseDSN), &gorm.Config{})
	if err != nil {
		log.Fatalln("Open database:", err.Error())
	}
	db.AutoMigrate(&Account{})

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(MiddWhitelistRequest(config.Whitelist, config.AllowMethods))

	e.Static("/", "./")
	e.GET("/welcome", welcome)
	e.GET("/accounts", accounts(db))
	// Just For Simulation
	dummySimulation(e, config.Listen)

	go func() {
		if err := e.Start(fmt.Sprintf(":%s", config.Listen)); err != nil && err != http.ErrServerClosed {
			log.Fatalln("Failed start HTTP server:", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalln("Failed shuttdown HTTP server:", err.Error())
	}
}

func dummySimulation(e *echo.Echo, port string) {
	// Simulasi Issues Cross-Site Request Forgrey (CSRF)
	// Domain yang digunakan attacker untuk melakukan redirect kembali
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://172.17.0.1:"+port+"/transfer?amount=1000&account=attacker-account", http.StatusMovedPermanently)
	})
	go http.ListenAndServe(":8080", nil)

	e.GET("/transfer", func(c echo.Context) error {
		params := c.QueryParams()
		return c.JSON(http.StatusCreated, echo.Map{"destination": params, "message": "Transfer Berhasil"})
	})
}

func CheckValue(val uint) bool {
	if val == 0 {
		return false
	}
	return true
}
