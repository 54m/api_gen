// Code generated by server_generator. DO NOT EDIT.
// generated version: unknown
package static

import (
	"net/http"

	"github.com/labstack/echo/v4"

	props "github.com/go-generalize/api_gen/server_generator/sample/props"
)

type Routes struct {
	router *echo.Group
}

func NewRoutes(p *props.ControllerProps, router *echo.Group) *Routes {
	r := &Routes{
		router: router,
	}

	router.GET("static_page", r.GetStaticPage(p))

	return r
}

func (r *Routes) GetStaticPage(p *props.ControllerProps) echo.HandlerFunc {
	i := NewGetStaticPageController(p)
	return func(c echo.Context) error {
		req := new(GetStaticPageRequest)
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "invalid request.",
			})
		}
		res, err := i.GetStaticPage(c, req)
		if err != nil {
			return err
		}
		if res == nil {
			return nil
		}

		return c.JSON(http.StatusOK, res)
	}
}

type IGetStaticPageController interface {
	GetStaticPage(c echo.Context, req *GetStaticPageRequest) (res *GetStaticPageResponse, err error)
}
