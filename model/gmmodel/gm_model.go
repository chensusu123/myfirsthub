package gmmodel

import (
	"github.com/iancoleman/orderedmap"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager"
)

type ShowSheet struct {
	ErrorCode uint64                          `json:"errorCode"`
	ErrorMsg  string                          `json:"errorMsg"`
	Data      []config_manager.ConfigShowItem `json:"data"`
}

type Output struct {
	Status int         `json:"status"`
	Desc   string      `json:"desc"`
	Data   DynamicData `json:"data"`
}

// 调整DynamicData以使用有序map
type DynamicData struct {
	List  []*orderedmap.OrderedMap `json:"list,omitempty"` // 有序map切片，保证对象格式和顺序
	Total int                      `json:"total,omitempty"`
}

func NewSheet() *ShowSheet {
	return &ShowSheet{}
}

func NewOutPut(status int, desc string, data DynamicData) *Output {
	return &Output{
		Status: status,
		Desc:   desc,
		Data:   data,
	}
}

func NewDynamicData(list []*orderedmap.OrderedMap, total int) *DynamicData {
	return &DynamicData{
		List:  list,
		Total: total,
	}
}

type AllianceInfo struct {
	AllianceID   int    `json:"alliance_id,omitempty"`
	AllianceName string `json:"alliance_name,omitempty"`
	CreateAt     int64  `json:"create_at,omitempty"`
}

type ListInfo struct {
	List     []*AllianceInfo `json:"list,omitempty"`
	Count    int             `json:"count,omitempty"`
	Total    int             `json:"total,omitempty"`
	Page     int             `json:"page,omitempty"`
	PageSize int             `json:"page_size,omitempty"`
}

type ReturnMsg struct {
	Code int    `json:"status"`
	Msg  string `json:"desc"`
	Data any    `json:"data,omitempty"`
}

func MakeSuccessReturnMsg(data any) ReturnMsg {
	return ReturnMsg{
		Code: 200,
		Msg:  "success",
		Data: data,
	}
}
