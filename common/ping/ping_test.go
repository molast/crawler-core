package ping

import (
	"testing"
)

func TestPing(t *testing.T) {

	t.Log(Ping("www.baidu.com", 5e9))
}
