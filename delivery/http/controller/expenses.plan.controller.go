package controller

import (
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/request"
	"github.com/benebobaa/amikom-bri-api/delivery/http/middleware"
	"github.com/benebobaa/amikom-bri-api/domain/usecase"
	"github.com/benebobaa/amikom-bri-api/util"
	"github.com/benebobaa/amikom-bri-api/util/token"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type ExpensesPlanController interface {
	Create(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	FindAllFilter(ctx *fiber.Ctx) error
	FindAllExportPdf(ctx *fiber.Ctx) error
}

type expensesPlanControllerImpl struct {
	expensesPlanUsecase usecase.ExpensesPlanUsecase
}

func NewExpensesPlanController(expensesPlanUsecase usecase.ExpensesPlanUsecase) ExpensesPlanController {
	return &expensesPlanControllerImpl{
		expensesPlanUsecase: expensesPlanUsecase,
	}
}

func (e *expensesPlanControllerImpl) Create(ctx *fiber.Ctx) error {
	requestBody := new(request.ExpensesPlanRequest)
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

	result, err := e.expensesPlanUsecase.ExpensesPlanCreate(ctx.UserContext(), requestBody, authPayload.UserID)

	if err != nil {
		resp, statusCode := util.ConstructBaseResponse(
			util.BaseResponse{
				Code:   fiber.StatusBadRequest,
				Status: err.Error(),
			},
		)
		return ctx.Status(statusCode).JSON(resp)
	}

	return ctx.Status(fiber.StatusCreated).JSON(result)

}

func (e *expensesPlanControllerImpl) Update(ctx *fiber.Ctx) error {
	requestBody := new(request.ExpensesPlanUpdateRequest)
	authPayload := ctx.Locals(middleware.AuthorizationPayloadKey).(*token.Payload)

	err := ctx.BodyParser(requestBody)

	paramId := ctx.Params("id")

	id, err := strconv.Atoi(paramId)
	if err != nil {
		resp, statusCode := util.ConstructBaseResponse(
			util.BaseResponse{
				Code:   fiber.StatusBadRequest,
				Status: err.Error(),
			},
		)
		return ctx.Status(statusCode).JSON(resp)
	}

	result, err := e.expensesPlanUsecase.UpdatePlan(ctx.UserContext(), requestBody, int64(id), authPayload.UserID)

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

func (e *expensesPlanControllerImpl) Delete(ctx *fiber.Ctx) error {
	authPayload := ctx.Locals(middleware.AuthorizationPayloadKey).(*token.Payload)
	paramId := ctx.Params("id")

	id, err := strconv.Atoi(paramId)
	if err != nil {
		resp, statusCode := util.ConstructBaseResponse(
			util.BaseResponse{
				Code:   fiber.StatusBadRequest,
				Status: err.Error(),
			},
		)
		return ctx.Status(statusCode).JSON(resp)
	}

	err = e.expensesPlanUsecase.DeletePlan(ctx.UserContext(), int64(id), authPayload.UserID)

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
		},
	)
	return ctx.Status(statusCode).JSON(resp)

}

func (e *expensesPlanControllerImpl) FindAllFilter(ctx *fiber.Ctx) error {
	authPayload := ctx.Locals(middleware.AuthorizationPayloadKey).(*token.Payload)
	request := &request.SearchPaginationRequest{
		Keyword: ctx.Query("keyword", ""),
		Page:    ctx.QueryInt("page", 1),
		Size:    ctx.QueryInt("size", 10),
	}

	result, err := e.expensesPlanUsecase.FindAllFilter(ctx.UserContext(), request, authPayload.UserID)

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

func (e *expensesPlanControllerImpl) FindAllExportPdf(ctx *fiber.Ctx) error {
	authPayload := ctx.Locals(middleware.AuthorizationPayloadKey).(*token.Payload)

	filename, err := e.expensesPlanUsecase.ExpensesPlanExportPdf(ctx.UserContext(), authPayload.UserID)

	if err != nil {
		resp, statusCode := util.ConstructBaseResponse(
			util.BaseResponse{
				Code:   fiber.StatusBadRequest,
				Status: err.Error(),
			},
		)
		return ctx.Status(statusCode).JSON(resp)
	}

	return ctx.SendFile(filename, true)
}
