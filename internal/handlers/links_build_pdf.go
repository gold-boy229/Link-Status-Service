package handlers

import (
	"Link-Status-Service/internal/consts"
	"Link-Status-Service/internal/dto"
	"Link-Status-Service/internal/entity"
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func (h linkHandler) BuildPDF(c echo.Context) error {
	var reqDTO dto.LinkBuildPDF_Request
	if err := c.Bind(&reqDTO); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewError(consts.ERROR_CODE_BAD_REQUEST, err.Error()))
	}
	if err := c.Validate(reqDTO); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewError(consts.ERROR_CODE_BAD_REQUEST, err.Error()))
	}

	params := convertDTOToEntity_LinkBuildPDF(reqDTO)
	// call service to get link statuses of union of all linkSets
	result, err := h.LinkService.GetStatusesOfLinkSets(context.TODO(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			dto.NewError(consts.ERROR_CODE_INTERNAL_SERVER_ERROR, err.Error()))
	}

	linkStatuses := convertLinkStatesToLinkStatuses(result.LinkStates)
	pdf := h.PDFBuilder.BuildPDF(linkStatuses)

	c.Request().Header.Set("Content-Type", "application/pdf")
	// Заголовок Content-Disposition с параметром "attachment" заставит браузер скачать файл
	// Если хотите отобразить PDF прямо в браузере, используйте "inline"
	c.Request().Header.Set("Content-Disposition", "attachment; filename=products_report.pdf")

	// Выводим PDF в http.ResponseWriter. Используется метод Output(), который может принимать http.ResponseWriter
	err = pdf.Output(c.Response().Writer)
	if err != nil {
		err = fmt.Errorf("cannot generate PDF: %w", err)
		return c.JSON(http.StatusInternalServerError,
			dto.NewError(consts.ERROR_CODE_INTERNAL_SERVER_ERROR, err.Error()))
	}
	return nil
}

func convertDTOToEntity_LinkBuildPDF(reqDTO dto.LinkBuildPDF_Request) entity.LinkBuildPDS_Params {
	return entity.LinkBuildPDS_Params{
		LinkNums: reqDTO.LinkNums,
	}
}

func convertLinkStatesToLinkStatuses(linkStates []entity.LinkState) []entity.LinkStatus {
	result := make([]entity.LinkStatus, len(linkStates))
	for idx, linkState := range linkStates {
		result[idx] = entity.LinkStatus{
			Address: linkState.Link,
			Status:  convertStatusToText(linkState.IsAvailable),
		}
	}
	return result
}
