package response

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	messageSuccess string = "Success"
	messageFailed  string = "Failed"
)

type BaseResponse struct {
	ResponseMessage string `json:"response_message"`
	ResponseData    any    `json:"response_data,omitempty"`
}

func (s BaseResponse) Success(ctx *gin.Context) {
	s.ResponseMessage = messageSuccess
	ctx.JSON(http.StatusOK, s)
}

func (s BaseResponse) Failed(ctx *gin.Context, msg any) {
	s.ResponseMessage = fmt.Sprintf("%s : %+v", messageFailed, msg)
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, s)
}
