# JSONRPC
 
 you can ues it make self-test or imiate chain node.

###  support methods
- [x]  test_echo 
- [ ] test_peerInfo
- [ ] test_sleep
- [ ] test_block
- [ ] test_repeat


## How to use
 you can see in **main_test.go** that  there have 2 options : Http or Websocket.
#### For example
```shell
# method: test_echo
# params: [string, int, {object}]

curl -X POST -H 'Content-Type: application/json'  --data '{"jsonrpc":"2.0","method":"test_echo","params":["f",3,{"S":"v"}],"id":1}'  "http://localhost:3000/rpc"
# ws
echo 'test_echo ["f",3,{"S":"v"}]'|websocat 'ws://localhost:3000/ws' --jsonrpc -n -1
```