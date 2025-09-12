// @Author pangchenyang 2025/6/17 19:43:00
// @Desc: 
package main

import (
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/pkg/discovery"
	"os"
	"bytes"
	"mime/multipart"
	"io"
	"net/http"
	"encoding/json"
	"maze_game_server/servers/maze_main_server/process/userprofile"
)

const (
	resUrl     = "/uploadImgFile" // 资源服务器地址
	resSvrName = "http_img.web"
	nameSpace  = "minigame-test"
)

var avatarSvr = &AvatarSvr{}

type AvatarSvr struct {
	resolver discovery.Resolver
}

type HttpResult struct {
	Status int32       `json:"status"`
	Desc   string      `json:"desc"`
	Data   interface{} `json:"data,omitempty"`
}

func (a *AvatarSvr) UploadAvatar(userID uint64, filePath string) error {
	fmt.Println("uploadAvatar start, user:", userID)

	// 10.102.13.254:5777
	serverAddress := "10.102.13.254:5777"
	token, err := userprofile.GenerateAvatarToken()
	if err != nil {
		return fmt.Errorf("generate token failed: %v", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file failed: %v", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return fmt.Errorf("create form file failed: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("copy file to part failed: %v", err)
	}
	writer.Close()

	url := fmt.Sprintf("http://%s%s", serverAddress, resUrl)
	reqPost, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("create request failed: %v", err)
	}

	reqPost.Header.Set("Content-Type", writer.FormDataContentType())
	reqPost.Header.Set("X-Cltx-Jwt", token)

	client := &http.Client{}
	resp, err := client.Do(reqPost)
	if err != nil {
		return fmt.Errorf("send request failed: %v", err)
	}
	defer resp.Body.Close()

	var result HttpResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf("decode response failed: %v", err)
	}

	fmt.Printf("uploadAvatar end, reqPost:%v result:%v \n", reqPost, result)
	return nil
}

func main() {
	fmt.Println("avatar script start ")
	// 设置本地环境变量
	cpuRequest := os.Getenv("K8S_REQUEST_CPU")
	if cpuRequest == "" {
		cpuRequest = "100m" // 设置默认值
		fmt.Println("使用默认CPU请求值:", cpuRequest)
	}
	userID := uint64(10001)
	filepath := "./bianmu.png"
	err := avatarSvr.UploadAvatar(userID, filepath)
	if err != nil {
		fmt.Println("uploadAvatar failed:", err)
		return
	}

	fmt.Println("avatar script end ")
}
