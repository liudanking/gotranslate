package gotranslate

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/liudanking/goutil/netutil"
)

const (
	GOOGLE_TRANSLATE_ADDR = "http://translate.google.cn"
)

var _tkkReg *regexp.Regexp

func init() {
	_tkkReg = regexp.MustCompile(`TKK\=eval\('\(\(function\(\)\{var\s+a\\x3d(-?\d+);var\s+b\\x3d(-?\d+);return\s+(\d+)\+`)
}

func getTKK() (int, int, error) {
	data, err := httpRequest("GET", GOOGLE_TRANSLATE_ADDR, nil)
	if err != nil {
		return 0, 0, err
	}
	c, d := findTKK(string(data))
	log.Printf("%d.%d", c, d)
	return c, d, nil
}

func httpRequest(method, addr string, params map[string]string) ([]byte, error) {
	pf := func(r *http.Request) (*url.URL, error) {
		purl, _ := url.Parse("http://127.0.0.1:6152")
		return purl, nil
	}
	data, code, err := netutil.DefaultHttpClient().RequestForm(method, addr, params).UserAgent(netutil.UA_SAFARI).
		Proxy(pf).
		DoByte()
	if err != nil {
		log.Printf("http request failed:[%d] %v", code, err)
	}
	return data, err
}

func findTKK(s string) (int, int) {
	rets := _tkkReg.FindStringSubmatch(s)
	if len(rets) != 4 {
		return 0, 0
	}
	// log.Printf("rets:%+v", rets)
	for k, v := range rets {
		log.Printf("%d:%s", k, v)
	}
	a, _ := strconv.Atoi(rets[1])
	b, _ := strconv.Atoi(rets[2])
	c, _ := strconv.Atoi(rets[3])
	d := a + b
	// log.Printf("%d.%d", c, d)
	return c, d
}

// javascript:
// https://developer.mozilla.org/zh-CN/docs/Web/JavaScript/Reference/Operators/Operator_Precedence
// https://developer.mozilla.org/zh-CN/docs/Web/JavaScript/Guide/Expressions_and_Operators
// http://stackoverflow.com/questions/18729405/how-to-convert-utf8-string-to-byte-array

/*
var TKK = ((function() {
  var a = 561666268;
  var b = 1526272306;
  return 406398 + '.' + (a + b);
})());
*/
// function b(a, b) {
//   for (var d = 0; d < b.length - 2; d += 3) {
//       var c = b.charAt(d + 2),
//           c = "a" <= c ? c.charCodeAt(0) - 87 : Number(c),
//           c = "+" == b.charAt(d + 1) ? a >>> c : a << c;
//       a = "+" == b.charAt(d) ? a + c & 4294967295 : a ^ c
//   }
//   return a
// }

// function tk(a) {
//     for (var e = TKK.split("."), h = Number(e[0]) || 0, g = [], d = 0, f = 0; f < a.length; f++) {
//         var c = a.charCodeAt(f);
//         128 > c ? g[d++] = c : (2048 > c ? g[d++] = c >> 6 | 192 : (55296 == (c & 64512) && f + 1 < a.length && 56320 == (a.charCodeAt(f + 1) & 64512) ? (c = 65536 + ((c & 1023) << 10) + (a.charCodeAt(++f) & 1023), g[d++] = c >> 18 | 240, g[d++] = c >> 12 & 63 | 128) : g[d++] = c >> 12 | 224, g[d++] = c >> 6 & 63 | 128), g[d++] = c & 63 | 128)
//     }
//     a = h;
//     for (d = 0; d < g.length; d++) a += g[d], a = b(a, "+-a^+6");
//     a = b(a, "+-3^+b+-f");
//     a ^= Number(e[1]) || 0;
//     0 > a && (a = (a & 2147483647) + 2147483648);
//     a %= 1E6;
//     return a.toString() + "." + (a ^ h)
// }

// </script>

func tk(h, h2 int, q string) string {
	qRune := []rune(q)
	g := []rune{}
	for i := 0; i < len(qRune); i++ {
		c := qRune[i]
		if 128 > c {
			g = append(g, c)
			continue
		}
		if 2048 > c {
			g = append(g, (c>>6)|192)
			continue
		}
		if 55296 == (c&64512) && i+1 < len(qRune) && 56320 == (qRune[i+1]&64512) {
			c2 := 65536 + ((c & 1023) << 10) + (qRune[i+1] & 1023)
			g = append(g, (c2>>18)|240)
			g = append(g, ((c2>>12)&63)|128)
			i++
		} else {
			g = append(g, (c>>12)|224)
			g = append(g, (c>>6)&63|128)
			g = append(g, (c&63)|128)
		}
	}
	// log.Printf("g:%+v", g)
	a := h
	for i := 0; i < len(g); i++ {
		a += int(g[i])
		a = bf(a, "+-a^+6")
	}
	// log.Printf("a:%d", a)
	a = bf(a, "+-3^+b+-f")
	a ^= h2
	if 0 > a {
		a = (a & 2147483647) + 2147483648
	}
	a %= 1e6

	s := fmt.Sprintf("%d.%d", a, a^h)

	log.Printf("tk:%s", s)
	return s
}

func bf(a int, s string) int {
	b := []rune(s)
	for i := 0; i < len(b)-2; i += 3 {
		c := int(b[i+2])
		if 'a' <= b[i+2] {
			c = int(b[i+2]) - 87
		} else {
			c, _ = strconv.Atoi(string([]byte{byte(c)}))
		}
		if '+' == b[i+1] {
			c = int(uint(a) >> uint(c))
		} else {
			c = a << uint(c)
		}
		if '+' == b[i] {
			a = (a + c) & 4294967295
		} else {
			a = a ^ c
		}
		// log.Printf("c in:%d", c)
		// log.Printf("a:%d", a)
	}
	return a
}
