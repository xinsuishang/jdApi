package conf

import "sync"

var (
	ErrResp sync.Map
	WG      sync.WaitGroup
	ApiConf = map[string]string{
		"access_token": "",
		"app_key":      "",
		"method":       "",
		"app_secret":   "",
		"v":            "",
	}
)

const (
	URI = "https://api.jd.com/routerjson?"
)
