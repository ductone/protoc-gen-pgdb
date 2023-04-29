package v1

import (
	"testing"
)

var toValidUTF8Tests = []struct {
	in  string
	out string
}{
	{"", ""},
	{"abc", "abc"},
	{"\uFDDD", "\uFDDD"},
	{"a\xffb", "a�b"},
	{"a\xffb\ufffd", "a�b�"},
	{"a☺\xffb☺\xC0\xAFc☺\xff", "a☺�b☺��c☺�"},
	{"\xC0\xAF", "��"},
	{"\xE0\x80\xAF", "���"},
	{"\xF0\x80\x80\xaf", "����"},
	{"\xF8\x80\x80\x80\xAF", "�����"},
	{"\xFC\x80\x80\x80\x80\xAF", "������"},
	{"a\x00b", "ab"},
	{"a\x00b\ufffd", "ab�"},
	{"a☺\x00b☺\xC0\xAFc☺\x00\xC0\xAF", "a☺b☺��c☺��"},
}

func TestToValidUTF8(t *testing.T) {
	for _, tc := range toValidUTF8Tests {
		got := ToValidUTF8(tc.in)
		if got != tc.out {
			t.Errorf("ToValidUTF8(%q) = %q; want %q", tc.in, got, tc.out)
		}
	}
}
