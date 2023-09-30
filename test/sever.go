package test

type TestService struct{}
type Args struct {
	A, B int
}

func (t *TestService) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

// func (s *TestService) NoArgsRets() {}
func (s *TestService) Echo(str string, i int, args *EchoArgs) EchoResult {
	return EchoResult{str, i, args}
}

//
//func (s *TestService) EchoWithCtx(ctx context.Context, str string, i int, args *EchoArgs) EchoResult {
//	return EchoResult{str, i, args}
//}
//
////func (s *TestService) PeerInfo(ctx context.Context) rpc.PeerInfo {
////	return rpc.PeerInfoFromContext(ctx)
////}
//
//func (s *TestService) Sleep(ctx context.Context, duration time.Duration) {
//	time.Sleep(duration * time.Second)
//}
//
//func (s *TestService) Block(ctx context.Context) error {
//	<-ctx.Done()
//	return errors.New("context canceled in testservice_block")
//}
//
//func (s *TestService) Rets() (string, error) {
//	return "", nil
//}
//
////lint:ignore ST1008 returns error first on purpose.
//func (s *TestService) InvalidRets1() (error, string) {
//	return nil, ""
//}
//
//func (s *TestService) InvalidRets2() (string, string) {
//	return "", ""
//}
//
//func (s *TestService) InvalidRets3() (string, string, error) {
//	return "", "", nil
//}
//
//func (s *TestService) ReturnError() error {
//	return TestError{}
//}
//
//func (s *TestService) Repeat(r int, c string) string {
//	if c == "" {
//		return strings.Repeat("x", r)
//	}
//	return strings.Repeat(c, r)
//}
//
////
////func (s *TestService) CallMeBack(ctx context.Context, method string, args []inter{}) (inter{}, error) {
////	c, ok := rpc.ClientFromContext(ctx)
////	if !ok {
////		return nil, errors.New("no client")
////	}
////	var result inter{}
////	err := c.Call(&result, method, args...)
////	return result, err
////}
////
////func (s *TestService) CallMeBackLater(ctx context.Context, method string, args []inter{}) error {
////	c, ok := rpc.ClientFromContext(ctx)
////	if !ok {
////		return errors.New("no client")
////	}
////	go func() {
////		<-ctx.Done()
////		var result inter{}
////		c.Call(&result, method, args...)
////	}()
////	return nil
////}
////
////func (s *TestService) Subscription(ctx context.Context) (*rpc.Subscription, error) {
////	return nil, nil
////}
