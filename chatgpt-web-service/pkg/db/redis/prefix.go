package redis

import "strings"

const Prifix = "chatgpt_web_service"

func Prefix_key(key string, parts ...string) string {
	key = Prifix + key
	if len(parts) != 0 {
		key += strings.Join(parts, "_")
	}
	return key
}
