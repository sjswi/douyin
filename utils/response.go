package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

type Response struct {
	StatusCode uint   `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}
type UploadPictureResponse struct {
	Success int    `json:"success"`
	Message string `json:"message"`
	Url     string `json:"url"`
}

// Redirect 重定向
func Redirect(c *gin.Context, location string) {
	c.Redirect(http.StatusFound, location)
	return
}

// RedirectBack 重定向到上一次的页面
func RedirectBack(c *gin.Context) {
	referer := c.GetHeader("Referer")
	pathInfo := ""
	if referer == "" {
		pathInfo = "/"
	} else {
		u, _ := url.Parse(referer)
		pathInfo = u.Path
	}

	c.Redirect(http.StatusFound, pathInfo)
	return
}
