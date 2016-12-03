# gotranslate - Google Translate library written by Go

It's a pity that Google translate API has no free quota anymore. But the ajax translate API is always available. `gotranslate` is based on the ajax one, and provides a convenient way to use it. BTW, if you are using Google Translate in commercial project, purchasing Google translate service is still strongly recommended.

## Install

`go get github.com/liudanking/gotranslate`

## Usage

```
package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/liudanking/gotranslate"
)

func main() {
	// use default translate settings
	ret, err := gotranslate.Translate("zh-CN", "en", "皮克斯：关于童心、勇气、创意和传奇")
	if err != nil {
		log.Printf("translate failed:%v", err)
		return
	}
	log.Printf("%+v", ret)

	pf := func(r *http.Request) (*url.URL, error) {
		purl, _ := url.Parse("http://127.0.0.1:6152")
		return purl, nil
	}
	// create a new translate using your proxy to fxck through GFW
	gt, err := gotranslate.New(gotranslate.TRANSLATE_COM_ADDR, pf)
	if err != nil {
		log.Printf("create gotranslate failed:%v", err)
		return
	}

	ret, err = gt.Translate("auto", "zh-TW", "台湾宝岛，富饶可人")
	if err != nil {
		log.Printf("translate failed:%v", err)
		return
	}
	log.Printf("%+v", ret)
}

```

## TODO

* Cache TKK (maybe)