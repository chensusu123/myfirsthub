package common_value

import (
	"gitlab.ifreetalk.com/plate/freetk/fkserver/thrift_service"
	"gitlab.ifreetalk.com/plate/protodef/MazeCommonValueSvr"
)

func RegisterRpcPackProcessor() {
	thrift_service.RegisterTwowaySimple(131447, &MazeCommonValueSvr.MazeCommonValueQueryRQ{},
		131448, &MazeCommonValueSvr.MazeCommonValueQueryRQ{}, OnMazeCommonValueQueryRQ)
	thrift_service.RegisterTwowaySimple(131449, &MazeCommonValueSvr.MazeCommonValueAddRQ{},
		131450, &MazeCommonValueSvr.MazeCommonValueAddRS{}, OnMazeCommonValueAddRQ)
	thrift_service.RegisterTwowaySimple(131451, &MazeCommonValueSvr.MazeCommonValueSubRQ{},
		131452, &MazeCommonValueSvr.MazeCommonValueSubRS{}, OnMazeCommonValueSubRQ)
	thrift_service.RegisterTwowaySimple(131453, &MazeCommonValueSvr.MazeCommonValueSetRQ{},
		131454, &MazeCommonValueSvr.MazeCommonValueSetRQ{}, OnMazeCommonValueSetRQ)
}
