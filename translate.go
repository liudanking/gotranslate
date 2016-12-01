package gotranslate

import (
	"fmt"
	"log"
	"net/url"
)

func Translate(q string) {
	// h1, h2, err := getTKK()
	// if err != nil {
	// 	log.Printf("get tkk error:%v", err)
	// 	return
	// }
	h1 := 411279
	h2 := 2826815461
	tkstr := tk(h1, h2, q)

	log.Printf("tk:%s", tkstr)

	q = url.QueryEscape(q)
	addr := fmt.Sprintf("http://translate.google.cn/translate_a/single?client=t&sl=zh-CN&tl=zh-TW&hl=zh-CN&dt=at&dt=bd&dt=ex&dt=ld&dt=md&dt=qca&dt=rw&dt=rm&dt=ss&dt=t&ie=UTF-8&oe=UTF-8&otf=2&ssel=0&tsel=0&kc=1&tk=%s&q=%s", tkstr, q)
	data, err := httpRequest("GET", addr, nil)
	log.Printf("%s %v", data, err)

}
