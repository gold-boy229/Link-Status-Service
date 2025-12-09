package handlers

import (
	"net/http"

	"Link-Status-Service/internal/consts"
	"Link-Status-Service/internal/dto"
	"Link-Status-Service/internal/entity"

	"github.com/labstack/echo"
)

func (h linkHandler) GetStatus(c echo.Context) error {
	var reqDTO dto.LinksGetStatusRequest
	if err := c.Bind(&reqDTO); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewError(consts.ErrorCodeBadRequest, err.Error()))
	}
	if err := c.Validate(reqDTO); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewError(consts.ErrorCodeBadRequest, err.Error()))
	}

	params := convertDTOToEntityLinksGetStatus(reqDTO)
	result, err := h.LinkService.GetStatus(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			dto.NewError(consts.ErrorCodeInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, convertEntityToDTOLinksGetStatus(result))
}

func convertDTOToEntityLinksGetStatus(reqDTO dto.LinksGetStatusRequest) entity.LinkGetStatusParams {
	return entity.LinkGetStatusParams{Links: reqDTO.Links}
}

func convertEntityToDTOLinksGetStatus(res entity.LinkGetStatusResult) dto.LinksGetStatusResponse {
	return dto.LinksGetStatusResponse{
		Links:    convertEntityToDTOLinkStates(res.LinkStates),
		LinksNum: res.LinkNum,
	}
}

func convertEntityToDTOLinkStates(states []entity.LinkState) dto.LinksStatusResponse {
	result := make([]dto.LinkStatusResponse, 0, len(states))
	for _, state := range states {
		result = append(result, convertEntityToDTOLinkState(state))
	}
	return result
}

func convertEntityToDTOLinkState(state entity.LinkState) dto.LinkStatusResponse {
	return dto.LinkStatusResponse{
		Address: state.Link,
		Status:  convertStatusToText(state.IsAvailable),
	}
}

func convertStatusToText(isAvailable bool) string {
	if isAvailable {
		return consts.Available
	}
	return consts.NotAvailable
}
