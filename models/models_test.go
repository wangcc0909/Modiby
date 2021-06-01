package models

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestMarshal(t *testing.T) {
	u1 := User{"teststruct", 12}
	bytes, _ := u1.MarshalJSON()
	u2 := User{}
	_ = u2.UnmarshalJSON(bytes)
	fmt.Println(string(bytes), reflect.TypeOf(u2))
}

func BenchmarkAddress_MarshalEasyJSON(b *testing.B) {
	u1 := Address{"beijing", "beijing", "1049"}
	u2 := Address{}
	for i := 0; i < b.N; i++ {
		ujs, _ := u1.MarshalJSON()
		_ = u2.UnmarshalJSON(ujs)
	}
}

func BenchmarkAddress_MarshalJSON(b *testing.B) {
	u1 := Address{"beijing", "beijing", "1049"}
	u2 := Address{}
	for i := 0; i < b.N; i++ {
		ujs, _ := json.Marshal(u1)
		_ = json.Unmarshal(ujs, &u2)
	}
}
