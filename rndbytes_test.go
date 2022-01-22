package rndbytes_test

import (
	"github.com/SergeyMosin/rndbytes"
	"testing"
)

func TestGetBytes(t *testing.T) {
	mapSize := 10000000
	stringKeySize := 48
	strMap := make(map[string]int, mapSize)
	intMap := make(map[int]bool, mapSize)

	str := ""

	allowFirstDash := false

	for i := 0; i < mapSize; i++ {
		//goland:noinspection GoBoolExpressions
		str = string(rndbytes.GetBytes(stringKeySize, allowFirstDash))
		n := rndbytes.GetInt()

		//goland:noinspection GoBoolExpressions
		if !allowFirstDash && str[0] == '-' {
			t.Fatal("error: first dash is not allowed, str: "+str+", i:", i)
		}

		if _, ok := strMap[str]; ok {
			t.Fatal("error: duplicate str key")
		}
		strMap[str] = n

		if _, ok := intMap[n]; ok {
			t.Fatal("error: duplicate int key")
		}
		intMap[n] = true
	}
	t.Log("Done:", str, strMap[str], intMap[strMap[str]])
}

func TestGetBytesPrint(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Log(string(rndbytes.GetBytes(42, false)))
	}
	t.Log("------------------------------------------------------")
	for i := 0; i < 100; i++ {
		t.Log(string(rndbytes.GetBytes(42, true)))
	}
}

func BenchmarkGetBytesNoDash42(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_ = rndbytes.GetBytes(42, false)
	}
}

func BenchmarkGetBytes42(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_ = rndbytes.GetBytes(42, true)
	}
}

func BenchmarkGetBytes1500(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_ = rndbytes.GetBytes(1500, true)
	}
}

func BenchmarkGetBytesNoDash1500(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_ = rndbytes.GetBytes(1500, false)
	}
}
