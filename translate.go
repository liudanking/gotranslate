package gotranslate

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/liudanking/goutil/netutil"
	"github.com/liudanking/goutil/strutil"
)

const (
	TRANSLATE_COM_ADDR = "https://translate.google.com"
	TRANSLATE_CN_ADDR  = "http://translate.google.cn"
)

// ISO839-1 https://cloud.google.com/translate/docs/languages
var _supportedLangs = []string{"af", "sq", "am", "ar", "hy", "az", "eu", "be", "bn", "bs",
	"bg", "ca", "ceb", "ny", "zh-CN", "zh-TW", "co", "hr", "cs", "da", "nl", "en",
	"eo", "et", "tl", "fi", "fr", "fy", "gl", "ka", "de", "el", "gu", "ht", "ha",
	"haw", "iw", "hi", "hmn", "hu", "is", "ig", "id", "ga", "it", "ja", "jw", "kn",
	"kk", "km", "ko", "ku", "ky", "lo", "la", "lv", "lt", "lb", "mk", "mg", "ms",
	"ml", "mt", "mi", "mr", "mn", "my", "ne", "no", "ps", "fa", "pl", "pt", "ma",
	"ro", "ru", "sm", "gd", "sr", "st", "sn", "sd", "si", "sk", "sl", "so", "es",
	"su", "sw", "sv", "tg", "ta", "te", "th", "tr", "uk", "ur", "uz", "vi", "cy",
	"xh", "yi", "yo", "zu"}

var defaultGTranslate *GTranslate

func init() {
	defaultGTranslate, _ = New(TRANSLATE_CN_ADDR, nil)
}

type GTranslate struct {
	srvAddr string
	proxy   func(r *http.Request) (*url.URL, error)
	// TODO, maybe
	tkkMtx sync.RWMutex
	tkk    typeTKK
}

type typeTKK struct {
	h1 int
	h2 int
}

func New(addr string, proxy func(r *http.Request) (*url.URL, error)) (*GTranslate, error) {
	if addr != TRANSLATE_CN_ADDR && addr != TRANSLATE_COM_ADDR {
		return nil, errors.New("addr not supported")
	}
	return &GTranslate{
		srvAddr: addr,
		proxy:   proxy,
	}, nil
}

type TranslateRet struct {
	Sentences []struct {
		Trans   string `json:"trans"`
		Orig    string `json:"orig"`
		Backend int    `json:"backend"`
	} `json:"sentences"`
	Src        string  `json:"src"`
	Confidence float64 `json:"confidence"`
	LdResult   struct {
		Srclangs            []string  `json:"srclangs"`
		SrclangsConfidences []float64 `json:"srclangs_confidences"`
		ExtendedSrclangs    []string  `json:"extended_srclangs"`
	} `json:"ld_result"`
}

// Translate translate q from sl to tl using default GTranslate
func Translate(sl, tl, q string) (*TranslateRet, error) {
	return defaultGTranslate.Translate(sl, tl, q)
}

// SimpleTranslate translate q to tl without test q sentences
func SimpleTranslate(sl, tl, q string) (string, error) {
	return defaultGTranslate.SimpleTranslate(sl, tl, q)
}

func (gt *GTranslate) Translate(sl, tl, q string) (*TranslateRet, error) {
	if sl != "auto" {
		if !strutil.StringIn(_supportedLangs, sl) {
			return nil, errors.New("source language not supported")
		}
	}

	if !strutil.StringIn(_supportedLangs, tl) {
		return nil, errors.New("target language not supported")
	}

	h1, h2, err := gt.getTKK()
	if err != nil {
		log.Printf("get tkk error:%v", err)
		return nil, err
	}
	// h1 = 411508
	// h2 = 1550816266

	tkstr := tk(h1, h2, q)
	// fmt.Printf("tk:%s", tkstr)
	// return nil, nil

	// https://translate.google.com/translate_a/single?client=t&sl=zh-CN&tl=zh-TW&hl=zh-CN&dt=at&dt=bd&dt=ex&dt=ld&dt=md&dt=qca&dt=rw&dt=rm&dt=ss&dt=t&ie=UTF-8&oe=UTF-8&otf=2&ssel=0&tsel=0&kc=1&tk=%s&q=%s
	var data []byte
	if false {
		q = url.QueryEscape(q)
		addr := fmt.Sprintf("%s/translate_a/single?client=t&dj=1&sl=zh-CN&tl=zh-TW&hl=zh-CN&dt=at&dt=bd&dt=ex&dt=ld&dt=md&dt=qca&dt=rw&dt=rm&dt=ss&dt=t&ie=UTF-8&oe=UTF-8&otf=2&ssel=0&tsel=0&kc=1&tk=%s&q=%s", gt.srvAddr, tkstr, q)
		data, err = gt.httpRequest("GET", addr, nil)
	} else {
		addr := fmt.Sprintf("%s/translate_a/single", gt.srvAddr)
		data, err = gt.httpRequest("GET", addr, gt.reqParams(sl, tl, tkstr, q))
	}
	// fmt.Printf("%s %v", data, err)

	ret := &TranslateRet{}
	err = json.Unmarshal(data, ret)
	return ret, err
}

func (gt *GTranslate) SimpleTranslate(sl, tl, q string) (string, error) {
	rsp, err := gt.Translate(sl, tl, q)
	if err != nil {
		return "", err
	}
	s := ""
	for _, sentence := range rsp.Sentences {
		s += sentence.Trans
	}
	return s, nil
}

func (gt *GTranslate) reqParams(sl, tl, tk, q string) map[string]interface{} {
	return map[string]interface{}{
		"client": "t",     // or gtx
		"sl":     sl,      // source language
		"tl":     tl,      // translated language
		"dj":     1,       // ensure return json is GoogleRes structure
		"ie":     "UTF-8", // input string encoding
		"oe":     "UTF-8", // output string encoding
		"tk":     tk,
		"q":      q,
		"dt":     []string{"t", "bd"}, // a list to add content to return json
		// possible dt values: correspond return json key
		// t: sentences
		// rm: sentences[1]
		// bd: dict
		// at: alternative_translations
		// ss: synsets
		// rw: related_words
		// ex: examples
		// ld: ld_result
	}
}

func (gt *GTranslate) httpRequest(method, addr string, params map[string]interface{}) ([]byte, error) {

	data, code, err := netutil.DefaultHttpClient().RequestForm(method, addr, params).UserAgent(netutil.UA_SAFARI).Proxy(gt.proxy).DoByte()
	if err != nil {
		log.Printf("http request failed:[%d] %v", code, err)
	}
	return data, err
}
