-- 图鉴系统数据库初始化脚本

-- 创建图鉴分片表（64个分片）
-- 分片规则：根据用户ID的后8位进行分片 (userID >> 8) % 64

-- 创建图鉴表
CREATE TABLE IF NOT EXISTS npc_library_0 (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    host_user_id BIGINT UNSIGNED NOT NULL COMMENT '主人用户ID',
    npc_role_id BIGINT UNSIGNED NOT NULL COMMENT 'NPC角色ID',
    npc_user_id BIGINT UNSIGNED NOT NULL COMMENT 'NPC用户ID',
    npc_count INT NOT NULL DEFAULT 0 COMMENT 'NPC个数',
    npc_level INT NOT NULL DEFAULT 1 COMMENT '图鉴级别',
    create_time BIGINT UNSIGNED NOT NULL COMMENT '创建时间(初次解锁时间)',
    next_level_need_count BIGINT UNSIGNED DEFAULT 0 COMMENT '下一级需要数量',
    next_level_cost INT DEFAULT 0 COMMENT '下一级消耗',
    next_level_cost_type INT DEFAULT 0 COMMENT '消耗类型',
    compose_flag INT DEFAULT 0 COMMENT '合成标记',
    next_cost_strength INT DEFAULT 0 COMMENT '下一级消耗强度',
    next_required_level INT DEFAULT 0 COMMENT '下一级要求等级',
    hp INT DEFAULT 0 COMMENT '血量',
    fight_val INT DEFAULT 0 COMMENT '战斗力',
    prestige INT DEFAULT 0 COMMENT '声望',
    compose_need_cnt INT DEFAULT 0 COMMENT '合成需要数量',
    `order` INT DEFAULT 0 COMMENT '排序',
    quality INT DEFAULT 0 COMMENT '品质',
    group_order INT DEFAULT 0 COMMENT '组排序',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_npc (host_user_id, npc_role_id),
    INDEX idx_user_id (host_user_id),
    INDEX idx_npc_role_id (npc_role_id),
    INDEX idx_create_time (create_time),
    INDEX idx_quality (quality),
    INDEX idx_group_order (group_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='NPC图鉴表_分片0';

-- 创建其他63个分片表（这里只展示创建语句，实际需要循环创建）
-- 可以通过脚本动态生成所有分片表

-- 创建存储过程：添加用户图鉴
DELIMITER $$
CREATE PROCEDURE IF NOT EXISTS pr_add_user_npc_lib_0(
    IN p_host_user_id BIGINT UNSIGNED,
    IN p_npc_role_id BIGINT UNSIGNED,
    IN p_npc_user_id BIGINT UNSIGNED,
    IN p_increse_npc_lib_cnt INT,
    IN p_npc_level INT,
    IN p_from_type INT
)
BEGIN
    DECLARE v_existing_count INT DEFAULT 0;
    DECLARE v_new_count INT DEFAULT 0;
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;
        RESIGNAL;
    END;
    
    START TRANSACTION;
    
    -- 检查是否已存在
    SELECT npc_count INTO v_existing_count 
    FROM npc_library_0 
    WHERE host_user_id = p_host_user_id AND npc_role_id = p_npc_role_id;
    
    IF v_existing_count > 0 THEN
        -- 更新现有记录
        UPDATE npc_library_0 
        SET npc_count = npc_count + p_increse_npc_lib_cnt,
            updated_at = CURRENT_TIMESTAMP
        WHERE host_user_id = p_host_user_id AND npc_role_id = p_npc_role_id;
    ELSE
        -- 插入新记录
        INSERT INTO npc_library_0 (
            host_user_id, npc_role_id, npc_user_id, npc_count, npc_level, 
            create_time, next_level_need_count, next_level_cost, next_level_cost_type
        ) VALUES (
            p_host_user_id, p_npc_role_id, p_npc_user_id, p_increse_npc_lib_cnt, p_npc_level,
            UNIX_TIMESTAMP(), 0, 0, 0
        );
    END IF;
    
    COMMIT;
END$$
DELIMITER ;

-- 创建存储过程：升级图鉴
DELIMITER $$
CREATE PROCEDURE IF NOT EXISTS pr_upgrade_user_npc_lib_0(
    IN p_host_user_id BIGINT UNSIGNED,
    IN p_npc_role_id BIGINT UNSIGNED,
    IN p_new_level INT,
    IN p_cost INT,
    IN p_cost_type INT
)
BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;
        RESIGNAL;
    END;
    
    START TRANSACTION;
    
    -- 更新图鉴级别
    UPDATE npc_library_0 
    SET npc_level = p_new_level,
        next_level_cost = p_cost,
        next_level_cost_type = p_cost_type,
        updated_at = CURRENT_TIMESTAMP
    WHERE host_user_id = p_host_user_id AND npc_role_id = p_npc_role_id;
    
    COMMIT;
END$$
DELIMITER ;

-- 创建存储过程：查询用户图鉴
DELIMITER $$
CREATE PROCEDURE IF NOT EXISTS pr_query_user_npc_lib_0(
    IN p_host_user_id BIGINT UNSIGNED,
    IN p_cursor VARCHAR(255),
    IN p_limit INT
)
BEGIN
    -- 分页查询用户图鉴
    SELECT * FROM npc_library_0 
    WHERE host_user_id = p_host_user_id
    ORDER BY create_time DESC
    LIMIT p_limit OFFSET IFNULL(p_cursor, 0);
END$$
DELIMITER ;

-- 创建存储过程：检查图鉴是否存在
DELIMITER $$
CREATE PROCEDURE IF NOT EXISTS pr_check_user_npc_lib_0(
    IN p_host_user_id BIGINT UNSIGNED,
    IN p_npc_role_ids TEXT
)
BEGIN
    -- 检查多个NPC角色ID是否存在
    SELECT npc_role_id, 1 as exists_flag
    FROM npc_library_0 
    WHERE host_user_id = p_host_user_id 
    AND FIND_IN_SET(npc_role_id, p_npc_role_ids) > 0;
END$$
DELIMITER ;

-- 创建索引优化
-- 为常用查询创建复合索引
CREATE INDEX IF NOT EXISTS idx_user_quality ON npc_library_0 (host_user_id, quality);
CREATE INDEX IF NOT EXISTS idx_user_group ON npc_library_0 (host_user_id, group_order);
CREATE INDEX IF NOT EXISTS idx_user_level ON npc_library_0 (host_user_id, npc_level);

-- 创建视图：用户图鉴统计
CREATE OR REPLACE VIEW v_user_npc_lib_stats AS
SELECT 
    host_user_id,
    COUNT(*) as total_npc_count,
    COUNT(CASE WHEN npc_level > 1 THEN 1 END) as upgraded_npc_count,
    AVG(npc_level) as avg_npc_level,
    MAX(npc_level) as max_npc_level,
    SUM(npc_count) as total_npc_fragments
FROM npc_library_0
GROUP BY host_user_id;

-- 创建触发器：自动更新统计信息
DELIMITER $$
CREATE TRIGGER IF NOT EXISTS tr_npc_lib_after_insert
AFTER INSERT ON npc_library_0
FOR EACH ROW
BEGIN
    -- 可以在这里添加统计更新逻辑
    -- 例如更新用户图鉴统计表
END$$
DELIMITER ;

DELIMITER $$
CREATE TRIGGER IF NOT EXISTS tr_npc_lib_after_update
AFTER UPDATE ON npc_library_0
FOR EACH ROW
BEGIN
    -- 可以在这里添加统计更新逻辑
    -- 例如更新用户图鉴统计表
END$$
DELIMITER ;

-- 创建分区表（可选，用于大数据量场景）
-- ALTER TABLE npc_library_0 PARTITION BY RANGE (host_user_id) (
--     PARTITION p0 VALUES LESS THAN (1000000),
--     PARTITION p1 VALUES LESS THAN (2000000),
--     PARTITION p2 VALUES LESS THAN (3000000),
--     PARTITION p3 VALUES LESS THAN (4000000),
--     PARTITION p4 VALUES LESS THAN (5000000),
--     PARTITION p5 VALUES LESS THAN (6000000),
--     PARTITION p6 VALUES LESS THAN (7000000),
--     PARTITION p7 VALUES LESS THAN (8000000),
--     PARTITION p8 VALUES LESS THAN (9000000),
--     PARTITION p9 VALUES LESS THAN (10000000),
--     PARTITION p10 VALUES LESS THAN MAXVALUE
-- );


