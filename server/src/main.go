package main

import (
    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
)

func main() {
    // Echoのインスタンス作る
    e := echo.New()

    // 全てのリクエストで差し込みたいミドルウェア（ログとか）はここ
    e.Use(middleware.CORS())
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    // ルーティング
    e.POST("/", hello)

    // サーバー起動
    e.Start(":8000")
}

func hello(c echo.Context) error {
    return c.String(200, "Hello, World!")
}