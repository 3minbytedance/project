package graphdb

import (
	"fmt"
	"go.uber.org/zap"
	"strconv"
)

func IsFriend(userId, followId uint) bool {
	query := "MATCH (v2)-->(v:person)-->(v2) where id(v) == $userId and id(v2)== $followId " +
		"RETURN COUNT(*) > 0 AS result;"

	params := map[string]interface{}{
		"userId":   int(userId),
		"followId": int(followId),
	}

	resp, err := sessionPool.ExecuteWithParameter(query, params)
	if err != nil {
		zap.L().Error("获取IsFollow失败", zap.Error(err))
		return false
	}
	result, err := resp.GetValuesByColName("result")
	follow, _ := result[0].AsBool()
	return follow
}

func AddFollow(userId, followId uint) error {
	addUserVertex := " INSERT VERTEX IF NOT EXISTS person() VALUES " + strconv.Itoa(int(userId)) + ":(); "
	addFollowVertex := " INSERT VERTEX IF NOT EXISTS person() VALUES " + strconv.Itoa(int(followId)) + ":(); "
	query := "INSERT EDGE follow () VALUES $userId->$followId:();"

	params := map[string]interface{}{
		"userId":   int(userId),
		"followId": int(followId),
	}
	_, err := sessionPool.Execute(addUserVertex)
	if err != nil {
		zap.L().Error("AddFollow失败", zap.Error(err))
		return err
	}
	_, err = sessionPool.Execute(addFollowVertex)
	if err != nil {
		zap.L().Error("AddFollow失败", zap.Error(err))
		return err
	}
	_, err = sessionPool.ExecuteWithParameter(query, params)
	if err != nil {
		zap.L().Error("AddFollow失败", zap.Error(err))
		return err
	}
	return nil
}

func DeleteFollowById(userId, followId uint) error {

	query := "DELETE EDGE follow $userId->$followId;"
	params := map[string]interface{}{
		"userId":   int(userId),
		"followId": int(followId),
	}
	_, err := sessionPool.ExecuteWithParameter(query, params)
	if err != nil {
		zap.L().Error("DeleteFollowById失败", zap.Error(err))
		return err
	}
	return nil
}

func GetFollowCnt(userId uint) (int64, error) {
	query := "MATCH (v)-->() where id(v) == " + strconv.Itoa(int(userId)) +
		"RETURN COUNT(*) AS result;"

	resp, err := sessionPool.Execute(query)
	if err != nil {
		zap.L().Error("GetFollowCnt失败", zap.Error(err))
		return 0, err
	}
	result, err := resp.GetValuesByColName("result")
	count, _ := result[0].AsInt()
	return count, nil
}

func GetFollowerCnt(userId uint) (int64, error) {
	query := "MATCH ()-->(v) where id(v) == " + strconv.Itoa(int(userId)) +
		"RETURN COUNT(*) AS result;"

	resp, err := sessionPool.Execute(query)
	if err != nil {
		zap.L().Error("GetFollowerCnt失败", zap.Error(err))
		return 0, err
	}
	result, err := resp.GetValuesByColName("result")
	count, _ := result[0].AsInt()
	return count, nil
}

func GetFollowList(userId uint) ([]uint, error) {
	query := "MATCH (v)-->(v2) where id(v) ==" + strconv.Itoa(int(userId)) +
		"RETURN id(v2) AS result;"
	resp, err := sessionPool.Execute(query)
	if err != nil {
		return nil, err
	}
	followList := make([]uint, 0, resp.GetRowSize())
	result, _ := resp.GetValuesByColName("result")
	for _, res := range result {
		asInt, _ := res.AsInt()
		followList = append(followList, uint(asInt))
	}
	return followList, nil
}

func GetFollowerList(userId uint) ([]uint, error) {
	query := "MATCH (v2)-->(v) where id(v) ==" + strconv.Itoa(int(userId)) +
		"RETURN id(v2) AS result;"
	resp, err := sessionPool.Execute(query)
	if err != nil {
		return nil, err
	}
	followList := make([]uint, 0, resp.GetRowSize())
	result, _ := resp.GetValuesByColName("result")
	for _, res := range result {
		asInt, _ := res.AsInt()
		followList = append(followList, uint(asInt))
	}
	return followList, nil
}

func GetAllFollowEdge() error {
	query := "MATCH (s)-[e:follow]->(d) RETURN id(s), id(d);"
	resp, err := sessionPool.Execute(query)
	if err != nil {
		return err
	}
	result := resp.AsStringTable()
	// 第0行为表头，从[1:]开始遍历
	for _, res := range result[1:] {
		fmt.Println(res[0])
	}
	return nil
}
