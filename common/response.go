package common

import "github.com/gin-gonic/gin"

var (
	SuccessMessage = "Success"
)

func BadRequestResponse(err error) gin.H{
	return gin.H{"error" : err.Error()}
}

func SuccesResponse(msg string) gin.H {
	return gin.H{"message" : msg}
}

func SuccesWithData(msg string , data interface{}) gin.H {
	return gin.H{"message" : msg, "data": data }
}