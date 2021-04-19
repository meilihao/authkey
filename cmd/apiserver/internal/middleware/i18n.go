package middleware

import (
	"github.com/meilihao/water"
)

var (
	_i18n = "i18n"
)

func GetLang(c *water.Context) string {
	return c.Get(_i18n).(string)
}

func I18n() water.HandlerFunc {
	return func(c *water.Context) {
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
