package jwt

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"project/utils"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// Auth 鉴权中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		//fmt.Println("token", token)
		// 没携带token
		if len(token) == 0 {

			// 没有token, 阻止后面函数执行
			c.Abort()
			c.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				StatusMsg:  "Unauthorized",
			})

		} else {

			claims, err := utils.ParseToken(token)
			fmt.Println("jwt test:", claims.ID)
			if err != nil {
				// token有误，阻止后面函数执行
				c.Abort()
				c.JSON(http.StatusUnauthorized, Response{
					StatusCode: -1,
					StatusMsg:  "Token Error",
				})
			} else {
				log.Println("token correct")
			}
			c.Set(utils.ContextUserIDKey, claims.ID)
			c.Next()

		}
	}
}

// AuthWithoutLogin 未登录情况，若携带token,解析用户id放入context;如果没有携带，则将用户id默认为0
func AuthWithoutLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		var userId uint
		if len(token) == 0 {
			// 没有token, 设置userId为0
			userId = 0
		} else {
			claims, err := utils.ParseToken(token)
			if err != nil {
				// token有误，阻止后面函数执行
				c.Abort()
				c.JSON(http.StatusUnauthorized, Response{
					StatusCode: -1,
					StatusMsg:  "Token Error",
				})
			} else {
				log.Println("token correct")
				userId = claims.ID
			}
			c.Set(utils.ContextUserIDKey, userId)
			c.Next()
		}
	}
}

// AuthBody 若token在请求体里，解析token
func AuthBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.PostFormValue("token")
		// 没携带token
		if len(token) == 0 {
			// 没有token, 阻止后面函数执行
			c.Abort()
			c.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				StatusMsg:  "Unauthorized",
			})
		} else {
			claims, err := utils.ParseToken(token)
			if err != nil {
				// token有误，阻止后面函数执行
				c.Abort()
				c.JSON(http.StatusUnauthorized, Response{
					StatusCode: -1,
					StatusMsg:  "Token Error",
				})
			} else {
				log.Println("token correct")
			}
			c.Set(utils.ContextUserIDKey, claims.ID)
			c.Next()
		}
	}
}
