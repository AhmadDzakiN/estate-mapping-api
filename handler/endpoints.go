package handler

import (
	"errors"
	"fmt"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"gorm.io/gorm"
	"net/http"
)

func stringToUUID(uuidSTR string) (parsedUUID openapi_types.UUID) {
	parsedUUID, _ = uuid.Parse(uuidSTR)
	return
}

// This is just a test endpoint to get you started. Please delete this endpoint.
// (GET /hello)
func (s *Server) GetHello(ctx echo.Context, params generated.GetHelloParams) error {
	var resp generated.NotFoundErrorResponse
	resp.Error = fmt.Sprintf("Hello User %d", params.Id)
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) CreateEstate(ctx echo.Context) error {
	var createReq generated.CreateEstateJSONBody
	err := ctx.Bind(&createReq)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.InvalidInputErrorResponse{Error: "Invalid input"})
	}

	err = ctx.Validate(createReq)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.InvalidInputErrorResponse{Error: "Invalid input"})
	}

	newEstate := repository.Estate{
		Width:  createReq.Width,
		Length: createReq.Length,
	}

	err = s.Repository.CreateEstate(ctx.Request().Context(), &newEstate)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.InternalServerErrorResponse{Error: "Oops, something wrong with the server. Please try again later"})
	}

	resp := generated.CreateEstateResponse{
		Id: stringToUUID(newEstate.ID),
	}

	return ctx.JSON(http.StatusCreated, resp)
}

func (s *Server) CreateTree(ctx echo.Context, estateID openapi_types.UUID) error {
	var createReq generated.CreateTreeJSONBody
	err := ctx.Bind(&createReq)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.InvalidInputErrorResponse{Error: "Invalid input"})
	}

	err = ctx.Validate(createReq)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.InvalidInputErrorResponse{Error: "Invalid input"})
	}

	estate, err := s.Repository.GetEstateByID(ctx.Request().Context(), estateID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.JSON(http.StatusNotFound, generated.NotFoundErrorResponse{Error: "Estate not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, generated.InternalServerErrorResponse{Error: "Oops, something wrong with the server. Please try again later"})
	}

	// Tree position is out of the estate's area
	if createReq.X > estate.Length || createReq.Y > estate.Width {
		return ctx.JSON(http.StatusBadRequest, generated.InvalidInputErrorResponse{Error: "Tree position is out of the estate's area"})
	}

	newTree := repository.Tree{
		EstateID:           estate.ID,
		HorizontalPosition: createReq.X,
		VerticalPosition:   createReq.Y,
		Height:             createReq.Height,
	}

	err = s.Repository.CreateTree(ctx.Request().Context(), &newTree)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.InternalServerErrorResponse{Error: "Oops, something wrong with the server. Please try again later"})
	}

	resp := generated.CreateTreeResponse{
		Id: stringToUUID(newTree.ID),
	}

	return ctx.JSON(http.StatusCreated, resp)
}

func (s *Server) GetEstateStats(ctx echo.Context, estateID openapi_types.UUID) error {
	estate, err := s.Repository.GetEstateByID(ctx.Request().Context(), estateID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.JSON(http.StatusNotFound, generated.NotFoundErrorResponse{Error: "Estate not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, generated.InternalServerErrorResponse{Error: "Oops, something wrong with the server. Please try again later"})
	}

	trees, err := s.Repository.GetTreeHeightsByEstateID(ctx.Request().Context(), estate.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.InternalServerErrorResponse{Error: "Oops, something wrong with the server. Please try again later"})
	}

	var resp generated.GetEstateStatsResponse
	treeCount := len(trees)
	if treeCount > 0 {
		resp.Count = treeCount
		resp.Max = trees[treeCount-1]
		resp.Min = trees[0]

		if treeCount%2 == 0 {
			resp.Median = (trees[treeCount/2-1] + trees[treeCount/2]) / 2 // the div result will be floored
		} else {
			resp.Median = trees[treeCount/2]
		}
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) GetEstateDronePlan(ctx echo.Context, estateId openapi_types.UUID) error {
	estate, err := s.Repository.GetEstateByID(ctx.Request().Context(), estateId.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.JSON(http.StatusNotFound, generated.NotFoundErrorResponse{Error: "Estate not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, generated.InternalServerErrorResponse{Error: "Oops, something wrong with the server. Please try again later"})
	}

	trees, err := s.Repository.GetTreesByEstateIDAndPlotsLocations(ctx.Request().Context(), estate.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.InternalServerErrorResponse{Error: "Oops, something wrong with the server. Please try again later"})
	}

	var resp generated.GetEstateDronePlanResponse
	if len(trees) > 0 {
		var totalTreeHeight int
		for _, tree := range trees {
			totalTreeHeight += tree.Height
		}
		resp.Distance += totalTreeHeight
	}

	// 2 is from the drone flies above the plot and returns down to the ground (1m + 1m)
	// each axis -1 is from when the farthest plot (10x10 square meter) of each axis will not be counted (got this logic from the documentation of the test)
	resp.Distance += (estate.Length-1)*10 + (estate.Width-1)*10 + 2

	return ctx.JSON(http.StatusOK, resp)
}
