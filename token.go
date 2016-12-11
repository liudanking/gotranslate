package gotranslate

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"
)

var _tkkReg *regexp.Regexp

func init() {
	_tkkReg = regexp.MustCompile(`TKK\=eval\('\(\(function\(\)\{var\s+a\\x3d(-?\d+);var\s+b\\x3d(-?\d+);return\s+(\d+)\+`)
}

func (gt *GTranslate) getTKK() (int, int, error) {
	data, err := gt.httpRequest("GET", gt.srvAddr, nil)
	if err != nil {
		return 0, 0, err
	}
	return findTKK(string(data))
}

// TODO
func (gt *GTranslate) setTKK(tkk typeTKK) {
	gt.tkkMtx.Lock()
	gt.tkk = tkk
	gt.tkkMtx.Unlock()
}

func (gt *GTranslate) initTKK() error {
	data, err := gt.httpRequest("GET", gt.srvAddr, nil)
	if err != nil {
		return err
	}
	h1, h2, err := findTKK(string(data))
	if err != nil {
		return err
	}
	gt.setTKK(typeTKK{h1: h1, h2: h2})
	return nil
}

func (gt *GTranslate) updateTKK() {
	var data []byte
	var err error
	var h1 int
	var h2 int
	for {
		data, err = gt.httpRequest("GET", gt.srvAddr, nil)
		if err != nil {
			log.Printf("get tkk failed:%v", err)
			goto next
		}
		h1, h2, err = findTKK(string(data))
		if err != nil {
			log.Printf("try to find tkk from [%s] failed:%v", data, err)
		}
		gt.setTKK(typeTKK{h1: h1, h2: h2})
	next:
		time.Sleep(60 * time.Second)
	}
}

func findTKK(s string) (int, int, error) {
	rets := _tkkReg.FindStringSubmatch(s)
	if len(rets) != 4 {
		return 0, 0, errors.New("can't find TKK")
	}

	a, _ := strconv.Atoi(rets[1])
	b, _ := strconv.Atoi(rets[2])
	c, _ := strconv.Atoi(rets[3])
	d := a + b
	return c, d, nil
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
	g := make([]rune, len(q))
	for i := 0; i < len(q); i++ {
		g[i] = rune(q[i])
	}
	// qRune := []rune(q)
	// for i := 0; i < len(qRune); i++ {
	// 	c := qRune[i]
	// 	if 128 > c {
	// 		g = append(g, c)
	// 		continue
	// 	}
	// 	if 2048 > c {
	// 		g = append(g, (c>>6)|192)
	// 		continue
	// 	}
	// 	if 55296 == (c&64512) && i+1 < len(qRune) && 56320 == (qRune[i+1]&64512) {
	// 		c2 := 65536 + ((c & 1023) << 10) + (qRune[i+1] & 1023)
	// 		g = append(g, (c2>>18)|240)
	// 		g = append(g, ((c2>>12)&63)|128)
	// 		i++
	// 	} else {
	// 		g = append(g, (c>>12)|224)
	// 		g = append(g, (c>>6)&63|128)
	// 		g = append(g, (c&63)|128)
	// 	}
	// }

	a := int32(h)
	for i := 0; i < len(g); i++ {
		a += int32(g[i])
		a = bf(a, "+-a^+6")
	}
	a = bf(a, "+-3^+b+-f")
	a ^= int32(h2)
	// fmt.Printf("a^h2:%d\n", a)
	var aInt64 int64
	if 0 > a {
		aInt64 = int64((int(a) & 2147483647) + 2147483648)
	}
	// fmt.Printf("0 > a:%d\n", a)

	aInt64 %= 1e6

	s := fmt.Sprintf("%d.%d", aInt64, aInt64^int64(h))

	// log.Printf("tk:%s", s)
	return s
}

func bf(a int32, s string) int32 {
	// log.Printf("round a in:%d, %s", a, s)
	b := []rune(s)
	for i := 0; i < len(b)-2; i += 3 {
		c := int32(b[i+2])
		if 'a' <= b[i+2] {
			c = int32(b[i+2]) - 87
		} else {
			cInt, _ := strconv.Atoi(string([]byte{byte(c)}))
			c = int32(cInt)
		}
		// log.Printf("c1:%d", c)
		if '+' == b[i+1] {
			c = int32(uint32(a) >> uint32(c))
		} else {
			c = a << uint32(c)
		}
		// log.Printf("c2:%d", c)
		if '+' == b[i] {
			a = int32((int(a) + int(c)) & 4294967295)
		} else {
			a = a ^ c
		}
		// log.Printf("%d:c:%d", i, c)
		// log.Printf("%d:a:%d", i, a)
	}
	// log.Printf("round a out:%d", a)
	return a
}
