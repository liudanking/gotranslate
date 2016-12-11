package gotranslate

import (
	"fmt"
	"log"

	"testing"
)

func TestFindTkk(t *testing.T) {
	s := `MAX_ALTERNATIVES_ROUNDTRIP_RESULTS=1;TKK=eval('((function(){var a\x3d3951944662;var b\x3d-3371426310;return 411232+\x27.\x27+(a+b)})())');WEB_TRANSLATION_PATH='/translate';SIGNED_IN=false;USAGE='';`
	// TKK=eval('((function(){var a\x3d3951944662;var b\x3d-3371426310;return 411232+\x27.\x27+(a+b)})())');
	// "411232.580518352"
	c, d, err := findTKK(s)
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%d.%d", c, d) != "411232.580518352" {
		t.Fatalf("%d.%d != 411232.580518352", c, d)
	}
}

// var TKK = ((function() {
//   return "411232.580518352";
// })());

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

func TestTK(t *testing.T) {
	h := 411232
	h2 := 580518352
	ret := tk(h, h2, "abc")
	if ret != "911334.762246" {
		t.Fatalf("%s != 911334.762246", ret)
	}
	log.Printf("tk[%s]", ret)

	ret = tk(h, h2, "abc中国人")
	if ret != "557546.968586" {
		t.Fatalf("%s != 911334.762246", ret)
	}
	log.Printf("tk[%s]", ret)
}

func TestTK2(t *testing.T) {
	h := 411508
	h2 := 1550816266
	q := "逗斗车 - 四川愣娃闯帝都 逗比天团再聚首 - 余洋" // 467843.91383
	ret := tk(h, h2, q)
	if ret != "467843.91383" {
		t.Fatalf("%s != 467843.91383", ret)
	}
	log.Printf("tk[%s]", ret)
}

func TestBF(t *testing.T) {
	a := int32(1024)
	b := "+-a^+6"

	ret := bf(a, b)
	if ret != 1066000 {
		t.Fatalf("%d != 1066000", ret)
	}

	ret = bf(a, "+-3^+b+-f")
	if ret != 302130180 {
		t.Fatalf("%d != 302130180", ret)
	}
}
