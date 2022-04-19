package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// RegisterRabbitmqRouter 注册Rabbitmq路由
func RegisterRabbitmqRouter(r *gin.Engine) {
	// 认证
	r.POST("/rb/auth/user", handleRabbitmqUser)
	// 权鉴
	r.POST("/rb/auth/vhost", handleRabbitmqVhost)
	r.POST("/rb/auth/resource", handleRabbitmqResource)
	r.POST("/rb/auth/topic", handleRabbitmqTopic)
}

//rabbitmq_mqtt插件 提供MQTT协议的支持
//rabbitmq_web_mqtt插件 提供Web MQTT协议的支持
//rabbitmq_auth_backend_http插件 提供HTTP认证的支持

// vim /opt/bitnami/rabbitmq/etc/rabbitmq/rabbitmq.conf
// auth_backends.1 = internal
// auth_backends.2 = cache
// # 缓存后端指定为 http
// auth_cache.cached_backend = http
// # 缓存时间，单位毫秒
// auth_cache.cache_ttl = 60000
//
// auth_http.http_method   = post
// auth_http.user_path     = http://host.docker.internal:8100/rb/auth/user
// auth_http.vhost_path    = http://host.docker.internal:8100/rb/auth/vhost
// auth_http.resource_path = http://host.docker.internal:8100/rb/auth/resource
// auth_http.topic_path    = http://host.docker.internal:8100/rb/auth/topic

func handleRabbitmqUser(c *gin.Context) {
	var form struct {
		Username string `form:"username"`
		Password string `form:"password"`
		ClientId string `form:"client_id"`
		Vhost    string `form:"vhost"`
		Ip       string `form:"ip"`
	}
	BindAndValid(c, &form)
	fmt.Println("handleRabbitmqUser", form)

	if form.Username == "user" && form.Password == "bitnami" {
		c.String(200, "allow administrator")
	} else if form.Username == "admin" && form.Password == "bitnami" {
		c.String(200, "allow management")
	} else {
		c.String(200, "allow")
	}
	// c.String(200, "deny")
}
func handleRabbitmqVhost(c *gin.Context) {
	var form struct {
		Username string `form:"username"`
		Vhost    string `form:"vhost"`
		Ip       string `form:"ip"`
	}
	BindAndValid(c, &form)
	fmt.Println("handleRabbitmqVhost", form)
	c.String(200, "allow")
}

func handleRabbitmqResource(c *gin.Context) {
	var form struct {
		Username   string `form:"username"`
		Vhost      string `form:"vhost"`
		Resource   string `form:"resource"`
		Name       string `form:"name"`
		Permission string `form:"permission"`
	}
	BindAndValid(c, &form)
	fmt.Println("handleRabbitmqResource", form)
	c.String(200, "allow")
}

func handleRabbitmqTopic(c *gin.Context) {
	var form struct {
		Username   string `form:"username"`
		Vhost      string `form:"vhost"`
		Resource   string `form:"resource"`
		Name       string `form:"name"`
		Permission string `form:"permission"`
		RoutingKey string `form:"routing_key"`
	}
	BindAndValid(c, &form)
	fmt.Println("handleRabbitmqTopic", form)
	c.String(200, "allow")
}
