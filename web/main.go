package main

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/Team-Fruit/Emojicord-web/web/handler"
	"github.com/Team-Fruit/Emojicord-web/web/model"
	"github.com/Team-Fruit/Emojicord-web/web/discord"
)

func main() {
	db := sqlx.MustConnect("mysql", "emojicord:@tcp(db:3306)/emojicord_db?parseTime=true")
	defer db.Close()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	m := model.NewModel(db)
	d := discord.NewClient(os.Getenv("BOT_TOKEN"))
	h := handler.NewHandler(m, d)

	e.GET("/auth/login", h.Auth)
	e.GET("/auth/callback", h.Callback)

	u := e.Group("/user")
	config := middleware.JWTConfig{
		Claims: &handler.JWTClaims{},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}
	u.Use(middleware.JWTWithConfig(config))
	u.GET("guilds", h.GetGuilds)

	e.Logger.Fatal(e.Start(":8082"))
}
