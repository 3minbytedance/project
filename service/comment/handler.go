package main

import (
	"context"
	"douyin/kitex_gen/comment"
)

// CommentServiceImpl implements the last service interface defined in the IDL.
type CommentServiceImpl struct{}

// CommentAction implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CommentAction(ctx context.Context, req *comment.CommentActionRequest) (resp *comment.CommentActionResponse, err error) {
	// TODO: Your code here...
	resp = new(comment.CommentActionResponse)
	return
}

// GetCommentList implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) GetCommentList(ctx context.Context, req *comment.CommentListRequest) (resp *comment.CommentListRequest, err error) {
	// TODO: Your code here...
	return
}

// GetCommentCount implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) GetCommentCount(ctx context.Context, req *comment.CommentCountRequest) (resp *comment.CommentCountResponse, err error) {
	// TODO: Your code here...
	return
}
