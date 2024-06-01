package handler

import (
	"fmt"
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/labstack/echo/v4"
)

// This is just a test endpoint to get you started. Please delete this endpoint.
// (GET /hello)
func (s *Server) GetHello(ctx echo.Context, params generated.GetHelloParams) error {
	var resp generated.HelloResponse
	resp.Message = fmt.Sprintf("Hello User %d", params.Id)
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) CreateEstate(ctx echo.Context) (err error) {
	var createReq generated.CreateEstateJSONBody
	err = ctx.Bind(&createReq)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	estateID, err := s.Repository.CreateEstate(ctx.Request().Context(), createReq.Length, createReq.Width)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Error Cui"})
	}

	return ctx.JSON(http.StatusCreated, map[string]string{"id": estateID})
}
