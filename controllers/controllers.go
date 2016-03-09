package controllers

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"strings"

	"github.com/nitrous-io/rise-server/models/oauthtoken"
	"github.com/nitrous-io/rise-server/models/user"

	"github.com/gin-gonic/gin"

	log "github.com/Sirupsen/logrus"
)

const (
	CurrentTokenKey = "current_token"
	CurrentUserKey  = "current_user"
)

func CurrentToken(c *gin.Context) *oauthtoken.OauthToken {
	ti, exists := c.Get(CurrentTokenKey)
	if ti == nil || !exists {
		return nil
	}

	t, ok := ti.(*oauthtoken.OauthToken)
	if !ok {
		return nil
	}
	return t
}

func CurrentUser(c *gin.Context) *user.User {
	ui, exists := c.Get(CurrentUserKey)
	if ui == nil || !exists {
		return nil
	}

	u, ok := ui.(*user.User)
	if !ok {
		return nil
	}
	return u
}

func InternalServerError(c *gin.Context, err error) {
	var (
		errMsg  = "internal server error"
		errHash string
	)

	if err != nil {
		errMsg = err.Error()
		errHash = fmt.Sprintf("%x", sha1.Sum([]byte(errMsg)))
	}

	req := c.Request

	fields := log.Fields{
		"req": fmt.Sprintf("%s %s", req.Method, req.URL.String()),
		"ip":  c.ClientIP(),
	}

	j := gin.H{
		"error": "internal_server_error",
	}

	if errHash != "" {
		fields["hash"] = errHash
		j["error_hash"] = errHash
	}

	if (req.Method == "POST" || req.Method == "PUT" || req.Method == "PATCH") && strings.Contains(c.ContentType(), "application/x-www-form-urlencoded") {
		if err := req.ParseForm(); err == nil {
			fields["form"] = req.PostForm.Encode()
		}
	}

	log.WithFields(fields).Error(errMsg)
	c.JSON(http.StatusInternalServerError, j)
}
