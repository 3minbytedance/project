package service

import (
	"fmt"
	"log"
	"project/dao/mysql"
	"project/dao/redis"
	"project/models"
	"time"
)

func AddComment(videoId, userId uint, content string) (models.CommentResponse, error) {
	// 评论信息
	commentData := models.Comment{
		VideoId: videoId,
		UserId:  userId,
		Content: content,
	}
	// 新增评论，返回评论id
	_, err := mysql.AddComment(&commentData)
	if err != nil {
		return models.CommentResponse{}, err
	}

	go func() {
		isSetKey, _ := checkAndSetRedisCommentKey(videoId)
		if isSetKey {
			return
		}
		// 更新commentCount
		err = redis.IncrementCommentCountByVideoId(videoId)
		if err != nil {
			log.Printf("更新videoId为%v的评论数失败  %v\n", videoId, err.Error())
		}
	}()

	// 查询user
	user, exist := GetUserInfoByUserId(uint(userId))
	if !exist {
		fmt.Println("根据评论中的user_id找用户失败, 评论ID为：", commentData.ID)
		return models.CommentResponse{}, err
	}

	// 封装返回数据
	var commentResp models.CommentResponse
	commentResp.Id = int64(commentData.ID)
	commentResp.User = user
	commentResp.Content = content
	commentResp.CreateDate = TranslateTime(commentData.CreatedAt.Unix())

	return commentResp, nil
}

// GetCommentList isLogged参数是为了返回用户信息中是否和自己关注
func GetCommentList(videoId uint, isLogged bool, userId uint) ([]models.CommentResponse, error) {
	// 1、根据videoId查询数据库，获取comments信息
	comments, err := mysql.FindCommentsByVideoId(videoId)
	if err != nil {
		fmt.Println("根据视频ID取评论失败")
		return nil, err
	}

	commentList := make([]models.CommentResponse, 0)
	for _, comment := range comments {
		user, _ := GetUserInfoByUserId(comment.UserId)
		if isLogged {
			user.IsFollow = IsInMyFollowList(userId, comment.UserId)
		}
		commentResp := models.CommentResponse{
			Id:         int64(comment.ID),
			User:       user,
			Content:    comment.Content,
			CreateDate: TranslateTime(comment.CreatedAt.Unix()),
		}
		commentList = append(commentList, commentResp)
	}

	go func() {
		err = redis.SetCommentCountByVideoId(videoId, int64(len(commentList)))
		if err != nil {
			log.Println("将评论数写入redis失败：", err.Error())
		}
	}()

	return commentList, nil
}

// TranslateTime 返回mm-dd格式
func TranslateTime(createTime int64) string {
	t := time.Unix(createTime, 0)
	return t.Format("01-02")
}

func DeleteComment(videoId, userId, commentId uint) (models.CommentResponse, error) {

	// 查询冗余字段
	// 查询comment
	comment, err := mysql.FindCommentById(commentId)
	if err != nil {
		fmt.Println("查询评论失败")
		return models.CommentResponse{}, err
	}

	// 查询user
	user, exist := GetUserInfoByUserId(comment.UserId)
	if !exist {
		log.Println("根据评论中的user_id找用户失败")
	}

	// 封装返回数据
	var commentResp models.CommentResponse
	commentResp.Id = int64(comment.ID)
	commentResp.User = user
	commentResp.Content = comment.Content
	commentResp.CreateDate = TranslateTime(comment.CreatedAt.Unix())

	// 1、 redis评论数-1
	err = redis.DecrementCommentCountByVideoId(videoId)
	if err != nil {
		log.Println("redis评论数-1失败")
		return models.CommentResponse{}, err
	}

	// 2、 mysql删除comment
	err = mysql.DeleteCommentById(commentId)
	if err != nil {
		fmt.Println("删除Comment失败")
		return commentResp, err
	}

	return commentResp, nil
}

// GetCommentCount 根据视频ID获取视频的评论数
func GetCommentCount(videoId uint) int64 {
	isSetKey, count := checkAndSetRedisCommentKey(videoId)
	if isSetKey {
		return count
	}
	// 从redis中获取评论数
	count, err := redis.GetCommentCountByVideoId(videoId)
	if err != nil {
		log.Println("redis获取评论数失败", err)
	}
	return count
}

// checkAndSetRedisCommentKey
// 返回true表示不存在这个key，并设置key
// 返回false表示已存在这个key，cnt数返回0
func checkAndSetRedisCommentKey(videoId uint) (isSet bool, count int64) {
	var cnt int64
	if !redis.IsExistVideoField(videoId, redis.CommentCountField) {
		// 获取最新commentCount
		cnt, err := mysql.GetCommentCnt(videoId)
		if err != nil {
			log.Println("mysql获取评论数失败", err)
		}
		// 设置最新commentCount
		err = redis.SetCommentCountByVideoId(videoId, cnt)
		if err != nil {
			log.Println("redis更新评论数失败", err)
		}
		return true, cnt
	}
	return false, cnt
}
