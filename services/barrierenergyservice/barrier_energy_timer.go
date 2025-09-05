package barrierenergyservice

import (
	"context"
	"sync"
	"time"

	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// 审核版本用于自动恢复体力
var (
	mu      sync.Mutex
	userMap = make(map[uint64]*time.Timer)
)

// 启动用户自动恢复体力
func (s service) startUserRecoverEnergy(ctx context.Context, userId uint64, nextUpdateTime int64) {
	logger := fklog.ContextAppLogger(ctx)
	mu.Lock()
	defer mu.Unlock()

	// 若已有 timer，跳过
	if timer, ok := userMap[userId]; ok {
		// timer.Stop()
		_ = timer
		return
	}

	// 设置体力自动恢复时间
	recoverTime := nextUpdateTime - time.Now().Unix()
	timer := time.AfterFunc(time.Duration(recoverTime)*time.Second, func() {
		s.safeTimer(ctx, userId)
	})
	userMap[userId] = timer
	logger.CtxInfo(ctx, "startUserRecoverEnergy success", zap.Any("userId", userId), zap.Int64("recoverTime", recoverTime), zap.Any("timer", timer))
}

func (s service) safeTimer(ctx context.Context, userID uint64) {
	logger := fklog.AppLogger().Clone("barrierenergyservice")
	logger.SetUid(userID)
	ctx = fklog.ContextWithLogger(context.Background(), logger)
	tracer := otel.Tracer("barrierenergyservice")
	ctx, span := tracer.Start(ctx, "UserRecoverEnergyTimer")

	span.SetAttributes(
		attribute.Int64("enduser.id", int64(userID)),
	)
	span.AddEvent("UserRecoverEnergy")
	// logger := fklog.ContextAppLogger(ctx)
	defer func() {
		if r := recover(); r != nil {
			logger.CtxError(ctx, "handleRecoverUserEnergy panic.", zap.Any("r", r))
		}
		span.End()
	}()
	isOnline := online.IsOnline(userID)
	if !isOnline {
		span.AddEvent("stopUserRecoverTimer")
		s.stopUserRecoverTimer(ctx, userID)
		return
	}
	span.AddEvent("handleRecoverUserEnergy")
	s.handleRecoverUserEnergy(ctx, userID)
}

// 自动恢复体力
func (s service) handleRecoverUserEnergy(ctx context.Context, userID uint64) {
	logger := fklog.ContextAppLogger(ctx)
	defer func() {
		// 重新设置timer
		mu.Lock()
		defer mu.Unlock()
		if old, ok := userMap[userID]; ok {
			old.Stop()
		}

		nextTriggerTime := GetEnergyRecoverCfg()
		timer := time.AfterFunc(time.Duration(nextTriggerTime)*time.Second, func() {
			s.safeTimer(ctx, userID)
		})
		userMap[userID] = timer
		logger.CtxInfo(ctx, "handleRecoverUserEnergy add timer", zap.Any("userID", userID), zap.Int64("nextTriggerTime", nextTriggerTime), zap.Any("timer", timer))
	}()

	curEnergy, nextTime, err := s.calEnergy(ctx, userID)
	if err != nil {
		logger.CtxInfo(ctx, "handleRecoverUserEnergy calEnergy failed", zap.Any("userID", userID), zap.Error(err))
		return
	}

	err = s.SendEnergyChgPack(ctx, userID, curEnergy, nextTime)
	if err != nil {
		err = nil
		logger.CtxError(ctx, "handleRecoverUserEnergy SendEnergyChgPack fail", zap.Error(err), zap.Any("userID", userID), zap.Int32("curEnergy", curEnergy), zap.Int64("nextTime", nextTime))
		// return
	}
}

func (s service) stopUserRecoverTimer(ctx context.Context, userID uint64) {
	mu.Lock()
	defer mu.Unlock()
	logger := fklog.ContextAppLogger(ctx)
	if timer, ok := userMap[userID]; ok {
		timer.Stop()
		delete(userMap, userID)
		logger.CtxInfo(ctx, "StopUserRecoverTimer success", zap.Uint64("userId", userID), zap.Any("timer", timer))
	}
}
