package controller

import (
	"errors"
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/request"
	"github.com/benebobaa/amikom-bri-api/delivery/http/middleware"
	"github.com/benebobaa/amikom-bri-api/domain/usecase"
	"github.com/benebobaa/amikom-bri-api/util"
	"github.com/benebobaa/amikom-bri-api/util/token"
	"github.com/gofiber/fiber/v2"
)

type TransferController interface {
	Transfer(ctx *fiber.Ctx) error
}

type transferControllerImpl struct {
	transferUsecase usecase.TransferUsecase
}

func NewTransferController(transferUsecase usecase.TransferUsecase) TransferController {
	return &transferControllerImpl{
		transferUsecase: transferUsecase,
	}
}

func (t *transferControllerImpl) Transfer(ctx *fiber.Ctx) error {
	requestBody := new(request.TransferRequest)
	authPayload := ctx.Locals(middleware.AuthorizationPayloadKey).(*token.Payload)

	err := ctx.BodyParser(requestBody)
	if err != nil {
		resp, statusCode := util.ConstructBaseResponse(
			util.BaseResponse{
				Code:   fiber.StatusBadRequest,
				Status: err.Error(),
			},
		)
		return ctx.Status(statusCode).JSON(resp)
	}

	err = t.transferUsecase.TransferMoney(ctx.UserContext(), requestBody, authPayload.UserID)

	if err != nil {
		if errors.Is(err, util.AccountNotBelongToUser) {
			resp, statusCode := util.ConstructBaseResponse(
				util.BaseResponse{
					Code:   fiber.StatusBadRequest,
					Status: err.Error(),
				},
			)
			return ctx.Status(statusCode).JSON(resp)
		}
		if errors.Is(err, util.InsufficientBalance) {
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
				Code:   fiber.StatusBadRequest,
				Status: err.Error(),
			},
		)
		return ctx.Status(statusCode).JSON(resp)
	}

	resp, statusCode := util.ConstructBaseResponse(
		util.BaseResponse{
			Code:   fiber.StatusOK,
			Status: "Transfer success",
		},
	)

	return ctx.Status(statusCode).JSON(resp)
}
