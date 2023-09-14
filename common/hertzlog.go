package common

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func AccessLog() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		hlog.CtxTracef(c, "status=%d method=%s full_path=%s QueryString=%s post=%s",
			ctx.Response.StatusCode(),
			ctx.Request.Header.Method(), ctx.Request.URI().PathOriginal(), ctx.Request.QueryString(),
			ctx.Request.PostArgString())
	}
}
