package business

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkalert"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/freetk/protodef/MysqlParam"

	"go.uber.org/zap"
)

func init() {
}

type LoadConfigMsg struct {
	ServerName   string `json:"server_name,omitempty"`
	ServerTypeId uint32 `json:"server_type_id,omitempty"`
	ServerId     uint32 `json:"server_id,omitempty"`
	GroupId      uint32 `json:"group_id,omitempty"`
	FileName     string `json:"file_name,omitempty"`
	SheetName    string `json:"sheet_name,omitempty"`
	SheetDataMd5 string `json:"sheet_data_md5,omitempty"`
	FileMd5      string `json:"file_md5,omitempty"`
	Desc         string `json:"desc,omitempty"`        // git_version_for_config.txt文件中的版本信息
	Md5Changed   bool   `json:"md5_changed,omitempty"` // 配表数据是否发生变化
	PushTimeS    int64  `json:"push_time_s,omitempty"` // 消息生成秒级时间戳
}

// OnLoadCacheConfig 加载缓存配置请求
func (tb *tCustomBusiness) LoadCacheConfig(logger fklog.FKLogI, req *MysqlParam.DbTplRQ, res *MysqlParam.DbTplRS) error {
	fkprometheus.InfoPMT("OnLoadCacheConfig")()
	start := time.Now().UnixNano()

	// 回包大小
	pbSize := 0

	res.ErrorCode = proto.Int32(0)
	res.ReadAll = proto.Bool(true)

	// 回包
	excelRes := &excelReadResult{desc: tb.gitCfgVersion.Load()}
	defer func() {
		endt := (time.Now().UnixNano() - start) / 1e6
		logger.DebugWF("OnLoadCacheConfig end",
			zap.Int64("ms", endt),
			zap.Uint64("request_id", req.GetRqId()),
			zap.String("file", req.GetFileName()),
			zap.String("sheet", req.GetDbProcName()),
			zap.String("sheetDataMd5", excelRes.sheetDataMd5),
			zap.String("fileMd5", excelRes.fileMd5),
			zap.Any("desc", excelRes.desc),
			zap.Strings("Fields", req.GetParamList()),
			zap.String("rq_md5", req.GetMd5Check()),
			zap.Uint32("start_index", req.GetStartIndex()),
			zap.String("server", req.GetServerName()),
			zap.Int("result len", len(res.Results)),
			zap.Int("row count", int(res.GetRowCount())),
			zap.String("md5", res.GetMd5Check()),
			zap.Bool("IsEnd", res.GetReadAll()),
			zap.Int("calc", pbSize))

		StatIns.AddLoadNum()
		if excelRes.all {
			msg := &LoadConfigMsg{
				ServerName:   req.GetServerName(),
				ServerTypeId: req.GetServerTypeId(),
				ServerId:     req.GetServerId(),
				GroupId:      req.GetGroupId(),
				FileName:     req.GetFileName(),
				SheetName:    req.GetDbProcName(),
				SheetDataMd5: excelRes.sheetDataMd5,
				FileMd5:      excelRes.fileMd5,
				Desc:         excelRes.desc,
				Md5Changed:   excelRes.md5Changed,
				PushTimeS:    time.Now().Unix(),
			}
			select {
			case buffer <- msg:
				StatIns.AddLoadAllNum()
			default:
				logger.ErrorWF("OnLoadCacheConfig kafka buffer full.")
				StatIns.AddBufferFullNum()
			}
		} else {
			StatIns.AddNoLoadAllNum()
			logger.DebugWF("OnLoadCacheConfig no load all",
				zap.Any("pageSize", req.GetPageSize()), zap.Any("startIndex", req.GetStartIndex()))
		}
	}()

	setError := func(code int32, tip string) {
		res.ErrorCode = &code
		tip = fmt.Sprintf("%s.request from %s sheet:%s get fiels:%v", tip,
			req.GetFileName(), req.GetDbProcName(), req.GetParamList())
		res.ErrorInfo = &tip
		logger.ErrorWF("OnLoadCacheConfig load cache config failed.",
			zap.Int32("code", code), zap.String("tip", tip),
			zap.Uint64("request_id", req.GetRqId()), zap.String("file", req.GetFileName()),
			zap.String("sheet", req.GetDbProcName()), zap.Strings("Fields", req.GetParamList()),
			zap.String("rq_md5", req.GetMd5Check()), zap.Uint32("start_index", req.GetStartIndex()),
			zap.String("server", req.GetServerName()))
	}

	alertName := req.GetFileName() + " " + req.GetDbProcName()

	res.RsId = req.RqId
	res.ShardingId = req.ShardingId

	if req.GetDbType() != 0 {
		fkalert.Alert(alertRequestDbIDNotZero, alertName)
		setError(1, "invliad request. not support read from db.")
		return errors.New("invliad request. not support read from db.")
	}
	// 查找缓存数据
	sheetCache := tb.getConfigCache(req.GetDbProcName())
	if sheetCache == nil {
		fkalert.Alert(alertSheetNotExists, alertName)
		setError(2, "request invalid sheet. not found sheet cache data.")
		return errors.New("request invalid sheet. not found sheet cache data.")
	}

	res.FieldCount = proto.Int32(int32(len(req.GetParamList())))

	// 检测MD5
	if sheetCache.Base.FileMd5 == req.GetMd5Check() {
		// 检查更新请求. 没有更新. 返回成功
		if req.GetStartIndex() == 0 {
			res.Md5Check = req.Md5Check
			excelRes.all = true // 拉取时缓存无更新直接返回，属于加载全部
			return nil
		}
	} else {
		excelRes.md5Changed = true
		// 如果md5值不一样. 检测开始索引. start index 还不是0. 可能是拉取数据过程中,excel文件更新了.
		// NOTE: 不要这个错误码(3). 客户端会特殊处理
		if req.GetStartIndex() != 0 {
			// 报错.
			setError(3, "maybe change when client pull config data.")
			return nil
		}
	}
	// 查找对应列的缓存数据
	configCache := sheetCache.getCacheData(req.GetParamList())
	if configCache == nil {
		fkalert.Alert(alertGetColumnDataFailed, alertName)
		setError(4, fmt.Sprintf("request invalid fields. has fields:%v md5:%s.",
			sheetCache.Base.Fields, sheetCache.Base.FileMd5))
		return nil
	}
	if len(req.GetParamList()) != len(configCache.Fields) {
		setError(5, fmt.Sprintf("get invalid data. fields num not match. has fields:%v.",
			configCache.Fields))
		return nil
	}
	res.FieldCount = proto.Int32(int32(len(configCache.Fields)))
	res.Md5Check = &sheetCache.Base.FileMd5
	res.GitVersion = &sheetCache.Base.Desc

	// 添加数据
	for i := int(req.GetStartIndex()); i < len(configCache.Data); i++ {
		row := &MysqlParam.DbRow{}
		row.Fields = configCache.Data[i]
		res.Results = append(res.Results, row)
		// pbSize += configCache.Size[i]
		// pbSize = proto.Size(res)
		// // 包大小限制
		// if req.GetSizeLimit() != 0 && pbSize >= int(req.GetSizeLimit()) {
		// 	// 一条数据大小就超出了包大小限制.
		// 	if len(res.Results) == 1 {
		// 		fkalert.Alert(alertDataTooBig, alertName)
		// 		setError(7, "too large config. one row data size out of limit.")
		// 		return nil
		// 	}
		// 	res.ReadAll = proto.Bool(false)
		// 	res.Results = res.Results[:len(res.Results)-1]
		// 	// pbSize -= configCache.Size[i]
		// 	pbSize = proto.Size(res)
		// 	break
		// }
		// // 条数限制
		// if req.GetPageSize() != 0 && i-int(req.GetStartIndex())+1 >= int(req.GetPageSize()) {
		// 	// 不是最后一条,就设置不是全部
		// 	if i+1 != len(configCache.Data) {
		// 		res.ReadAll = proto.Bool(false)
		// 	}
		// 	break
		// }
	}
	res.RowCount = proto.Int32(int32(len(res.GetResults())))

	excelRes.file = strings.Replace(sheetCache.File, flagConfigPath, "", -1)
	excelRes.sheet = sheetCache.Sheet
	excelRes.md5 = sheetCache.Base.Md5Sum
	excelRes.last = time.Now().Unix()
	excelRes.all = res.GetReadAll()
	excelRes.server = req.GetServerName()
	excelRes.sheetDataMd5 = sheetCache.Base.SheetDataMd5
	excelRes.fileMd5 = sheetCache.Base.FileMd5
	excelRes.desc = sheetCache.Base.Desc
	tb.rdCh <- excelRes

	// 回包
	logger.InfoWF("OnLoadCacheConfig load cache config success.", zap.Uint64("request_id", req.GetRqId()),
		zap.String("file", req.GetFileName()),
		zap.String("sheet", req.GetDbProcName()),
		zap.String("sheetDataMd5", excelRes.sheetDataMd5),
		zap.String("fileMd5", excelRes.fileMd5),
		zap.Any("desc", excelRes.desc),
		zap.Strings("Fields", req.GetParamList()),
		zap.String("rq_md5", req.GetMd5Check()), zap.String("rs_md5", res.GetMd5Check()),
		zap.Int("ret_count", len(res.GetResults())), zap.Uint32("start_index", req.GetStartIndex()),
		zap.String("git_version", res.GetGitVersion()), zap.String("server", req.GetServerName()))
	return nil
}
