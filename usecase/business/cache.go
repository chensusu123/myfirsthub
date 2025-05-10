package business

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"sync"
	"time"

	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"

	"go.uber.org/zap"
)

// excell 缓存
type excelCache struct {
	file   *excellFileCache // 文件缓存
	logger fklog.FKLogI     // 日志接口
	Md5Sum string           // sheet fields md5
	Desc   string           // version desc
	Fields []string         // fields
	Data   [][]string       // config data
}

type excelFullCache struct {
	excelCache                  // normal cache
	FieldMap     map[string]int // file fields map
	FileMd5      string         // file md5
	SheetDataMd5 string
}

func (cache *excelFullCache) buildCache(sum string, fields []string) (data *excelCache) {
	logger := cache.logger
	// 检测数据列
	var fieldIDs []int
	failed := false
	for _, field := range fields {
		// 有的列在全部数据中不存在.抛弃
		if id, ok := cache.FieldMap[field]; ok {
			fieldIDs = append(fieldIDs, id)
			continue
		}
		failed = true
		logger.ErrorWF("buildCache not found field.", zap.String("file", cache.file.File),
			zap.String("sheet", cache.file.Sheet),
			zap.String("field", field))
	}
	if failed {
		logger.ErrorWF("buildCache not found field.", zap.String("file", cache.file.File),
			zap.String("sheet", cache.file.Sheet),
			zap.String("md5", cache.FileMd5))
		return
	}
	// 构建缓存
	data = &excelCache{}
	for _, row := range cache.Data {
		var tmp []string
		for _, index := range fieldIDs {
			tmp = append(tmp, row[index])
		}
		data.Data = append(data.Data, tmp)
	}
	data.logger = cache.logger
	data.Fields = fields
	data.Md5Sum = sum
	data.file = cache.file
	logger.InfoWF("buildCache build cache config.", zap.String("file", cache.file.File),
		zap.String("sheet", cache.file.Sheet),
		zap.String("md5", cache.FileMd5), zap.Strings("fields", fields))
	return
}

// excell 文件缓存
type excellFileCache struct {
	logger fklog.FKLogI    // log info
	File   string          // file name
	Sheet  string          // sheet name
	Base   *excelFullCache // all fields info
	Data   sync.Map        //map[string]*excelCache - custom fields info
}

// 更新缓存数据
func (cache *excellFileCache) updateExcellCache(newCache *excellFileCache) {
	cache.Data.Range(func(k, v interface{}) bool {
		data := v.(*excelCache)
		newItem := newCache.Base.buildCache(data.Md5Sum, data.Fields)
		if newItem != nil {
			newCache.Data.Store(newItem.Md5Sum, newItem)
		} else {
			cache.logger.ErrorWF("updateExcellCache rebuild cache item failed.",
				zap.String("file", cache.File), zap.String("sheet", cache.Sheet),
				zap.String("last_md5", cache.Base.FileMd5), zap.String("cur_md5", newCache.Base.FileMd5))
		}
		return true
	})
}

func (cache *excellFileCache) getCacheData(fields []string) *excelCache {
	md5sum := calcFieldsMd5(fields)
	if cache.Base.Md5Sum == md5sum {
		return &cache.Base.excelCache
	}
	cache.logger.WarnWF("getCacheData get file default sheet failed.", zap.String("file", cache.File),
		zap.Strings("fields", cache.Base.Fields), zap.Strings("get", fields),
		zap.String("src_md5", cache.Base.Md5Sum), zap.String("get_md5", md5sum))
	if val, ok := cache.Data.Load(md5sum); ok {
		return val.(*excelCache)
	}
	// build cache
	newItem := cache.Base.buildCache(md5sum, fields)
	if newItem != nil {
		cache.Data.Store(newItem.Md5Sum, newItem)
	}
	cache.logger.WarnWF("getCacheData build sheet cache.", zap.Bool("success", newItem != nil),
		zap.String("file", cache.File), zap.Strings("fields", cache.Base.Fields),
		zap.Strings("get", fields), zap.String("get_md5", md5sum))
	return newItem
}

// 获取文件缓存
func (tb *tCustomBusiness) getConfigCache(sheet string) *excellFileCache {
	v, ok := tb.configCache.Load(sheet)
	if !ok {
		return nil
	}
	cache, ok := v.(*excellFileCache)
	if !ok {
		return nil
	}
	return cache
}

// 保存/更新文件缓存
func (tb *tCustomBusiness) saveConfigCache(sheet string, v *excellFileCache) {
	tb.configCache.Store(sheet, v)
}

// 保存/更新文件缓存
func (tb *tCustomBusiness) saveConfigCacheForChange(sheet string, v *excellFileCache) {
	tb.configCacheForChange.Store(sheet, v)
}

func calcFieldsMd5(fields []string) string {
	buf := &bytes.Buffer{}
	for _, v := range fields {
		buf.WriteString(v)
	}
	return fmt.Sprintf("%x", md5.Sum(buf.Bytes()))
}

func calcSheetDataMd5(logger fklog.FKLogI, data [][]string) string {
	start := time.Now().UnixNano() / 1e6
	buf := &bytes.Buffer{}
	for _, v1 := range data {
		for _, v2 := range v1 {
			buf.WriteString(v2)
		}
	}
	res := fmt.Sprintf("%x", md5.Sum(buf.Bytes()))
	StatIns.AddSheetMd5CalNum()
	StatIns.AddSheetMd5CalCost(uint64(time.Now().UnixNano()/1e6 - start))
	return res
}
