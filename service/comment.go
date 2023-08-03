package service

import (
	"fmt"
	"log"
	"project/dao/mysql"
	"project/dao/redis"
	"project/models"
	"time"
)

func AddComment(videoId, userId int64, content string) (models.CommentResponse, error) {
	// 评论信息
	commentResp := models.CommentResponse{}
	commentData := models.Comment{
		VideoId: videoId,
		UserId:  userId,
		Content: content,
	}
	// 新增评论，返回评论id
	_, err := mysql.AddComment(&commentData)
	if err != nil {
		return commentResp, err
	}
	// 查询user
	user, exist := models.FindUserByID(mysql.DB, int(userId))
	if !exist {
		fmt.Println("根据评论中的user_id找用户失败")
	}

	// 封装返回数据
	commentResp.Id = int64(commentData.ID)
	commentResp.User = user
	commentResp.Content = content
	commentResp.CreateDate = models.TranslateTime(commentData.CreatedAt.Unix(), time.Now().Unix())

	//commentResp.User
	//commentResp.Content
	//commentResp.CreateDate
	// TODO 更新视频信息，由于未确定视频表设计，延后再写
	//video, b := models.FindVideoByVideoId(videoId, content)
	//if !b {
	//	fmt.Println("未找到对应的视频")
	//} else {
	//	num := video.CommentCount + 1
	//	mysql.DB.Model(&video).Update("comment_count", strconv.Itoa(int(num)))
	//}
	return commentResp, nil

}

func GetCommentList(videoId int64) ([]models.CommentResponse, error) {
	comments, err := mysql.FindCommentsByVideoId(videoId)
	if err != nil {
		fmt.Println("根据视频ID取评论失败")
		return nil, err
	}
	commentList := make([]models.CommentResponse, 0)
	for _, comment := range comments {
		user, exist := models.FindUserByID(mysql.DB, int(comment.UserId))
		if !exist {
			fmt.Println("根据评论中的user_id找用户失败")
		}
		commentResp := models.CommentResponse{
			Id:         int64(comment.ID),
			User:       user,
			Content:    comment.Content,
			CreateDate: models.TranslateTime(comment.CreatedAt.Unix(), time.Now().Unix()),
		}
		commentList = append(commentList, commentResp)
	}
	return commentList, nil
}

func DeleteComment(videoId, userId, commentId int64) (models.CommentResponse, error) {
	// TODO：等后续视频表建立完成，再看是否需要进行其他操作
	commentResp := models.CommentResponse{}

	// 查询comment
	comment, err := mysql.FindCommentById(commentId)
	if err != nil {
		fmt.Println("查询评论失败")
		return commentResp, err
	}

	// 查询user
	user, exist := models.FindUserByID(mysql.DB, int(comment.UserId))
	if !exist {
		fmt.Println("根据评论中的user_id找用户失败")
	}

	// 封装返回数据
	commentResp.Id = int64(comment.ID)
	commentResp.User = user
	commentResp.Content = comment.Content
	commentResp.CreateDate = models.TranslateTime(comment.CreatedAt.Unix(), time.Now().Unix())

	// 删除comment
	err = mysql.DeleteCommentById(commentId)
	if err != nil {
		fmt.Println("删除Comment失败")
		return commentResp, err
	}

	return commentResp, nil

}

// GetCommentCount 根据视频ID获取视频的评论数
func GetCommentCount(videoId int64) (int64, error) {
	// 从redis中获取评论数
	count, err := redis.GetCommentCountByVideoId(videoId)
	if err != nil {
		log.Println("从redis中获取评论数失败：", err)
		return 0, err
	}
	// 缓存中有数据, 直接返回
	if count > 0 {
		log.Println("从redis中获取评论数成功：", count)
		return count, nil
	}
	// 缓存中没有数据，从数据库中获取
	count, err = mysql.GetCommentCnt(videoId)
	if err != nil {
		log.Println("从数据库中获取评论数失败：", err.Error())
		return 0, nil
	}
	log.Println("从数据库中获取评论数成功：", count)
	// FIXME 按道理, 这里获取评论数时, 不应该把评论的内容放在redis中,
	// 但是能够从redis获取评论数的前提是评论内容在redis中，
	// 这就会导致如果评论内容一直没在redis中, 那么每次都要从数据库中计算评论数
	// 所以，应该在redis中存储评论数，而不仅仅是评论内容
	return count, nil
}
