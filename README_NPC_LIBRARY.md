# 图鉴系统 (NPC Library System)

## 概述

图鉴系统是一个完整的NPC收集、管理和升级系统，支持多种获取方式、分片存储、缓存优化和实时通知。

## 系统特性

- **多源获取**: 支持战斗抓取、抽卡掉落、任务奖励等多种获取方式
- **分片存储**: 64个分片表，支持大规模数据存储
- **缓存优化**: Redis缓存，提升查询性能
- **实时通知**: 支持图鉴变化实时通知
- **升级系统**: 支持图鉴等级提升和属性增强
- **合成功能**: 支持NPC强制合成

## 快速开始

### 1. 环境准备

```bash
# 安装依赖
go mod tidy

# 安装MySQL客户端
sudo apt-get install mysql-client

# 安装Redis客户端
sudo apt-get install redis-tools
```

### 2. 数据库初始化

```bash
# 设置环境变量
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=your_password
export DB_NAME=maze_game

# 执行初始化脚本
chmod +x scripts/npc_library_setup.sh
./scripts/npc_library_setup.sh
```

### 3. 启动服务

```bash
# 启动主服务器
go run servers/maze_main_server/main.go
```

### 4. 运行测试

```bash
# 运行图鉴系统测试
go test ./tester/npc_library_t/... -v
```

## 系统架构

### 文件结构

```
maze_game_server/
├── model/npclibrarymodel/           # 数据模型层
│   ├── npc_library_info.go         # 基础数据结构
│   └── npc_library_model.go        # 数据访问接口
├── services/npclibraryservice/      # 服务层
│   ├── npc_library_service.go      # 业务逻辑
│   └── npc_library_redis.go        # Redis缓存
├── servers/maze_main_server/process/npclibrary/  # 处理器层
│   ├── npc_library_handler.go      # RPC请求处理
│   └── npc_library_register.go    # 处理器注册
├── pb/common/NPCLibrary/           # 协议定义
│   └── NPCLibrary.pb.go           # Protobuf消息
├── config/                         # 配置文件
│   └── npc_library_config.yaml    # 系统配置
├── scripts/                        # 脚本文件
│   ├── npc_library_init.sql       # 数据库初始化
│   └── npc_library_setup.sh       # 系统安装脚本
├── tester/npc_library_t/           # 测试文件
│   └── npc_library_test.go        # 单元测试
└── docs/                          # 文档
    └── npc_library_system.md       # 系统文档
```

### 核心组件

1. **数据模型** (`model/npclibrarymodel/`)
   - `NPCLibraryInfo`: 基础图鉴信息
   - `NPCLibraryInfoEx`: 扩展图鉴信息
   - `NPCLibraryModel`: 数据访问接口

2. **服务层** (`services/npclibraryservice/`)
   - `NPCLibraryService`: 业务逻辑接口
   - `NPCLibraryRedis`: Redis缓存接口

3. **处理器** (`servers/maze_main_server/process/npclibrary/`)
   - `NPCLibraryHandler`: RPC请求处理器
   - `NPCLibraryRegister`: 处理器注册

4. **协议定义** (`pb/common/NPCLibrary/`)
   - Protobuf消息定义
   - RPC接口规范

## 数据模型

### 核心数据结构

#### NPCLibraryInfo (基础信息)
```go
type NPCLibraryInfo struct {
    UserID            uint64 `json:"user_id"`             // 用户ID
    NPCRoleID         uint64 `json:"npc_role_id"`        // NPC角色ID
    NPCUserID         uint64 `json:"npc_user_id"`        // NPC用户ID
    NPCCount          int32  `json:"npc_count"`          // NPC数量
    NPCLevel          int32  `json:"npc_level"`          // NPC等级
    CreateTime        uint64 `json:"create_time"`        // 创建时间
    ChangeToken       uint64 `json:"change_token"`       // 变化令牌
}
```

#### NPCLibraryInfoEx (扩展信息)
```go
type NPCLibraryInfoEx struct {
    NPCLibraryInfo
    NextLevelNeedCount uint64 `json:"next_level_need_count"` // 下一级需要数量
    NextLevelCost      int32  `json:"next_level_cost"`       // 下一级消耗
    NextLevelCostType  int32  `json:"next_level_cost_type"`  // 消耗类型
    ComposeFlag        int32  `json:"compose_flag"`          // 合成标记
    NextCostStrength   int32  `json:"next_cost_strength"`    // 下一级消耗强度
    NextRequiredLevel  int32  `json:"next_required_level"`   // 下一级要求等级
    HP                 int32  `json:"hp"`                    // 血量
    FightVal           int32  `json:"fight_val"`             // 战斗力
    Prestige           int32  `json:"prestige"`              // 声望
    ComposeNeedCnt     int32  `json:"compose_need_cnt"`      // 合成需要数量
    Order              int32  `json:"order"`                 // 排序
    Quality            int32  `json:"quality"`               // 品质
    GroupOrder         int32  `json:"group_order"`           // 组排序
}
```

## 数据库设计

### 分片策略

- **分片数量**: 64个分片
- **分片规则**: `(userID >> 8) % 64`
- **表名格式**: `npc_library_%d` (0-63)

### 表结构

```sql
CREATE TABLE npc_library_%d (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    host_user_id BIGINT UNSIGNED NOT NULL,
    npc_role_id BIGINT UNSIGNED NOT NULL,
    npc_user_id BIGINT UNSIGNED NOT NULL,
    npc_count INT NOT NULL DEFAULT 0,
    npc_level INT NOT NULL DEFAULT 1,
    create_time BIGINT UNSIGNED NOT NULL,
    next_level_need_count BIGINT UNSIGNED DEFAULT 0,
    next_level_cost INT DEFAULT 0,
    next_level_cost_type INT DEFAULT 0,
    compose_flag INT DEFAULT 0,
    next_cost_strength INT DEFAULT 0,
    next_required_level INT DEFAULT 0,
    hp INT DEFAULT 0,
    fight_val INT DEFAULT 0,
    prestige INT DEFAULT 0,
    compose_need_cnt INT DEFAULT 0,
    `order` INT DEFAULT 0,
    quality INT DEFAULT 0,
    group_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_npc (host_user_id, npc_role_id),
    INDEX idx_user_id (host_user_id),
    INDEX idx_npc_role_id (npc_role_id),
    INDEX idx_create_time (create_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## RPC接口

### 添加图鉴

**请求**: `DEF_SVR_NPC_LIBRARY_ADD_RQ`
```protobuf
message NPCLibraryAddRequest {
    uint64 user_id = 1;
    uint64 npc_role_id = 2;
    uint64 npc_user_id = 3;
    int32 count = 4;
    int32 from_type = 5;
}
```

**响应**: `DEF_SVR_NPC_LIBRARY_ADD_RS`
```protobuf
message NPCLibraryAddResponse {
    int32 result = 1;
    string message = 2;
    NPCLibraryInfoEx npc_library = 3;
}
```

### 查询图鉴

**请求**: `DEF_SVR_NPC_LIBRARY_QUERY_RQ`
```protobuf
message NPCLibraryQueryRequest {
    uint64 user_id = 1;
    string cursor = 2;
    int32 limit = 3;
}
```

**响应**: `DEF_SVR_NPC_LIBRARY_QUERY_RS`
```protobuf
message NPCLibraryQueryResponse {
    int32 result = 1;
    string message = 2;
    repeated NPCLibraryInfoEx npc_libraries = 3;
    string next_cursor = 4;
    bool has_more = 5;
}
```

### 升级图鉴

**请求**: `DEF_SVR_NPC_LIBRARY_UPGRADE_RQ`
```protobuf
message NPCLibraryUpgradeRequest {
    uint64 user_id = 1;
    uint64 npc_role_id = 2;
}
```

**响应**: `DEF_SVR_NPC_LIBRARY_UPGRADE_RS`
```protobuf
message NPCLibraryUpgradeResponse {
    int32 result = 1;
    string message = 2;
    NPCLibraryInfoEx npc_library = 3;
}
```

### 检查图鉴

**请求**: `DEF_SVR_NPC_LIBRARY_CHECK_RQ`
```protobuf
message NPCLibraryCheckRequest {
    uint64 user_id = 1;
    repeated uint64 npc_role_ids = 2;
}
```

**响应**: `DEF_SVR_NPC_LIBRARY_CHECK_RS`
```protobuf
message NPCLibraryCheckResponse {
    int32 result = 1;
    string message = 2;
    repeated NPCLibraryInfoEx npc_libraries = 3;
}
```

### 获取下一级信息

**请求**: `DEF_SVR_NPC_LIBRARY_NEXTLEVEL_INFO_RQ`
```protobuf
message NPCLibraryNextLevelInfoRequest {
    uint64 user_id = 1;
    uint64 npc_role_id = 2;
}
```

**响应**: `DEF_SVR_NPC_LIBRARY_NEXTLEVEL_INFO_RS`
```protobuf
message NPCLibraryNextLevelInfoResponse {
    int32 result = 1;
    string message = 2;
    NPCLibraryInfoEx npc_library = 3;
}
```

### 强制合成NPC

**请求**: `DEF_SVR_NPC_LIBRARYS_CHECK_RQ`
```protobuf
message NPCLibraryForceComposeRequest {
    uint64 user_id = 1;
    uint64 npc_role_id = 2;
}
```

**响应**: `DEF_SVR_NPC_LIBRARYS_CHECK_RS`
```protobuf
message NPCLibraryForceComposeResponse {
    int32 result = 1;
    string message = 2;
    NPCLibraryInfoEx npc_library = 3;
}
```

### 图鉴变化通知

**通知**: `DEF_SVR_NPC_LIBRARY_NOTIFY`
```protobuf
message NPCLibraryNotify {
    uint64 user_id = 1;
    uint64 npc_role_id = 2;
    int32 change_type = 3;
    NPCLibraryInfoEx npc_library = 4;
}
```

## 配置说明

### 系统配置

```yaml
npc_library:
  # 分片配置
  sharding:
    shard_count: 64
    shard_rule: "user_id_mod"
  
  # 缓存配置
  cache:
    expire_time: 3600
    key_prefix: "npc_library"
    enabled: true
  
  # 数据库配置
  database:
    max_connections: 100
    connection_timeout: 30
    query_timeout: 10
  
  # 业务配置
  business:
    default_level: 1
    max_level: 10
    upgrade_cost:
      base_cost: 100
      level_multiplier: 1.5
      cost_type: 1
```

### 环境变量

| 变量名 | 描述 | 默认值 |
|--------|------|--------|
| DB_HOST | 数据库主机 | localhost |
| DB_PORT | 数据库端口 | 3306 |
| DB_USER | 数据库用户 | root |
| DB_PASSWORD | 数据库密码 | 空 |
| DB_NAME | 数据库名称 | maze_game |
| REDIS_HOST | Redis主机 | localhost |
| REDIS_PORT | Redis端口 | 6379 |
| REDIS_PASSWORD | Redis密码 | 空 |

## 业务流程

### 添加图鉴流程

1. 客户端发送添加图鉴请求
2. Handler接收请求并验证参数
3. Service处理业务逻辑
4. Model访问数据库
5. 更新Redis缓存
6. 发送变化通知
7. 返回结果给客户端

### 查询图鉴流程

1. 客户端发送查询请求
2. Handler接收请求
3. Service处理业务逻辑
4. 检查Redis缓存
5. 如果缓存未命中，查询数据库
6. 更新缓存
7. 返回结果给客户端

### 升级图鉴流程

1. 客户端发送升级请求
2. Handler接收请求
3. Service检查升级条件
4. 计算升级消耗
5. 扣除消耗
6. 更新图鉴等级
7. 更新缓存
8. 发送通知
9. 返回结果

## 监控和运维

### 关键指标

1. **性能指标**
   - 响应时间
   - 吞吐量
   - 错误率

2. **业务指标**
   - 图鉴添加成功率
   - 升级成功率
   - 查询命中率

3. **系统指标**
   - 数据库连接数
   - Redis内存使用
   - CPU使用率

### 日志配置

```yaml
logging:
  level: "info"
  file_path: "logs/npc_library.log"
  max_size: 100
  max_age: 7
  max_backups: 3
```

### 告警配置

```yaml
monitoring:
  enabled: true
  collect_interval: 30
  alerts:
    error_rate_threshold: 0.05
    response_time_threshold: 1000
    db_connection_threshold: 80
```

## 故障排查

### 常见问题

1. **数据库连接失败**
   - 检查数据库服务状态
   - 验证连接参数
   - 检查网络连通性

2. **Redis连接失败**
   - 检查Redis服务状态
   - 验证连接参数
   - 检查内存使用情况

3. **分片表不存在**
   - 检查初始化脚本执行结果
   - 验证分片规则
   - 重新执行初始化

4. **缓存不一致**
   - 检查缓存更新逻辑
   - 验证数据同步
   - 清理缓存重新加载

### 调试工具

```bash
# 检查数据库表
mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASSWORD $DB_NAME -e "SHOW TABLES LIKE 'npc_library_%';"

# 检查Redis连接
redis-cli -h$REDIS_HOST -p$REDIS_PORT ping

# 查看日志
tail -f logs/npc_library.log
```

## 扩展说明

### 功能扩展

1. **新增获取来源**
   - 在`LibFromType`中添加新类型
   - 更新业务逻辑处理

2. **新增图鉴属性**
   - 扩展`NPCLibraryInfoEx`结构
   - 更新数据库表结构
   - 更新缓存逻辑

3. **新增业务规则**
   - 在Service层添加新方法
   - 更新Handler层处理
   - 添加新的RPC接口

### 性能优化

1. **缓存优化**
   - 调整缓存过期时间
   - 优化缓存键设计
   - 实现缓存预热

2. **数据库优化**
   - 添加索引优化查询
   - 调整分片策略
   - 优化SQL语句

3. **并发优化**
   - 调整连接池大小
   - 优化锁机制
   - 实现异步处理

## 版本历史

| 版本 | 描述 |
|------|------|------|
| v1.0.0 | 初始版本，基础功能实现 |
| v1.1.0 | 添加缓存优化 |
| v1.2.0 | 添加分片支持 |
| v1.3.0 |  添加监控和告警 |




