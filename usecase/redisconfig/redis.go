// @Author pangchenyang 2025/6/12 17:56:00
// @Desc: 
package redisconfig

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	cfg2 "maze_game_server/lib/nano/cfg"
	"fmt"
)

const (
	RedisConfigName  = "Redis"
	DefaultRedisName = "default"
)

// RedisConfig 存储Redis连接配置
type RedisConfig struct {
	Addr         string // Redis服务器地址
	Password     string // 密码
	DB           int    // 数据库索引
	PoolSize     int    // 连接池大小
	MinIdleConns int    // 最小空闲连接数
	Timeout      time.Duration
}

// RedisService Redis管理服务类
type RedisService struct {
	ctx            context.Context
	cancel         context.CancelFunc
	config         RedisConfig
	cfg            cfg2.CfgSvr // todo 配置更新
	client         *redis.Client
	clientInstance map[string]redis.Client // key:clientName
	initialized    bool
	once           sync.Once
	mutex          sync.RWMutex
}

// 单例实例
var (
	redisService *RedisService
	onceInstance sync.Once
)

// GetRedisService 获取Redis服务
func GetRedisService() *RedisService {
	onceInstance.Do(func() {
		redisService = &RedisService{}
	})
	return redisService
}

// Init 初始化Redis配置
func (rs *RedisService) Init(cfg cfg2.CfgSvr) error {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	if rs.initialized {
		return errors.New("redis service already initialized")
	}
	rs.cfg = cfg
	var redisCfg = RedisConfig{}
	ret, err := rs.cfg.LoadConfig(RedisConfigName, redisCfg)
	if err != nil {
		fmt.Println("redis service init failed, err:", err)
		return err
	}
	c, ok := ret.(map[string]interface{})
	if !ok {
		fmt.Println("redis service init failed, err:", err)
		return err
	}
	redisCfg.Addr = c["addr"].(string)
	redisCfg.Password = c["pwd"].(string)

	rs.config = redisCfg
	rs.initialized = true
	rs.ctx, rs.cancel = context.WithCancel(context.Background())
	fmt.Printf("Redis service initialized with config:%v redisService:%v \n", redisCfg, rs)
	return nil
}

// Start 启动Redis连接
func (rs *RedisService) Start() error {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	// 检查是否已初始化
	if !rs.initialized {
		return errors.New("redis service not initialized, please call Init first")
	}

	// 检查是否已启动
	if rs.client != nil {
		return errors.New("redis service already started")
	}

	// 创建Redis客户端
	rs.client = redis.NewClient(&redis.Options{
		Addr:         rs.config.Addr,
		Password:     rs.config.Password,
		DB:           rs.config.DB,
		PoolSize:     rs.config.PoolSize,
		MinIdleConns: rs.config.MinIdleConns,
	})

	_, err := rs.client.Ping(rs.ctx).Result()
	if err != nil {
		rs.client = nil
		return err
	}

	fmt.Println("Redis connection started successfully")
	return nil
}

// Stop 停止Redis连接
func (rs *RedisService) Stop() error {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	// 检查是否已启动
	if rs.client == nil {
		return errors.New("redis service not started")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭连接
	err := rs.client.Close()
	if err != nil {
		log.Printf("Error closing Redis connection: %v \n", err)
	} else {
		log.Println("Redis connection closed successfully")
	}
	rs.ctx = ctx
	rs.client = nil
	return err
}

// GetClient 获取Redis客户端
func (rs *RedisService) GetClient() (*redis.Client, error) {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	if rs.client == nil {
		return nil, errors.New("redis service not started")
	}
	return rs.client, nil
}
