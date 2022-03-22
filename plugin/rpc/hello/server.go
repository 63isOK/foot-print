package hello

// NOTE: step2: 基于rpc/grpc实现对接口的封装
// 服务端代码,封装了接口对具体实现的调用
// 这层封装会在Plugin.Server实现中使用
// 这里套多层的目的是为了简化实现插件时的写法法

// Server is a real impl
type Server struct {
	Impl Greeter
}

func (s *Server) Greet(args interface{}, resp *string) (err error) {
	*resp, err = s.Impl.Greet()
	return
}

func (s *Server) GreetAgain(args interface{}, resp *string) error {
	*resp = s.Impl.GreetAgain()
	return nil
}
