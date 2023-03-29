package common

import "github.com/gin-gonic/gin"

func SuccessResponse() gin.H {
	return gin.H{
		"message": "ok",
	}
}

func MsgResponse(msg string) gin.H {
	return gin.H{
		"message": msg,
	}
}

func SingleResponse(key, value string) gin.H {
	return gin.H{
		key: value,
	}
}

func ListResponse[T any](data []T, total int64) gin.H {
	return gin.H{
		"list":  data,
		"total": total,
	}
}
