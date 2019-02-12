package request

import (
	"fmt"
	"runtime"

	"github.com/isayme/go-ddns/src/util"
)

// UserAgent user agent
var UserAgent string

func init() {
	UserAgent = fmt.Sprintf("%s/%v Go/%s", util.Name, util.Version, runtime.Version())
}
