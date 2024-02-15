package controller

import (
	"errors"
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/request"
	"github.com/benebobaa/amikom-bri-api/domain/usecase"
	"github.com/benebobaa/amikom-bri-api/util"
	"github.com/gofiber/fiber/v2"
)

type AuthController interface {
	Register(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	RenewAccessToken(ctx *fiber.Ctx) error
}

type authControllerImpl struct {
	userUsecase    usecase.UserUsecase
	loginUsecase   usecase.LoginUseCase
	sessionUsecase usecase.SessionUseCase
}

func NewAuthController(userUsecase usecase.UserUsecase, useCase usecase.LoginUseCase, sessionUseCase usecase.SessionUseCase) AuthController {
	return &authControllerImpl{
		userUsecase:    userUsecase,
		loginUsecase:   useCase,
		sessionUsecase: sessionUseCase,
	}
}

func (a *authControllerImpl) Register(ctx *fiber.Ctx) error {
	requestBody := new(request.UserRegisterRequest)
	baseUrl := ctx.BaseURL()
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

	result, err := a.userUsecase.RegisterNewUser(ctx.UserContext(), requestBody, baseUrl)

	if err != nil {
		if errors.Is(err, util.UsernameAlreadyExists) {
			resp, statusCode := util.ConstructBaseResponse(
				util.BaseResponse{
					Code:   fiber.StatusBadRequest,
					Status: err.Error(),
				},
			)
			return ctx.Status(statusCode).JSON(resp)
		}
		if errors.Is(err, util.EmailAlreadyExists) {
			resp, statusCode := util.ConstructBaseResponse(
				util.BaseResponse{
					Code:   fiber.StatusBadRequest,
					Status: err.Error(),
				},
			)
			return ctx.Status(statusCode).JSON(resp)

		} else {
			resp, statusCode := util.ConstructBaseResponse(
				util.BaseResponse{
					Code:   fiber.StatusBadRequest,
					Status: err.Error(),
				},
			)
			return ctx.Status(statusCode).JSON(resp)
		}
	}

	resp, statusCode := util.ConstructBaseResponse(
		util.BaseResponse{
			Code:   fiber.StatusCreated,
			Status: "Success",
			Data:   result,
		},
	)

	return ctx.Status(statusCode).JSON(resp)
}

func (a *authControllerImpl) Login(ctx *fiber.Ctx) error {
	requestBody := new(request.LoginRequest)
	userAgent := ctx.Get("User-Agent")
	clientIP := ctx.IP()

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

	loginResult, err := a.loginUsecase.LoginUser(ctx.UserContext(), requestBody, userAgent, clientIP)

	if err != nil {
		if errors.Is(err, util.InvalidPassword) {
			resp, statusCode := util.ConstructBaseResponse(
				util.BaseResponse{
					Code:   fiber.StatusUnauthorized,
					Status: err.Error(),
				},
			)
			return ctx.Status(statusCode).JSON(resp)
		}
		if errors.Is(err, util.UsernameOrEmailNotFound) {
			resp, statusCode := util.ConstructBaseResponse(
				util.BaseResponse{
					Code:   fiber.StatusNotFound,
					Status: err.Error(),
				},
			)
			return ctx.Status(statusCode).JSON(resp)
		}
		if errors.Is(err, util.EmailNotVerified) {
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
			Data:   loginResult,
		},
	)

	return ctx.Status(statusCode).JSON(resp)
}

func (a *authControllerImpl) RenewAccessToken(ctx *fiber.Ctx) error {
	requestBody := new(request.TokenRenewRequest)

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

	result, err := a.sessionUsecase.RenewAccessToken(ctx.UserContext(), requestBody)
	if err != nil {
		if errors.Is(err, util.SessionNotFound) {
			resp, statusCode := util.ConstructBaseResponse(
				util.BaseResponse{
					Code:   fiber.StatusNotFound,
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

	return ctx.Status(statusCode).JSON(resp)
}
