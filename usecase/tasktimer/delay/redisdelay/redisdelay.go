package redisdelay

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/usecase/tasktimer/delay/redisclient"
)

// define bucket ticker
type BucketTicker struct {
	Ticker       *time.Ticker
	Interval     time.Duration
	Name         string
	CallbackFunc func(interface{}) bool
	fklog.FKLogI
}

// define task
type Task struct {
	Id        string        // task id global uniqueness
	Data      string        // data of task
	Delay     time.Duration // delay time, 30 means after 30 second
	Timestamp int
}

// new ticker
func New(logger fklog.FKLogI, interval time.Duration, bucketName string, redisAddr string, callbackFunc func(interface{}) bool) (*BucketTicker, error) {
	if interval <= 0 || callbackFunc == nil {
		return nil, errors.New("create bucket ticker instance fail")
	}
	redisclient.Init(redisAddr)
	bucket := &BucketTicker{
		Interval:     interval,
		Name:         bucketName,
		CallbackFunc: callbackFunc,
		FKLogI:       logger.Clone(bucketName),
	}
	return bucket, nil
}

// add task
func (bucket *BucketTicker) AddTask(task *Task) error {
	// task id and delay time in redis zset
	timestamp := time.Now().Add(task.Delay).Unix()
	bucket.InfoWF("add task",
		zap.String("taskId", task.Id),
		zap.Int64("timestamp", timestamp),
	)
	err := redisclient.ZAdd(bucket.Name, int(timestamp), task.Id)
	if err != nil {
		return err
	}
	// task body in redis string
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	err = redisclient.Set(task.Id, string(data))
	if err != nil {
		return err
	}
	return nil
}

// start ticker
func (bucket *BucketTicker) Start(logger fklog.FKLogI, ctx context.Context) {
	timer := time.NewTicker(bucket.Interval) // interval
	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.WarnWF("bucket ticker done")
				return
			case t := <-timer.C:
				bucket.tickHandler(t, bucket.Name)
			}
		}
	}()
}

// tick handler
func (bucket *BucketTicker) tickHandler(currentTime time.Time, bucketName string) {
	for {
		task, err := getTask(bucketName)
		if err != nil {
			bucket.ErrorWF("tickHandler get task failed", zap.Error(err))
			return
		}
		if task == nil { // no task
			return
		}
		// not arrival execution time
		if task.Timestamp > int(currentTime.Unix()) {
			return
		}
		// do task
		taskDetail, err := getTaskDetail(task.Id)
		if err != nil { // retry
			bucket.ErrorWF("get task detail failed", zap.String("taskId", task.Id))
			continue
		}
		// if callback success, remove finish task
		if ok := bucket.CallbackFunc(taskDetail.Data); ok {
			err = removeTask(bucketName, task.Id)
			if err != nil {
				continue
			}
		} else {
			bucket.ErrorWF("error happen!", zap.Any("taskID", task.Id))
			continue // retry
		}
		return
	}
}

// get task from redis zset
func getTask(bucketName string) (*Task, error) {
	value, err := redisclient.ZRangeFirst(bucketName) // ZRANGE key 0 0 WITHSCORES
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}
	timestamp := int(value[0].(float64))
	taskId := value[1].(string)
	task := Task{
		Id:        taskId,
		Timestamp: timestamp,
	}
	return &task, nil
}

// get task detail by taskId
func getTaskDetail(taskId string) (*Task, error) {
	v, err := redisclient.Get(taskId)
	if err != nil {
		return nil, err
	}
	if v == "" {
		return nil, nil
	}
	task := Task{}
	err = json.Unmarshal([]byte(v), &task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// remove the task
func removeTask(bucketName string, taskId string) error {
	err := redisclient.ZRem(bucketName, taskId)
	if err != nil {
		return err
	}
	err = redisclient.Del(taskId)
	if err != nil {
		return err
	}
	return nil
}

// add task
func (bucket *BucketTicker) DelTask(task *Task) error {
	return removeTask(bucket.Name, task.Id)
}
