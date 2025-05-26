package business

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig/param"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager/loadconfigapi"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil/filemonitor"
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

// 监控文件路径
var (
	flagConfigPath        string
	flagConfigHistoryPath string
)

// 配表拉取消息缓冲区长度
var bufferLen int

// 定时同步变化缓存，单位ms，最小时间500ms,默认1000ms
var cacheApplyCheck int

// 是否记录程序统计信息(stat info)，1:打开 0:关闭
var statOpen int

var buffer chan *LoadConfigMsg

func SetConfigPath(dir string) {
	flagConfigPath = dir
}

func init() {
	param.StringP(&flagConfigPath, "config:path", "./conf.d/data/")
	param.StringP(&flagConfigHistoryPath, "config:history:path", "./conf.d/history/")
	param.IntP(&cacheApplyCheck, "cache:apply:check", 1000, "配置同步时间间隔.单位毫秒")
	param.IntP(&bufferLen, "buffer:len", 5000, "配表拉取消息缓冲区长度")
	param.IntP(&statOpen, "stat:open", 0, "是否打开程序统计")

	loadconfigapi.SetLoadConfigFunc(GCustomBusiness.LoadCacheConfig)
	loadconfigapi.SetInitConfigCacheFunc(GCustomBusiness.Init)
}

type excelReadResult struct {
	sheet, file  string // excel sheet,file
	server       string // server name
	md5          string // file md5
	last         int64  // last read time
	all          bool   // 是否全部加载完
	sheetDataMd5 string
	fileMd5      string
	desc         string // git_version_for_config.txt文件中的版本信息
	md5Changed   bool   // 配表数据是否发生变化
}

const (
	StatusNoCache int32 = 0 // 无缓存
	StatusCaching int32 = 1 // 缓存中
	StatusCached  int32 = 2 // 缓存完成
)

type readInfo struct {
	time map[string]time.Time // key:最后读取时间
	lock sync.RWMutex         // 读写锁
}

type tCustomBusiness struct {
	closeMonitor  func()                // 关闭文件监控
	gitFile       string                // git版本文件
	gitCfgVersion atomic.String         // 文件版本
	configCache   sync.Map              // 所有文件缓存数据 map[string]*excellFileCache
	rdCh          chan *excelReadResult // 读取通道
	loadCh        chan *excelReadResult // 加载通道
	serverData    sync.Map              // server:readInfo
	fileData      sync.Map              // sheet:readInfo
	loadTime      sync.Map              // file:int64

	// 配置缓存加载控制逻辑
	isFirstLoad          bool     // 程序启动特殊逻辑控制
	isSameDir            bool     // 线上配置目录与发布目录是否一致
	configCacheForChange sync.Map // 非第一次加载变化做预缓存
	cacheForChangeStatus int32    // 变化缓存状态，0：变化缓存已消费，无缓存；1：缓存进行中；2：缓存完成，待消费
}

// FKServiceI 服务接口
func (tb *tCustomBusiness) Name() string {
	return "ConfigCacheServer"
}

func (tb *tCustomBusiness) OnInit(logger fklog.FKLogI, cfg fkconfig.FkConfigerI) (err error) {
	return nil
	// 如果指定没有文件,就递归从上层目录查找.如果到根目录还查找不到.就直接使用输入的目录
	gitPath := flagConfigPath
	for !filemonitor.IsFileExists(filepath.Clean(gitPath + "/" + versionFile)) {
		gitPath = filepath.Dir(gitPath)
		if len(gitPath) < 2 {
			gitPath = flagConfigPath
			break
		}
	}
	tb.gitFile = filepath.Clean(gitPath + "/" + versionFile)
	tb.rdCh = make(chan *excelReadResult, 1024)
	tb.loadCh = make(chan *excelReadResult, 1024)
	go tb.dealReadFile(logger)

	// 服务启动，第一次加载标记设置为true
	tb.isFirstLoad = true
	tb.isSameDir = true
	tb.cacheForChangeStatus = StatusNoCache

	if cacheApplyCheck < 500 {
		cacheApplyCheck = 500
	}

	go tb.applyChangeForCache(logger)

	// 业务服务远程读配表缓存消息打点
	buffer = make(chan *LoadConfigMsg, bufferLen)
	go func() {
		for msg := range buffer {
			// loadConfigKafka.PushLoadCfgInfo(logger, msg)
			_ = msg
		}
	}()

	// 检测线上配置目录和发布目录是否一致,默认一致
	tb.isSameDir = true
	readDir := ""

	// usingVersion, err := ConfigCacheInfo.GetCurrentUsingGitVersion(logger, int32(fkconfig.GetServerConfig().GroupID))
	// if err != nil {
	// 	logger.ErrorWF("OnInit GetCurrentUsingGitVersion failed")
	// 	return
	// }

	// logger.WarnWF("OnInit GetCurrentUsingGitVersion", zap.String("usingVersion", usingVersion))
	// 如果当前使用版本存在
	// if len(usingVersion) > 0 {
	// 	logger.WarnWF("OnInit current using version is ", zap.String("version", usingVersion))
	// 	//判断是否与当前current目录下version一致
	// 	data, err := ioutil.ReadFile(tb.gitFile)
	// 	if err != nil {
	// 		logger.ErrorWF("OnInit read git file failed", zap.String("filename", tb.gitFile))
	// 		return err
	// 	}
	// 	//fileMd5 := fmt.Sprintf("%x", md5.Sum(data))

	// 	//如果不相等
	// 	if string(data) != usingVersion {
	// 		logger.WarnWF("OnInit using dir is not same deploy dir", zap.String("usingVersion", usingVersion))
	// 		//遍历查找git文件md5等于usingVersion的文件夹路径，如果找不到，直接读取最新发布目录
	// 		dir, found, git_version_dir := tb.GetUsingConfigDir(logger, usingVersion)
	// 		//如果找到了，则使用目录和发布目录不一致,其余情况都直接读取最新发布目录
	// 		if found {
	// 			tb.isSameDir = false
	// 			readDir = dir + "config" + "/" + "svn" + "/"
	// 			tb.gitFile = git_version_dir
	// 			logger.WarnWF("OnInit get using config dir ", zap.String("readDir", readDir), zap.String("git_version_dir", git_version_dir))
	// 		}
	// 	} else {
	// 		logger.WarnWF("OnInit using dir is same as deploy dir", zap.String("usingVersion", usingVersion))
	// 	}
	// }

	data, err := ioutil.ReadFile(tb.gitFile)
	if err != nil {
		logger.ErrorWF("OnInit read git file failed", zap.String("filename", tb.gitFile))
		return err
	}
	// err = ConfigCacheInfo.SetCurrentServerStatusInfo(logger, context.TODO(), StatusNoCache, int32(fkconfig.GetServerConfig().ServerTypeID),
	// 	int32(fkconfig.GetServerConfig().ShardingID), int32(fkconfig.GetServerConfig().ServerID), int32(fkconfig.GetServerConfig().GroupID), fkconfig.GetServerConfig().Address, string(data), nil)
	// if err != nil {
	// 	logger.ErrorWF("OnInit SetCurrentServerStatusInfo failed", zap.Error(err), zap.Int32("status", StatusNoCache))
	// }
	_ = data
	// 监控文件变化
	tb.closeMonitor, err = filemonitor.MonitorChangeForCacheServer(logger, ".xlsx",
		tb, tb.isSameDir, flagConfigPath, readDir)
	if err != nil {
		logger.ErrorWF("OnInit monitor file path failed.", zap.String("path", flagConfigPath), zap.String("version", tb.gitFile),
			zap.Error(err))
		return
	}
	logger.InfoWF("OnInit monitor file path success.", zap.String("path", flagConfigPath), zap.String("version", tb.gitFile))
	return
}

func (tb *tCustomBusiness) OnStart(logger fklog.FKLogI, cfg fkconfig.FkConfigerI) (err error) {
	return
}

func (tb *tCustomBusiness) OnStop(logger fklog.FKLogI) (err error) {
	// 关闭文件监控
	tb.closeMonitor()
	return
}

func (tb *tCustomBusiness) OnFinish(logger fklog.FKLogI) (err error) {
	// 关闭buffer channel
	close(buffer)
	return
}

func (tb *tCustomBusiness) dealReadFile(logger fklog.FKLogI) {
	ctx := fkserver.GetContext()
	var val interface{}
	var ok bool
	var res *excelReadResult
	var ri *readInfo
	for {
		select {
		case <-ctx.Done():
			logger.WarnWF("dealReadFile exit file history.")
			return
		case res = <-tb.rdCh:
			val, ok = tb.serverData.Load(res.server)
			if !ok {
				ri = &readInfo{time: make(map[string]time.Time)}
				tb.serverData.Store(res.server, ri)
			} else {
				ri = val.(*readInfo)
			}
			ri.lock.Lock()
			ri.time[res.sheet] = time.Unix(res.last, 0)
			ri.lock.Unlock()

			val, ok = tb.fileData.Load(res.sheet)
			if !ok {
				ri = &readInfo{time: make(map[string]time.Time)}
				tb.fileData.Store(res.sheet, ri)
			} else {
				ri = val.(*readInfo)
			}
			ri.lock.Lock()
			ri.time[res.server] = time.Unix(res.last, 0)
			ri.lock.Unlock()
		case res = <-tb.loadCh:
			tb.loadTime.Store(res.file, res)
		}
	}
}

func (tb *tCustomBusiness) applyChangeForCache(logger fklog.FKLogI) {
	logger.WarnWF("applyChangeForCache  begin")
	timer := time.NewTicker(time.Duration(cacheApplyCheck) * time.Millisecond)
	defer func() {
		timer.Stop()
		logger.WarnWF("applyChangeForCache  end")
	}()
	ctx := fkserver.GetContext()
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			// 如果缓存不是处于已完成状态，直接返回
			if tb.cacheForChangeStatus != StatusCached {
				// logger.WarnWF("cacheForChangeStatus not equal StatusCached", zap.Int32("cacheForChangeStatus", tb.cacheForChangeStatus))
				continue
			}

			// 如果同步标记没有打开，直接返回
			// canApply, err := ConfigCacheInfo.IsCanApplyCacheChange(logger)
			// if canApply == false || err != nil {
			// 	logger.WarnWF("applyChangeForCache IsCanApplyCacheChange false", zap.Bool("canApply", canApply), zap.Error(err))
			// 	continue
			// }

			// 应用缓存
			tb.configCacheForChange.Range(func(key, value interface{}) bool {
				tb.configCache.Store(key, value)
				return true
			})
			tb.configCacheForChange.Range(func(key, value interface{}) bool {
				tb.configCacheForChange.Delete(key)
				return true
			})

			gitPath := flagConfigPath
			for !filemonitor.IsFileExists(filepath.Clean(gitPath + "/" + versionFile)) {
				gitPath = filepath.Dir(gitPath)
				if len(gitPath) < 2 {
					gitPath = flagConfigPath
					break
				}
			}
			tb.gitFile = filepath.Clean(gitPath + "/" + versionFile)

			data, err := ioutil.ReadFile(tb.gitFile)
			if err != nil {
				logger.ErrorWF("applyChangeForCache read config git version failed.", zap.String("gitCfg", tb.gitFile))
				gitPath := filepath.Dir(tb.gitFile)
				for !filemonitor.IsFileExists(filepath.Clean(gitPath + "/" + versionFile)) {
					gitPath = filepath.Dir(gitPath)
					if len(gitPath) < 2 {
						gitPath = ""
						break
					}
				}
				if gitPath != "" {
					newFile := filepath.Clean(gitPath + "/" + versionFile)
					data, err = ioutil.ReadFile(tb.gitFile)
					logger.InfoWF("applyChangeForCache retry find and read config git version.", zap.String("gitCfg", tb.gitFile),
						zap.String("file", newFile), zap.String("git_version", string(data)))
					tb.gitFile = newFile
				}
			}

			cfgVersion := string(data)
			tb.gitCfgVersion.Store(cfgVersion)

			tb.configCache.Range(func(k, v interface{}) bool {
				cache := v.(*excellFileCache)
				cache.Base.Desc = cfgVersion
				cache.Data.Range(func(k, v2 interface{}) bool {
					item := v2.(*excelCache)
					item.Desc = cfgVersion
					return true
				})
				return true
			})
			// fileMd5 := fmt.Sprintf("%x", md5.Sum([]byte(cfgVersion)))

			// err = ConfigCacheInfo.SetCurrentUsingGitVersion(logger, cfgVersion, int32(fkconfig.GetServerConfig().GroupID))
			// if err != nil {
			// 	logger.ErrorWF("applyChangeForCache SetCurrentUsingGitVersion failed", zap.String("version", cfgVersion), zap.Error(err))
			// }
			logger.WarnWF("applyChangeForCache SetCurrentUsingGitVersion success", zap.String("version", cfgVersion))

			tb.cacheForChangeStatus = StatusNoCache
			filemonitor.ApplyFileCache(logger)
			// err = ConfigCacheInfo.SetCurrentServerStatusInfo(logger, ctx, StatusNoCache, int32(fkconfig.GetServerConfig().ServerTypeID),
			// 	int32(fkconfig.GetServerConfig().ShardingID), int32(fkconfig.GetServerConfig().ServerID), int32(fkconfig.GetServerConfig().GroupID), fkconfig.GetServerConfig().Address, cfgVersion, nil)
			// if err != nil {
			// 	logger.ErrorWF("applyChangeForCache SetCurrentServerStatusInfo failed", zap.Error(err), zap.Int32("status", StatusNoCache))
			// }
			logger.WarnWF("applyChangeForCache SetCurrentServerStatusInfo success", zap.Int32("status", StatusNoCache))
		}
	}
}

func (tb *tCustomBusiness) GetUsingConfigDir(logger fklog.FKLogI, usingVersion string) (configDir string, found bool, git_version_file string) {
	configDir = ""
	found = false
	git_version_file = ""

	s, err := os.Stat(flagConfigHistoryPath)
	if err != nil {
		return
	}
	if !s.IsDir() {
		return
	}

	rd, err := ioutil.ReadDir(flagConfigHistoryPath)
	if err != nil {
		return
	}
	sort.Slice(rd, func(i, j int) bool { return rd[i].ModTime().Unix() > rd[j].ModTime().Unix() })
	for _, fi := range rd {
		if fi.IsDir() {
			historyDir := flagConfigHistoryPath + fi.Name() + "/"
			historyVersionFile := historyDir + versionFile
			logger.WarnWF("GetUsingConfigDir get version file", zap.String("historyDir", historyDir), zap.String("historyVersionFile", historyVersionFile))

			data, err := ioutil.ReadFile(historyVersionFile)
			if err != nil {
				logger.ErrorWF("GetUsingConfigDir read git file failed", zap.String("filename", historyVersionFile))
				continue
			}

			// fileMd5 := fmt.Sprintf("%x", md5.Sum(data))
			// 如果不相等
			if string(data) == usingVersion {
				configDir = historyDir
				found = true
				git_version_file = historyVersionFile
				logger.WarnWF("GetUsingConfigDir get using version dir", zap.String("configDir", configDir), zap.String("git_version_file", git_version_file))
				return
			}
		}
	}

	return
}

func (tb *tCustomBusiness) GetAllSheet() (info *AllSheetsInfo) {
	info = &AllSheetsInfo{}
	tb.configCache.Range(func(key, value interface{}) bool {
		excelInfo := value.(*excellFileCache)
		info.SheetsInfo = append(info.SheetsInfo, &SheetInfo{SheetName: excelInfo.Sheet, FileName: filepath.Base(excelInfo.File)})
		return true
	})
	return
}

type SheetInfo struct {
	SheetName string `json:"sheet_name"`
	FileName  string `json:"file_name"`
}

type AllSheetsInfo struct {
	SheetsInfo []*SheetInfo `json:"sheets_info"`
}

// GCustomBusiness 自定义服务逻辑
var GCustomBusiness tCustomBusiness

func (tb *tCustomBusiness) Init(logger fklog.FKLogI) (err error) {
	// 如果指定没有文件,就递归从上层目录查找.如果到根目录还查找不到.就直接使用输入的目录
	gitPath := flagConfigPath
	for !filemonitor.IsFileExists(filepath.Clean(gitPath + "/" + versionFile)) {
		gitPath = filepath.Dir(gitPath)
		if len(gitPath) < 2 {
			gitPath = flagConfigPath
			break
		}
	}
	tb.gitFile = filepath.Clean(gitPath + "/" + versionFile)
	tb.rdCh = make(chan *excelReadResult, 1024)
	tb.loadCh = make(chan *excelReadResult, 1024)
	go tb.dealReadFile(logger)

	// 服务启动，第一次加载标记设置为true
	tb.isFirstLoad = true
	tb.isSameDir = true
	tb.cacheForChangeStatus = StatusNoCache

	if cacheApplyCheck < 500 {
		cacheApplyCheck = 500
	}

	go tb.applyChangeForCache(logger)

	// 业务服务远程读配表缓存消息打点
	buffer = make(chan *LoadConfigMsg, bufferLen)
	go func() {
		for msg := range buffer {
			// loadConfigKafka.PushLoadCfgInfo(logger, msg)
			_ = msg
		}
	}()

	// 检测线上配置目录和发布目录是否一致,默认一致
	tb.isSameDir = true
	readDir := ""

	// usingVersion, err := ConfigCacheInfo.GetCurrentUsingGitVersion(logger, int32(fkconfig.GetServerConfig().GroupID))
	// if err != nil {
	// 	logger.ErrorWF("OnInit GetCurrentUsingGitVersion failed")
	// 	return
	// }

	// logger.WarnWF("OnInit GetCurrentUsingGitVersion", zap.String("usingVersion", usingVersion))
	// 如果当前使用版本存在
	// if len(usingVersion) > 0 {
	// 	logger.WarnWF("OnInit current using version is ", zap.String("version", usingVersion))
	// 	//判断是否与当前current目录下version一致
	// 	data, err := ioutil.ReadFile(tb.gitFile)
	// 	if err != nil {
	// 		logger.ErrorWF("OnInit read git file failed", zap.String("filename", tb.gitFile))
	// 		return err
	// 	}
	// 	//fileMd5 := fmt.Sprintf("%x", md5.Sum(data))

	// 	//如果不相等
	// 	if string(data) != usingVersion {
	// 		logger.WarnWF("OnInit using dir is not same deploy dir", zap.String("usingVersion", usingVersion))
	// 		//遍历查找git文件md5等于usingVersion的文件夹路径，如果找不到，直接读取最新发布目录
	// 		dir, found, git_version_dir := tb.GetUsingConfigDir(logger, usingVersion)
	// 		//如果找到了，则使用目录和发布目录不一致,其余情况都直接读取最新发布目录
	// 		if found {
	// 			tb.isSameDir = false
	// 			readDir = dir + "config" + "/" + "svn" + "/"
	// 			tb.gitFile = git_version_dir
	// 			logger.WarnWF("OnInit get using config dir ", zap.String("readDir", readDir), zap.String("git_version_dir", git_version_dir))
	// 		}
	// 	} else {
	// 		logger.WarnWF("OnInit using dir is same as deploy dir", zap.String("usingVersion", usingVersion))
	// 	}
	// }

	data, err := ioutil.ReadFile(tb.gitFile)
	if err != nil {
		logger.ErrorWF("OnInit read git file failed", zap.String("filename", tb.gitFile))
		return err
	}
	// err = ConfigCacheInfo.SetCurrentServerStatusInfo(logger, context.TODO(), StatusNoCache, int32(fkconfig.GetServerConfig().ServerTypeID),
	// 	int32(fkconfig.GetServerConfig().ShardingID), int32(fkconfig.GetServerConfig().ServerID), int32(fkconfig.GetServerConfig().GroupID), fkconfig.GetServerConfig().Address, string(data), nil)
	// if err != nil {
	// 	logger.ErrorWF("OnInit SetCurrentServerStatusInfo failed", zap.Error(err), zap.Int32("status", StatusNoCache))
	// }
	_ = data
	// 监控文件变化
	tb.closeMonitor, err = filemonitor.MonitorChangeForCacheServer(logger, ".xlsx",
		tb, tb.isSameDir, flagConfigPath, readDir)
	if err != nil {
		logger.ErrorWF("OnInit monitor file path failed.", zap.String("path", flagConfigPath), zap.String("version", tb.gitFile),
			zap.Error(err))
		return
	}
	logger.InfoWF("OnInit monitor file path success.", zap.String("path", flagConfigPath), zap.String("version", tb.gitFile))

	return
}
