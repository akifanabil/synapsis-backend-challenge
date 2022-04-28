package handlers

import (
	"github.com/akifanabil/synapsis-backend-challenge/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get Product List
// @Description Get list of product
// @ID GetProducts
// @Produce json
// @Param category path string true "Category"
// @Success 200 {object} interfaces.Products "Product List"
// @Failure 500 {object} interfaces.ErrorResponse "Error Response"
// @Router /product/{category} [get]
func (h handler) GetProducts(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	category := c.Param("category")

	var result []interfaces.ProductResponse
	if res := h.DB.Table("products").Select("id, name, category, description, amount, price").Where("category = ?",category).Scan(&result); res.Error != nil {
		var response = &interfaces.MessageResponse{Message: "error while getting products",}
		c.JSON(http.StatusInternalServerError, &interfaces.ErrorResponse{
			Error : *response,
		})
		return
	}

	var products = &interfaces.Products{Products : result,}

	// Setup response
	c.JSON(http.StatusOK, products)
}


