CREATE DATABASE  IF NOT EXISTS maze DEFAULT CHARACTER SET utf8mb4;
use maze;


CREATE TABLE `t_maze_attr_chg_record`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `server_id` int(0) NOT NULL DEFAULT 0 COMMENT '服ID',
  `user_id` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户id',
  `attr_id` int(0) NOT NULL DEFAULT 0 COMMENT '属性ID',
  `attr_type` int(0) NOT NULL DEFAULT 0 COMMENT '属性类型',
  `new_val` bigint(0) NOT NULL DEFAULT 0 COMMENT '新值',
  `old_val` bigint(0) NOT NULL DEFAULT 0 COMMENT '旧值',
  `chg_type` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '变化类型',
  `chg_sub_type` int(0) NOT NULL DEFAULT 0 COMMENT '变化子类型',
  `chg_desc` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '原因描述',
  `extra` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '扩展信息',
  `group_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '组id',
  `create_time` bigint(0) NOT NULL DEFAULT 0 COMMENT '操作时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `user_id`(`user_id`) USING BTREE,
  INDEX `create_time`(`create_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '迷宫属性变化流水' ROW_FORMAT = Dynamic;

CREATE TABLE `t_maze_barrier_user_record`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `server_id` int(0) NOT NULL DEFAULT 0 COMMENT '服ID',
  `user_id` bigint(0) NOT NULL DEFAULT 0 COMMENT '用户id',
  `barrier` int(0) NOT NULL DEFAULT 0 COMMENT '关卡id',
  `game_ret` bigint(0) NOT NULL DEFAULT 0 COMMENT '用户闯关结果 1-通关成功 2-死亡失败 3-扫荡',
  `awards` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '本次获得的奖励',
  `group_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '组id',
  `create_time` bigint(0) NOT NULL DEFAULT 0 COMMENT '操作时间 毫秒',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `user_id`(`user_id`) USING BTREE,
  INDEX `create_time`(`create_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户迷宫闯关流水' ROW_FORMAT = Dynamic;

CREATE TABLE `t_maze_collect_chg_record`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `server_id` int(0) NOT NULL DEFAULT 0 COMMENT '服ID',
  `user_id` bigint(0) UNSIGNED NOT NULL COMMENT '用户id',
  `op_type` int(0) NOT NULL COMMENT '操作类型 1:初始化 2:定时收集 3:领取',
  `start_time` bigint(0) NOT NULL COMMENT '开始时间',
  `last_time` bigint(0) NOT NULL COMMENT '上次收集结算时间',
  `new_last_time` bigint(0) NOT NULL COMMENT '本次收集结算时间',
  `available_time` bigint(0) NOT NULL COMMENT '可领取时间',
  `end_time` bigint(0) NOT NULL COMMENT '结束时间',
  `period_time` int(0) NOT NULL COMMENT '产出周期（秒）',
  `collect_times` bigint(0) NOT NULL COMMENT '产出周期数',
  `barrier_id` int(0) NOT NULL COMMENT '关卡id',
  `trade_no` bigint(0) UNSIGNED NOT NULL COMMENT '加物品流水号',
  `add_items` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '收集的道具/领取的道具',
  `remain_items` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '领取后遗留的道具/累计产出道具',
  `ret_code` bigint(0) NOT NULL COMMENT '结果 0:成功 其他:失败',
  `group_id` int(0) UNSIGNED NOT NULL COMMENT '组id',
  `create_time` bigint(0) NOT NULL COMMENT '操作时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user`(`user_id`, `create_time`) USING BTREE,
  INDEX `idx_dt`(`create_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '迷宫挂机变化流水' ROW_FORMAT = Dynamic;

CREATE TABLE `t_maze_equip_assemble_record`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `server_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '服ID',
  `user_id` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户id',
  `equip_pos` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '装备位ID',
  `op_type` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '穿戴装备/更换装备/卸下装备',
  `new_equip_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '穿戴装备配置ID',
  `new_guid` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '穿戴装备guid',
  `old_equip_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '卸下装备配置ID',
  `old_guid` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '卸下装备guid',
  `old_f_elem` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '变化前激活信息',
  `new_f_elem` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '变化后激活信息',
  `ret_code` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '0:成功  其他失败',
  `code_mask` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '业务掩码',
  `trans_id` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '事务Id',
  `group_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '组id',
  `create_time` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '操作时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user`(`user_id`, `create_time`) USING BTREE,
  INDEX `idx_dt`(`create_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '装备穿戴流水' ROW_FORMAT = Dynamic;

CREATE TABLE `t_maze_equip_bag_record`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `server_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '服ID',
  `user_id` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户id',
  `chg_type` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '变化原因 1 添加 2 删除 3 更新 4 锁定 5 解锁 6 实例化装备 7 删除实例化装备',
  `trade_num` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '交易单号',
  `add_equip_guids` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '新增装备guid列表',
  `del_equip_guids` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '删除装备guid列表',
  `op_type` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '业务类型 挂机/锻造/购买',
  `is_fail` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '操作是否失败 0-成功 1-失败',
  `group_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '组id',
  `create_time` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '操作时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user`(`user_id`, `create_time`) USING BTREE,
  INDEX `idx_dt`(`create_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '装备分解流水' ROW_FORMAT = Dynamic;

CREATE TABLE `t_maze_equip_dismant_record`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `user_id` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户id',
  `equip_guids` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '装备guid列表',
  `trade_num` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '交易单号',
  `award` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '分解获得的材料',
  `op_type` int(0) NOT NULL DEFAULT 0 COMMENT '分解的操作来源',
  `is_fail` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '操作是否失败 0-成功 1-失败',
  `group_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '组id',
  `server_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '服务器id',
  `create_time` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '操作时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user`(`user_id`, `create_time`) USING BTREE,
  INDEX `idx_dt`(`create_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '装备分解流水' ROW_FORMAT = Dynamic;

CREATE TABLE `t_maze_equip_pos_level_up_record`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `server_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '服ID',
  `user_id` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户id',
  `pos_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '装备位id',
  `old_pos_lv` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '强化前装备位等级',
  `new_pos_lv` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '强化后装备位等级',
  `old_pos_suit_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '强化前装备位套装id',
  `new_pos_suit_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '强化后装备位套装id',
  `trade_no` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '扣物品流水号',
  `cost_items` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '扣物品',
  `result` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '结果 0:成功 1:强化失败 2:存储武力值属性失败 3:存储非武力值属性失败',
  `group_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '组id',
  `create_time` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '毫秒时间戳',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user`(`user_id`, `create_time`) USING BTREE,
  INDEX `idx_dt`(`create_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '装备位强化流水' ROW_FORMAT = Dynamic;

CREATE TABLE `t_maze_foe_record`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `server_id` int(0) NOT NULL DEFAULT 0 COMMENT '服ID',
  `user_id` bigint(0) UNSIGNED NOT NULL COMMENT '用户id',
  `barrier` int(0) NOT NULL COMMENT '关卡id',
  `area` int(0) NOT NULL COMMENT '区域id',
  `level` int(0) NOT NULL COMMENT '用户等级',
  `master_id` bigint(0) NOT NULL COMMENT '怪物id',
  `award_list` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '打怪获得的奖励 1银子 2装备积分 3经验',
  `equips` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '打怪获得的装备',
  `equip_points` int(0) NOT NULL COMMENT '本次打怪后的当前装备积分',
  `group_id` int(0) UNSIGNED NOT NULL COMMENT '组id',
  `create_time` bigint(0) NOT NULL COMMENT '操作时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_id`(`user_id`, `create_time`) USING BTREE,
  INDEX `idx_dt`(`create_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户打怪变化流水' ROW_FORMAT = Dynamic;

CREATE TABLE `t_maze_money_record`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `server_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '服ID',
  `user_id` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户id',
  `old_money_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '旧货币id',
  `old_money_count` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '旧货币数量',
  `new_money_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '新货币id',
  `new_money_count` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '新货币数量',
  `trade_no` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '交易号',
  `chg_reason` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '变化原因',
  `group_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '组id',
  `create_time` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '操作时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user`(`user_id`, `create_time`) USING BTREE,
  INDEX `idx_dt`(`create_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户货币变化流水' ROW_FORMAT = Dynamic;

CREATE TABLE `t_maze_sweep_record`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `server_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '服ID',
  `user_id` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户id',
  `barrier` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '关卡id',
  `awards` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '本次获得的奖励',
  `group_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '组id',
  `create_time` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '操作时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user`(`user_id`, `create_time`) USING BTREE,
  INDEX `idx_dt`(`create_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户迷宫扫荡流水' ROW_FORMAT = Dynamic;

CREATE TABLE `t_maze_temp_buff_change_record`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `server_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '服ID',
  `user_id` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户Id',
  `stage_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '关卡id',
  `chg_attrs` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '变化的属性（序列化的属性变更列表）',
  `chg_type` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '变化类型',
  `chg_desc` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '原因描述',
  `group_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '分组Id',
  `create_time` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '毫秒时间戳',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user`(`user_id`, `create_time`) USING BTREE,
  INDEX `idx_dt`(`create_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '临时buff流水' ROW_FORMAT = Dynamic;

CREATE TABLE `t_maze_user_level_record`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `server_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '服ID',
  `user_id` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户id',
  `old_level` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '旧等级',
  `old_total_exp` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '旧经验总值',
  `new_level` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '新等级',
  `new_total_exp` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '新经验总值',
  `group_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '组id',
  `create_time` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '操作时间 毫秒',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user`(`user_id`, `create_time`) USING BTREE,
  INDEX `idx_dt`(`create_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户等级变化流水' ROW_FORMAT = Dynamic;

CREATE TABLE `t_maze_user_reborn_record`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `server_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '服ID',
  `user_id` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户id',
  `barrier` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '关卡id',
  `reborn_count` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '复活次数 第n次复活',
  `reborn_cost` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '复活的消耗',
  `group_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '组id',
  `create_time` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '操作时间 毫秒',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user`(`user_id`, `create_time`) USING BTREE,
  INDEX `idx_dt`(`create_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户复活流水' ROW_FORMAT = Dynamic;

CREATE TABLE `t_maze_user_energy_record`  (
 `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
 `server_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '服ID',
 `group_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '组id',
 `user_id` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户id',
 `old_val` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '旧值',
 `new_val` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '新值',
 `op_type` int(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '操作类型\r\n100=时间恢复\r\n101=初始化体力\r\n102=体力瓶\r\n103=进入关卡\r\n104=扫荡关卡\r\n105=GM添加\r\n106=重置',
 `last_time` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '上次恢复时间',
 `create_time` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '操作时间',
 PRIMARY KEY (`id`) USING BTREE,
 INDEX `idx_user`(`user_id`, `create_time`) USING BTREE,
 INDEX `idx_dt`(`create_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户体力变化流水' ROW_FORMAT = Dynamic;

