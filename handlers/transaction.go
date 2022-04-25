package handlers

import (
	"fmt"
	"errors"
	"strconv"

	"github.com/akifanabil/synapsis-backend-challenge/interfaces"
	"github.com/akifanabil/synapsis-backend-challenge/middleware"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"

)

// @Summary Delete cart Item
// @Security JWT
// @ID Deletecart
// @Consume application/x-www-form-urlencoded
// @Produce json
// @Param product_id formData int true "Product ID"
// @Param amount formData int true "Amount"
// @Success 200 {object} interfaces.SuccessResponse "Success Response"
// @Failure 400 {object} interfaces.ErrorResponse "Error Response"
// @Failure 500 {object} interfaces.ErrorResponse "Error Response"
// @Router /cart [delete]
func (h handler) Checkout(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	token, _ := middleware.VerifyToken(c.Request)

	claims, _ := token.Claims.(jwt.MapClaims)
	userID64, _ := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	productID64, _ := strconv.ParseUint(c.Request.FormValue("product_id"), 10, 64)
	amount, _ := strconv.Atoi(c.Request.FormValue("amount"))

	userID := uint(userID64)
	productID := uint(productID64)

	// get amount of product  and check whether amount cart request valid (<= amount of product) or not
	product := &interfaces.Product{}
	if result := h.DB.Where("id = ?", productID).First(&product); errors.Is(result.Error, gorm.ErrRecordNotFound) {
		var response = &interfaces.MessageResponse{Message: "Product not found in cart",}
		c.JSON(http.StatusBadRequest, &interfaces.ErrorResponse{
			Error : *response,
		})
		return
	}

	if (product.Amount < amount) {
		// Number of product available < requested amount to buy
		var response = &interfaces.MessageResponse{Message: "Product quantity is not enough",}
		c.JSON(http.StatusInternalServerError, &interfaces.ErrorResponse{
			Error : *response,
		})
		return	
	}

	// check whether the user balance is sufficient or not
	totals := amount * product.Price
	user := &interfaces.Customer{}
	if result := h.DB.Where("id = ?", userID).First(&user); errors.Is(result.Error, gorm.ErrRecordNotFound) {
		var response = &interfaces.MessageResponse{Message: "User not registered",}
		c.JSON(http.StatusBadRequest, &interfaces.ErrorResponse{
			Error : *response,
		})
		return
	}

	if (user.Balance < totals){
		var response = &interfaces.MessageResponse{Message: "Balance is not enough",}
		c.JSON(http.StatusBadRequest, &interfaces.ErrorResponse{
			Error : *response,
		})
		return
	}

	// Reduce user's balance
	user.Balance -= totals
	if save := h.DB.Save(&user); save.Error != nil {
		var response = &interfaces.MessageResponse{Message: "Failed in updating customer's balance",}
		c.JSON(http.StatusInternalServerError, &interfaces.ErrorResponse{
			Error : *response,
		})
		return
	}

	// Reduce product stock
	product.Amount -= amount
	if save := h.DB.Save(&product); save.Error != nil {
		var response = &interfaces.MessageResponse{Message: "Failed in updating product's amount",}
		c.JSON(http.StatusInternalServerError, &interfaces.ErrorResponse{
			Error : *response,
		})
		return
	}

	// Adjust the number of products in the cart
	cart := &interfaces.Cart{}
	if result := h.DB.Where("product_id = ? and customer_id = ?", productID, userID).First(&cart); errors.Is(result.Error, gorm.ErrRecordNotFound) {
	} else if (cart.Amount > amount) {
		// Reduce amount of product in cart
		cart.Amount -= amount
		h.DB.Save(&cart);
	} else{
		// cart.Amount <= amount -> delete item from cart
		h.DB.Where("product_id = ? and customer_id = ?",productID, userID).Unscoped().Delete(&interfaces.Cart{});
	}

	// Add to transactions table
	transaction := &interfaces.Transaction{
		CustomerID 	: userID,
		ProductID 	: productID,
		Amount		: amount,
		Totals  	: totals,	
	}
	h.DB.Create(&transaction)

	var response = &interfaces.MessageResponse{Message: "Successfully buy item",}
	c.JSON(http.StatusOK, &interfaces.SuccessResponse{
		Response : *response,
	})
}