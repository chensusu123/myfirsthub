package idgenerator

const (
	// 起始时间戳，2021-01-01 00:00:00 的毫秒时间戳
	epoch = 1609459200000
	// 时间戳所占位数
	timestampBits = 41
	// 序列号所占位数
	sequenceBits = 22
	// 时间戳向左移位数
	timestampShift = sequenceBits
	// 序列号最大值
	maxSequence = -1 ^ (-1 << sequenceBits)
	// 时间戳最大值
	maxTimestamp = -1 ^ (-1 << timestampBits)

	// 消息 ID 最大值
	MaxMessageID = (maxTimestamp << timestampShift) | maxSequence
)

// MessageID 基于传入的毫秒时间戳和序列号生成63位唯一消息ID
func MessageID(timestamp int64, sequence uint64) uint64 {
	// 计算相对时间戳
	relative := uint64(timestamp - epoch)

	// 确保时间戳和序列号在有效范围内
	if relative > maxTimestamp {
		relative = maxTimestamp
	}
	if sequence > maxSequence {
		sequence = maxSequence
	}

	// 组合时间戳和序列号生成消息 ID
	messageID := (relative << timestampShift) | sequence
	return messageID
}
