package controller

import (
	"errors"
	"fmt"
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/request"
	"github.com/benebobaa/amikom-bri-api/domain/usecase"
	"github.com/benebobaa/amikom-bri-api/util"
	"github.com/gofiber/fiber/v2"
)

type WebController interface {
	VerifyEmail(ctx *fiber.Ctx) error
	ResetPasswordForm(ctx *fiber.Ctx) error
	ResetPasswordSubmit(ctx *fiber.Ctx) error
	ResetPasswordSuccess(ctx *fiber.Ctx) error
}

type webControllerImpl struct {
	userUsecase           usecase.UserUsecase
	forgotPasswordUsecase usecase.ForgotPasswordUsecase
}

func NewWebController(userUsecase usecase.UserUsecase, passwordUsecase usecase.ForgotPasswordUsecase) WebController {
	return &webControllerImpl{
		userUsecase:           userUsecase,
		forgotPasswordUsecase: passwordUsecase,
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
	secretToken := ctx.Query("secret", "")

	fmt.Println("secretToken", secretToken)
	err := w.forgotPasswordUsecase.CheckResetTokenIsUsed(ctx.UserContext(), secretToken)

	if err != nil {
		if errors.Is(err, util.ResetTokenAlreadyUsed) {
			return ctx.Render("resource/views/reset_password_expired", fiber.Map{})
		}
		resp, statusCode := util.ConstructBaseResponse(
			util.BaseResponse{
				Code:   fiber.StatusInternalServerError,
				Status: err.Error(),
			},
		)
		return ctx.Status(statusCode).JSON(resp)
	}

	return ctx.Render("resource/views/reset_password", fiber.Map{
		"SecretToken": secretToken,
	})
}

func (w *webControllerImpl) ResetPasswordSubmit(ctx *fiber.Ctx) error {
	requestBody := new(request.ResetPasswordRequest)

	secretToken := ctx.Query("secret", "")

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

	err = w.forgotPasswordUsecase.ResetPasswordRequest(ctx.UserContext(), requestBody, secretToken)
	if err != nil {
		if errors.Is(err, util.InvalidResetToken) {
			resp, statusCode := util.ConstructBaseResponse(
				util.BaseResponse{
					Code:   fiber.StatusUnauthorized,
					Status: err.Error(),
				},
			)
			return ctx.Status(statusCode).JSON(resp)
		}

		if errors.Is(err, util.ResetTokenAlreadyUsed) {
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

	return ctx.Redirect("/users/reset-password-success")
}

func (w *webControllerImpl) ResetPasswordSuccess(ctx *fiber.Ctx) error {
	return ctx.Render("resource/views/reset_password_success", fiber.Map{})
}
