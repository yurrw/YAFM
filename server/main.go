package main

// TODO IMPLEMENTAR SSO
import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"yafm/pkg/router"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var IsLoggedIn = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey: []byte(jwtKey),
})

//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJOYW1lIjoiSm9uIFNub3ciLCJSb2xlIjoiQWRtaW4iLCJleHAiOjE1ODc1ODE1OTR9.5ULmgWiXZcAgZ5cHHP0Mp1sceWrF9a4M71U0boes2xk
//Seta rotas
func routes(e *echo.Echo) {
	// e.GET("/", Welcome())

	e.POST("/login", Login)
	e.GET("/is-loggedin", restricted, IsLoggedIn)
	e.GET("/is-admin", restricted, IsLoggedIn, isAdmin)
	e.POST("/rnw-token", RenewAccessToken, IsLoggedIn)
	e.POST("/logout", LogoutUser, IsLoggedIn)

}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	// map[Name:Jon Snow Role:Admin exp:1.587582519e+09]
	name := claims["Name"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func isAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		isAdmin := claims["admin"].(bool) // TODO MUDAR
		if isAdmin == false {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}

func TokenValid(c echo.Context) error {
	token, err := VerifyToken(c)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

// router.POST("/logout", isAuthenthicated(), Logout)
// funcaozinha pra autenticar se nao ta na blackalist
func isAuthenthicated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		err := TokenValid(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		c.Next()
	}
}

// curl -X POST -d 'username=jon' -d 'password=shhh!' localhost:7777/login
func Login(c echo.Context) error {
	// TODO : mudar para dados requisitados do bd
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username != "jon" || password != "shhh!" {
		return echo.ErrUnauthorized
	}
	tokenPair, _ := setUpToken()

	return c.JSON(http.StatusOK, tokenPair)
}

type Claims struct {
	Name string
	Role string
	// ExpiresAt int64
	jwt.StandardClaims
}

// Create the JWT key used to create the signature
// The signing string should be secret (a generated UUID works too)

var jwtKey = []byte("VAI_TOMAR_NO_CU_ALEX")

// CREATING TOKENs
func setUpToken() (tokenPair map[string]string, err error) {

	// stackoverflow.com/questions/27726066/jwt-refresh-token-flow
	//Claims
	expireTime := time.Now().Add(time.Hour * 72).Unix()

	// pwc = payload with claims
	pwc := &Claims{
		Name: "Jon Snow",
		Role: "Admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, pwc)
	// Create the JWT string
	t, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		// log.Fatalln(err)
		return nil, err
	}
	expireTimeRefresh := time.Now().Add(time.Hour * 24).Unix()
	pwcRefresh := &Claims{
		Name: "Jon Snow",
		Role: "Admin",
		StandardClaims: jwt.StandardClaims{
			Subject:   "1", // change this for
			ExpiresAt: expireTimeRefresh,
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, pwcRefresh)
	rt, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		// log.Fatalln(err)
		return nil, err

	}
	tokenPair = map[string]string{
		"access_token":  t,
		"refresh_token": rt,
	}
	return tokenPair, nil
	// claims = payload

}

var client *redis.Client

func initRedis() {
	//Initializing redis
	//TODO
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:63	79"
	}
	client = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

}

func main() {
	e := router.New()
	routes(e)
	initRedis()
	//Starting server
	setUpToken()
	if err := e.Start(":7777"); err != nil {
		log.Fatalln(err)
	}

}

func ExtractToken(c echo.Context) string {
	bearToken := c.Request().Header.Get("Authorization")
	// bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")

	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

// This is the api to refresh tokens
// Most of the code is taken from the jwt-go package's sample codes
// https://godoc.org/github.com/dgrijalva/jwt-go#example-Parse--Hmac
func RenewAccessToken(c echo.Context) error {
	// TODO : POR ISSO OUTSIDE
	type tokenReqBody struct {
		RefreshToken string `json:"refresh_token"`
	}
	// t := c.FormValue("refresh_token")
	// fmt.Println(t)

	tokenReq := tokenReqBody{}
	c.Bind(&tokenReq)
	fmt.Println("-------------------------------------------------------")
	fmt.Println("TOKEN")
	fmt.Println(tokenReq.RefreshToken)
	fmt.Println("-------------------------------------------------------")

	// Parse takes the token string and a function for looking up the key.
	// The latter is especially useful if you use multiple keys for your application.
	// The standard is to use 'kid' in the head of the token to identify
	// which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenReq.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return jwtKey, nil
	})
	fmt.Println("AQUI")
	// fmt.Println(token.Claims.(jwt.MapClaims))
	// fmt.Println(token.Valid)
	fmt.Println("AQUI")

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get the user record from database or
		// run through your business logic to verify if the user can log in
		if int(claims["sub"].(float64)) == 1 {
			fmt.Println("AQUI dentro")

			newTokenPair, err := setUpToken()
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, newTokenPair)
		}

		return echo.ErrUnauthorized
	}

	return err
}

//----------------------------------------------------------------------
type AccessDetails struct {
	token  string
	UserId uint64
}

func VerifyToken(c echo.Context) (*jwt.Token, error) {
	tokenString := ExtractToken(c)

	// verify reddis

	//

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// =============================
// REDIS =======================
// =============================

func FetchAuthReddis(reddisAuth AccessDetails) (uint64, error) {
	userid, err := client.Get(reddisAuth.token).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}
func SaveInReddis() {

	saveErr := CreateAuth(userId, ts)
	if saveErr != nil {
		c.JSON(http.StatusForbidden, saveErr.Error())
		return
	}

}
func CreateAuth(userid uint64, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := client.Set(td.AccessUuid, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := client.Set(td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func DeleteAuthReddis(token string) (int64, error) {
	deleted, err := client.Del(token).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
func BlackListIt(token string) (int64, error) {
	deleted, err := client.Del(token).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func LogoutUser(c echo.Context) {

	token := ExtractToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Authorization Token is required"})
		// c.Abort()
		log.Fatalln("Authorization token was not provided")

		return
	}

	extractedToken := strings.Split(token, "Bearer ")
	fmt.Println(extractedToken[1])
	// This method will add the token to the redis db
	_, err := BlackListIt(extractedToken[1])
	if err != nil {
		// c.Abort()
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": http.StatusAccepted, "message": "Done"})

}
