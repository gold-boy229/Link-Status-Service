package handlers

import (
	"Link-Status-Service/internal/consts"
	"Link-Status-Service/internal/dto"
	"Link-Status-Service/internal/entity"
	"net/http"

	"github.com/labstack/echo"
)

func (h linkHandler) GetStatus(c echo.Context) error {
	var reqDTO dto.LinksGetStatus_Request
	if err := c.Bind(&reqDTO); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewError(consts.ERROR_CODE_BAD_REQUEST, err.Error()))
	}
	if err := c.Validate(reqDTO); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewError(consts.ERROR_CODE_BAD_REQUEST, err.Error()))
	}

	params := convertDTOToEntity_LinksGetStatus(reqDTO)
	result, err := h.LinkService.GetStatus(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			dto.NewError(consts.ERROR_CODE_INTERNAL_SERVER_ERROR, err.Error()))
	}

	return c.JSON(http.StatusOK, convertEntityToDTO_LinksGetStatus(result))
}

func convertDTOToEntity_LinksGetStatus(reqDTO dto.LinksGetStatus_Request) entity.LinkGetStatus_Params {
	return entity.LinkGetStatus_Params{Links: reqDTO.Links}
}

func convertEntityToDTO_LinksGetStatus(res entity.LinkGetStatus_Result) dto.LinksGetStatus_Response {
	return dto.LinksGetStatus_Response{
		Links:    convertEntityToDTO_LinkStates(res.LinkStates),
		LinksNum: res.LinkNum,
	}
}

func convertEntityToDTO_LinkStates(states []entity.LinkState) dto.LinksStatus_Response {
	result := make([]dto.LinkStatus_Response, 0, len(states))
	for _, state := range states {
		result = append(result, convertEntityToDTO_LinkState(state))
	}
	return result
}

func convertEntityToDTO_LinkState(state entity.LinkState) dto.LinkStatus_Response {
	return dto.LinkStatus_Response{
		Address: state.Link,
		Status:  convertStatusToText(state.IsAvailable),
	}
}

func convertStatusToText(isAvailable bool) string {
	if isAvailable {
		return consts.AVAILABLE
	}
	return consts.NOT_AVAILABLE
}
