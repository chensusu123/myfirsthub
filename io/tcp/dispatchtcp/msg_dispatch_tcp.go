package dispatchtcp

import (
	"fmt"
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/pb/server/KafkaMsgNotify"
	"strings"
	"sync"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/common/fkfmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig/param"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fktcpclient"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var svrTypeToClient sync.Map     // svrType(int32):*fktcpclient.SyncClient
var topicTypeToSvrTypes sync.Map // topicType(int32):svrType,svrType,svrType topic类型:需要消费的服务类型
var svrTypeToCfg sync.Map        // svrType(int32):*fkconfig.CGKConfigNode

var paramRetryTimes int // 重试次数
var paramSleepTime int  // 等待时间

const separator = "_" // map类型的cgk参数中多个服务类型的分隔符

func init() {
	fkconfig.RegCGKVersionChange(onCgkParamChgTcp)
	param.SafeMapParam(&topicTypeToSvrTypes, "topic:to:server:type:map:tcp", func() {}, param.ConvInt32, param.ConvString, "topic映射tcp代理服务类型")
	param.IntP(&paramRetryTimes, "tcp:retry:time", 3, "重试次数")
	param.IntP(&paramSleepTime, "tcp:sleep:time", 500, "重试等待时间,单位毫秒")
}

func onCgkParamChgTcp(log fklog.FKLogI) error {
	log.InfoWF("onCgkParamChgTcp start")
	fkconfig.RangeConfigNode(func(cfg *fkconfig.CGKConfigNode) bool {
		log.InfoWF("onCgkParamChgTcp depend list", zap.Uint32("svrType", cfg.ServerTypeID))
		svrTypeToCfg.Store(int32(cfg.ServerTypeID), cfg)
		return true
	})

	var err error
	topicTypeToSvrTypes.Range(func(key, value interface{}) bool {
		fkutil.CaptureException()
		topicType := key.(int32)
		svrTypes := strings.Split(value.(string), separator)
		for _, v := range svrTypes {
			svrType := fkutil.ToInt32(v)
			log.InfoWF("onCgkParamChgTcp range map", zap.Int32("topicType", topicType), zap.Int32("svrType", svrType))
			if svrType == 0 {
				continue
			}

			cfg, ok := svrTypeToCfg.Load(int32(svrType))
			if !ok {
				err = errors.New(fmt.Sprintf("no found node cgk config topic:%d svrType:%d", topicType, svrType))
				log.ErrorWF("onCgkParamChgTcp cannot find config",
					zap.Int32("topicType", topicType),
					zap.Int32("svrType", svrType))
				return false
			} else {
				log.InfoWF("onCgkParamChgTcp 1",
					zap.Int32("topicType", topicType),
					zap.Int32("svrType", svrType))
				node := cfg.(*fkconfig.CGKConfigNode)
				if node != nil {
					log.InfoWF("onCgkParamChgTcp 2", zap.Int32("topicType", topicType),
						zap.Int32("svrType", svrType))
					e := rangeParamAsync(log, topicType, int32(svrType), node)
					if e != nil {
						err = e
						return false
					}
				} else {
					log.InfoWF("onCgkParamChgTcp node nil", zap.Int32("topicType", topicType), zap.Int32("svrType", svrType))
				}
			}
		}
		return true
	})
	if err != nil {
		log.ErrorWF("onCgkParamChgTcp end with err", zap.Error(err))
	} else {
		log.InfoWF("onCgkParamChgTcp end")
	}

	return err
}

func rangeParamAsync(log fklog.FKLogI, sailType, svrType int32, svrNode *fkconfig.CGKConfigNode) error {
	log.DebugWF("rangeParamAsync start with", zap.Int32("svrType", svrType), zap.Int32("sailType", sailType), zap.Any("svrNode", svrNode))
	fkutil.CaptureException()
	if sailType == 0 || svrType == 0 {
		log.WarnWF("rangeParamAsync config == 0")
		return nil
	}
	// 已存在. 就不需要手动添加了.
	if _, loaded := svrTypeToClient.Load(svrType); loaded {
		log.WarnWF("rangeParamAsync already init", zap.Int32("sailType", sailType), zap.Int32("svrType", svrType))
		return nil
	}
	c := &fktcpclient.SyncClient{MutilClient: &fknet.MutilClient{RegisterType: 1}}
	svrTypeToClient.Store(svrType, c)

	err := fkconfig.RegisterNameNodeWithOpen(svrNode.ServiceName, uint32(svrType), c, svrNode, log)
	if err != nil {
		log.ErrorWF("rangeParamAsync open rpc fail", zap.Error(err), zap.Int32("sailType", sailType),
			zap.Int32("svrType", svrType), zap.String("serviceName", svrNode.ServiceName))
		fkfmt.Println("rangeParamAsync open rpc fail", "err:", err, "sailType:", sailType, "svrType:", svrType, "serviceName:", svrNode.ServiceName)
		return err
	}

	log.InfoWF("rangeParamAsync open rpc succ", zap.Int32("sailType", sailType), zap.Int32("svrType", svrType),
		zap.String("serviceName", svrNode.ServiceName))
	return nil
}

func getClientByTopic(agent fklog.FKLogI, topic int32) (typeAndClients map[int32]*fktcpclient.SyncClient) {
	svrs, ok := topicTypeToSvrTypes.Load(topic)
	if !ok {
		agent.WarnWF("getClientByTopic cannot find svr type", zap.Int32("topic", topic))
		return
	}
	svrTypeStrs := strings.Split(svrs.(string), separator)
	if len(svrTypeStrs) == 0 {
		agent.WarnWF("getClientByTopic svrType == 0", zap.Int32("topic", topic))
		return
	}
	typeAndClients = make(map[int32]*fktcpclient.SyncClient)
	for _, svrTypeStr := range svrTypeStrs {

		svrType := fkutil.ToInt32(svrTypeStr)
		if svrType <= 0 {
			agent.WarnWF("getClientByTopic svrTypeStr error", zap.String("svrTypeStr", svrTypeStr))
			continue
		}
		client, ok := svrTypeToClient.Load(int32(svrType))
		if !ok {
			agent.ErrorWF("getClientByTopic cannot find client", zap.Int32("svrType", svrType))
			continue
		}
		typeAndClients[svrType] = client.(*fktcpclient.SyncClient)
	}
	return
}

func DispatchKafkaMsgTcp(logger fklog.FKLogI, uid uint64, kafkaName string, data []byte) (err error) {
	topicType := constdef.GetMsgType(kafkaName)
	if topicType == 0 {
		logger.ErrorWF("DispatchKafkaMsgTcp cannot find typ",
			zap.Uint64("uid", uid),
			zap.String("kafkaName", kafkaName),
			zap.String("msg", string(data)))
		return errors.New("cannot find topicType")
	}

	rq := &KafkaMsgNotify.KafkaMsgDistributeRQ{}
	rq.UserId = proto.Uint64(uid)
	rq.MsgData = data
	rq.MsgType = proto.Int32(topicType)
	tcpClients := getClientByTopic(logger, topicType)
	if len(tcpClients) == 0 {
		logger.WarnWF("DispatchKafkaMsgTcp 未找到后端业务", zap.Int32("topicType", topicType), zap.String("kafkaName", kafkaName))
		return nil
	}
	var wg sync.WaitGroup
	wg.Add(len(tcpClients))
	for svrType, cli := range tcpClients {
		tmpSvrType := svrType
		tmpCli := cli
		go func() {
			defer wg.Done()
			rs := &KafkaMsgNotify.KafkaMsgDistributeRS{}
			count := 1

			resp, _, err := tmpCli.Call(logger, uid, 20989, rq, 20990, rs)
			res := resp.(*KafkaMsgNotify.KafkaMsgDistributeRS)
			defer func() {
				logger.InfoWF("DispatchKafkaMsgTcp end",
					zap.Int32("svrType", tmpSvrType),
					zap.Int32("topicType", topicType),
					zap.Int("cnt", count),
					zap.Any("rq", rq), zap.Any("rs", res))
			}()
			if err == nil && res != nil && res.GetErrInfo() != nil && res.GetErrInfo().GetErrCode() == errors.NO_ERROR_CODE {
				return
			}

			// 需要重试
			for ; count <= paramRetryTimes; count++ {
				resp, _, err := tmpCli.Call(logger, uid, 20989, rq, 20990, rs)
				res := resp.(*KafkaMsgNotify.KafkaMsgDistributeRS)
				if err == nil && res != nil && res.GetErrInfo() != nil && res.GetErrInfo().GetErrCode() == errors.NO_ERROR_CODE {
					return
				}
				if paramSleepTime == 0 {
					paramSleepTime = 10
				}
				time.Sleep(time.Duration(paramSleepTime) * time.Millisecond)
			}
		}()
	}

	wg.Wait()
	return
}
