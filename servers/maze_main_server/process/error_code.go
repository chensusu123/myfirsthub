package process

const (
	ERROR_CODE_NOT_IN_HIGH_AREA      int64 = 80001 //不在最高区域
	ERROR_CODE_USER_CANT_REBORN      int64 = 80002 //不能复活
	ERROR_CODE_REBORN_COST_NOT_MATCH int64 = 80003 //复活消耗不匹配
)

var errMsg = map[int64]string{
	ERROR_CODE_NOT_IN_HIGH_AREA: "非最高区域",
	ERROR_CODE_USER_CANT_REBORN: "用户不能复活",
}
