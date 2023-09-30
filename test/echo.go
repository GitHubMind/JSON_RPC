package test

type EchoArgs struct {
	S string `json:"s"`
}

type EchoResult struct {
	String string
	Int    int
	Args   *EchoArgs
}
type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
