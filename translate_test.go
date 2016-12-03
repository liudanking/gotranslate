package gotranslate

import (
	"net/http"
	"net/url"
	"testing"
)

func TestTranslate(t *testing.T) {
	pf := func(r *http.Request) (*url.URL, error) {
		purl, _ := url.Parse("http://127.0.0.1:6152")
		return purl, nil
	}
	gt, err := New(TRANSLATE_CN_ADDR, pf)
	if err != nil {
		t.Fatal(err)
	}
	ret, err := gt.Translate("zh-CN", "zh-TW", "abc 中国人")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("ret:%+v", ret)
}
