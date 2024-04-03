package main

import (
	"github.com/golfz/fun-exercise-api/postgres"
	"github.com/golfz/fun-exercise-api/wallet"
	"github.com/labstack/echo/v4"

	_ "github.com/golfz/fun-exercise-api/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title			Wallet API
// @version		1.0
// @description	Sophisticated Wallet API
// @host			localhost:1323
func main() {
	p, err := postgres.New()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	handler := wallet.New(p)

	g := e.Group("/api/v1")
	g.GET("/wallets", handler.GetWalletsHandler)           // challenge 3
	g.GET("/users/:id/wallets", handler.UserWalletHandler) // challenge 4
	g.POST("/wallets", handler.CreateWalletHandler)
	g.PUT("/wallets", handler.UpdateWalletHandler)
	g.DELETE("/users/:id/wallets", handler.DeleteUserWalletHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
