package utils

import "github.com/microcosm-cc/bluemonday"

var p *bluemonday.Policy = bluemonday.UGCPolicy()

func FilterHTML(input string) string {
	// 过滤不受信任的 HTML
	sanitizedInput := p.Sanitize(input)

	return sanitizedInput
}
