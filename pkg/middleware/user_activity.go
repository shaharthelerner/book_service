package user_activity_middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"pkg/service/pkg/consts"
	"pkg/service/pkg/interfaces"
	"pkg/service/pkg/models/request"
)

func Middleware(usersHandler interfaces.UsersHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Skip to the next handler if the path is the user activity endpoint
		if ctx.FullPath() == consts.GetUserActivityUrlPath {
			ctx.Next()
			return
		}

		if ctx.Request.Body == nil || ctx.Request.Body == http.NoBody {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no request body"})
			return
		}

		origBody, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(origBody)) // Return the original body for the next read

		req := request.Common{}
		if err = ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(origBody)) // Return the original body for the next read

		userAction := request.CreateUserAction{
			Username: req.Username,
			Method:   ctx.Request.Method,
			Route:    ctx.FullPath(),
		}

		if err = usersHandler.SaveUserAction(userAction); err != nil {
			log.Printf("failed to save user action: %s", err.Error())
		}

		ctx.Next()
	}
}
