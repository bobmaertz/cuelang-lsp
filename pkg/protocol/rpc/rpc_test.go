package rpc

import "testing"

type EncodingExample struct {
	Testing bool
}

func TestEncode(t *testing.T) {
	expected := "Content-Length: 16\r\n\r\n{\"Testing\":true}"

	actual := EncodeMessage(EncodingExample{Testing: true})

	if actual != expected {
		t.Errorf("Not equal: %v != %v", actual, expected)
		t.Fail()
	}
}

func TestDecode(t *testing.T) {
	msg := "Content-Length: 16\r\n\r\n{\"method\":\"hel\"}"

	method, content, err := DecodeMessage([]byte(msg))
	if err != nil {
		t.Error(err)
		return
	}

	if method != "hel" {
		t.Errorf("Not equal: hel != %v", method)
	}

	if string(content) != string("{\"method\":\"hel\"}") {
		t.Fatal("decode content not equal")
	}
}
