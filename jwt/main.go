package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
)

type post struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Tag   string `json:"tag"`
	URL   string `json:"url"`
}

var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// headerのセット
	token := jwt.New(jwt.SigningMethodHS256)

	// claimsのセット
	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = true
	claims["sub"] = "19950510"
	claims["name"] = "shuya"
	claims["iat"] = time.Now()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// 電子署名
	tokenString, _ := token.SignedString([]byte(os.Getenv("SIGNINGKEY")))

	// JWTを返却
	w.Write([]byte(tokenString))
})

// JwtMiddleware check token
var JwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SIGNINGKEY")), nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

func main() {
	r := mux.NewRouter()
	// localhost:8080/publicでpublicハンドラーを実行
	r.Handle("/public", public)

	r.Handle("/private", JwtMiddleware.Handler(private))
	r.Handle("/auth", GetTokenHandler)
	//サーバー起動
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("ListenAndServe:", nil)
	}
}

var public = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	post := &post{
		Id:    "1",
		Title: "公開しているやつ",
		Tag:   "公開",
		URL:   "https://nwiizo.github.io/",
	}
	json.NewEncoder(w).Encode(post)
})

var private = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	post := &post{
		Id:    "2",
		Title: "秘密にしているやつ",
		Tag:   "秘密",
		URL:   "https://twitter.com/nwiizo",
	}
	json.NewEncoder(w).Encode(post)
})
