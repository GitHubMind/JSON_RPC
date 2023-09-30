package test

import (
	"fmt"
	"sync"
	"time"

	"github.com/valyala/fastjson"
)

var (
	jsonParserPool = fastjson.ParserPool{}
)

/*
*
a simplified version that will be cached and kept until ws response is received
*/
type RPCHeader struct {
	ID     string `json:"id"`
	Method string `json:"method"`
	// for redis cache
	CacheKey string        `json:"cacheKey"`
	CacheTTL time.Duration `json:"cacheTTL"`

	Weight int `json:"weight"`
}

type RPC struct {
	RPCHeader `mapstructure:",squash"`
	Params    string `json:"-"`
	Error     error  `json:"-"`
	BatchId   int    `json:"-"`
	Cached    bool   `json:"-"`
	//Length    int
}

type RPCPushMsg struct {
	Method string
	Params string
	Length int
	Seq    int
}

type RPCResponse struct {
	ID     string // it is a raw jsonstring, means quotes are included when type is string
	Error  error
	Result string // it is a raw jsonstring, means quotes are included when type is string
	Cached bool
	RpcReq *RPCHeader
}

//type RPCPair struct {
//	WsReq  *RPCHeader
//	Resp *RPCResponse
//}

func decodeRPC(v *fastjson.Value) *RPC {
	rpc := &RPC{}
	idValue := v.Get("id")
	if idValue != nil {
		rpc.ID = idValue.String()
	}
	rpc.Method = string(v.GetStringBytes("method"))
	params := v.Get("params")
	if params != nil {
		rpc.Params = params.String()
	}
	return rpc
}

func DecodeJSONRPC(b []byte) (rpcs []*RPC, err error) {
	parser := jsonParserPool.Get()
	defer jsonParserPool.Put(parser)
	v, err := parser.ParseBytes(b)
	if err != nil {
		rpc := &RPC{
			Error: err,
		}
		rpcs = append(rpcs, rpc)
		return
	}
	if v.Type() == fastjson.TypeArray {
		a := v.GetArray()
		l := len(a)
		for _, item := range a {
			rpc := decodeRPC(item)
			rpc.BatchId = l
			rpcs = append(rpcs, rpc)
		}
		return
	}
	rpc := decodeRPC(v)
	rpcs = append(rpcs, rpc)
	return
}

func (rpc RPC) IDString() string {
	return rpc.ID
}

func isPushMsg(v *fastjson.Value) bool {
	return !v.Exists("id")
}

func DecodeResponse(msg []byte) (err error, resps []*RPCResponse, pushMsg *RPCPushMsg) {
	parser := jsonParserPool.Get()
	defer func() {
		jsonParserPool.Put(parser)
	}()
	v, err := parser.ParseBytes(msg)
	if err != nil {
		return
	}
	if v.Type() == fastjson.TypeArray {
		a := v.GetArray()
		for _, item := range a {
			resps = append(resps, decodeRPCResponse(item))
		}
		return
	}
	if isPushMsg(v) {
		pushMsg = parsePushMsg(v)
		pushMsg.Length = len(msg)
		return nil, nil, pushMsg
	}

	resps = append(resps, decodeRPCResponse(v))
	return
}

func parsePushMsg(v *fastjson.Value) (push *RPCPushMsg) {
	push = &RPCPushMsg{}
	method := v.Get("method")
	if method != nil {
		push.Method = string(v.GetStringBytes("method"))
	}
	params := v.Get("params")
	if params != nil {
		push.Params = params.String()
	}
	return
}

func decodeRPCResponse(v *fastjson.Value) (resp *RPCResponse) {
	resp = &RPCResponse{Cached: false}
	idValue := v.Get("id")
	if idValue != nil {
		resp.ID = idValue.String()
	}
	result := v.Get("result")
	if result != nil {
		resp.Result = result.String()
	} else {
		obj := v.GetObject("error")
		if obj != nil {
			err := GetErrJSONRPC(obj.Get("code").GetInt(), obj.Get("message").String())
			resp.Error = err
		}
	}
	return
}

type ErrJSONRPC struct {
	code int
	msg  string
}

// TODO why jsonrpc has no id
func (err ErrJSONRPC) Error() string {
	return fmt.Sprintf("JSONRPC error: %s, Code %d", err.msg, err.code)
}

var errPool = sync.Pool{
	New: func() interface{} {
		return new(ErrJSONRPC)
	},
}

func GetErrJSONRPC(code int, msg string) *ErrJSONRPC {
	ret := errPool.Get().(*ErrJSONRPC)
	ret.code = code
	ret.msg = msg
	return ret
}

func PutErrJSONRPC(err *ErrJSONRPC) {
	errPool.Put(err)
}
