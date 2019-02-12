package dnspod

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/isayme/go-ddns/src/conf"
	"github.com/isayme/go-ddns/src/request"
	logger "github.com/isayme/go-logger"
)

// DNSPod ...
type DNSPod struct {
	config conf.DNSPod
}

// NewDNSPod ...
func NewDNSPod(cfg conf.DNSPod) *DNSPod {
	return &DNSPod{
		config: cfg,
	}
}

// Sleep ...
func (dnspod *DNSPod) Sleep() {
	time.Sleep(time.Duration(dnspod.config.Interval))
}

// UpdateIP update ip
func (dnspod *DNSPod) UpdateIP(ip string) error {
	recordID, err := dnspod.getRecordID()
	if err != nil {
		return err
	}

	resp, err := request.Request(request.Option{
		Method: http.MethodPost,
		URL:    "https://dnsapi.cn/Record.DDns",
		Form: request.M{
			"login_token": dnspod.config.Token,
			"domain":      dnspod.config.Domain,
			"sub_domain":  dnspod.config.SubDomain,
			"record_id":   recordID,
			"record_line": "默认",
			"format":      "json",
			"value":       ip,
		},
	})
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	logger.Debugf("update ip response: %s", string(respBody))

	if resp.StatusCode >= 300 {
		return fmt.Errorf("update ip fail: %s", string(respBody))
	}

	return nil
}

type recordListResponse struct {
	Records []struct {
		ID string `json:"id"`
	}
}

func (dnspod *DNSPod) getRecordID() (string, error) {
	resp, err := request.Request(request.Option{
		Method: http.MethodPost,
		URL:    "https://dnsapi.cn/Record.List",
		Form: request.M{
			"login_token": dnspod.config.Token,
			"domain":      dnspod.config.Domain,
			"sub_domain":  dnspod.config.SubDomain,
			"format":      "json",
		},
	})
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	logger.Debugf("get record id response: %s", string(respBody))

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("get record id fail: %s", string(respBody))
	}

	var result recordListResponse
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return "", err
	}

	if len(result.Records) > 0 {
		logger.Debugw("record id found", "subdomain", dnspod.config.SubDomain, "recordId", result.Records[0].ID)
		return result.Records[0].ID, nil
	}
	return "", fmt.Errorf("record id not found for subdomain: %s", dnspod.config.SubDomain)
}
