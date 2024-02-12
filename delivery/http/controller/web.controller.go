package controller

import (
	"errors"
	"github.com/benebobaa/amikom-bri-api/domain/usecase"
	"github.com/benebobaa/amikom-bri-api/util"
	"github.com/gofiber/fiber/v2"
)

type WebController interface {
	VerifyEmail(ctx *fiber.Ctx) error
	ResetPasswordForm(ctx *fiber.Ctx) error
	ResetPasswordSubmit(ctx *fiber.Ctx) error
}

type webControllerImpl struct {
	userUsecase usecase.UserUsecase
}

func NewWebController(userUsecase usecase.UserUsecase) WebController {
	return &webControllerImpl{
		userUsecase: userUsecase,
	}
}

func (w *webControllerImpl) VerifyEmail(ctx *fiber.Ctx) error {

	secretCode := ctx.Query("secret", "")

	result, err := w.userUsecase.VerifyUserEmail(ctx.UserContext(), secretCode)

	if err != nil {
		if errors.Is(err, util.EmailVerifyExpired) {
			resp, statusCode := util.ConstructBaseResponse(
				util.BaseResponse{
					Code:   fiber.StatusUnauthorized,
					Status: err.Error(),
				},
			)
			return ctx.Status(statusCode).JSON(resp)
		}
		if errors.Is(err, util.EmailVerifyCodeNotValid) {
			resp, statusCode := util.ConstructBaseResponse(
				util.BaseResponse{
					Code:   fiber.StatusUnauthorized,
					Status: err.Error(),
				},
			)
			return ctx.Status(statusCode).JSON(resp)
		}
		if errors.Is(err, util.EmailVerifyAlreadyUsed) {
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
			Status: "Success",
			Data:   result,
		},
	)

	return ctx.Render("resource/views/verify_success", fiber.Map{
		"response": resp,
		"status":   statusCode,
	})
}

func (w *webControllerImpl) ResetPasswordForm(ctx *fiber.Ctx) error {
	return ctx.Render("resource/views/test", fiber.Map{})
}

func (w *webControllerImpl) ResetPasswordSubmit(ctx *fiber.Ctx) error {
	return ctx.Render("resource/views/reset_password", fiber.Map{})
}
