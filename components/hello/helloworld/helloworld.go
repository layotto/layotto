package helloworld

import "github.com/layotto/layotto/components/hello"

type HelloWorld struct {
	Say string
}

var _ hello.HelloService = &HelloWorld{}

func NewHelloWorld() hello.HelloService {
	return &HelloWorld{}
}

func (hw *HelloWorld) Init(config *hello.HelloConfig) error {
	hw.Say = config.HelloString
	return nil
}

func (hw *HelloWorld) Hello(req *hello.HelloRequest) (*hello.HelloReponse, error) {
	return &hello.HelloReponse{
		HelloString: hw.Say,
	}, nil
}
