package httpserver

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/name5566/leaf/chanrpc"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
)

// Server Server
type Server struct {
	*http.Server
	ln          net.Listener
	wgLn        sync.WaitGroup
	Gate        *gate.Gate
	Addr        string
	CertFile    string
	KeyFile     string
	HTTPTimeout time.Duration
}

// Close Close
func (its *Server) Close() {
	its.ln.Close()
	its.wgLn.Wait()
}

// Start Start
func (its *Server) Start() {
	its.init()
	go its.run()
}

func (its *Server) init() {
	ln, err := net.Listen("tcp", its.Addr)
	if err != nil {
		log.Fatal("%v", err)
	}

	if its.HTTPTimeout <= 0 {
		its.HTTPTimeout = 10 * time.Second
		log.Release("invalid HTTPTimeout, reset to %v", its.HTTPTimeout)
	}

	if its.CertFile != "" && its.KeyFile != "" {
		config := &tls.Config{}
		config.NextProtos = []string{"http/1.1"}

		var err error
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(its.CertFile, its.KeyFile)
		if err != nil {
			log.Fatal("%v", err)
		}

		ln = tls.NewListener(ln, config)
	}

	its.ln = ln
	its.Server = &http.Server{
		Addr:           its.Addr,
		ReadTimeout:    its.HTTPTimeout,
		WriteTimeout:   its.HTTPTimeout,
		MaxHeaderBytes: 1024,
	}
}

func (its *Server) run() {
	its.wgLn.Add(1)
	defer its.wgLn.Done()
	log.Debug("启动http服=====================")
	router := its.NewRouter()
	its.Handler = router.Handler()
	its.Serve(its.ln)

	router.OnClose()
}

// NewRouter NewRouter
func (its *Server) NewRouter() IRouter {
	result := &Router{
		Router: httprouter.New(),
		gate:   its.Gate,
		chans:  make(map[string]*chanrpc.Server),
	}
	if its.Gate.AgentChanRPC != nil {
		result.isDone.Add(1)
		its.Gate.AgentChanRPC.Go("NewServer", result)
		result.isDone.Wait()
	}
	return result
}

// IRouter IRouter
type IRouter interface {
	OnClose()
	Handler() http.Handler
	GET(string, interface{})
	POST(string, interface{})
	Done()
}

// Router Router
type Router struct {
	*httprouter.Router
	gate   *gate.Gate
	chans  map[string]*chanrpc.Server
	isDone sync.WaitGroup
	wg     sync.WaitGroup
}

// Done Done
func (its *Router) Done() {
	its.isDone.Done()
}

// Handler Handler
func (its *Router) Handler() http.Handler {
	return its.Router
}

// OnClose OnClose
func (its *Router) OnClose() {
	its.wg.Wait()
}

// POST POST
func (its *Router) POST(path string, f interface{}) {
	chanRPC := its.gate.AgentChanRPC
	id := fmt.Sprintf("POST%s", path)
	log.Debug("注册httpID: %v", id)
	// 注册到`chanrpc`
	if _, ok := its.chans[id]; !ok && chanRPC != nil {
		its.chans[id] = chanRPC
		chanRPC.Register(id, f)
	}

	// 注册到`httprouter`
	its.Router.POST(path, func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		log.Debug("http接收到信息----------------------------")
		its.wg.Add(1)
		defer its.wg.Done()

		w.Header().Set("Access-Control-Allow-Origin", "*")
		// if e := route(id, chanRPC, w, ReadPost(req)); e != nil {
		if e := route(id, chanRPC, w, ReadGet(req)); e != nil {
			fmt.Fprintf(w, "Failed: %v", e)

		}
	})
}

// GET GET
func (its *Router) GET(path string, f interface{}) {
	chanRPC := its.gate.AgentChanRPC
	id := fmt.Sprintf("GET%s", path)

	// 注册到`chanrpc`
	if _, ok := its.chans[id]; !ok && chanRPC != nil {
		its.chans[id] = chanRPC
		chanRPC.Register(id, f)
	}

	// 注册到`httprouter`
	its.Router.GET(path, func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		its.wg.Add(1)
		defer its.wg.Done()

		w.Header().Set("Access-Control-Allow-Origin", "*")

		if e := route(id, chanRPC, w, ReadGet(req)); e != nil {
			fmt.Fprintf(w, "Failed: %v", e)
		}
	})
}

// ReadGet ReadGet
func ReadGet(r *http.Request) map[interface{}]interface{} {
	defer func() {
		if e := recover(); e != nil {
			log.Error("read get failed: %v", e)
		}
	}()

	result := make(map[interface{}]interface{})
	form := r.URL.Query()
	for k, v := range form {
		result[k] = v
	}
	return result
}

// ReadPost ReadPost
func ReadPost(r *http.Request) []interface{} {
	defer func() {
		if e := recover(); e != nil {
			log.Error("read post failed: %v", e)
		}
	}()

	result := make([]interface{}, 0)
	if e := r.ParseForm(); e != nil {
		panic(e.Error())
	}

	for k := range r.Form {
		result = append(result, r.Form.Get(k))
	}
	return result
}

func route(id string, chanRPC *chanrpc.Server, w http.ResponseWriter, args ...interface{}) error {
	if chanRPC == nil {
		return fmt.Errorf("chanRPC must not be nil")
	}

	result, e := chanRPC.Call1(id, args...)
	if e != nil {
		return e
	}

	// 获取结构类型
	t := reflect.TypeOf(result)

	// 类型判断（仅支持`string`和`ptr`）
	if k := t.Kind(); k == reflect.String {
		fmt.Fprintf(w, "%v", result)
	} else if k == reflect.Ptr {
		if e = json.NewEncoder(w).Encode(result); e != nil {
			return e
		}
	} else {
		return fmt.Errorf("invalid result type %v", t)
	}
	return nil
}
