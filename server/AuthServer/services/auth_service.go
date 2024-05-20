//Business logic layer containing services for authentication and other related operations.

package services

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	authDao "github.com/naman2607/netflixClone/dao"
	user "github.com/naman2607/netflixClone/models"
)

type ServiceResponse struct {
	Message    string
	StatusCode int
	JwtToken   string
}

type CustomClaim struct {
	Email string `json:"email"`
	jwt.StandardClaims
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
	// claims := CustomClaim{
	// 	email,
	// 	jwt.StandardClaims{
	// 		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	// 		Issuer:    "netflix",
	// 	},
	// }
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// conf := config.GetInstance()
	// secret := []byte(conf.GetSecretKey())
	// tokenString, _ := token.SignedString(secret)
	// log.Println("jwt token : ", tokenString)
	// response.JwtToken = tokenString
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
	log.Println("Signin service called", email, password)
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

	authDao.GetUser(email)
	hashedPassword := getHashedPassword(password)
	log.Println(err, hashedPassword)
	return &response
}
