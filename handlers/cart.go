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

// @Summary Get User's cart
// @Security JWT
// @Description Get list of cart
// @ID Getcart
// @Produce json
// @Success 200 {object} interfaces.carts "List of cart"
// @Failure 500 {object} interfaces.ErrorResponse "Error Response"
// @Router /cart [get]
func (h handler) GetCart(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	token, _ := middleware.VerifyToken(c.Request)
	claims, _ := token.Claims.(jwt.MapClaims)
	id_from_payload := claims["user_id"]
	user_id, _ := strconv.Atoi(fmt.Sprintf("%v", id_from_payload))

	cartResponse := []interfaces.CartResponse{}

	if result := h.DB.Table("products").Select("product_id, products.name, products.price, products.category, products.description,carts.amount").
		Joins("JOIN carts on carts.customer_id = ? AND carts.product_id = products.id",user_id).
		Scan(&cartResponse); errors.Is(result.Error, gorm.ErrRecordNotFound) {
			var response = &interfaces.MessageResponse{Message: "error while getting carts",}
			c.JSON(http.StatusInternalServerError, &interfaces.ErrorResponse{
				Error : *response,
			})
			return
		}

	
	var carts = &interfaces.Carts{Carts : cartResponse,}

	c.JSON(http.StatusOK, carts)
}

// @Summary Add cart
// @Security JWT
// @Description Add product to user's cart
// @ID Addcart
// @Produce json
// @Param product_id formData int true "Product ID"
// @Param amount formData int true "Amount of Product"
// @Success 200 {object} interfaces.SuccessResponse "Success"
// @Failure 400 {object} interfaces.ErrorResponse "Error Response"
// @Failure 500 {object} interfaces.ErrorResponse "Error Response"
// @Router /cart [post]
func (h handler) AddCart(c *gin.Context) {
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

	product := &interfaces.Product{}
	// get amount of product  and check whether amount cart request valid (<= amount of product) or not
	if result := h.DB.Where("id = ?", productID).First(&product); errors.Is(result.Error, gorm.ErrRecordNotFound) {
		var response = &interfaces.MessageResponse{Message: "Product not found in cart",}
		c.JSON(http.StatusBadRequest, &interfaces.ErrorResponse{
			Error : *response,
		})
		return
	}

	// check whether amount of product > requested amount
	if (amount > product.Amount){
		var response = &interfaces.MessageResponse{Message: "Invalid product amount",}
		c.JSON(http.StatusBadRequest, &interfaces.ErrorResponse{
			Error : *response,
		})
		return
	}

	// check whether the product is already in the cart or not
	cart := &interfaces.Cart{}
	if result := h.DB.Where("product_id = ? and customer_id = ?", productID, userID).First(&cart); errors.Is(result.Error, gorm.ErrRecordNotFound) {
		cart = &interfaces.Cart{
			CustomerID	: userID,
			ProductID	: productID,
			Amount		: amount,
		}
	
		if result := h.DB.Create(&cart); result.Error != nil {
			var response = &interfaces.MessageResponse{Message: "Error while adding carts",}
			c.JSON(http.StatusInternalServerError, &interfaces.ErrorResponse{
				Error : *response,
			})
			return
		}
	} else{
		if (cart.Amount + amount > product.Amount){
			var response = &interfaces.MessageResponse{Message: "Invalid product amount",}
			c.JSON(http.StatusBadRequest, &interfaces.ErrorResponse{
				Error : *response,
			})
			return
		}

		cart.Amount += amount

		if save := h.DB.Save(&cart); save.Error != nil {
			var response = &interfaces.MessageResponse{Message: "Failed update cart item amount",}
			c.JSON(http.StatusInternalServerError, &interfaces.ErrorResponse{
				Error : *response,
			})
			return
		}

	}

	var response = &interfaces.MessageResponse{Message: "Successed add product to cart",}
	c.JSON(http.StatusOK, &interfaces.SuccessResponse{
		Response : *response,
	})
}

// @Summary Delete cart Item
// @Security JWT
// @ID Deletecart
// @Consume application/x-www-form-urlencoded
// @Produce json
// @Param product_id formData int true "Product ID"
// @Param amount formData int true "Amount"
// @Success 200 {object} interfaces.SuccessResponse "Deleted cart Response"
// @Failure 400 {object} interfaces.ErrorResponse "Error Response"
// @Failure 500 {object} interfaces.ErrorResponse "Error Response"
// @Router /cart [delete]
func (h handler) DeleteCartItem(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	token, _ := middleware.VerifyToken(c.Request)

	claims, _ := token.Claims.(jwt.MapClaims)
	userID64, _ := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	productID64, _ := strconv.ParseUint(c.Request.FormValue("product_id"), 10, 64)
	amount, _ := strconv.Atoi(c.Request.FormValue("amount"))

	userID := uint(userID64)
	productID := uint(productID64)

	cart := &interfaces.Cart{}

	if result := h.DB.Where("product_id = ? and customer_id = ?", productID, userID).First(&cart); errors.Is(result.Error, gorm.ErrRecordNotFound) {
		var response = &interfaces.MessageResponse{Message: "Product not found in cart",}
		c.JSON(http.StatusBadRequest, &interfaces.ErrorResponse{
			Error : *response,
		})
		return
	}

	if (cart.Amount > amount) {
		// Reduce amount of product in cart
		cart.Amount -= amount

		if save := h.DB.Save(&cart); save.Error != nil {
			var response = &interfaces.MessageResponse{Message: "Failed update cart item amount",}
			c.JSON(http.StatusInternalServerError, &interfaces.ErrorResponse{
				Error : *response,
			})
			return
		}
	} else if (cart.Amount == amount){
		// delete item from cart
		if result := h.DB.Where("product_id = ? and customer_id = ?",productID, userID).Unscoped().Delete(&interfaces.Cart{}); result.Error != nil {
			var response = &interfaces.MessageResponse{Message: "Failed delete cart item",}
			c.JSON(http.StatusInternalServerError, &interfaces.ErrorResponse{
				Error : *response,
			})
			return
		}	
	} else{
		// requested item amount to be deleted > cart item amount
		var response = &interfaces.MessageResponse{Message: "Invalid product quantity",}
		c.JSON(http.StatusBadRequest, &interfaces.ErrorResponse{
			Error : *response,
		})
		return
	}

	var response = &interfaces.MessageResponse{Message: "Successfully delete cart item from user's cart",}
	c.JSON(http.StatusOK, &interfaces.SuccessResponse{
		Response : *response,
	})

}
