package ip

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/isayme/go-ddns/src/conf"
	"github.com/isayme/go-ddns/src/request"
	logger "github.com/isayme/go-logger"
)

// Get get current ip
func Get() (string, error) {
	urls := conf.Get().IPURLs

	for _, url := range urls {
		resp, err := request.Request(request.Option{
			URL: url,
		})
		if err != nil {
			logger.Warnw("get ip fail", "err", err, "url", url)
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode >= 300 {
			logger.Warnw("get ip fail", "err", err, "url", url, "statusCode", resp.StatusCode)
			continue
		}

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		return string(bytes.TrimSpace(respBody)), nil
	}

	return "", fmt.Errorf("get ip fail")
}
