package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

const SECRETKEY = "My$3cr3TKey"

//Generate JWT token
func GenerateJWT(username string, age int) (string, error) {
	var mySigningKey = []byte(SECRETKEY)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["age"] = age
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Fatal(err.Error())
		return "", err
	}
	return tokenString, nil
}

func IsAuthorized(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if strings.HasPrefix(r.URL.Path, "/weather") {
			handler.ServeHTTP(w, r)
			return
		}

		if strings.HasSuffix(r.URL.Path, "/users/register") ||
			strings.HasSuffix(r.URL.Path, "/users/login") {
			handler.ServeHTTP(w, r)
			return
		}

		clientToken := r.Header.Get("Authorization")
		if clientToken == "" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("No Authorization Header Provided"))
			return
		}

		splitToken := strings.Split(clientToken, "Bearer ")
		if len(splitToken) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid Token"))
			return
		}

		accessToken := splitToken[1]
		if len(accessToken) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid Token"))
			return
		}

		var mySigningKey = []byte(SECRETKEY)

		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing token.")
			}
			return mySigningKey, nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Your Token has been expired!"))
			return
		}

		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			handler.ServeHTTP(w, r)
			return
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}
	})
}
