# gotranslate - Google Translate library written by Go

It's a pity that Google translate API has no free quota anymore. But the ajax translate API is always available. `gotranslate` is based on the ajax one, and provides a convenient way to use it. BTW, if you are using Google Translate in commercial project, purchasing Google translate service is still strongly recommended.

## Install

`go get github.com/liudanking/gotranslate`

## Usage

```go
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

## Supported Languages

From [ISO839-1](https://cloud.google.com/translate/docs/languages)

```go
var _supportedLangs = []string{"af", "sq", "am", "ar", "hy", "az", "eu", "be", "bn", "bs",
	"bg", "ca", "ceb", "ny", "zh-CN", "zh-TW", "co", "hr", "cs", "da", "nl", "en",
	"eo", "et", "tl", "fi", "fr", "fy", "gl", "ka", "de", "el", "gu", "ht", "ha",
	"haw", "iw", "hi", "hmn", "hu", "is", "ig", "id", "ga", "it", "ja", "jw", "kn",
	"kk", "km", "ko", "ku", "ky", "lo", "la", "lv", "lt", "lb", "mk", "mg", "ms",
	"ml", "mt", "mi", "mr", "mn", "my", "ne", "no", "ps", "fa", "pl", "pt", "ma",
	"ro", "ru", "sm", "gd", "sr", "st", "sn", "sd", "si", "sk", "sl", "so", "es",
	"su", "sw", "sv", "tg", "ta", "te", "th", "tr", "uk", "ur", "uz", "vi", "cy",
	"xh", "yi", "yo", "zu"}
```

## TODO

* ~~Cache TKK (maybe)~~