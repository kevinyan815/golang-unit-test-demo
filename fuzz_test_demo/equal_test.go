package fuzz_test_demo

import (
	"strconv"
	"testing"
)

func TestEqual(t *testing.T) {
	if !Equal([]byte{'f', 'u', 'z', 'z'}, []byte{'f', 'u', 'z', 'z'}) {
		t.Error("expected true, got false")
	}
}

func TestEqualWithTable(t *testing.T) {
	// 定义测试表格
	// 这里使用匿名结构体定义了若干个测试用例
	// 并且为每个测试用例设置了一个名称
	tests := []struct {
		name   string
		inputA []byte
		inputB []byte
		want   bool
	}{
		{"right case", []byte{'f', 'u', 'z', 'z'}, []byte{'f', 'u', 'z', 'z'}, true},
		{"right case", []byte{'a', 'b', 'c'}, []byte{'b', 'c', 'd'}, false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := Equal(tt.inputA, tt.inputB); got != tt.want {
				t.Error("expected " + strconv.FormatBool(tt.want) + ",  got " + strconv.FormatBool(got))
			}
		})
	}
}

func FuzzEqual(f *testing.F) {
	//f.Add([]byte{'a', 'b', 'c'}, []byte{'a', 'b', 'c'})
	f.Fuzz(func(t *testing.T, a []byte, b []byte) {
		Equal(a, b)
	})
}
