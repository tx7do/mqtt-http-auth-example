package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// RegisterEmqxRouter 注册emqx路由
func RegisterEmqxRouter(r *gin.Engine) {
	r.POST("/emqx/auth", handleEmqxAuth)

	r.POST("/emqx/hook", handleEmqxHook)
}

// 启用插件: emqx_auth_http
// vi etc/plugins/emqx_auth_http.conf
// emqx_auth_http 插件同时包含 ACL 功能，可通过注释禁用。
//
// auth.http.auth_req = http://127.0.0.1:8100/emqx/auth
// auth.http.auth_req_method = POST
// auth.http.auth_req_content_type = application/json
// auth.http.auth_req.params = clientid=%c,username=%u,password=%P,ipaddress=%a
//
// 认证失败：API 返回非 200 状态码
// 认证成功：API 返回 200 状态码
// 忽略认证：API 返回 200 状态码且消息体 ignore

func handleEmqxAuth(c *gin.Context) {
	var form struct {
		Username  string `form:"username"`
		Password  string `form:"password"`
		ClientId  string `form:"clientid"`
		IpAddress string `form:"ipaddress"`
		Protocol  string `form:"protocol"`
		SockPort  string `form:"sockport"`
	}
	BindAndValid(c, &form)
	fmt.Println("handleEmqxAuth", form)

	c.String(200, "")
}

// 启用插件: emqx_web_hook
//
// EMQ X 并不关心 Web 服务的返回
//
// vi etc/plugins/emqx_web_hook.conf

// web.hook.url = http://host.docker.internal:8100/emqx/hook
// web.hook.rule.client.connected.1     = {"action": "on_client_connected"}
// web.hook.rule.client.disconnected.1  = {"action": "on_client_disconnected"}

func handleEmqxHook(c *gin.Context) {
	//printRequestHeader(c)
	//printRequestBody(c)
	var form struct {
		Action      string `json:"action,omitempty"`    // 事件名称
		Username    string `json:"username,omitempty"`  // 客户端 Username，不存在时该值为 "undefined"
		ClientId    string `json:"clientid,omitempty"`  // 客户端 ClientId
		IpAddress   string `json:"ipaddress,omitempty"` // 客户端源 IP 地址
		Reason      string `json:"reason,omitempty"`
		KeepAlive   int64  `json:"keepalive,omitempty"`
		ProtoVer    int64  `json:"proto_ver,omitempty"` // 协议版本号
		ConnectedAt int64  `json:"connected_at,omitempty"`
	}
	BindAndValid(c, &form)

	fmt.Println("handleEmqxHook", form)

	switch form.Action {
	case "client_connected":
		fmt.Println("设备连接成功", form.Username, form.ClientId)
	case "client_disconnected":
		fmt.Println("设备断开连接成功", form.Username, form.ClientId)
	}
}
