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

type UserController interface {
	Register(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	Profile(ctx *fiber.Ctx) error
	FindAll(ctx *fiber.Ctx) error
	ForgotPassword(ctx *fiber.Ctx) error
}

type userControllerImpl struct {
	userUsecase           usecase.UserUsecase
	loginUsecase          usecase.LoginUseCase
	forgotPasswordUsecase usecase.ForgotPasswordUsecase
}

func NewUserController(userUsecase usecase.UserUsecase, useCase usecase.LoginUseCase, passwordUsecase usecase.ForgotPasswordUsecase) UserController {
	return &userControllerImpl{
		userUsecase:           userUsecase,
		loginUsecase:          useCase,
		forgotPasswordUsecase: passwordUsecase,
	}
}

func (u *userControllerImpl) Register(ctx *fiber.Ctx) error {
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

	result, err := u.userUsecase.RegisterNewUser(ctx.UserContext(), requestBody, baseUrl)

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

func (u *userControllerImpl) Login(ctx *fiber.Ctx) error {
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

	loginResult, err := u.loginUsecase.LoginUser(ctx.UserContext(), requestBody, userAgent, clientIP)

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

func (u *userControllerImpl) Delete(ctx *fiber.Ctx) error {

	username := ctx.Params("username")
	authPayload := ctx.Locals(middleware.AuthorizationPayloadKey).(*token.Payload)

	err := u.userUsecase.DeleteUser(ctx.UserContext(), username, authPayload.Username)

	if err != nil {
		if errors.Is(err, util.UsernameNotFound) {
			resp, statusCode := util.ConstructBaseResponse(
				util.BaseResponse{
					Code:   fiber.StatusNotFound,
					Status: err.Error(),
				},
			)
			return ctx.Status(statusCode).JSON(resp)
		}
		if errors.Is(err, util.UnauthorizedDeleteUser) {
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
		},
	)

	return ctx.Status(statusCode).JSON(resp)
}

func (u *userControllerImpl) Profile(ctx *fiber.Ctx) error {
	authPayload := ctx.Locals(middleware.AuthorizationPayloadKey).(*token.Payload)

	result, err := u.userUsecase.ProfileUser(ctx.UserContext(), authPayload.UserID)

	if err != nil {
		if errors.Is(err, util.UserNotFound) {
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

func (u *userControllerImpl) FindAll(ctx *fiber.Ctx) error {

	request := &request.SearchPaginationRequest{
		Keyword: ctx.Query("keyword", ""),
		Page:    ctx.QueryInt("page", 1),
		Size:    ctx.QueryInt("size", 10),
	}

	result, err := u.userUsecase.GetAllUsers(ctx.UserContext(), request)

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

func (u *userControllerImpl) ForgotPassword(ctx *fiber.Ctx) error {
	requestBody := new(request.ForgotPasswordRequest)
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

	err = u.forgotPasswordUsecase.ForgotPasswordRequest(ctx.UserContext(), requestBody, baseUrl)

	if err != nil {
		if errors.Is(err, util.UsernameNotExist) {
			resp, statusCode := util.ConstructBaseResponse(
				util.BaseResponse{
					Code:   fiber.StatusNotFound,
					Status: err.Error(),
				},
			)
			return ctx.Status(statusCode).JSON(resp)
		}
		if errors.Is(err, util.EmailNotExists) {
			resp, statusCode := util.ConstructBaseResponse(
				util.BaseResponse{
					Code:   fiber.StatusNotFound,
					Status: err.Error(),
				},
			)
			return ctx.Status(statusCode).JSON(resp)
		}

		if errors.Is(err, util.UsernameAndEmailNotMatch) {
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
		},
	)

	return ctx.Status(statusCode).JSON(resp)
}
