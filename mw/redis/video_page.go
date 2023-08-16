package redis

const (
	Page     = "page"
	PageSize = "page_size"
)

func SetPageAndPageSize(tempId string, page int, pageSize int) error {
	key := VideoPage + tempId
	err := Rdb.HMSet(Ctx, key, Page, page, PageSize, pageSize).Err()
	return err
}

func GetPageAndPageSize(tempId string) (page int, pageSize int, err error) {
	key := VideoPage + tempId
	result, err := Rdb.HMGet(Ctx, key, Page, PageSize).Result()
	if err != nil {
		return 0, 0, err
	}
	page = result[0].(int)
	pageSize = result[1].(int)
	return page, pageSize, nil
}
