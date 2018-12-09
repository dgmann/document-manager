package http

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int
	Data       interface{}
	context    *gin.Context
}

func (r *Response) JSON() {
	r.context.Status(r.StatusCode)
	r.context.Header("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(r.context.Writer).Encode(r.Data); err != nil {
		r.context.AbortWithError(400, err)
	}
}
