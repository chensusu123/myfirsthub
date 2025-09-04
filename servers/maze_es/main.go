package main

import (
	"fmt"
	"os"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/hertzhttpservice"
)

func GetEnv(name string) (string, error) {
	val, exists := os.LookupEnv(name)
	if !exists {
		return "", fmt.Errorf("environment %s not set", name)
	}
	return val, nil
}

// 19987	UN_CGK_SVR_TYPE_MAZE_MAIN_SERVER 小程序版迷宫主服务
func main() {
	fkserver.Version = "1.0.0"

	httpSvr_2 := hertzhttpservice.NewHttpService("htttp_demo_2", &HttpDemo2{})

	fkserver.AddBusiness(httpSvr_2)

	fkserver.Run()
}
