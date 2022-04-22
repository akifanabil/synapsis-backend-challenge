package handlers

import (
	"fmt"
	"errors"
	"github.com/akifanabil/synapsis-backend-challenge/helpers"
	"github.com/akifanabil/synapsis-backend-challenge/interfaces"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
)


// @Summary Provides a JSON Web Token for Authenticated Customers
// @Description Authenticates a user and provides a JWT to Authorize API calls
// @ID Authentication
// @Consume application/x-www-form-urlencoded
// @Produce json
// @Param email formData string true "Email"
// @Param password formData string true "Password"
// @Success 200 {object} interfaces.AuthResponse "Login Response"
// @Failure 401 {object} interfaces.ErrorResponse "Error Response"
// @Router /login [post]
func (h handler) Login(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")

	user := &interfaces.Customer{}

	if result := h.DB.Where("email = ?", email).First(&user); errors.Is(result.Error, gorm.ErrRecordNotFound) {
		var response = &interfaces.MessageResponse{Message: "wrong email or password",}
		c.JSON(http.StatusUnauthorized, &interfaces.ErrorResponse{
			Error : *response,
		})
		return
	}

	passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
		var response = &interfaces.MessageResponse{Message: "wrong email or password",}
		c.JSON(http.StatusUnauthorized, &interfaces.ErrorResponse{
			Error : *response,
		})
		return
	}

	// Setup response
	customerResponse := &interfaces.CustomerResponse{
		ID:     user.ID,
		Name:	user.Name,
		Email:  user.Email,
		Gender: user.Gender,
	}

	// Sign token
	tokenContent := jwt.MapClaims{
		"user_id": customerResponse.ID,
		"role": "user"    ,
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte("TokenPassword"))
	helpers.HandleError(err)

	// Prepare response
	var authData = &interfaces.AuthResponse{
		Message: "Ok",
		Jwt: token,
		ExpiresIn_hour: 720,
		Data: *customerResponse,
	}

	c.JSON(http.StatusOK, authData)
}


// @Summary Register New User
// @Description Register New User and Generate JWT Token
// @ID Registration
// @Consume application/x-www-form-urlencoded
// @Produce json
// @Param email formData string true "Email"
// @Param password formData string true "Password"
// @Param name formData string true "Name"
// @Param gender formData string true "Gender (w/m)"
// @Success 200 {object} interfaces.AuthResponse "Login Response"
// @Failure 400 {object} interfaces.ErrorResponse "Error Response"
// @Router /register [post]
func (h handler) Register(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")
	name := c.Request.FormValue("name")
	gender := c.Request.FormValue("gender")

	user := &interfaces.Customer{
		Email:    email,
		Password: helpers.HashAndSalt([]byte(password)),
		Name:     name,
		Gender:   gender,
	}

	var email_unique = false;

	// check if email already used
	if result := h.DB.Where("email = ?", email).First(&user); errors.Is(result.Error, gorm.ErrRecordNotFound) {
		email_unique = true;
	}

	if (!email_unique){
		var response = &interfaces.MessageResponse{Message: "email already used",}
		c.JSON(http.StatusBadRequest, &interfaces.ErrorResponse{
			Error : *response,
		})
		return
	}

	if result := h.DB.Create(&user); result.Error != nil {
		fmt.Println(result.Error)
		var response = &interfaces.MessageResponse{Message: "failed creating user",}
		c.JSON(http.StatusBadRequest, &interfaces.ErrorResponse{
			Error : *response,
		})
		return
	}

	// Setup response
	customerResponse := &interfaces.CustomerResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Gender: user.Gender,
	}

	// Sign token
	tokenContent := jwt.MapClaims{
		"user_id": customerResponse.ID,
		"role":    "user",
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte("TokenPassword"))
	helpers.HandleError(err)

	// Prepare response
	var authData = &interfaces.AuthResponse{
		Message: "Ok",
		Jwt: token,
		ExpiresIn_hour: 720,
		Data: *customerResponse,
	}

	c.JSON(http.StatusOK, authData)
}
