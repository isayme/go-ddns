package ip

import (
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/isayme/go-ddns/src/conf"
	"github.com/isayme/go-ddns/src/request"
	logger "github.com/isayme/go-logger"
)

var ipReg *regexp.Regexp

func init() {
	// https://stackoverflow.com/questions/13386461/regex-for-ip-v4-address
	ipReg = regexp.MustCompile("(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)")
}

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

		ip := ipReg.Find(respBody)
		if len(ip) == 0 {
			logger.Warnw("get ip fail", "err", err, "url", url, "respBody", string(respBody))
			continue
		}

		return string(ip), nil
	}

	return "", fmt.Errorf("get ip fail")
}
