// Package user2 ...
// generated version: devel
package user2

import (
	"github.com/go-generalize/api_gen/server_generator/sample/props"
	"github.com/labstack/echo/v4"
)

// GetUserController ...
type GetUserController struct {
	*props.ControllerProps
}

// NewGetUserController ...
func NewGetUserController(cp *props.ControllerProps) *GetUserController {
	g := &GetUserController{
		ControllerProps: cp,
	}
	return g
}

// GetUser ...
// @Summary WIP
// @Description WIP
// @Accept json
// @Produce json
// @Param userID path string WIP:${isRequire} WIP:${description}
// @Param search_request query string WIP:${isRequire} WIP:${description}
// @Success 200 {object} GetUserResponse
// @Failure 400 {object} wrapper.APIError
// @Failure 500 {object} wrapper.APIError
// @Router /service/user2/{userID} [GET]
func (g *GetUserController) GetUser(
	_ echo.Context, req *GetUserRequest,
) (res *GetUserResponse, err error) {
	return &GetUserResponse{
		ID:            req.ID,
		SearchRequest: req.SearchRequest,
	}, nil
}
