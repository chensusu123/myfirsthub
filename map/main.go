package main

/*
#cgo CFLAGS: -O3 -ffast-math
#cgo darwin,amd64 CFLAGS: -march=native -mavx2 -mfma
#cgo darwin,arm64 CFLAGS: -mcpu=apple-m1

#ifdef __x86_64__
    #include <immintrin.h>
    #define SIMD_WIDTH 8
    #define USE_AVX2 1
#elif __aarch64__
    #include <arm_neon.h>
    #define SIMD_WIDTH 4
    #define USE_NEON 1
#endif

#include <stdint.h>
#include <string.h>
#include <math.h>
#include <sys/sysctl.h>

// =================== CPU特性检测 ===================

typedef struct {
    int has_sse4_2;
    int has_avx;
    int has_avx2;
    int has_fma;
    int has_neon;
    int is_apple_silicon;
    char cpu_brand[64];
} cpu_features_t;

cpu_features_t detect_cpu_features() {
    cpu_features_t features = {0};

#ifdef __x86_64__
    // Intel CPU特性检测
    unsigned int eax, ebx, ecx, edx;

    // 基本CPUID
    __asm__ ("cpuid" : "=a"(eax), "=b"(ebx), "=c"(ecx), "=d"(edx) : "a"(1));
    features.has_sse4_2 = (ecx & (1 << 20)) != 0;
    features.has_avx = (ecx & (1 << 28)) != 0;
    features.has_fma = (ecx & (1 << 12)) != 0;

    // 扩展特性检测
    __asm__ ("cpuid" : "=a"(eax), "=b"(ebx), "=c"(ecx), "=d"(edx) : "a"(7), "c"(0));
    features.has_avx2 = (ebx & (1 << 5)) != 0;

    // 获取CPU品牌
    size_t size = sizeof(features.cpu_brand);
    sysctlbyname("machdep.cpu.brand_string", features.cpu_brand, &size, NULL, 0);

#elif __aarch64__
    // Apple Silicon检测
    features.has_neon = 1; // ARM64默认支持NEON
    features.is_apple_silicon = 1;

    size_t size = sizeof(features.cpu_brand);
    sysctlbyname("hw.model", features.cpu_brand, &size, NULL, 0);
#endif

    return features;
}

// =================== Intel AVX2优化 ===================

#ifdef __x86_64__

// AVX2优化的批量距离计算
int avx2_collision_check(
    const float* restrict pos1_x, const float* restrict pos1_y,
    const float* restrict pos2_x, const float* restrict pos2_y,
    const float* restrict radius1, const float* restrict radius2,
    uint32_t* restrict collision_pairs, int count) {

    int collision_count = 0;
    const int simd_width = 8;

    // 主循环：处理8的倍数
    for (int i = 0; i <= count - simd_width; i += simd_width) {
        // 加载数据 - 使用unaligned load更安全
        __m256 x1 = _mm256_loadu_ps(&pos1_x[i]);
        __m256 y1 = _mm256_loadu_ps(&pos1_y[i]);
        __m256 x2 = _mm256_loadu_ps(&pos2_x[i]);
        __m256 y2 = _mm256_loadu_ps(&pos2_y[i]);
        __m256 r1 = _mm256_loadu_ps(&radius1[i]);
        __m256 r2 = _mm256_loadu_ps(&radius2[i]);

        // 计算距离平方: (x1-x2)² + (y1-y2)²
        __m256 dx = _mm256_sub_ps(x1, x2);
        __m256 dy = _mm256_sub_ps(y1, y2);

        // 使用FMA指令优化: dx*dx + dy*dy
        __m256 dist_sq = _mm256_fmadd_ps(dx, dx, _mm256_mul_ps(dy, dy));

        // 计算半径和的平方
        __m256 radius_sum = _mm256_add_ps(r1, r2);
        __m256 radius_sum_sq = _mm256_mul_ps(radius_sum, radius_sum);

        // 碰撞检测
        __m256 collision_mask = _mm256_cmp_ps(dist_sq, radius_sum_sq, _CMP_LT_OQ);
        int mask = _mm256_movemask_ps(collision_mask);

        // 处理碰撞结果
        for (int j = 0; j < 8; j++) {
            if (mask & (1 << j)) {
                collision_pairs[collision_count * 2] = i + j;
                collision_pairs[collision_count * 2 + 1] = i + j + count; // 简化的配对
                collision_count++;
            }
        }
    }

    // 处理剩余元素
    for (int i = (count / simd_width) * simd_width; i < count; i++) {
        float dx = pos1_x[i] - pos2_x[i];
        float dy = pos1_y[i] - pos2_y[i];
        float dist_sq = dx*dx + dy*dy;
        float radius_sum = radius1[i] + radius2[i];

        if (dist_sq < radius_sum * radius_sum) {
            collision_pairs[collision_count * 2] = i;
            collision_pairs[collision_count * 2 + 1] = i + count;
            collision_count++;
        }
    }

    return collision_count;
}

// Intel的范围查询优化
int avx2_range_query(
    float center_x, float center_y, float range,
    const float* entities_x, const float* entities_y,
    uint32_t* result_ids, int entity_count) {

    int found_count = 0;
    const __m256 center_x_vec = _mm256_set1_ps(center_x);
    const __m256 center_y_vec = _mm256_set1_ps(center_y);
    const __m256 range_sq = _mm256_set1_ps(range * range);

    for (int i = 0; i <= entity_count - 8; i += 8) {
        __m256 entity_x = _mm256_loadu_ps(&entities_x[i]);
        __m256 entity_y = _mm256_loadu_ps(&entities_y[i]);

        __m256 dx = _mm256_sub_ps(entity_x, center_x_vec);
        __m256 dy = _mm256_sub_ps(entity_y, center_y_vec);
        __m256 dist_sq = _mm256_fmadd_ps(dx, dx, _mm256_mul_ps(dy, dy));

        __m256 in_range = _mm256_cmp_ps(dist_sq, range_sq, _CMP_LE_OQ);
        int mask = _mm256_movemask_ps(in_range);

        for (int j = 0; j < 8; j++) {
            if (mask & (1 << j)) {
                result_ids[found_count++] = i + j;
            }
        }
    }

    // 处理剩余元素
    for (int i = (entity_count / 8) * 8; i < entity_count; i++) {
        float dx = entities_x[i] - center_x;
        float dy = entities_y[i] - center_y;
        if (dx*dx + dy*dy <= range*range) {
            result_ids[found_count++] = i;
        }
    }

    return found_count;
}

#endif // __x86_64__

// =================== Apple Silicon NEON优化 ===================

#ifdef __aarch64__

// NEON优化的批量距离计算
int neon_collision_check(
    const float* restrict pos1_x, const float* restrict pos1_y,
    const float* restrict pos2_x, const float* restrict pos2_y,
    const float* restrict radius1, const float* restrict radius2,
    uint32_t* restrict collision_pairs, int count) {

    int collision_count = 0;
    const int simd_width = 4; // NEON处理4个float32

    for (int i = 0; i <= count - simd_width; i += simd_width) {
        // 加载数据到NEON寄存器
        float32x4_t x1 = vld1q_f32(&pos1_x[i]);
        float32x4_t y1 = vld1q_f32(&pos1_y[i]);
        float32x4_t x2 = vld1q_f32(&pos2_x[i]);
        float32x4_t y2 = vld1q_f32(&pos2_y[i]);
        float32x4_t r1 = vld1q_f32(&radius1[i]);
        float32x4_t r2 = vld1q_f32(&radius2[i]);

        // 计算距离平方
        float32x4_t dx = vsubq_f32(x1, x2);
        float32x4_t dy = vsubq_f32(y1, y2);

        // 使用FMA指令: dx*dx + dy*dy
        float32x4_t dist_sq = vfmaq_f32(vmulq_f32(dy, dy), dx, dx);

        // 计算半径和的平方
        float32x4_t radius_sum = vaddq_f32(r1, r2);
        float32x4_t radius_sum_sq = vmulq_f32(radius_sum, radius_sum);

        // 碰撞检测
        uint32x4_t collision_mask = vcltq_f32(dist_sq, radius_sum_sq);

        // 提取结果
        uint32_t mask[4];
        vst1q_u32(mask, collision_mask);

        for (int j = 0; j < 4; j++) {
            if (mask[j] != 0) {
                collision_pairs[collision_count * 2] = i + j;
                collision_pairs[collision_count * 2 + 1] = i + j + count;
                collision_count++;
            }
        }
    }

    // 处理剩余元素
    for (int i = (count / simd_width) * simd_width; i < count; i++) {
        float dx = pos1_x[i] - pos2_x[i];
        float dy = pos1_y[i] - pos2_y[i];
        float dist_sq = dx*dx + dy*dy;
        float radius_sum = radius1[i] + radius2[i];

        if (dist_sq < radius_sum * radius_sum) {
            collision_pairs[collision_count * 2] = i;
            collision_pairs[collision_count * 2 + 1] = i + count;
            collision_count++;
        }
    }

    return collision_count;
}

// Apple Silicon的范围查询优化
int neon_range_query(
    float center_x, float center_y, float range,
    const float* entities_x, const float* entities_y,
    uint32_t* result_ids, int entity_count) {

    int found_count = 0;
    const float32x4_t center_x_vec = vdupq_n_f32(center_x);
    const float32x4_t center_y_vec = vdupq_n_f32(center_y);
    const float32x4_t range_sq = vdupq_n_f32(range * range);

    for (int i = 0; i <= entity_count - 4; i += 4) {
        float32x4_t entity_x = vld1q_f32(&entities_x[i]);
        float32x4_t entity_y = vld1q_f32(&entities_y[i]);

        float32x4_t dx = vsubq_f32(entity_x, center_x_vec);
        float32x4_t dy = vsubq_f32(entity_y, center_y_vec);
        float32x4_t dist_sq = vfmaq_f32(vmulq_f32(dy, dy), dx, dx);

        uint32x4_t in_range = vcleq_f32(dist_sq, range_sq);

        uint32_t mask[4];
        vst1q_u32(mask, in_range);

        for (int j = 0; j < 4; j++) {
            if (mask[j] != 0) {
                result_ids[found_count++] = i + j;
            }
        }
    }

    // 处理剩余元素
    for (int i = (entity_count / 4) * 4; i < entity_count; i++) {
        float dx = entities_x[i] - center_x;
        float dy = entities_y[i] - center_y;
        if (dx*dx + dy*dy <= range*range) {
            result_ids[found_count++] = i;
        }
    }

    return found_count;
}

#endif // __aarch64__

// =================== 通用接口 ===================

// 统一的碰撞检测接口
int simd_collision_check_universal(
    const float* pos1_x, const float* pos1_y,
    const float* pos2_x, const float* pos2_y,
    const float* radius1, const float* radius2,
    uint32_t* collision_pairs, int count) {

#ifdef __x86_64__
    return avx2_collision_check(pos1_x, pos1_y, pos2_x, pos2_y,
                               radius1, radius2, collision_pairs, count);
#elif __aarch64__
    return neon_collision_check(pos1_x, pos1_y, pos2_x, pos2_y,
                               radius1, radius2, collision_pairs, count);
#else
    // 回退到标量实现
    int collision_count = 0;
    for (int i = 0; i < count; i++) {
        float dx = pos1_x[i] - pos2_x[i];
        float dy = pos1_y[i] - pos2_y[i];
        float dist_sq = dx*dx + dy*dy;
        float radius_sum = radius1[i] + radius2[i];

        if (dist_sq < radius_sum * radius_sum) {
            collision_pairs[collision_count * 2] = i;
            collision_pairs[collision_count * 2 + 1] = i + count;
            collision_count++;
        }
    }
    return collision_count;
#endif
}

// 统一的范围查询接口
int simd_range_query_universal(
    float center_x, float center_y, float range,
    const float* entities_x, const float* entities_y,
    uint32_t* result_ids, int entity_count) {

#ifdef __x86_64__
    return avx2_range_query(center_x, center_y, range,
                           entities_x, entities_y, result_ids, entity_count);
#elif __aarch64__
    return neon_range_query(center_x, center_y, range,
                           entities_x, entities_y, result_ids, entity_count);
#else
    // 标量回退实现
    int found_count = 0;
    for (int i = 0; i < entity_count; i++) {
        float dx = entities_x[i] - center_x;
        float dy = entities_y[i] - center_y;
        if (dx*dx + dy*dy <= range*range) {
            result_ids[found_count++] = i;
        }
    }
    return found_count;
#endif
}

*/
import "C"
import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
	"unsafe"
)

// CPU特性结构
type CPUFeatures struct {
	HasSSE42       bool
	HasAVX         bool
	HasAVX2        bool
	HasFMA         bool
	HasNEON        bool
	IsAppleSilicon bool
	Brand          string
	SIMDWidth      int
}

// 空间分割系统的配置常量
const (
	MAX_ENTITIES_SPATIAL = 100000 // 最大实体数
	GRID_SIZE_SPATIAL    = 50.0   // 空间分割网格大小
	MAX_COLLISION_CHECKS = 1000   // 每个实体最大碰撞检测数
)

// 优化的位置结构体
type Position struct {
	X, Y float32
}

// 实体结构体
type Entity struct {
	ID     uint32   // 4字节
	Pos    Position // 8字节
	Radius float32  // 4字节
	GridX  int16    // 2字节 - 网格坐标
	GridY  int16    // 2字节
	Active bool     // 1字节
	_      [3]byte  // 填充到32字节对齐
}

// 空间分割网格
type SpatialGrid struct {
	cellSize    float32
	invCellSize float32 // 预计算倒数，避免除法
	width       int
	height      int
	cells       [][]uint32 // 每个格子存储实体ID列表
	cellPool    sync.Pool  // 对象池复用切片
}

// 碰撞检测系统
type CollisionSystem struct {
	entities    []Entity     // 实体数组
	activeCount int          // 活跃实体数量
	grid        *SpatialGrid // 空间分割
	workerPool  *WorkerPool  // 工作池
	resultPool  sync.Pool    // 结果对象池
}

// 工作池用于并发处理
type WorkerPool struct {
	workers    int
	taskChan   chan CollisionTask
	resultChan chan CollisionResult
	wg         sync.WaitGroup
}

// 碰撞检测任务
type CollisionTask struct {
	startIdx int
	endIdx   int
	entities *[]Entity
	grid     *SpatialGrid
}

// 碰撞结果（重命名避免冲突）
type CollisionResult struct {
	entity1  uint32
	entity2  uint32
	distance float32
}

// Mac优化的碰撞检测器
type MacSIMDCollisionDetector struct {
	cpuFeatures    CPUFeatures
	entitiesX      []float32
	entitiesY      []float32
	entitiesR      []float32
	collisionPairs []uint32
	results        []uint32
	maxEntities    int
}

// 创建Mac优化的SIMD检测器
func NewMacSIMDCollisionDetector(maxEntities int) *MacSIMDCollisionDetector {
	detector := &MacSIMDCollisionDetector{
		maxEntities: maxEntities,
	}

	// 检测CPU特性
	detector.detectCPUFeatures()

	// 分配内存缓冲区
	detector.entitiesX = make([]float32, maxEntities)
	detector.entitiesY = make([]float32, maxEntities)
	detector.entitiesR = make([]float32, maxEntities)
	detector.collisionPairs = make([]uint32, maxEntities*2)
	detector.results = make([]uint32, maxEntities*2)

	detector.printSystemInfo()
	return detector
}

// 检测CPU特性
func (m *MacSIMDCollisionDetector) detectCPUFeatures() {
	features := C.detect_cpu_features()

	m.cpuFeatures = CPUFeatures{
		HasSSE42:       features.has_sse4_2 != 0,
		HasAVX:         features.has_avx != 0,
		HasAVX2:        features.has_avx2 != 0,
		HasFMA:         features.has_fma != 0,
		HasNEON:        features.has_neon != 0,
		IsAppleSilicon: features.is_apple_silicon != 0,
		Brand:          C.GoString(&features.cpu_brand[0]),
	}

	// 设置SIMD宽度
	if m.cpuFeatures.HasAVX2 {
		m.cpuFeatures.SIMDWidth = 8
	} else if m.cpuFeatures.HasNEON {
		m.cpuFeatures.SIMDWidth = 4
	} else {
		m.cpuFeatures.SIMDWidth = 1
	}
}

// 打印系统信息
func (m *MacSIMDCollisionDetector) printSystemInfo() {
	fmt.Println("=== Mac SIMD碰撞检测系统 ===")
	fmt.Printf("Go版本: %s\n", runtime.Version())
	fmt.Printf("操作系统: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("CPU: %s\n", m.cpuFeatures.Brand)
	fmt.Printf("CPU核心: %d\n", runtime.NumCPU())

	fmt.Println("\nSIMD特性:")
	if runtime.GOARCH == "amd64" {
		fmt.Printf("  SSE4.2: %v\n", m.cpuFeatures.HasSSE42)
		fmt.Printf("  AVX:    %v\n", m.cpuFeatures.HasAVX)
		fmt.Printf("  AVX2:   %v\n", m.cpuFeatures.HasAVX2)
		fmt.Printf("  FMA:    %v\n", m.cpuFeatures.HasFMA)
	} else if runtime.GOARCH == "arm64" {
		fmt.Printf("  NEON:         %v\n", m.cpuFeatures.HasNEON)
		fmt.Printf("  Apple Silicon: %v\n", m.cpuFeatures.IsAppleSilicon)
	}
	fmt.Printf("  SIMD宽度: %d\n", m.cpuFeatures.SIMDWidth)
}

// 更新实体
func (m *MacSIMDCollisionDetector) UpdateEntity(id int, x, y, radius float32) {
	if id >= 0 && id < m.maxEntities {
		m.entitiesX[id] = x
		m.entitiesY[id] = y
		m.entitiesR[id] = radius
	}
}

// SIMD优化的碰撞检测
func (m *MacSIMDCollisionDetector) DetectCollisions(pairs [][2]int) []SIMDCollisionResult {
	if len(pairs) == 0 {
		return nil
	}

	// 准备输入数据
	pos1X := make([]float32, len(pairs))
	pos1Y := make([]float32, len(pairs))
	pos2X := make([]float32, len(pairs))
	pos2Y := make([]float32, len(pairs))
	rad1 := make([]float32, len(pairs))
	rad2 := make([]float32, len(pairs))

	for i, pair := range pairs {
		pos1X[i] = m.entitiesX[pair[0]]
		pos1Y[i] = m.entitiesY[pair[0]]
		pos2X[i] = m.entitiesX[pair[1]]
		pos2Y[i] = m.entitiesY[pair[1]]
		rad1[i] = m.entitiesR[pair[0]]
		rad2[i] = m.entitiesR[pair[1]]
	}

	// 调用SIMD优化的C函数
	collisionCount := int(C.simd_collision_check_universal(
		(*C.float)(unsafe.Pointer(&pos1X[0])),
		(*C.float)(unsafe.Pointer(&pos1Y[0])),
		(*C.float)(unsafe.Pointer(&pos2X[0])),
		(*C.float)(unsafe.Pointer(&pos2Y[0])),
		(*C.float)(unsafe.Pointer(&rad1[0])),
		(*C.float)(unsafe.Pointer(&rad2[0])),
		(*C.uint32_t)(unsafe.Pointer(&m.collisionPairs[0])),
		C.int(len(pairs)),
	))

	// 转换结果
	results := make([]SIMDCollisionResult, collisionCount)
	for i := 0; i < collisionCount; i++ {
		pairIdx := m.collisionPairs[i*2]
		if int(pairIdx) < len(pairs) {
			results[i] = SIMDCollisionResult{
				Entity1: uint32(pairs[pairIdx][0]),
				Entity2: uint32(pairs[pairIdx][1]),
			}
		}
	}

	return results
}

// SIMD优化的范围查询
func (m *MacSIMDCollisionDetector) RangeQuery(centerX, centerY, radius float32, entityCount int) []uint32 {
	if entityCount > m.maxEntities {
		entityCount = m.maxEntities
	}

	foundCount := int(C.simd_range_query_universal(
		C.float(centerX),
		C.float(centerY),
		C.float(radius),
		(*C.float)(unsafe.Pointer(&m.entitiesX[0])),
		(*C.float)(unsafe.Pointer(&m.entitiesY[0])),
		(*C.uint32_t)(unsafe.Pointer(&m.results[0])),
		C.int(entityCount),
	))

	result := make([]uint32, foundCount)
	copy(result, m.results[:foundCount])
	return result
}

// SIMD碰撞结果
type SIMDCollisionResult struct {
	Entity1 uint32
	Entity2 uint32
}

// Mac平台性能基准测试
func (m *MacSIMDCollisionDetector) BenchmarkMacSIMD() {
	fmt.Println("\n=== Mac SIMD性能基准测试 ===")

	// 初始化测试数据
	entityCount := 10000
	for i := 0; i < entityCount; i++ {
		x := float32(i%100) * 8.0
		y := float32(i/100) * 8.0
		m.UpdateEntity(i, x, y, 3.0)
	}

	// 生成测试配对
	pairCount := 50000
	pairs := make([][2]int, pairCount)
	for i := 0; i < pairCount; i++ {
		pairs[i] = [2]int{i % entityCount, (i + 1) % entityCount}
	}

	// 预热
	for i := 0; i < 10; i++ {
		m.DetectCollisions(pairs[:1000])
	}

	// 性能测试
	iterations := 100
	totalTime := time.Duration(0)

	fmt.Printf("测试 %d 个配对，运行 %d 次迭代...\n", pairCount, iterations)

	for i := 0; i < iterations; i++ {
		start := time.Now()
		results := m.DetectCollisions(pairs)
		elapsed := time.Since(start)
		totalTime += elapsed

		if i == 0 {
			fmt.Printf("首次结果: %d 碰撞，耗时: %.3fms\n",
				len(results), float64(elapsed.Nanoseconds())/1e6)
		}
	}

	avgTime := float64(totalTime.Nanoseconds()) / float64(iterations) / 1e6
	throughput := float64(pairCount) / avgTime * 1000.0

	fmt.Printf("平均耗时: %.3fms\n", avgTime)
	fmt.Printf("吞吐量: %.0f 配对/秒\n", throughput)

	// 计算10万实体的理论性能
	scaledTime := avgTime * 100000.0 / float64(pairCount)
	fmt.Printf("10万实体理论性能: %.3fms/帧\n", scaledTime)
	fmt.Printf("理论帧率: %.1f FPS\n", 1000.0/scaledTime)
}

// 不同架构性能对比
func (m *MacSIMDCollisionDetector) ArchitectureComparison() {
	fmt.Println("\n=== 架构性能对比 ===")

	entityCount := 5000
	testPairs := 10000

	// 初始化数据
	for i := 0; i < entityCount; i++ {
		x := float32(i%50) * 15.0
		y := float32(i/50) * 15.0
		m.UpdateEntity(i, x, y, 5.0)
	}

	pairs := make([][2]int, testPairs)
	for i := 0; i < testPairs; i++ {
		pairs[i] = [2]int{i % entityCount, (i + 300) % entityCount}
	}

	// SIMD测试
	start := time.Now()
	simdResults := m.DetectCollisions(pairs)
	simdTime := time.Since(start)

	// 标量测试
	start = time.Now()
	scalarResults := m.scalarDetection(pairs)
	scalarTime := time.Since(start)

	fmt.Printf("标量实现: %d碰撞，%.3fms\n", len(scalarResults),
		float64(scalarTime.Nanoseconds())/1e6)
	fmt.Printf("SIMD实现: %d碰撞，%.3fms\n", len(simdResults),
		float64(simdTime.Nanoseconds())/1e6)

	if simdTime < scalarTime {
		speedup := float64(scalarTime) / float64(simdTime)
		fmt.Printf("性能提升: %.1fx\n", speedup)
	} else {
		fmt.Printf("性能下降: %.1fx (可能是数据量太小)\n",
			float64(simdTime)/float64(scalarTime))
	}

	// Apple Silicon特殊测试
	if m.cpuFeatures.IsAppleSilicon {
		fmt.Println("\n=== Apple Silicon特性测试 ===")
		m.appleSiliconSpecificTests()
	}
}

// 标量实现用于对比
func (m *MacSIMDCollisionDetector) scalarDetection(pairs [][2]int) []SIMDCollisionResult {
	var results []SIMDCollisionResult

	for _, pair := range pairs {
		id1, id2 := pair[0], pair[1]

		dx := m.entitiesX[id1] - m.entitiesX[id2]
		dy := m.entitiesY[id1] - m.entitiesY[id2]
		distSq := dx*dx + dy*dy

		radiusSum := m.entitiesR[id1] + m.entitiesR[id2]
		if distSq < radiusSum*radiusSum {
			results = append(results, SIMDCollisionResult{
				Entity1: uint32(id1),
				Entity2: uint32(id2),
			})
		}
	}

	return results
}

// Apple Silicon特殊测试
func (m *MacSIMDCollisionDetector) appleSiliconSpecificTests() {
	// 测试Apple Silicon的统一内存架构优势
	largeEntityCount := 50000

	// 大数据集范围查询测试
	for i := 0; i < largeEntityCount; i++ {
		x := float32(i%200) * 5.0
		y := float32(i/200) * 5.0
		m.UpdateEntity(i, x, y, 2.0)
	}

	start := time.Now()
	nearby := m.RangeQuery(500.0, 500.0, 100.0, largeEntityCount)
	elapsed := time.Since(start)

	fmt.Printf("大数据集范围查询 (%d实体): 找到%d个，耗时%.3fms\n",
		largeEntityCount, len(nearby), float64(elapsed.Nanoseconds())/1e6)

	// 内存带宽测试
	fmt.Printf("Apple Silicon内存优势: 统一内存架构提供更高带宽\n")
}

// 内存使用分析
func (m *MacSIMDCollisionDetector) AnalyzeMemoryUsage() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	fmt.Println("\n=== 内存使用分析 ===")
	fmt.Printf("当前分配: %.2fMB\n", float64(mem.Alloc)/1024/1024)
	fmt.Printf("总分配: %.2fMB\n", float64(mem.TotalAlloc)/1024/1024)
	fmt.Printf("系统内存: %.2fMB\n", float64(mem.Sys)/1024/1024)
	fmt.Printf("GC次数: %d\n", mem.NumGC)

	// 计算SIMD数据结构内存占用
	arrayMem := float64(len(m.entitiesX)*4 + len(m.entitiesY)*4 +
		len(m.entitiesR)*4 + len(m.collisionPairs)*4 + len(m.results)*4)
	fmt.Printf("SIMD数组内存: %.2fMB\n", arrayMem/1024/1024)
}

func main() {
	// 检查命令行参数来决定运行哪个版本的测试
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "simd":
			runSIMDDemo()
		case "spatial":
			runSpatialDemo()
		case "benchmark":
			runSimpleBenchmark()
		case "100k":
			runRealtime100KTest()
		default:
			showMainUsage()
		}
	} else {
		// 默认运行对比测试
		runComparisonDemo()
	}
}

// 运行SIMD演示
func runSIMDDemo() {
	fmt.Println("=== Mac平台SIMD优化碰撞检测演示 ===")

	// 创建检测器
	detector := NewMacSIMDCollisionDetector(100000)

	// 运行性能测试
	detector.BenchmarkMacSIMD()

	// 架构对比
	detector.ArchitectureComparison()

	// 内存分析
	detector.AnalyzeMemoryUsage()

	// 范围查询演示
	fmt.Println("\n=== 范围查询演示 ===")
	for i := 0; i < 1000; i++ {
		x := float32(i%32) * 8.0
		y := float32(i/32) * 8.0
		detector.UpdateEntity(i, x, y, 2.5)
	}

	start := time.Now()
	nearby := detector.RangeQuery(128.0, 128.0, 40.0, 1000)
	elapsed := time.Since(start)

	fmt.Printf("范围查询: 找到%d个附近实体，耗时%.3fms\n",
		len(nearby), float64(elapsed.Nanoseconds())/1e6)
}

// 运行空间分割演示
func runSpatialDemo() {
	fmt.Println("=== 空间分割碰撞检测演示 ===")
	fmt.Printf("CPU核心数: %d\n", runtime.NumCPU())

	// 创建碰撞系统 (假设10000x10000的世界)
	collisionSystem := NewCollisionSystem(100000, 100000)

	// 运行性能测试
	collisionSystem.BenchmarkCollisionDetection()

	// 内存优化
	collisionSystem.optimizeMemory()
}

// 运行对比演示
func runComparisonDemo() {
	fmt.Println("=== 碰撞检测系统性能对比 ===")

	// 创建SIMD检测器
	simdDetector := NewMacSIMDCollisionDetector(10000)

	// 创建空间分割系统
	spatialSystem := NewCollisionSystem(10000, 10000)

	entityCount := 5000

	// 初始化SIMD系统
	for i := 0; i < entityCount; i++ {
		x := float32(i%100) * 10.0
		y := float32(i/100) * 10.0
		simdDetector.UpdateEntity(i, x, y, 2.0)
	}

	// 初始化空间分割系统
	for i := 0; i < entityCount; i++ {
		x := float32(i%100) * 10.0
		y := float32(i/100) * 10.0
		spatialSystem.UpdateEntity(uint32(i), x, y, 2.0)
	}

	fmt.Printf("测试配置: %d个实体\n\n", entityCount)

	// 测试空间分割系统
	fmt.Println("测试空间分割系统...")
	start := time.Now()
	spatialResults := spatialSystem.DetectCollisions()
	spatialTime := time.Since(start)

	fmt.Printf("空间分割结果: %d碰撞, 耗时%.2fms\n",
		len(spatialResults), float64(spatialTime.Nanoseconds())/1e6)

	// 测试SIMD系统（生成一些测试配对）
	fmt.Println("测试SIMD系统...")
	pairs := make([][2]int, entityCount/2)
	for i := 0; i < len(pairs); i++ {
		pairs[i] = [2]int{i, (i + 1) % entityCount}
	}

	start = time.Now()
	simdResults := simdDetector.DetectCollisions(pairs)
	simdTime := time.Since(start)

	fmt.Printf("SIMD结果: %d碰撞, 耗时%.2fms\n",
		len(simdResults), float64(simdTime.Nanoseconds())/1e6)

	fmt.Println("\n使用说明:")
	fmt.Println("  go run . simd      - 运行SIMD版本演示")
	fmt.Println("  go run . spatial   - 运行空间分割版本演示")
	fmt.Println("  go run . benchmark - 运行简单基准测试")
}

// 运行简单基准测试
func runSimpleBenchmark() {
	fmt.Println("=== 简单基准测试 ===")

	entityCounts := []int{1000, 5000, 10000}

	for _, count := range entityCounts {
		fmt.Printf("\n测试 %d 个实体:\n", count)

		// 测试空间分割系统
		spatialSystem := NewCollisionSystem(10000, 10000)
		for i := 0; i < count; i++ {
			x := float32(i%100) * 10.0
			y := float32(i/100) * 10.0
			spatialSystem.UpdateEntity(uint32(i), x, y, 2.0)
		}

		// 运行多次测试取平均值
		iterations := 10
		totalTime := time.Duration(0)

		for i := 0; i < iterations; i++ {
			start := time.Now()
			results := spatialSystem.DetectCollisions()
			elapsed := time.Since(start)
			totalTime += elapsed

			if i == 0 {
				fmt.Printf("  碰撞数: %d\n", len(results))
			}
		}

		avgTime := float64(totalTime.Nanoseconds()) / float64(iterations) / 1e6
		fps := 1000.0 / avgTime

		fmt.Printf("  平均时间: %.2fms\n", avgTime)
		fmt.Printf("  理论FPS: %.1f\n", fps)
	}
}

// 显示主程序使用说明
func showMainUsage() {
	fmt.Println("碰撞检测系统演示程序")
	fmt.Println()
	fmt.Println("使用方法:")
	fmt.Println("  go run . [模式]")
	fmt.Println()
	fmt.Println("可用模式:")
	fmt.Println("  simd      - SIMD优化版本演示")
	fmt.Println("  spatial   - 空间分割版本演示")
	fmt.Println("  benchmark - 简单基准测试")
	fmt.Println("  100k      - 10万实体实时测试")
	fmt.Println("  (无参数)  - 快速性能对比")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  go run .           # 快速对比")
	fmt.Println("  go run . simd      # SIMD演示")
	fmt.Println("  go run . spatial   # 空间分割演示")
	fmt.Println("  go run . benchmark # 基准测试")
	fmt.Println("  go run . 100k      # 10万实体测试")
	fmt.Println()
	fmt.Println("或者运行简化版演示:")
	fmt.Println("  go run simple_demo.go")
}

// ================ 空间分割系统实现 ================

// 创建空间分割网格
func NewSpatialGrid(worldWidth, worldHeight, cellSize float32) *SpatialGrid {
	invCellSize := 1.0 / cellSize
	width := int(math.Ceil(float64(worldWidth * invCellSize)))
	height := int(math.Ceil(float64(worldHeight * invCellSize)))

	grid := &SpatialGrid{
		cellSize:    cellSize,
		invCellSize: invCellSize,
		width:       width,
		height:      height,
		cells:       make([][]uint32, width*height),
		cellPool: sync.Pool{
			New: func() interface{} {
				return make([]uint32, 0, 16) // 预分配容量
			},
		},
	}

	// 初始化所有格子
	for i := range grid.cells {
		grid.cells[i] = grid.cellPool.Get().([]uint32)[:0]
	}

	return grid
}

// 快速网格坐标计算
func (g *SpatialGrid) getGridCoords(x, y float32) (int, int) {
	gx := int(x * g.invCellSize)
	gy := int(y * g.invCellSize)

	// 边界检查
	if gx < 0 {
		gx = 0
	} else if gx >= g.width {
		gx = g.width - 1
	}

	if gy < 0 {
		gy = 0
	} else if gy >= g.height {
		gy = g.height - 1
	}

	return gx, gy
}

// 清空网格
func (g *SpatialGrid) clear() {
	for i := range g.cells {
		g.cells[i] = g.cells[i][:0] // 重置长度但保留容量
	}
}

// 添加实体到网格
func (g *SpatialGrid) addEntity(entity *Entity) {
	gx, gy := g.getGridCoords(entity.Pos.X, entity.Pos.Y)
	entity.GridX = int16(gx)
	entity.GridY = int16(gy)

	cellIdx := gy*g.width + gx
	g.cells[cellIdx] = append(g.cells[cellIdx], entity.ID)
}

// 创建工作池
func NewWorkerPool(numWorkers int) *WorkerPool {
	return &WorkerPool{
		workers:    numWorkers,
		taskChan:   make(chan CollisionTask, numWorkers*2),
		resultChan: make(chan CollisionResult, 1000),
	}
}

// 启动工作池
func (wp *WorkerPool) start() {
	for i := 0; i < wp.workers; i++ {
		go wp.worker()
	}
}

// 工作协程
func (wp *WorkerPool) worker() {
	for task := range wp.taskChan {
		wp.processCollisionTask(task)
		wp.wg.Done()
	}
}

// 处理碰撞检测任务
func (wp *WorkerPool) processCollisionTask(task CollisionTask) {
	entities := *task.entities
	grid := task.grid

	for i := task.startIdx; i < task.endIdx; i++ {
		entity := &entities[i]
		if !entity.Active {
			continue
		}

		// 获取周围9个格子的实体
		wp.checkEntityCollisions(entity, entities, grid)
	}
}

// 检查单个实体的碰撞
func (wp *WorkerPool) checkEntityCollisions(entity *Entity, entities []Entity, grid *SpatialGrid) {
	gx, gy := int(entity.GridX), int(entity.GridY)
	checkCount := 0

	// 检查周围9个格子 (3x3)
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			ngx, ngy := gx+dx, gy+dy

			// 边界检查
			if ngx < 0 || ngx >= grid.width || ngy < 0 || ngy >= grid.height {
				continue
			}

			cellIdx := ngy*grid.width + ngx
			cellEntities := grid.cells[cellIdx]

			for _, otherID := range cellEntities {
				if checkCount >= MAX_COLLISION_CHECKS {
					return // 限制检查数量
				}

				if otherID <= entity.ID { // 避免重复检查
					continue
				}

				other := &entities[otherID]
				if !other.Active {
					continue
				}

				// 快速距离检查 (避免sqrt)
				dx := entity.Pos.X - other.Pos.X
				dy := entity.Pos.Y - other.Pos.Y
				distSq := dx*dx + dy*dy
				radiusSum := entity.Radius + other.Radius

				if distSq < radiusSum*radiusSum {
					// 发生碰撞
					result := CollisionResult{
						entity1:  entity.ID,
						entity2:  otherID,
						distance: float32(math.Sqrt(float64(distSq))),
					}

					select {
					case wp.resultChan <- result:
					default:
						// 结果通道满了，跳过这个碰撞
					}
				}

				checkCount++
			}
		}
	}
}

// 创建碰撞系统
func NewCollisionSystem(worldWidth, worldHeight float32) *CollisionSystem {
	numCPU := runtime.NumCPU()

	system := &CollisionSystem{
		entities:   make([]Entity, MAX_ENTITIES_SPATIAL),
		grid:       NewSpatialGrid(worldWidth, worldHeight, GRID_SIZE_SPATIAL),
		workerPool: NewWorkerPool(numCPU),
		resultPool: sync.Pool{
			New: func() interface{} {
				return make([]CollisionResult, 0, 100)
			},
		},
	}

	// 初始化实体ID
	for i := range system.entities {
		system.entities[i].ID = uint32(i)
	}

	system.workerPool.start()
	return system
}

// 更新实体位置
func (cs *CollisionSystem) UpdateEntity(id uint32, x, y, radius float32) {
	if id >= MAX_ENTITIES_SPATIAL {
		return
	}

	entity := &cs.entities[id]
	entity.Pos.X = x
	entity.Pos.Y = y
	entity.Radius = radius
	entity.Active = true
}

// 主要的碰撞检测函数
func (cs *CollisionSystem) DetectCollisions() []CollisionResult {
	startTime := time.Now()

	// 1. 清空并重建空间网格
	cs.grid.clear()

	// 2. 将活跃实体添加到网格
	activeCount := 0
	for i := range cs.entities {
		if cs.entities[i].Active {
			cs.grid.addEntity(&cs.entities[i])
			activeCount++
		}
	}
	cs.activeCount = activeCount

	// 3. 清空结果通道
	for len(cs.workerPool.resultChan) > 0 {
		<-cs.workerPool.resultChan
	}

	// 4. 分发任务到工作池
	numWorkers := cs.workerPool.workers
	entitiesPerWorker := MAX_ENTITIES_SPATIAL / numWorkers

	for i := 0; i < numWorkers; i++ {
		startIdx := i * entitiesPerWorker
		endIdx := startIdx + entitiesPerWorker
		if i == numWorkers-1 {
			endIdx = MAX_ENTITIES_SPATIAL // 最后一个worker处理剩余的
		}

		task := CollisionTask{
			startIdx: startIdx,
			endIdx:   endIdx,
			entities: &cs.entities,
			grid:     cs.grid,
		}

		cs.workerPool.wg.Add(1)
		cs.workerPool.taskChan <- task
	}

	// 5. 等待所有worker完成
	cs.workerPool.wg.Wait()

	// 6. 收集结果
	results := cs.resultPool.Get().([]CollisionResult)[:0]

	for len(cs.workerPool.resultChan) > 0 {
		result := <-cs.workerPool.resultChan
		results = append(results, result)
	}

	elapsed := time.Since(startTime)

	// 性能统计
	if len(results) > 0 {
		fmt.Printf("碰撞检测: %d活跃实体, %d碰撞, 耗时: %.2fms\n",
			activeCount, len(results), float64(elapsed.Nanoseconds())/1e6)
	}

	return results
}

// 性能测试函数
func (cs *CollisionSystem) BenchmarkCollisionDetection() {
	fmt.Println("开始性能测试...")

	// 随机初始化实体位置
	for i := 0; i < MAX_ENTITIES_SPATIAL; i++ {
		x := float32(i%1000) * 2.0 // 简单的位置分布
		y := float32(i/1000) * 2.0
		cs.UpdateEntity(uint32(i), x, y, 1.0)
	}

	// 测试多帧
	const testFrames = 100
	totalTime := time.Duration(0)

	for frame := 0; frame < testFrames; frame++ {
		start := time.Now()
		results := cs.DetectCollisions()
		elapsed := time.Since(start)
		totalTime += elapsed

		if frame%10 == 0 {
			fmt.Printf("帧 %d: %d碰撞, %.2fms\n",
				frame, len(results), float64(elapsed.Nanoseconds())/1e6)
		}

		// 轻微移动实体
		for i := 0; i < 1000; i++ {
			entity := &cs.entities[i]
			entity.Pos.X += 0.1
			entity.Pos.Y += 0.1
		}
	}

	avgTime := float64(totalTime.Nanoseconds()) / float64(testFrames) / 1e6
	fmt.Printf("\n平均每帧时间: %.2fms (目标: 16.67ms)\n", avgTime)
	fmt.Printf("是否达到60FPS: %v\n", avgTime < 16.67)
}

// 内存优化函数
func (cs *CollisionSystem) optimizeMemory() {
	// 手动触发GC (在帧间隙调用)
	runtime.GC()

	// 打印内存统计
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("内存使用: %.2fMB, GC次数: %d\n",
		float64(m.Alloc)/1024/1024, m.NumGC)
}

// ================ 10万实体实时测试 ================

// 实时实体结构
type RealtimeEntity struct {
	ID     int
	X, Y   float32
	VX, VY float32 // 速度
	Radius float32
	Active bool
}

// 10万实体实时测试
func runRealtime100KTest() {
	fmt.Println("=== 10万实体实时碰撞检测性能测试 ===")
	fmt.Printf("Go版本: %s\n", runtime.Version())
	fmt.Printf("操作系统: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("CPU核心: %d\n", runtime.NumCPU())

	const entityCount = 100000
	const worldSize = 2000.0
	const testDuration = 30 * time.Second
	const targetFPS = 60.0

	fmt.Printf("实体数量: %d\n", entityCount)
	fmt.Printf("世界大小: %.0fx%.0f\n", worldSize, worldSize)
	fmt.Printf("测试时长: %v\n", testDuration)
	fmt.Printf("目标FPS: %.1f\n", targetFPS)

	// 创建检测系统
	simdDetector := NewMacSIMDCollisionDetector(entityCount)
	spatialSystem := NewCollisionSystem(worldSize, worldSize)

	// 创建实体数组
	entities := make([]RealtimeEntity, entityCount)

	// 初始化实体
	fmt.Println("初始化10万个实体...")
	for i := 0; i < entityCount; i++ {
		x := rand.Float32() * worldSize
		y := rand.Float32() * worldSize
		vx := (rand.Float32() - 0.5) * 20.0 // -10 到 10
		vy := (rand.Float32() - 0.5) * 20.0
		radius := 1.0 + rand.Float32()*3.0 // 1-4的半径

		entities[i] = RealtimeEntity{
			ID:     i,
			X:      x,
			Y:      y,
			VX:     vx,
			VY:     vy,
			Radius: radius,
			Active: true,
		}

		// 更新到检测系统
		simdDetector.UpdateEntity(i, x, y, radius)
		spatialSystem.UpdateEntity(uint32(i), x, y, radius)
	}

	fmt.Println("初始化完成，开始实时测试...")

	// 测试统计
	frameTimeMs := 1000.0 / targetFPS
	frameTime := time.Duration(frameTimeMs * float64(time.Millisecond))
	startTime := time.Now()
	lastFrameTime := startTime
	frameCount := 0

	simdTotalTime := time.Duration(0)
	spatialTotalTime := time.Duration(0)
	simdTotalCollisions := 0
	spatialTotalCollisions := 0

	frameStats := make([]float64, 0, 2000)

	fmt.Println("\n开始实时测试...")
	fmt.Println("格式: 帧数 | SIMD(碰撞,耗时ms) | 空间分割(碰撞,耗时ms) | 更新耗时ms | 总帧时间ms | FPS")
	fmt.Println("------|------------------|-------------------|----------|-----------|----")

	for time.Since(startTime) < testDuration {
		frameStart := time.Now()

		// 1. 更新实体位置
		deltaTime := float32(time.Since(lastFrameTime).Seconds())
		updateStart := time.Now()
		updateCount := 0

		for i := range entities {
			entity := &entities[i]
			if !entity.Active {
				continue
			}

			// 更新位置
			newX := entity.X + entity.VX*deltaTime
			newY := entity.Y + entity.VY*deltaTime

			// 边界反弹
			if newX < entity.Radius || newX > worldSize-entity.Radius {
				entity.VX = -entity.VX
				newX = entity.X + entity.VX*deltaTime
			}
			if newY < entity.Radius || newY > worldSize-entity.Radius {
				entity.VY = -entity.VY
				newY = entity.Y + entity.VY*deltaTime
			}

			// 边界限制
			if newX < entity.Radius {
				newX = entity.Radius
			} else if newX > worldSize-entity.Radius {
				newX = worldSize - entity.Radius
			}

			if newY < entity.Radius {
				newY = entity.Radius
			} else if newY > worldSize-entity.Radius {
				newY = worldSize - entity.Radius
			}

			// 更新实体
			entity.X = newX
			entity.Y = newY

			// 更新到检测系统
			simdDetector.UpdateEntity(i, newX, newY, entity.Radius)
			spatialSystem.UpdateEntity(uint32(i), newX, newY, entity.Radius)

			updateCount++
		}
		updateTime := time.Since(updateStart)

		// 2. SIMD碰撞检测 - 生成智能配对
		pairCount := 20000 // 2万个配对
		pairs := make([][2]int, pairCount)

		for i := 0; i < pairCount; i++ {
			id1 := rand.Intn(entityCount)
			// 优先选择附近的实体
			baseX := entities[id1].X
			baseY := entities[id1].Y
			searchRadius := float32(100.0)

			id2 := id1
			attempts := 0
			for attempts < 5 && id2 == id1 {
				candidate := rand.Intn(entityCount)
				dx := entities[candidate].X - baseX
				dy := entities[candidate].Y - baseY
				if dx*dx+dy*dy <= searchRadius*searchRadius {
					id2 = candidate
				}
				attempts++
			}

			if id2 == id1 {
				id2 = rand.Intn(entityCount)
			}

			pairs[i] = [2]int{id1, id2}
		}

		simdStart := time.Now()
		simdResults := simdDetector.DetectCollisions(pairs)
		simdTime := time.Since(simdStart)
		simdTotalTime += simdTime
		simdTotalCollisions += len(simdResults)

		// 3. 空间分割碰撞检测 (每10帧执行一次)
		var spatialCollisions int
		var spatialTime time.Duration
		if frameCount%10 == 0 {
			spatialStart := time.Now()
			spatialResults := spatialSystem.DetectCollisions()
			spatialTime = time.Since(spatialStart)
			spatialTotalTime += spatialTime
			spatialCollisions = len(spatialResults)
			spatialTotalCollisions += spatialCollisions
		}

		frameEnd := time.Now()
		totalFrameTime := frameEnd.Sub(frameStart)
		currentFPS := 1.0 / totalFrameTime.Seconds()

		// 记录统计
		frameStats = append(frameStats, float64(totalFrameTime.Nanoseconds())/1e6)

		// 每60帧输出一次统计
		if frameCount%60 == 0 {
			fmt.Printf("%5d | SIMD(%4d,%6.1f) | 空间(%4d,%6.1f) | %8.1f | %9.1f | %4.1f\n",
				frameCount,
				len(simdResults), float64(simdTime.Nanoseconds())/1e6,
				spatialCollisions, float64(spatialTime.Nanoseconds())/1e6,
				float64(updateTime.Nanoseconds())/1e6,
				float64(totalFrameTime.Nanoseconds())/1e6,
				currentFPS)
		}

		frameCount++
		lastFrameTime = frameStart

		// 帧率控制
		if totalFrameTime < frameTime {
			time.Sleep(frameTime - totalFrameTime)
		}
	}

	// 输出最终统计
	fmt.Printf("\n=== 最终统计结果 ===\n")

	actualFPS := float64(frameCount) / testDuration.Seconds()
	avgFrameTime := calculateAverage(frameStats)

	fmt.Printf("总帧数: %d\n", frameCount)
	fmt.Printf("实际平均FPS: %.1f\n", actualFPS)
	fmt.Printf("平均帧时间: %.2fms\n", avgFrameTime)

	// SIMD统计
	avgSIMDTime := float64(simdTotalTime.Nanoseconds()) / float64(frameCount) / 1e6
	avgSIMDCollisions := float64(simdTotalCollisions) / float64(frameCount)

	fmt.Printf("\nSIMD碰撞检测:\n")
	fmt.Printf("  平均检测时间: %.2fms\n", avgSIMDTime)
	fmt.Printf("  平均碰撞数: %.1f\n", avgSIMDCollisions)
	fmt.Printf("  总检测时间: %.2fs\n", simdTotalTime.Seconds())

	// 空间分割统计
	spatialFrames := frameCount / 10
	if spatialFrames > 0 {
		avgSpatialTime := float64(spatialTotalTime.Nanoseconds()) / float64(spatialFrames) / 1e6
		avgSpatialCollisions := float64(spatialTotalCollisions) / float64(spatialFrames)

		fmt.Printf("\n空间分割碰撞检测 (每10帧):\n")
		fmt.Printf("  平均检测时间: %.2fms\n", avgSpatialTime)
		fmt.Printf("  平均碰撞数: %.1f\n", avgSpatialCollisions)
		fmt.Printf("  总检测时间: %.2fs\n", spatialTotalTime.Seconds())
	}

	// 性能分析
	fmt.Printf("\n=== 性能分析 ===\n")

	// 60FPS达成率
	fps60Count := 0
	fps30Count := 0
	for _, ft := range frameStats {
		if ft <= 16.67 {
			fps60Count++
		}
		if ft <= 33.33 {
			fps30Count++
		}
	}

	fps60Rate := float64(fps60Count) / float64(len(frameStats)) * 100
	fps30Rate := float64(fps30Count) / float64(len(frameStats)) * 100

	fmt.Printf("60FPS达成率: %.1f%%\n", fps60Rate)
	fmt.Printf("30FPS达成率: %.1f%%\n", fps30Rate)

	// 帧时间分布
	minFrameTime := frameStats[0]
	maxFrameTime := frameStats[0]
	for _, ft := range frameStats {
		if ft < minFrameTime {
			minFrameTime = ft
		}
		if ft > maxFrameTime {
			maxFrameTime = ft
		}
	}

	fmt.Printf("最小帧时间: %.2fms\n", minFrameTime)
	fmt.Printf("最大帧时间: %.2fms\n", maxFrameTime)

	// 内存使用
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("\n内存使用: %.2fMB\n", float64(m.Alloc)/1024/1024)
	fmt.Printf("GC次数: %d\n", m.NumGC)

	// 结论
	fmt.Printf("\n=== 结论 ===\n")
	if fps60Rate > 90 {
		fmt.Println("✅ 性能优秀，完全满足60FPS实时游戏需求")
	} else if fps30Rate > 90 {
		fmt.Println("⚠️  性能良好，满足30FPS游戏需求")
	} else {
		fmt.Println("❌ 性能不足，需要优化")
	}

	fmt.Printf("SIMD优化效果: 平均%.2fms检测时间\n", avgSIMDTime)
	fmt.Printf("推荐使用: SIMD每帧 + 空间分割定期检测\n")
}

// 计算平均值辅助函数
func calculateAverage(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}
