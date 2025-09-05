package svrheadertrace

import (
	"context"
	"time"

	"maze_game_server/lib/codec"
	"maze_game_server/lib/codec/raw_pkg"
	"maze_game_server/lib/nano/frame"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/protocol/svrheader"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func Serialize(v interface{}) ([]byte, error) {
	if data, ok := v.([]byte); ok {
		return data, nil
	}

	data, err := codec.NewProtobufSerializer().Marshal(v)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func XXX(packetType uint16, pack interface{}) (uint64, []byte) {
	data, _ := Serialize(pack)
	return codec.ToMessageID(uint32(time.Now().Unix()), 0, packetType), data
}

func PacketWithSvrHeader(ctx context.Context, p []byte) (data []byte, err error) {
	ak := &raw_pkg.StruSvrEsRawBaseHead{}
	err = ak.UnPack(p)
	if err != nil {
		return
	}
	fklog.AppLogger().InfoWF("PacketWithSvrHeader", zap.Any("ak", ak))

	span := trace.SpanFromContext(ctx)
	xxx := &raw_pkg.StruSvrEsRawBaseHead{
		Header: make(map[string]string),
	}
	span.AddEvent("svrheader.Encode")
	carrier := otel.GetTextMapPropagator()
	carrier.Inject(ctx, codec.NewBaseHeader(xxx))
	hhh := svrheader.SvrHeader{
		Header: xxx.Header,
		Body:   p,
	}
	yyy, err := hhh.Encode()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return yyy, err
	}
	return yyy, err
}

func YYY(ctx context.Context, packetType uint16, pack interface{}, cc *codec.EsPacketCodec) []byte {
	mid, data := XXX(packetType, pack)
	m := &frame.Message{
		Data: data,
		ID:   mid,
	}
	xData, _ := cc.Encode(m)
	data, _ = PacketWithSvrHeader(ctx, xData)
	return data
}
