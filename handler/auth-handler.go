package handler

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

const USERNAME = "belajar-golang"
const PASSWORD = "Password123"
const SECUREPASSWORD = "$2a$14$t4CQ67jNbN3EdaBKfa3Gpe5UMqtXB45IFK8MQrAK.d43ftW4UbMp2"

func SecureMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if strings.HasPrefix(r.URL.Path, "/users") ||
		// 	strings.HasPrefix(r.URL.Path, "/orders") {
		// 	next.ServeHTTP(w, r)
		// }
		username, password, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`Authentication Failed`))
			return
		}

		if username != USERNAME {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`Invalid Username or Password`))
			return
		}

		err := bcrypt.CompareHashAndPassword([]byte(SECUREPASSWORD), []byte(password))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`Invalid Username or Password`))
			return
		}

		next.ServeHTTP(w, r)
	})
}
