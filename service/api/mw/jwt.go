package mw

import (
	"context"
	"douyin/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// Auth 鉴权中间件
func Auth() app.HandlerFunc {
	return func(ctx context.Context, rc *app.RequestContext) {
		token := rc.Query("token")
		// 没携带token
		if len(token) == 0 {
			// 没有token, 阻止后面函数执行
			rc.Abort()
			rc.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				StatusMsg:  "Unauthorized",
			})

		} else {
			claims, err := utils.ParseToken(token)
			if err != nil {
				// token有误，阻止后面函数执行
				rc.Abort()
				rc.JSON(http.StatusUnauthorized, Response{
					StatusCode: -1,
					StatusMsg:  "Token Error",
				})
			}
			rc.Set(utils.ContextUserIDKey, claims.ID)
			rc.Next(ctx)

		}
	}
}

// AuthWithoutLogin 未登录情况，若携带token,解析用户id放入context;如果没有携带，则将用户id默认为0
func AuthWithoutLogin() app.HandlerFunc {
	return func(ctx context.Context, rc *app.RequestContext) {
		token := rc.Query("token")
		var userId uint
		if len(token) == 0 {
			// 没有token, 设置userId为0
			userId = 0
		} else {
			claims, err := utils.ParseToken(token)
			if err != nil {
				// token有误，阻止后面函数执行
				rc.Abort()
				rc.JSON(http.StatusUnauthorized, Response{
					StatusCode: -1,
					StatusMsg:  "Token Error",
				})
			} else {
				userId = claims.ID
			}
			rc.Set(utils.ContextUserIDKey, userId)
			rc.Next(ctx)
		}
	}
}

// AuthBody 若token在请求体里，解析token
func AuthBody() app.HandlerFunc {
	return func(ctx context.Context, rc *app.RequestContext) {
		token := rc.PostForm("token")
		// 没携带token
		if len(token) == 0 {
			// 没有token, 阻止后面函数执行
			rc.Abort()
			rc.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				StatusMsg:  "Unauthorized",
			})
		} else {
			claims, err := utils.ParseToken(token)
			if err != nil {
				// token有误，阻止后面函数执行
				rc.Abort()
				rc.JSON(http.StatusUnauthorized, Response{
					StatusCode: -1,
					StatusMsg:  "Token Error",
				})
			}
			rc.Set(utils.ContextUserIDKey, claims.ID)
			rc.Next(ctx)
		}
	}
}

//var (
//	JwtMiddleware *jwt.HertzJWTMiddleware
//	identityKey   = "userId"
//	secretKey     = "secret-key"
//)
//
//func InitJWT() {
//	JwtMiddleware, _ = jwt.New(&jwt.HertzJWTMiddleware{
//		//设置签名密钥
//		Key: []byte(secretKey),
//		//设置token的获取源，可选header、query、cookie、param、form 默认header: Authorization
//		TokenLookup: "header: Authorization, query: token, form: token, cookie: jwt",
//		//用于设置从header中获取token时的前缀，默认为Bearer
//		TokenHeadName: "Bearer",
//		//设置获取当前时间的函数
//		TimeFunc: time.Now,
//		//token的过期时间
//		Timeout: time.Hour,
//		//最大token刷新时间，允许客户端在tokenTime+MaxRefresh内刷新token的有效时间，追加一个timeout时长
//		MaxRefresh: 24 * time.Hour,
//		//用于设置检索身份的键 默认identity
//		IdentityKey: identityKey,
//		//设置获取身份信息的函数,与PayloadFunc一致
//		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
//			claims := jwt.ExtractClaims(ctx, c)
//			return &model.User{
//				ID: uint(claims[identityKey].(float64)),
//			}
//		},
//		//设置登陆时认证用户信息的函数，这个函数的返回值 user.ID 将为后续生成 jwt token 提供 payload 数据源。
//		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
//			var loginStruct struct {
//				Username string `form:"username" json:"username" query:"username" vd:"(len($) > 5 && len($) < 30); msg:'Illegal format'"`
//				Password string `form:"password" json:"password" query:"password" vd:"(len($) > 5 && len($) < 30); msg:'Illegal format'"`
//			}
//			if err := c.BindAndValidate(&loginStruct); err != nil {
//				return nil, err
//			}
//			user, err := mysql.CheckUser(loginStruct.Username, loginStruct.Password)
//			if err != nil {
//				return nil, err
//			}
//			if user.ID > 0 {
//				return nil, errors.New("user already exists or wrong password")
//			}
//
//			return user.ID, nil
//		},
//		//登陆成功后为向token中添加自定义负载信息的函数,额外存储了用户id，如不设置则只存储token的过期时间和创建时间
//		PayloadFunc: func(data interface{}) jwt.MapClaims {
//			if v, ok := data.(uint); ok {
//				return jwt.MapClaims{
//					identityKey: v,
//				}
//			}
//			return jwt.MapClaims{}
//		},
//		//设置登陆的响应函数
//		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
//			c.JSON(http.StatusOK, utils.H{
//				"status_code": code,
//				"status_msg":  "success",
//				"user_id":     nil,
//				"token":       token,
//			})
//		},
//
//		//验证流程失败的响应函数
//		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
//			c.JSON(http.StatusOK, utils.H{
//				"status_code": code,
//				"status_msg":  "unauthorized",
//				"user_id":     nil,
//				"token":       nil,
//			})
//		},
//		//设置jwt校验流程发生错误时相应所包含的错误信息
//		HTTPStatusMessageFunc: func(e error, ctx context.Context, c *app.RequestContext) string {
//			hlog.CtxErrorf(ctx, "jwt biz err = %+v", e.Error())
//			return e.Error()
//		},
//	})
//}
