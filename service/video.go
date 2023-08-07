package service

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/spf13/viper"
	cos "github.com/tencentyun/cos-go-sdk-v5"
)

func Upload_video(name string, path string) error {
	req_url := viper.GetString("oss.tencent")
	u, _ := url.Parse(req_url)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     "SECRETID",
			SecretKey:    "SECRETKEY",
			SessionToken: "SECRETTOKEN",
		},
	})

	// 通过文件流上传对象
	fd, err := os.Open("./test")
	if err != nil {
		return err
	}
	defer fd.Close()
	_, err = c.Object.Put(context.Background(), name, fd, nil)
	if err != nil {
		return err
	}
	return nil
}

func Download_video(name string) (interface{}, error) {
	req_url := viper.GetString("oss.tencent")
	u, _ := url.Parse(req_url)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     "SECRETID",
			SecretKey:    "SECRETKEY",
			SessionToken: "SECRETTOKEN",
		},
	})

	resp, err := c.Object.Get(context.Background(), name, nil)
	if err != nil {
		return nil, err
	}
	bs, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return bs, nil
}
