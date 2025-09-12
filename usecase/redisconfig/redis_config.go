// @Author pangchenyang 2025/6/10 14:14:00
// @Desc: 
package redisconfig

// type RedisService struct {
// 	cfg  cfg2.CfgSvr
// 	name string
// 	fklog.FKLogI
// 	ctx       context.Context
// 	Client    *redis.Client
// 	instances map[string]*redis.Client // 支持多实例
// 	cancel    context.CancelFunc
// 	mu        sync.RWMutex
// }
//
// const (
// 	RedisConfig      = "Redis"
// 	DefaultRedisName = "default"
// )
//
// type RedisCfg struct {
// 	Addr string `yaml:"addr"`
// 	Pwd  string `yaml:"pwd"`
// 	Db   int    `yaml:"db"`
// }
//
// func NewRedisService(config cfg2.CfgSvr) *RedisService {
// 	return &RedisService{
// 		cfg:       config,
// 		name:      "redis_service",
// 		instances: make(map[string]*redis.Client),
// 	}
// }
//
// func (rs *RedisService) Name() string {
// 	return rs.name
// }
//
// func (rs *RedisService) OnInit(logger fklog.FKLogI, config fkcore.FkConfigerI) (err error) {
// 	rs.FKLogI = logger.Clone("RedisService")
// 	ctx, cancel := context.WithCancel(context.Background())
// 	rs.ctx, rs.cancel = ctx, cancel
// 	return
// }
//
// func (rs *RedisService) OnStart(logger fklog.FKLogI, config fkcore.FkConfigerI) (err error) {
// 	logger.InfoWF("RedisService start begin")
// 	redisCfg := RedisCfg{}
// 	ret, err := rs.cfg.LoadConfig(RedisConfig, redisCfg)
// 	if err != nil {
// 		logger.ErrorWF("RedisService LoadConfig failed.", zap.Error(err))
// 		return
// 	}
// 	c, ok := ret.(map[string]interface{})
// 	if !ok {
// 		logger.ErrorWF("RedisService LoadConfig unmarshal failed.", zap.Any("ret", ret))
// 		return
// 	}
// 	redisCfg.Addr = c["addr"].(string)
// 	redisCfg.Pwd = c["pwd"].(string)
// 	// redisCfg.Db = c["Db"].(int)
// 	client, err := rs.InitRedis(redisCfg, DefaultRedisName)
// 	rs.Client = client
// 	logger.InfoWF("RedisService start success.", zap.Any("redisCfg", redisCfg), zap.Any("globalRedisService", globalRedisService))
// 	return
// }
//
// func (rs *RedisService) OnStop(logger fklog.FKLogI) (err error) {
// 	logger.WarnWF("RedisService stop begin")
//
// 	rs.mu.Lock()
// 	defer rs.mu.Unlock()
//
// 	for name, client := range rs.instances {
// 		if err := client.Close(); err != nil {
// 			logger.ErrorWF("Close Redis connection failed",
// 				zap.String("instance", name),
// 				zap.Error(err))
// 		}
// 		delete(rs.instances, name)
// 	}
//
// 	if rs.cancel != nil {
// 		rs.cancel()
// 	}
// 	globalRedisService = nil
// 	logger.InfoWF("RedisService stop success")
// 	return nil
// }
//
// func (rs *RedisService) OnFinish(logger fklog.FKLogI) (err error) {
// 	logger.InfoWF("RedisService finish success.")
// 	return
// }
//
// func GetClient(name ...string) *redis.Client {
// 	return globalRedisService.GetClient(name...)
// }
//
// // GetClient 获取Redis客户端
// func (rs *RedisService) GetClient(name ...string) *redis.Client {
// 	instanceName := DefaultRedisName
// 	if len(name) > 0 {
// 		instanceName = name[0]
// 	}
//
// 	rs.mu.RLock()
// 	defer rs.mu.RUnlock()
//
// 	return rs.instances[instanceName]
// }
//
// // InitRedis 初始化Redis连接
// func (rs *RedisService) InitRedis(cfg RedisCfg, name ...string) (*redis.Client, error) {
// 	instanceName := DefaultRedisName
// 	if len(name) > 0 {
// 		instanceName = name[0]
// 	}
//
// 	client := redis.NewClient(&redis.Options{
// 		Addr:     cfg.Addr,
// 		Password: cfg.Pwd,
// 		DB:       cfg.Db,
// 	})
//
// 	if _, err := client.Ping(rs.ctx).Result(); err != nil {
// 		return nil, err
// 	}
//
// 	rs.mu.Lock()
// 	defer rs.mu.Unlock()
// 	rs.instances[instanceName] = client
//
// 	return client, nil
// }
