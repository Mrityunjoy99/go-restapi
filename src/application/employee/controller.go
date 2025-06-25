package employee

import (
	"net/http"

	"github.com/Mrityunjoy99/sample-go/src/common/constant"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type controller struct {
	service Service
}

type Controller interface {
	GetEmployeeById(c *gin.Context)
}

func NewController(service Service) Controller {
	return &controller{
		service: service,
	}
}
func (ctrl *controller) GetEmployeeById(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invelid id",
		})

		return
	}

	employee, gerr := ctrl.service.GetEmployeeById(id)
	if gerr != nil {
		if gerr.GetCode() == constant.ErrorCodeResourceNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "record not found",
			})

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gerr.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, employee)
}
