package mysql

import (
	"context"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

type GormLogger struct {
	Logger        fklog.FKLogI
	LogLevel      logger.LogLevel
	SlowThreshold time.Duration
}

func NewGormLogger(logger fklog.FKLogI, level logger.LogLevel) *GormLogger {
	return &GormLogger{
		Logger:        logger,
		LogLevel:      level,
		SlowThreshold: 200 * time.Millisecond, // 可自定义慢查询阈值
	}
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.Logger.InfoWF(msg, zap.Any("data", data))
	}
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.Logger.WarnWF(msg, zap.Any("data", data))
	}
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.Logger.ErrorWF(msg, zap.Any("data", data))
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && l.LogLevel >= logger.Error:
		l.Logger.WarnWF("SQL error",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("elapsed", elapsed),
			// zap.String("file", utils.FileWithLineNum()),
			zap.Error(err),
		)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		l.Logger.WarnWF("SLOW SQL",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("elapsed", elapsed),
			// zap.String("file", utils.FileWithLineNum()),
		)
	case l.LogLevel >= logger.Info:
		l.Logger.WarnWF("SQL executed",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("elapsed", elapsed),
			// zap.String("file", utils.FileWithLineNum()),
		)
	}
}
