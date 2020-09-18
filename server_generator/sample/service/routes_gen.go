// Code generated by server_generator. DO NOT EDIT.
// generated version: 0.4.0

package service

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

	router.GET("article", r.GetArticle(p))

	return r
}

func (r *Routes) GetArticle(p *props.ControllerProps) echo.HandlerFunc {
	i := NewGetArticleController(p)
	return func(c echo.Context) error {
		req := new(GetArticleRequest)
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "invalid request.",
			})
		}
		res, err := i.GetArticle(c, req)
		if err != nil {
			return err
		}
		if res == nil {
			return nil
		}

		return c.JSON(http.StatusOK, res)
	}
}

type IGetArticleController interface {
	GetArticle(c echo.Context, req *GetArticleRequest) (res *GetArticleResponse, err error)
}
