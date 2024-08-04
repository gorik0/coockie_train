package main

import (
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

func main() {

	gi := gin.New()
	gi.Use(cook)

	// ---------------- HANDLERS ----------------

	//						:::GET COOKIES::::
	gi.GET("/getCookie/:cook", func(c *gin.Context) {
		cookName := c.Param("cook")
		err, coo := Read(c.Request, cookName)
		if err != nil {
			c.Error(err)

		}
		c.JSON(200, gin.H{
			"cookies ": coo,
		})
	})

	//						:::SET COOKIES::::
	gi.GET("/setCook/:cook", func(c *gin.Context) {
		coo := c.Param("cook")
		cookieToSet := http.Cookie{
			Name:  "coo",
			Value: coo,
			Path:  "/",
		}

		writer := c.Writer
		http.SetCookie(writer, &cookieToSet)
		if err := Writer(writer, cookieToSet); err != nil {
			c.Error(err)
			return
		}
		c.JSON(200, gin.H{
			"acitons": "setting cookies of",
		})

	})

	// ---------------- HANDLERS ----------------

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	err = http.Serve(lis, gi)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

}

var ErrGetCookie = errors.New("Errors to get cookie")

var ErrDecoedCookie = errors.New("Errors to decode cookie")

func Read(r *http.Request, cookieName string) (error, string) {

	getCookie, err := r.Cookie(cookieName)
	if err != nil {
		return ErrGetCookie, ""

	}
	value, err := base64.URLEncoding.DecodeString(getCookie.Value)
	if err != nil {
		return ErrDecoedCookie, ""
	}

	return nil, string(value)
}

var ErrValueTooLong = errors.New("Cookie value length is too long")

func Writer(writer gin.ResponseWriter, coo http.Cookie) error {

	coo.Value = base64.URLEncoding.EncodeToString([]byte(coo.Value))

	if len(coo.String()) > 4096 {
		return ErrValueTooLong
	}
	http.SetCookie(writer, &coo)
	return nil
}

func cook(context *gin.Context) {

	cookExtract := context.Request.Header.Get("Cookie")

	println("cookExtract:", cookExtract)
	context.Next()

}
