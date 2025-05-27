package codec

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/lonng/nano/component"
)

var (
	ErrInvalidMethod = errors.New("invalid method")
)

var (
	reg = regexp.MustCompile(`^[0-9a-zA-Z\_]+\_([0-9]+)\_([0-9]+)$`)
)

type Routes struct {
	comps []component.Component
}

// Register 注册nano接口格式的组件实例
func (rts *Routes) Register(comp component.Component) {
	rts.comps = append(rts.comps, comp)
}

type target struct {
	rsID    uint16
	handler string
}

// route 解析路由处理函数，返回请求RqID与RsID，如：Game.Welcome_1_2
func route(fnName string) (rqID, rsID uint16, ok bool) {
	match := reg.FindStringSubmatch(fnName)
	// Matched
	if ok = len(match) > 0; !ok {
		return
	}
	rq, err := strconv.Atoi(match[1])
	if err != nil || rq <= 0 {
		return 0, 0, false
	}
	rs, err := strconv.Atoi(match[2])
	if err != nil || rs <= 0 {
		return 0, 0, false
	}
	return uint16(rq), uint16(rs), true
}
