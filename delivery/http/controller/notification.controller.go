package controller

import (
	"github.com/benebobaa/amikom-bri-api/delivery/http/middleware"
	"github.com/benebobaa/amikom-bri-api/domain/usecase"
	"github.com/benebobaa/amikom-bri-api/util"
	"github.com/benebobaa/amikom-bri-api/util/token"
	"github.com/gofiber/fiber/v2"
)

type NotificationController interface {
	FindALL(ctx *fiber.Ctx) error
}

type notificationControllerImpl struct {
	notificationUsecase usecase.NotificationUsecase
}

func NewNotificationController(notificationUsecase usecase.NotificationUsecase) NotificationController {
	return &notificationControllerImpl{
		notificationUsecase: notificationUsecase,
	}
}

func (n *notificationControllerImpl) FindALL(ctx *fiber.Ctx) error {
	authPayload := ctx.Locals(middleware.AuthorizationPayloadKey).(*token.Payload)

	result, err := n.notificationUsecase.GetAllNotifHistory(ctx.UserContext(), authPayload.UserID)

	if err != nil {
		resp, statusCode := util.ConstructBaseResponse(
			util.BaseResponse{
				Code:   fiber.StatusInternalServerError,
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
