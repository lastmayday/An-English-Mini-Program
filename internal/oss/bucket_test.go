package oss

import (
	"testing"
)

func TestUpload(t *testing.T) {
	filePath := "/Users/lastmayday/test.txt"
	fileName := "test.txt"
	_, err := Upload(fileName, filePath)
	if err != nil {
		t.Error("upload failed")
	}
}
