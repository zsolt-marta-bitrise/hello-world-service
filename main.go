package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/hello", HelloHandler)
	e.POST("/echo", EchoHandler)
	e.GET("/headers", HeadersHandler)

	e.GET("/livez", HelloHandler)
	e.GET("/readyz", HelloHandler)

	fmt.Println("Starting")
	if err := e.Start(":3000"); err != nil {
		fmt.Println(err)
	}
}

func HelloHandler(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Hello, world!")
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
