//Handlers for different HTTP routes and their corresponding business logic.

package authHandler

import (
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/naman2607/netflixClone/services"
)

type SignupRequestBody struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SigninRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupResponse struct {
	Message string `json:"message"`
}

func validateEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

func validatePassword(password string) bool {
	passLength := len(password)
	if passLength == 0 || passLength <= 5 {
		return false
	}
	return true
}

func Signup(c *gin.Context) {

	var reqBody SignupRequestBody
	if err := c.BindJSON(&reqBody); err != nil {
		log.Println("failed to process request")
		return
	}

	validCred := validateEmail(reqBody.Email) && validatePassword(reqBody.Password)

	if !validCred {
		signupResponse := SignupResponse{Message: "Invalid Credentials"}
		c.JSON(http.StatusBadRequest, signupResponse)
		return
	}

	response := services.SignupService(reqBody.Username, reqBody.Email, reqBody.Password)
	signupResponse := SignupResponse{Message: response.Message}
	c.JSON(response.StatusCode, signupResponse)
}

func Signin(c *gin.Context) {
	var requestBody SigninRequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		log.Println("Failed to process login")
		return
	}
	services.Signin(requestBody.Email, requestBody.Password)
}
