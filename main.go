package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func ResponseJSON(c *gin.Context, httpCode, errCode int, msg string, data interface{}) {
	c.JSON(httpCode, Response{
		Code: errCode,
		Msg:  msg,
		Data: data,
	})
}

func BindAndValid(c *gin.Context, form interface{}) (int, int) {
	err := c.Bind(form)
	if err != nil {
		return http.StatusBadRequest, 400
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return http.StatusInternalServerError, 500
	}
	if !check {
		return http.StatusBadRequest, 400
	}

	return http.StatusOK, 200
}

func printRequest(c *gin.Context) {
	fmt.Println("Request:", c.Request.Method, c.Request.URL)
}

func printRequestHeader(c *gin.Context) {
	fmt.Println("Request Header:")
	for k, v := range c.Request.Header {
		fmt.Println(k, v)
	}
}

func printRequestBody(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println("---body/--- \r\n " + string(body))
}

func runHttpServer() {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	RegisterEmqxRouter(r)
	RegisterRabbitmqRouter(r)

	_ = r.Run(":7100")
}

func main() {
	runHttpServer()
}
