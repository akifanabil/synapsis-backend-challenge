package handlers

import (
	"fmt"
	"errors"
	"strconv"

	// "github.com/akifanabil/synapsis-backend-challenge/helpers"
	"github.com/akifanabil/synapsis-backend-challenge/interfaces"
	"github.com/akifanabil/synapsis-backend-challenge/middleware"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"

)

// @Summary Get User's Chart
// @Security JWT
// @Description Get list favorite tour area of a user
// @ID GetChart
// @Produce json
// @Success 200 {object} interfaces.Charts "List of Chart"
// @Failure 204 {object} interfaces.ErrorResponse "Error Response"
// @Router /chart [get]
func (h handler) GetChart(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	token, _ := middleware.VerifyToken(c.Request)
	claims, _ := token.Claims.(jwt.MapClaims)
	id_from_payload := claims["user_id"]
	user_id, _ := strconv.Atoi(fmt.Sprintf("%v", id_from_payload))

	chartResponse := []interfaces.ChartResponse{}

	if result := h.DB.Table("products").Select("charts.id as ChartID, products.id as ProductID, products.name, products.price, products.category, products.description,charts.amount").
		Joins("JOIN charts on charts.customer_id = ? AND charts.product_id = products.id",user_id).
		Scan(&chartResponse); errors.Is(result.Error, gorm.ErrRecordNotFound) {
			var response = &interfaces.MessageResponse{Message: "error while getting charts",}
			c.JSON(http.StatusInternalServerError, &interfaces.ErrorResponse{
				Error : *response,
			})
			return
		}

	
	var charts = &interfaces.Charts{Charts : chartResponse,}

	c.JSON(http.StatusOK, charts)
}

// @Summary Add Chart
// @Security JWT
// @Description Add product to user's chart
// @ID AddChart
// @Produce json
// @Param product_id formData int true "Product ID"
// @Param amount formData int true "Amount of Product"
// @Success 200 {object} interfaces.SuccessResponse "Success"
// @Failure 400 {object} interfaces.ErrorResponse "Error Response"
// @Router /chart [post]
func (h handler) AddChart(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	token, _ := middleware.VerifyToken(c.Request)

	claims, _ := token.Claims.(jwt.MapClaims)
	role_payload := claims["role"]
	role := fmt.Sprintf("%v", role_payload)

	// not user
	if role != "user" {
		var response = &interfaces.MessageResponse{Message: "Not authorized",}
		c.JSON(http.StatusInternalServerError, &interfaces.ErrorResponse{
			Error : *response,
		})
		return
	}

	userID64, _ := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	productID64, _ := strconv.ParseUint(c.Request.FormValue("product_id"), 10, 64)
	userID := uint(userID64)
	productID := uint(productID64)
	amount, _ := strconv.Atoi(c.Request.FormValue("amount"))

	// TODO : get amount of product id and check whether amount chart request valid (<= amount of product) or not
	
	chart := &interfaces.Chart{
		CustomerID	: userID,
		ProductID	: productID,
		Amount		: amount,
	}

	if result := h.DB.Create(&chart); result.Error != nil {
		var response = &interfaces.MessageResponse{Message: "Error while adding charts",}
		c.JSON(http.StatusInternalServerError, &interfaces.ErrorResponse{
			Error : *response,
		})
		return
	}

	var response = map[string]interface{}{"message": "Add product to chart success"}
	c.JSON(http.StatusOK, gin.H{
		"response": response,
	})
}