package test

type TestError struct{}

func (TestError) Error() string          { return "testError" }
func (TestError) ErrorCode() int         { return 444 }
func (TestError) ErrorData() interface{} { return "testError data" }
