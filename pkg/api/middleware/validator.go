package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// MyClaims represents custom claims for JWT
type MyClaims struct {
	UserID uint   `json:"userID"`
	Phone  string `json:"phone"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// we define the expiration time of JWT, taking 2 hours
const TokenExpireDuration = time.Hour * 2

func GenToken(userID uint, phone string, c *gin.Context) (string, error) {
	MySecretKEY := os.Getenv("MySecretKEY")
	claims := MyClaims{
		UserID: userID, // Custom field
		Phone:  phone,
		Role:   "user", // Set the user's role
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // Expiration time
			Issuer:    "my-project",                               // Issuer
		},
	}
	// Creates a signed object using the specified signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(MySecretKEY))
	if err != nil {
		return "", err
	}

	// Calculate the number of seconds until the token expires
	expireSeconds := int(TokenExpireDuration.Seconds())

	// Set the cookie with the specified name and expiration
	c.SetCookie("Authorize", tokenString, expireSeconds, "", "", false, true)

	return tokenString, nil
}

func ValidateCookie(c *gin.Context) {
	tokenString, err := c.Cookie("Authorize")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"StatusCode": 401,
			"msg":        "Unauthorized User",
		})
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if err != nil {
			fmt.Println("error in parsing", err)
		}
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method:%v", token.Header["alg"])
		}
		return []byte(os.Getenv("MySecretKEY")), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"StatusCode": 401,
				"msg":        "Jwt session expired",
			})

			return
		}
		fmt.Println("ClAIMS", token.Claims.(jwt.MapClaims))
		c.Set("userID", fmt.Sprint(claims["userID"]))
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Statuscode": 401,
			"Msg":        "Invalid claims",
		})
		return
	}
}

func DeleteCookie(c *gin.Context) error {
	c.SetCookie("Authorize", "", 0, "", "", false, true)
	fmt.Println("Cookie deleted")
	return nil
}
