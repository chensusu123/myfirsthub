package errors

// 道具错误提示
func WrapItemKnowErr(err error, defaultErr *CodeError) error {
	errInfo, ok := err.(*CodeError)
	if !ok {
		return defaultErr
	}

	switch errInfo {
	case SILVER_REACH_MAX, ERR_SHIP_EXP_FULL, ERR_MERITORIOUS_SERVICE_FULL, STOREHOUSE_FULL_ERROR,
		ERR_SHIP_FULL_EQUIP_BAG, ERR_SHIP_LESS_EQUIP_BAG, BAG_FULL_ERROR, BAG_ITEM_LIMIT_MAX_ERROR:
		// 绿钞 , //航海术,// 功勋,//仓库，//装备，//装备 //背包满 //背包道具上限
		return errInfo
	case NotProcessItemType:
		// 未注册处理的道具类型
		return errInfo
	case ItemCheckFailure:
		return errInfo
	default:
		return defaultErr
	}

	return defaultErr
}
