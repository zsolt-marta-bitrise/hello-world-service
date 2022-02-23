package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/hello", HelloHandler)
	e.POST("/echo", EchoHandler)
	e.GET("/headers", HeadersHandler)
	e.GET("/db/ping", DBHandler)
	e.GET("/env", EnvHandler)

	e.GET("/livez", HelloHandler)
	e.GET("/readyz", HelloHandler)

	fmt.Println("Starting")
	if err := e.Start(":3000"); err != nil {
		fmt.Println(err)
	}
}

func HelloHandler(ctx echo.Context) error {
    text := os.Getenv("HELLO_TEXT")
    if text == "" {
      text = "Hello, world!"
    }
	return ctx.String(http.StatusOK, text)
}

func EchoHandler(ctx echo.Context) error {
	body := ctx.Request().Body
	defer body.Close()

	bodybytes, err := ioutil.ReadAll(body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Cannot read body")
	}

	return ctx.Blob(http.StatusOK, ctx.Request().Header.Get("Content-type"), bodybytes)
}

func HeadersHandler(ctx echo.Context) error {
	resp := make(map[string]interface{})

	headers := make(map[string]string)
	for key, values := range ctx.Request().Header {
		headers[key] = strings.Join(values, ",")
	}
	resp["headers"] = headers
	resp["requestUri"] = ctx.Request().RequestURI
	resp["remoteAddress"] = ctx.Request().RemoteAddr

	return ctx.JSON(http.StatusOK, resp)
}

func DBHandler(ctx echo.Context) error {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	defer conn.Close(context.Background())
	if err = conn.Ping(context.Background()); err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.String(http.StatusOK, "OK")
}

func EnvHandler(ctx echo.Context) error {
	resp := make(map[string]interface{})

	for _, env := range os.Environ() {
		split := strings.Split(env, "=")
		if len(split) < 2 {
			continue
		}
		key := split[0]
		value := os.Getenv(key)
		resp[key] = value
	}

	return ctx.JSON(http.StatusOK, resp)
}
