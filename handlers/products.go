package handlers

import (
	// "github.com/akifanabil/synapsis-backend-challenge/helpers"
	"github.com/akifanabil/synapsis-backend-challenge/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/golang-jwt/jwt/v4"
)

// @Summary Get Product List
// @Description Get list of product
// @ID GetProducts
// @Produce json
// @Success 200 {object} interfaces.Products "Product List"
// @Failure 500 {object} interfaces.ErrorResponse "Error Response"
// @Router /tour [get]
func (h handler) GetProducts(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	var result []interfaces.ProductResponse
	if res := h.DB.Table("products").Select("id", "name", "description", "amount", "price").Scan(&result); res.Error != nil {
		var response = &interfaces.MessageResponse{Message: "error while getting products",}
		c.JSON(http.StatusInternalServerError, &interfaces.ErrorResponse{
			Error : *response,
		})
		return
	}

	var products = &interfaces.Products{Products : result}

	// Setup response
	c.JSON(http.StatusOK, products)
}