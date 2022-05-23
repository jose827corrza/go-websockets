package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jose827corrza/go-websockets/models"
	"github.com/jose827corrza/go-websockets/repository"
	"github.com/jose827corrza/go-websockets/server"
	"github.com/jose827corrza/go-websockets/tokenization"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	HASH_COST = 9
)

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	UserName string `json:"username"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type SignUpResponse struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	UserName string `json:"username"`
}

func SignUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request = SignUpRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), HASH_COST)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		id, err := ksuid.NewRandom()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var user = models.User{
			Email:    request.Email,
			Password: string(hashedPassword),
			UserName: request.UserName,
			Id:       id.String(),
		}
		err = repository.InsertUser(r.Context(), &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SignUpResponse{
			Id:       user.Id,
			Email:    user.Email,
			UserName: user.UserName,
		})
	}
}

func LoginHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request = LoginRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		user, err := repository.GetUserByEmail(r.Context(), request.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if user == nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		// claims := models.AppClaims{
		// 	UserId: user.Id,
		// 	StandardClaims: jwt.StandardClaims{
		// 		ExpiresAt: time.Now().Add(2 * time.Hour * 2).Unix(),
		// 	},
		// }
		// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// tokenString, err := token.SignedString([]byte(s.Config().JwtSecret))
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		tokenString, err := tokenization.NewAuthorization().SignToken(s.Config().JwtSecret, user.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResponse{
			Token: tokenString,
		})
	}
}

func MeHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
		// token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 	return []byte(s.Config().JwtSecret), nil
		// })
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusUnauthorized)
		// 	return
		// }
		// if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
		// 	user, err := repository.GetUserById(r.Context(), claims.UserId)
		// 	if err != nil {
		// 		http.Error(w, err.Error(), http.StatusInternalServerError)
		// 		return
		// 	}
		claims, err := tokenization.NewAuthorization().ParseAndVerifyToken(s.Config().JwtSecret, tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user, err := repository.GetUserById(r.Context(), claims.UserId)
		if user.Id == "" {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
		// } else {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
	}
}

////////////
func GettingUserById(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := repository.GetUserById(r.Context(), r.URL.Path)
		fmt.Print(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&SignUpResponse{
			Id:       user.Id,
			Email:    user.Email,
			UserName: user.UserName,
		})
	}
}
