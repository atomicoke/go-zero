package reg

import (
	"dm.com/toolx/arr"
	"regexp"
	"testing"
)

var regExp = regexp.MustCompile(`(\w+):"(.*?)"`)

func TestFindAll(t *testing.T) {

	all := regExp.FindAllStringSubmatch(`json:"title" validate:"required,lt=60" label:"会议标题"`, -1)

	var mp = arr.NewMapX[string, string]()

	for _, v := range all {
		mp.Set(v[1], v[2])
	}

	t.Log(mp.Get("label", ""))
}
