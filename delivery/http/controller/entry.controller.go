package controller

import (
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/request"
	"github.com/benebobaa/amikom-bri-api/delivery/http/middleware"
	"github.com/benebobaa/amikom-bri-api/domain/usecase"
	"github.com/benebobaa/amikom-bri-api/util"
	"github.com/benebobaa/amikom-bri-api/util/token"
	"github.com/gofiber/fiber/v2"
)

type EntryController interface {
	FindAllEntries(ctx *fiber.Ctx) error
	FindAllFilterDate(ctx *fiber.Ctx) error
}

type entryControllerImpl struct {
	entryUsecase usecase.EntryUsecase
}

func NewEntryController(entryUsecase usecase.EntryUsecase) EntryController {
	return &entryControllerImpl{
		entryUsecase: entryUsecase,
	}
}

func (e *entryControllerImpl) FindAllEntries(ctx *fiber.Ctx) error {

	authPayload := ctx.Locals(middleware.AuthorizationPayloadKey).(*token.Payload)

	request := &request.SearchPaginationRequest{
		Keyword: ctx.Query("entry_type", ""),
		Page:    ctx.QueryInt("page", 1),
		Size:    ctx.QueryInt("size", 10),
	}

	result, err := e.entryUsecase.FindAllHistoryTransfer(ctx.UserContext(), request, authPayload.UserID)

	if err != nil {
		resp, statusCode := util.ConstructBaseResponse(
			util.BaseResponse{
				Code:   fiber.StatusBadRequest,
				Status: err.Error(),
			},
		)
		return ctx.Status(statusCode).JSON(resp)
	}

	resp, statusCode := util.ConstructBaseResponse(
		util.BaseResponse{
			Code:   fiber.StatusOK,
			Status: "Success",
			Data:   result,
		},
	)

	return ctx.Status(statusCode).JSON(resp)
}

func (e *entryControllerImpl) FindAllFilterDate(ctx *fiber.Ctx) error {

	authPayload := ctx.Locals(middleware.AuthorizationPayloadKey).(*token.Payload)

	request := &request.SearchPaginationRequest{
		Filter: ctx.Query("filter", ""),
		Date:   ctx.Query("date", ""),
		Page:   ctx.QueryInt("page", 1),
		Size:   ctx.QueryInt("size", 10),
	}

	result, err := e.entryUsecase.FindAllFilterDate(ctx.UserContext(), request, authPayload.UserID)

	if err != nil {
		resp, statusCode := util.ConstructBaseResponse(
			util.BaseResponse{
				Code:   fiber.StatusBadRequest,
				Status: err.Error(),
			},
		)
		return ctx.Status(statusCode).JSON(resp)
	}

	resp, statusCode := util.ConstructBaseResponse(
		util.BaseResponse{
			Code:   fiber.StatusOK,
			Status: "Success",
			Data:   result,
		},
	)

	return ctx.Status(statusCode).JSON(resp)
}
