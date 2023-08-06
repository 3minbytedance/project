package service

import (
	"fmt"
	"project/dao/mysql"
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
	user, exist := mysql.FindUserByID(int(userId))
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
		user, exist := mysql.FindUserByID(int(comment.UserId))
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
	user, exist := mysql.FindUserByID(int(comment.UserId))
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
