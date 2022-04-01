package main

import (
	"fmt"
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

func handleEmqxMqttAuth(c *gin.Context) {
	// 启用插件: emqx_auth_http
	// emqx_auth_http 插件同时包含 ACL 功能，可通过注释禁用。
	//
	// auth.http.auth_req = http://127.0.0.1:8100/mqtt/auth
	// auth.http.auth_req_method = POST
	// auth.http.auth_req_content_type = application/json
	// auth.http.auth_req.params = clientid=%c,username=%u,password=%P,ipaddress=%a
	//
	// 认证失败：API 返回非 200 状态码
	// 认证成功：API 返回 200 状态码
	// 忽略认证：API 返回 200 状态码且消息体 ignore

	var form struct {
		Username  string `form:"username"`
		Password  string `form:"password"`
		ClientId  string `form:"clientid"`
		IpAddress string `form:"ipaddress"`
		Protocol  string `form:"protocol"`
		SockPort  string `form:"sockport"`
	}
	BindAndValid(c, &form)
	fmt.Println("handleEmqxMqttAuth", form)

	c.String(200, "")
}

// vim /opt/bitnami/rabbitmq/etc/rabbitmq/rabbitmq.conf
// auth_backends.1 = internal
// auth_backends.2 = cache
// # 缓存后端指定为 http
// auth_cache.cached_backend = http
// # 缓存时间，单位毫秒
// auth_cache.cache_ttl = 60000

// auth_http.http_method   = post
// auth_http.user_path     = http://host.docker.internal:8100/auth/user
// auth_http.vhost_path    = http://host.docker.internal:8100/auth/vhost
// auth_http.resource_path = http://host.docker.internal:8100/auth/resource
// auth_http.topic_path    = http://host.docker.internal:8100/auth/topic

func handleRabbitMqUser(c *gin.Context) {
	// fmt.Println("handleRabbitMqUser==============")

	// body, _ := ioutil.ReadAll(c.Request.Body)
	// fmt.Println("---body/--- \r\n " + string(body))
	// for k, v := range c.Request.Header {
	// 	fmt.Println(k, v)
	// }

	var form struct {
		Username string `form:"username"`
		Password string `form:"password"`
		ClientId string `form:"client_id"`
		Vhost    string `form:"vhost"`
		Ip       string `form:"ip"`
	}
	BindAndValid(c, &form)
	fmt.Println("handleRabbitMqUser", form)

	if form.Username == "user" && form.Password == "bitnami" {
		c.String(200, "allow administrator")
	} else if form.Username == "admin" && form.Password == "bitnami" {
		c.String(200, "allow management")
	} else {
		c.String(200, "allow")
	}
	// c.String(200, "deny")
}
func handleRabbitMqVhost(c *gin.Context) {
	var form struct {
		Username string `form:"username"`
		Vhost    string `form:"vhost"`
		Ip       string `form:"ip"`
	}
	BindAndValid(c, &form)
	fmt.Println("handleRabbitMqVhost", form)
	c.String(200, "allow")
}
func handleRabbitMqResource(c *gin.Context) {
	var form struct {
		Username   string `form:"username"`
		Vhost      string `form:"vhost"`
		Resource   string `form:"resource"`
		Name       string `form:"name"`
		Permission string `form:"permission"`
	}
	BindAndValid(c, &form)
	fmt.Println("handleRabbitMqResource", form)
	c.String(200, "allow")
}
func handleRabbitMqTopic(c *gin.Context) {
	var form struct {
		Username   string `form:"username"`
		Vhost      string `form:"vhost"`
		Resource   string `form:"resource"`
		Name       string `form:"name"`
		Permission string `form:"permission"`
		RoutingKey string `form:"routing_key"`
	}
	BindAndValid(c, &form)
	fmt.Println("handleRabbitMqTopic", form)
	c.String(200, "allow")
}

func runHttpServer() {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Emqx
	r.POST("/mqtt/auth", handleEmqxMqttAuth)

	// RabbitMq
	r.POST("/auth/user", handleRabbitMqUser)
	r.POST("/auth/vhost", handleRabbitMqVhost)
	r.POST("/auth/resource", handleRabbitMqResource)
	r.POST("/auth/topic", handleRabbitMqTopic)

	_ = r.Run(":8100")
}

func main() {
	runHttpServer()
}
