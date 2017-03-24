package gotranslate

import (
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestTranslate(t *testing.T) {
	pf := func(r *http.Request) (*url.URL, error) {
		purl, _ := url.Parse("http://127.0.0.1:6152")
		return purl, nil
	}
	gt, err := New(TRANSLATE_COM_ADDR, pf)
	if err != nil {
		t.Fatal(err)
	}
	q := "逗斗车 - 四川愣娃闯帝都 逗比天团再聚首 - 余洋"
	// q := "abc中国人"
	ret, err := gt.Translate("zh-CN", "zh-TW", q)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("ret:%+v", ret)

}

func TestSimpleTranslate(t *testing.T) {
	pf := func(r *http.Request) (*url.URL, error) {
		purl, _ := url.Parse("http://127.0.0.1:6152")
		return purl, nil
	}
	gt, err := New(TRANSLATE_COM_ADDR, pf)
	if err != nil {
		t.Fatal(err)
	}
	q := "丈母娘唠车 别再加价买丰田埃尔法了！开车如开房的奔驰... 胡永平"
	// q := "abc中国人"
	for i := 0; i < 10; i++ {
		start := time.Now()

		ret, err := gt.SimpleTranslate("zh-CN", "zh-TW", q)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("ret:%+v", ret)

		log.Printf("cost:%v", time.Since(start))
	}

}
