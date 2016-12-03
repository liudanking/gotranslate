package gotranslate

import (
	"testing"
)

func TestTranslate(t *testing.T) {
	Translate("zh-CN", "zh-TW", "abc 中国人")
}
