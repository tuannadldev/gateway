package wrapper

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

type ctxKey int

const (
	wrapperKey ctxKey = iota
)

// context wrapper of gin.Context
type Context struct {
	*gin.Context
}

// ToContext ...
func (c *Context) ToContext() context.Context {
	return context.WithValue(context.Background(), wrapperKey, c)
}

type HandlerFunc func(ctx *Context)

// WithContext HOC for wrapping the gin context
func WithContext(handler HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		wrappedContext := &Context{
			ctx,
		}
		handler(wrappedContext)
	}
}

func ContextWithWrapper(ctx context.Context, wrapperCtx *Context) context.Context {
	return context.WithValue(ctx, wrapperKey, wrapperCtx)
}

const (
	hiddenContent         = "<HIDDEN>"
	ignoreContent         = "<IGNORE>"
	emptyContentTag       = "<EMPTY>"
	contentSizeLimitation = 10000
)

func isIgnoreRequestBody(ctx *gin.Context) bool {
	contentSize := ctx.Request.ContentLength
	// Ingore content too large
	if contentSize == -1 || contentSize >= contentSizeLimitation {
		return true
	}

	contentType := ctx.ContentType()
	return contentType == gin.MIMEMultipartPOSTForm
}

func GetRequestBody(ctx *gin.Context) string {
	log := ctxzap.Extract(ctx).Sugar()

	requestBody := hiddenContent

	if isIgnoreRequestBody(ctx) {
		return requestBody
	}

	buf, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.With("err", err).Error("can't read body content from request")
		return ignoreContent
	}
	readCloser1 := io.NopCloser(bytes.NewBuffer(buf))
	readCloser2 := io.NopCloser(bytes.NewBuffer(buf))
	ctx.Request.Body = readCloser2

	// convert readCloser1 to string
	bytesBuffer := new(bytes.Buffer)
	_, err = bytesBuffer.ReadFrom(readCloser1)
	if err != nil {
		log.Error("can't read byte array from reader", err)
		return ignoreContent
	}
	requestBody = bytesBuffer.String()
	if requestBody == "" {
		// return tag to easy filter
		return emptyContentTag
	}
	return requestBody
}

// JSONData ...
func (c *Context) JSONData(statusCode int, data interface{}) {
	c.JSONResponse(statusCode, data, "", "", "")
}

// JSONError ...
func (c *Context) JSONError(statusCode int, errorCode string, errorMessage string) {
	c.JSONResponse(statusCode, nil, errorCode, errorMessage, "")
}

// JSONResponse ...
func (c *Context) JSONResponse(statusCode int, data interface{}, errorCode string, errorMessage string, errorField string) {
	resp := gin.H{}

	success := true
	if data != nil {
		resp["data"] = data
	}

	if errorCode != "" && errorMessage != "" {
		success = false
		resp["errors"] = []gin.H{
			{
				"message": errorMessage,
				"extensions": gin.H{
					"status_code": statusCode,
					"code":        errorCode,
					"field":       errorField,
				},
			},
		}
	}
	resp["success"] = success

	c.JSON(statusCode, resp)

	if statusCode != http.StatusOK {
		c.Abort()
	}
}
