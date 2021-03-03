// Code generated by server_generator. DO NOT EDIT.
// generated version: devel

package apigen

import (
	"io"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/go-generalize/api_gen/samples/empty_root/server/apigen/props"
	serverFooBar "github.com/go-generalize/api_gen/samples/empty_root/server/foo/bar"
	"github.com/labstack/echo/v4"
)

// MiddlewareList ...
type MiddlewareList []*MiddlewareSet

// MiddlewareMap ...
type MiddlewareMap map[string][]echo.MiddlewareFunc

// MiddlewareSet ...
type MiddlewareSet struct {
	Path           string
	MiddlewareFunc []echo.MiddlewareFunc
}

// ToMap ...
func (m MiddlewareList) ToMap() MiddlewareMap {
	mf := make(map[string][]echo.MiddlewareFunc)
	for _, middleware := range m {
		mf[middleware.Path] = middleware.MiddlewareFunc
	}
	return mf
}

// Bootstrap ...
func Bootstrap(p *props.ControllerProps, e *echo.Echo, middlewareList MiddlewareList, opts ...io.Writer) {
	if len(opts) > 0 {
		if w := opts[0]; w != nil {
			log.SetOutput(w)
		}
	}

	middleware := middlewareList.ToMap()

	// error handling
	e.Use(func(before echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				recoverErr := recover()
				if recoverErr == nil {
					return
				}

				debug.PrintStack()

				if httpErr, ok := recoverErr.(*echo.HTTPError); ok {
					err = c.JSON(httpErr.Code, httpErr.Message)
				}

				err = c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"code":    http.StatusInternalServerError,
					"message": "internal server error.",
				})
			}()

			return before(c)
		}
	})

	rootGroup := e.Group("/")
	setMiddleware(rootGroup, "/", middleware)

	serverFooBarGroup := rootGroup.Group("foo/bar/")
	setMiddleware(serverFooBarGroup, "/foo/bar/", middleware)
	serverFooBar.NewRoutes(p, serverFooBarGroup, opts...)
}

func setMiddleware(group *echo.Group, path string, list MiddlewareMap) {
	if ms, ok := list[path]; ok {
		for i := range ms {
			group.Use(ms[i])
		}
	}
}
