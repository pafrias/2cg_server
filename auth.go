package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var keyString = []byte("maybeitshouldbeshorter")

func (s *server) checkAuth(authLevel int, next http.HandlerFunc) http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {
		fmt.Println(req.Header.Get("authorization"))

		fmt.Println("Required Auth level ", authLevel)

		next(res, req)
	}
}

func (s *server) signin() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		req.ParseForm()

		username := req.Form.Get("username")
		password := req.Form.Get("password")

		fmt.Printf("user: %v, pw: %v", username, password)

		if err := s.DB.Ping(); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(err.Error()))
			fmt.Println(password)
			return
		}

		//query := "SELECT salt, hash, user_type FROM Users WHERE name = ? LIMIT 1"

		// var salt, hash string
		var userType int

		//err := s.db.QueryRow(query, username).Scan(&salt, &hash, &userType)

		userType = 3
		username = "diablo"

		// if err != nil {
		// 	res.WriteHeader(http.StatusBadRequest)
		// 	res.Write([]byte("Username does not exist"))
		// }

		// test if hash(password + salt) = hash

		// make token, send to front end

		tokenString, err := generateJWT(username, userType)

		if err != nil {
			fmt.Println(err.Error())
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(err.Error()))
			return
		}

		res.Header().Set("authorization", tokenString)
		res.WriteHeader(200)

	}
}

func generateJWT(user string, userType int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["user"] = user
	claims["user_type"] = userType
	claims["exp"] = time.Now().Add(time.Hour * 6).Unix()
	claims["iat"] = time.Now().Unix()

	return token.SignedString(keyString)

}
