package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"

	"mock_server/test"

	"github.com/gorilla/websocket"
)

var listen string
var id string
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var registerMap map[string]any

func init() {
	registerMap = map[string]any{
		"test": new(test.TestService),
	}
}

// set header
func headerAdd(w *http.ResponseWriter) {

	(*w).Header().Add("Onf-Endpoint", id)

}
func handleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	headerAdd(&w)
	conn, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()
	conn.PongHandler()
	for true {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
		}
		jsonRpcArr, err := test.DecodeJSONRPC(message)
		//reflect.TypeOf()
		//对应 map xxx.xxx
		if err != nil {
			log.Println("decode is err", err)
			return
		}
		//batch
		for _, RpcArr := range jsonRpcArr {
			// TODO
			val, err := register(RpcArr, false)
			if err != nil {
				fmt.Fprint(w, err)

			} else {
				conn.WriteJSON(val)
			}
		}

	}
}
func register(RpcArr *test.RPC, Json bool) (any, error) {
	methodsStrArr := strings.Split(RpcArr.RPCHeader.Method, "_")
	if len(methodsStrArr) < 2 {
		return nil, errors.New("method must be  xxx_xxx")

	}
	end, ok := registerMap[methodsStrArr[0]]
	if !ok {
		//return nil, errors.New("no register router")
	}
	valueReflect := reflect.TypeOf(end)
	v := reflect.ValueOf(end)
	for i := 0; i < valueReflect.NumMethod(); i++ {
		method := valueReflect.Method(i)
		if strings.ToLower(method.Name) == strings.ToLower(methodsStrArr[1]) {
			//method.Func
			var data []interface{}
			err := json.Unmarshal([]byte(RpcArr.Params), &data)
			if err != nil {
				return nil, fmt.Errorf("json unmarshall err:%v", err)
			}
			if len(data) != method.Type.NumIn()-1 {
				return nil, errors.New("input params of quantity != the method  of parameter")
			}
			args := make([]reflect.Value, len(data))
			for i := 0; i < method.Type.NumIn(); i++ {
				paramType := method.Type.In(i)
				// -1 name of methods
				if i == 0 {
					continue
				}

				// more layer structure
				if reflect.TypeOf(map[string]interface{}{}) == reflect.TypeOf(data[i-1]) {
					jsonBytes, err := json.Marshal(data[i-1])
					if err != nil {
						return nil, fmt.Errorf("Error: %v", err)
					}
					value := reflect.New(paramType)
					err = json.Unmarshal(jsonBytes, value.Interface())
					if err != nil {
						return nil, fmt.Errorf("Error: %v", err)
					}
					args[i-1] = value.Elem()
				} else {
					//base type
					args[i-1] = reflect.ValueOf(data[i-1]).Convert(paramType)
				}
				//fmt.Printf("Parameter %d: %s\n", i-1, paramType)
			}

			result := v.MethodByName(method.Name).Call(args)
			results := make([]interface{}, len(result))
			for i, v := range result {
				results[i] = v.Interface()
			}
			if !Json {
				return results, nil
			}
			val, err := json.Marshal(results)
			if err != nil {
				log.Println("marshal err:", err)
			}
			return val, nil

		}
	}
	return nil, errors.New("can't find this method ,but struct is existed")

}
func handleJSONRPCRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	headerAdd(&w)
	body, _ := ioutil.ReadAll(r.Body)
	//[]byte
	jsonRpcArr, err := test.DecodeJSONRPC(body)
	//reflect.TypeOf()
	//对应 map xxx.xxx
	if err != nil {
		log.Println("decode is err", err)
		return
	}
	//batch
	for _, RpcArr := range jsonRpcArr {
		// TODO
		val, err := register(RpcArr, true)
		if err != nil {
			fmt.Fprint(w, err)
		} else {
			byptVal := val.([]byte)
			fmt.Fprint(w, string(byptVal))
		}
	}
}

func main() {
	listen = "3000"
	//TODO compatible command line startup
	flag.StringVar(&listen, "listen", "3000", "address to listen to")
	flag.StringVar(&id, "id", "", "id for rpc server ")
	flag.Parse()
	if id == "" {
		id = listen
	}
	//rpc
	http.HandleFunc("/ws", handleWebSocketConnection)
	http.HandleFunc("/rpc", handleJSONRPCRequest)

	log.Fatal(http.ListenAndServe(":"+listen, nil))
}
