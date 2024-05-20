//Business logic layer containing services for authentication and other related operations.

package services

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	authDao "github.com/naman2607/netflixClone/dao"
	user "github.com/naman2607/netflixClone/models"
	jwtHelper "github.com/naman2607/netflixClone/utils"
)

type ServiceResponse struct {
	Message    string
	StatusCode int
	JwtToken   string
}

func getHashedPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	hashedBytes := hash.Sum(nil)
	hashedString := hex.EncodeToString(hashedBytes)
	return hashedString
}

func SignupService(username string, email string, password string) *ServiceResponse {
	response := ServiceResponse{}
	userAlreadyExists, err := authDao.CheckIfUserExists(email)
	if err != nil {
		log.Println("failed to check the user in the db ", err)
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Sign up failed"
		return &response
	}
	if userAlreadyExists {
		response.StatusCode = http.StatusAccepted
		response.Message = "User Already Exists"
		return &response
	}

	var user user.UserBasicDetails
	user.Username = username
	user.Email = email

	user.Password = getHashedPassword(password)
	if err := authDao.InsertUserIntoDatabase(user); err != nil {
		log.Println("failed to insert the user in the db ", err)
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Sign up failed"
		return &response
	}
	response.StatusCode = http.StatusOK
	response.Message = "Successfully Registered"
	return &response
}

func Signin(email string, password string) *ServiceResponse {
	userExists, err := authDao.CheckIfUserExists(email)
	var response ServiceResponse

	if err != nil {
		log.Println("failed to check the user in the db ", err)
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Sign up failed"
		return &response
	}

	if !userExists {
		response.StatusCode = http.StatusOK
		response.Message = "User does not Exist"
		return &response
	}

	user, err := authDao.GetUser(email)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		return &response
	}
	hashedPassword := getHashedPassword(password)
	if user.Password != hashedPassword {
		response.StatusCode = http.StatusUnauthorized
		response.Message = "Passwod is incorrect"
		return &response
	}
	claims := jwtHelper.CustomClaim{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "netflix",
		},
	}
	token := jwtHelper.CreateNewToken(claims)
	response.Message = "Login Successful"
	response.JwtToken = token
	return &response
}
