package utils

import (
	"net/http"
	"strings"
)

func GetLanguageFromRequest(req *http.Request) string {
	lang := req.Header.Get("Accept-Language")
	lang = strings.ToLower(lang)
	if lang == "bn" || lang == "ban" || lang == "bang" || lang == "bangla" || lang == "bengali" {
		return "bn"
	} else {
		return "en"
	}
}
