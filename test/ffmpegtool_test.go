package test

import (
	"project/service"
	"testing"
)

func TestUploadVideo(t *testing.T) {
	service.UploadVideo("bear.mp4", "./bear.mp4")
	//log.Fatal(err)

}
