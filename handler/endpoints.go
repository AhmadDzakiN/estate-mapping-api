package handler

import (
	"database/sql"
	"fmt"
	openapi_types "github.com/oapi-codegen/runtime/types"
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

func (s *Server) CreateEstate(ctx echo.Context) error {
	var createReq generated.CreateEstateJSONBody
	err := ctx.Bind(&createReq)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	err = ctx.Validate(createReq)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	estateID, err := s.Repository.CreateEstate(ctx.Request().Context(), createReq.Length, createReq.Width)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Oops, something wrong with the server. Please try again later"})
	}

	return ctx.JSON(http.StatusCreated, map[string]string{"id": estateID})
}

func (s *Server) CreateTree(ctx echo.Context, estateID openapi_types.UUID) error {
	var createReq generated.CreateTreeJSONBody
	err := ctx.Bind(&createReq)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	err = ctx.Validate(createReq)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	estate, err := s.Repository.GetEstate(ctx.Request().Context(), estateID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Estate not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Oops, something wrong with the server. Please try again later"})
	}

	// Tree position is not in the estate's area
	if createReq.X > estate.Length || createReq.Y > estate.Width {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	treeID, err := s.Repository.CreateTree(ctx.Request().Context(), estate.ID, createReq.X, createReq.Y, createReq.Height)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Oops, something wrong with the server. Please try again later"})
	}

	return ctx.JSON(http.StatusCreated, map[string]string{"id": treeID})
}
