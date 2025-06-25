package jwt

import (
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/proto"
	"maze_game_server/pb/common/MazePay"
	"testing"
)

func TestGenPayJwt(t *testing.T) {
	orderJwt, err := GeneratePayJWT(123, "123", "wx123", 1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(orderJwt)
}

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEyMywidW5pcXVlSWQiOiIxMjMiLCJ3ZWNoYXRPcGVuSWQiOiJ3eDEyMyIsImV4cCI6MTc1MDA0MTYzNn0.9NA7NFIopgp0z0zWeQ3Ng4qjGSw47HIjm2AJAuFgYaY
func TestValidatePayJwt(t *testing.T) {
	order, err := ValidatePayJWT("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEyMywidW5pcXVlSWQiOiIxMjMiLCJ3ZWNoYXRPcGVuSWQiOiJ3eDEyMyIsImV4cCI6MTc1MDA0MTYzNn0.9NA7NFIopgp0z0zWeQ3Ng4qjGSw47HIjm2AJAuFgYaY")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(order)
}

func Test(t *testing.T) {
	req := &MazePay.MazePayTokenRQ{}
	req.UniqueId = proto.String("xxx")
	b, _ := json.Marshal(req)
	s := string(b)
	fmt.Println(s)
}
