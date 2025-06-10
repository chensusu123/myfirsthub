package business

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkalert"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil/filemonitor"

	"github.com/xuri/excelize/v2"

	"go.uber.org/zap"
)

// versionFile custom git version file
var versionFile string = "git_version_for_config.txt"

func (tb *tCustomBusiness) OnPrev(logger fklog.FKLogI, file string) {
	//不是同一个目录，主动先读取一次，不用走首次逻辑,首次也会回调
	//if tb.isSameDir == false {
	//	tb.isFirstLoad = false
	//}
	//如果不是第一次加载，设置预缓存状态为1，正在进行缓存
	if !tb.isFirstLoad {
		// 清除已有缓存
		tb.configCacheForChange.Range(func(key, value interface{}) bool {
			tb.configCacheForChange.Delete(key)
			return true
		})
		filemonitor.ClearFileCache(logger)
		tb.cacheForChangeStatus = StatusCaching
		// cfgVersion := tb.gitCfgVersion.Load()
		// deployVersion := tb.getDeployGitVersion(logger)
		// if deployVersion == cfgVersion {
		// 	err := ConfigCacheInfo.SetCurrentServerStatusInfo(logger, context.TODO(), StatusNoCache, int32(fkconfig.GetServerConfig().ServerTypeID),
		// 		int32(fkconfig.GetServerConfig().ShardingID), int32(fkconfig.GetServerConfig().ServerID), int32(fkconfig.GetServerConfig().GroupID), fkconfig.GetServerConfig().Address, deployVersion, nil)
		// 	if err != nil {
		// 		logger.ErrorWF("OnPrev SetCurrentServerStatusInfo failed", zap.Error(err), zap.Int32("status", StatusNoCache))
		// 	}
		// 	logger.WarnWF("OnPrev version is same, no need to cache", zap.String("cfgVersion", cfgVersion))
		// 	return
		// }
		// err := ConfigCacheInfo.SetCurrentServerStatusInfo(logger, fkserver.GetContext(), StatusCaching, int32(fkconfig.GetServerConfig().ServerTypeID),
		// 	int32(fkconfig.GetServerConfig().ShardingID), int32(fkconfig.GetServerConfig().ServerID), int32(fkconfig.GetServerConfig().GroupID), fkconfig.GetServerConfig().Address, deployVersion, nil)
		// if err != nil {
		// 	logger.ErrorWF("OnPrev SetCurrentServerStatusInfo failed", zap.Error(err), zap.Int32("status", StatusCaching))
		// }

		logger.WarnWF("OnPrev SetCurrentServerStatusInfo success", zap.Int32("status", StatusCaching))
	} else {
		data, err := ioutil.ReadFile(tb.gitFile)
		if err != nil {
			logger.ErrorWF("OnPrev read config git version failed.", zap.String("gitCfg", tb.gitFile), zap.String("file", file))
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
				logger.InfoWF("OnPrev retry find and read config git version.", zap.String("gitCfg", tb.gitFile),
					zap.String("file", newFile), zap.String("git_version", string(data)))
				tb.gitFile = newFile
			}
		}
		tb.gitCfgVersion.Store(string(data))
		logger.DebugWF("OnPrev ready to update git config version.", zap.String("git_version", string(data)), zap.String("file", tb.gitFile))
	}
	return
}

func (tb *tCustomBusiness) OnEnd(logger fklog.FKLogI, file string) {
	// 非第一次加载，设置预加载状态为完成，待消费
	if !tb.isFirstLoad {
		deployGitVersion := tb.getDeployGitVersion(logger)
		// cfgVersion := tb.gitCfgVersion.Load()

		// if cfgVersion == deployGitVersion {
		// 	err := ConfigCacheInfo.SetCurrentServerStatusInfo(logger, context.TODO(), StatusNoCache, int32(fkconfig.GetServerConfig().ServerTypeID),
		// 		int32(fkconfig.GetServerConfig().ShardingID), int32(fkconfig.GetServerConfig().ServerID), int32(fkconfig.GetServerConfig().GroupID), fkconfig.GetServerConfig().Address, deployGitVersion, nil)
		// 	if err != nil {
		// 		logger.ErrorWF("OnEnd SetCurrentServerStatusInfo failed", zap.Error(err), zap.Int32("status", StatusNoCache))
		// 	}
		// 	logger.WarnWF("OnEnd version is same, no need to cache")
		// 	return
		// }

		changeDesc := make([]string, 0)
		tb.configCacheForChange.Range(func(key, value interface{}) bool {
			cache := value.(*excellFileCache)
			sheet := key.(string)
			desc := filepath.Base(cache.File) + ":"
			desc += sheet + ","
			oldCacheTmp, ok := tb.configCache.Load(key)
			// 使用缓存中没有
			if !ok {
				desc += "新增表,字段:"
				newField := cache.Base.FieldMap
				for key := range newField {
					desc += key + "  "
				}
				changeDesc = append(changeDesc, desc)
				return true
			} else {
				newField := cache.Base.FieldMap
				oldCache := oldCacheTmp.(*excellFileCache)
				oldField := oldCache.Base.FieldMap

				addField := "新增字段:"
				hasAddField := false
				for key := range newField {
					if _, ok := oldField[key]; !ok {
						addField += key + "  "
						hasAddField = true
					}
				}
				if hasAddField {
					desc += addField
				}

				delField := ",删除字段:"
				hasDelField := false
				for key := range oldField {
					if _, ok := newField[key]; !ok {
						delField += key + " "
						hasDelField = true
					}
				}
				if hasDelField {
					desc += delField
				}
				if hasAddField || hasDelField {
					changeDesc = append(changeDesc, desc)
				}
				return true
			}
		})

		// err := ConfigCacheInfo.SetCurrentServerStatusInfo(logger, fkserver.GetContext(), StatusCached, int32(fkconfig.GetServerConfig().ServerTypeID),
		// 	int32(fkconfig.GetServerConfig().ShardingID), int32(fkconfig.GetServerConfig().ServerID), int32(fkconfig.GetServerConfig().GroupID), fkconfig.GetServerConfig().Address, deployGitVersion, changeDesc)
		// if err != nil {
		// 	logger.ErrorWF("OnEnd SetCurrentServerStatusInfo failed", zap.Error(err), zap.Int32("status", StatusCached))
		// }
		tb.cacheForChangeStatus = StatusCached
		logger.WarnWF("OnEnd SetCurrentServerStatusInfo success", zap.Int32("status", StatusCached), zap.String("git_version", deployGitVersion), zap.Any("change", changeDesc))
	} else {
		cfgVersion := tb.gitCfgVersion.Load()

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
		// err := ConfigCacheInfo.SetCurrentUsingGitVersion(logger, cfgVersion, int32(fkconfig.GetServerConfig().GroupID))
		// if err != nil {
		// 	logger.ErrorWF("OnEnd SetCurrentUsingGitVersion failed", zap.String("version", cfgVersion), zap.Error(err))
		// }
		logger.WarnWF("OnEnd SetCurrentUsingGitVersion success", zap.String("version", cfgVersion))
	}
	tb.configCache.Range(func(key, value interface{}) bool {
		sheet := key.(string)
		logger.WarnWF("OnEnd current cache file", zap.String("sheet", sheet))
		return true
	})
	// 只要加载成功过一次，就不能直接加载
	tb.isFirstLoad = false
}

// MonitorFileModify 监听文件变动
func (tb *tCustomBusiness) OnChange(logger fklog.FKLogI, file string, data []byte) error {
	// 监控解析excel文件时间
	// endf := monitorLoadFile.Start()
	// defer endf()

	cfgVersion := tb.gitCfgVersion.Load()
	xlsxFile, err := excelize.OpenFile(file)
	if err != nil {
		fkalert.Alert(alertOpenExcelFailed, file)
		logger.ErrorWF("OnChange open excel by binary failed.", zap.String("file", file), zap.Int("len", len(data)), zap.Error(err))
		return err
	}
	fileMd5 := fmt.Sprintf("%x", md5.Sum(data))
	// 编译所有sheet
	for _, sheetSrcName := range xlsxFile.GetSheetMap() {
		sheetName, ok := checkSheetName(sheetSrcName)
		logger.DebugWF("OnChange recv file.", zap.String("file", file), zap.String("md5", fileMd5),
			zap.String("src_sheet", sheetSrcName), zap.String("sheet", sheetName))
		if !ok {
			logger.WarnWF("OnChange ignore sheet. check sheet name.", zap.String("file", file), zap.String("md5", fileMd5),
				zap.String("src_sheet", sheetSrcName), zap.String("sheet", sheetName))
			continue
		}
		rows, err := xlsxFile.GetRows(sheetSrcName)
		if err != nil {
			logger.WarnWF("OnChange ignore sheet. get rows data.", zap.String("file", file), zap.String("md5", fileMd5),
				zap.String("src_sheet", sheetSrcName), zap.String("sheet", sheetName))
			continue
		}

		full, err := parseExcellSheet(logger, file, sheetSrcName, rows)
		if err != nil || full == nil {
			logger.ErrorWF("OnChange load sheet data failed.", zap.String("file", file), zap.String("md5", fileMd5), zap.String("src", sheetSrcName),
				zap.String("sheet", sheetName), zap.Bool("full", full == nil), zap.Error(err))
			continue
		}
		full.FileMd5 = fileMd5
		full.SheetDataMd5 = calcSheetDataMd5(logger, full.Data)
		full.Md5Sum = calcFieldsMd5(full.Fields)
		newCache := &excellFileCache{}
		newCache.Base = full
		newCache.File = file
		newCache.Sheet = sheetSrcName
		newCache.logger = logger
		full.file = newCache
		// last := getConfigCache(sheetName)
		// if last != nil {
		// 	last.updateExcellCache(newCache)
		// }
		full.Desc = cfgVersion

		// 第一次加载，直接更新缓存，否则预存储
		if tb.isFirstLoad {
			tb.saveConfigCache(sheetName, newCache)
		} else {
			tb.saveConfigCacheForChange(sheetName, newCache)
		}

		logger.DebugWF("OnChange update config.", zap.String("file", file),
			zap.String("src_sheet", sheetSrcName), zap.String("sheet", sheetName),
			zap.String("fileMd5", full.FileMd5), zap.String("sheetDataMd5", full.SheetDataMd5), zap.Any("desc", full.Desc),
			zap.Strings("fields", full.Fields), zap.String("git_version", cfgVersion), zap.Any("rows", len(rows)))
	}
	// 加载记录
	res := &excelReadResult{}
	res.file = strings.Replace(file, flagConfigPath, "", -1)
	res.md5 = fileMd5
	res.last = time.Now().Unix()
	tb.loadCh <- res
	return nil
}

func (tb *tCustomBusiness) IsNeedCache(logger fklog.FKLogI, fileName string) (needCache bool) {
	if tb.isFirstLoad == true {
		needCache = false
	} else {
		needCache = true
	}
	return needCache
}

func (tb *tCustomBusiness) getDeployGitVersion(logger fklog.FKLogI) (git_version string) {
	gitPath := flagConfigPath
	for !filemonitor.IsFileExists(filepath.Clean(gitPath + "/" + versionFile)) {
		gitPath = filepath.Dir(gitPath)
		if len(gitPath) < 2 {
			gitPath = flagConfigPath
			break
		}
	}
	gitFile := filepath.Clean(gitPath + "/" + versionFile)

	data, err := ioutil.ReadFile(gitFile)
	if err != nil {
		logger.ErrorWF("getDeployGitVersion read config git version failed.", zap.String("gitCfg", gitFile))
		gitPath := filepath.Dir(gitFile)
		for !filemonitor.IsFileExists(filepath.Clean(gitPath + "/" + versionFile)) {
			gitPath = filepath.Dir(gitPath)
			if len(gitPath) < 2 {
				gitPath = ""
				break
			}
		}
		if gitPath != "" {
			newFile := filepath.Clean(gitPath + "/" + versionFile)
			data, err = ioutil.ReadFile(gitFile)
			logger.InfoWF("getDeployGitVersion retry find and read config git version.", zap.String("gitCfg", tb.gitFile),
				zap.String("file", newFile), zap.String("git_version", string(data)))
			gitFile = newFile
		}
	}

	return string(data)
}
