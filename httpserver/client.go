package httpserver

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/name5566/leaf/chanrpc"
	"github.com/name5566/leaf/log"
)

// NewClient NewClient
func NewClient() *Client {
	return &Client{
		chans: make(map[string]*chanrpc.Server),
	}
}

// Client Client
type Client struct {
	chans map[string]*chanrpc.Server
}

// GET GET
func (its *Client) GET(url string, chanRPC *chanrpc.Server, f interface{}) {
	// 向服务器请求
	resp, err := http.Get(url)
	if err != nil {
		log.Error("GET URL %q failed: %v", url, err)
		return
	}

	// 注册到`chanrpc`
	id := fmt.Sprintf("GET%s", url)

	// 消息路由
	its.route(id, url, resp, chanRPC, f)
}

// POST POST
func (its *Client) POST(url string, values url.Values, chanRPC *chanrpc.Server, f interface{}) {
	resp, err := http.PostForm(url, values)
	if err != nil {
		log.Error("POST URL %q failed %v", url, err)
		return
	}

	// 注册到`chanrpc`
	id := fmt.Sprintf("POST%s%v", url, values)

	// 消息路由
	its.route(id, url, resp, chanRPC, f)
}

// SyncPost SyncPost
func (its *Client) SyncPost(url string, values url.Values, fun func(*http.Response)) {
	resp, err := http.PostForm(url, values)
	if err != nil {
		log.Error("POST URL %q failed %v", url, err)
		return
	}

	if nil != fun {
		fun(resp)
	}
}

func (its *Client) route(id, url string, resp *http.Response, chanRPC *chanrpc.Server, f interface{}) {
	if _, ok := its.chans[id]; !ok && chanRPC != nil {
		its.chans[url] = chanRPC
		chanRPC.Register(id, f)
	}
	if chanRPC != nil {
		chanRPC.Go(id, resp)
	}
}
