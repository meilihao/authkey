package util

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/i18n"
	"google.golang.org/grpc/metadata"
)

type CodeErr struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e CodeErr) Error() string {
	return fmt.Sprintf("code=%s, message=%s", e.Code, e.Message)
}

func I18nError(c *gin.Context, key string, args ...interface{}) error {
	if c.MustGet("i18n").(string) == "cn" {
		return CodeErr{
			Code:    key,
			Message: i18n.Tr("cn", key, args...),
		}
	} else {
		return CodeErr{
			Code:    key,
			Message: i18n.Tr("en", key, args...),
		}
	}
}

// for grpc server
func I18nError2(ctx context.Context, key string, args ...interface{}) error {
	if getLangFromCtx(ctx) == "cn" {
		return CodeErr{
			Code:    key,
			Message: i18n.Tr("cn", key, args...),
		}
	} else {
		return CodeErr{
			Code:    key,
			Message: i18n.Tr("en", key, args...),
		}
	}
}

func getLangFromCtx(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "en"
	}

	if len(md["lang"]) > 0 || md["lang"][0] == "cn" {
		return "cn"
	}

	return "en"
}
