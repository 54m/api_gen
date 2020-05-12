package main

import (
	"fmt"
	"log"

	"github.com/gcp-kit/api_gen/server_generator/sample"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	sample.Bootstrap(e, sample.MiddlewareList{
		{
			Path: "/service/user2/",
			MiddlewareFunc: []echo.MiddlewareFunc{
				func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
					return func(context echo.Context) error {
						fmt.Println(context.Request().URL)
						return handlerFunc(context)
					}
				},
			},
		},
	})

	if err := e.Start("localhost:7899"); err != nil {
		log.Fatal(err)
	}
}
