// Package user2 ...
// generated version: unknown
package user2

import (
	props "github.com/go-generalize/api_gen/server_generator/sample/props"
	"github.com/labstack/echo/v4"
)

// PostUpdateUserNameController ...
type PostUpdateUserNameController struct {
	*props.ControllerProps
}

// NewPostUpdateUserNameController ...
func NewPostUpdateUserNameController(props *props.ControllerProps) *PostUpdateUserNameController {
	p := &PostUpdateUserNameController{
		ControllerProps: props,
	}
	return p
}

// PostUpdateUserName ...
// @Summary WIP
// @Description WIP
// @Accept json
// @Produce json
// @Param Name body string WIP:${isRequire} WIP:${description}
// @Success 200 {object} PostUpdateUserNameResponse
// @Failure 400 {object} WIP
// @Router /service/user2/update_user_name [POST]
func (p *PostUpdateUserNameController) PostUpdateUserName(
	c echo.Context, req *PostUpdateUserNameRequest,
) (res *PostUpdateUserNameResponse, err error) {
	panic("require implements.")
}
