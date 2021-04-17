package middleware

import (
	"github.com/gin-gonic/gin"
)

var (
	_i18n = "i18n"
)

func GetLang(c *gin.Context) string {
	return c.MustGet(_i18n).(string)
}

func I18n() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.Request.Header.Get("Accept-Language")
		if lang == "" {
			lang = c.Query("lang")
		}
		if lang == "" || (lang != "zh" && lang != "en") {
			lang = "en"
		}

		c.Set(_i18n, lang)

		c.Next()
	}
}
