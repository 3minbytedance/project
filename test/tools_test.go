package test

import (
	"douyin/common"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"testing"
)

func TestBloom(t *testing.T) {
	common.InitBloomFilter()
	common.AddToBloom("user1")
	common.AddToBloom("user2")
	common.AddToBloom("user3")
	common.AddToBloom("user4")
	common.AddToBloom("use1")
	common.AddToBloom("use2")

	assert.True(t, common.TestBloom("usera1"))
	assert.False(t, common.TestBloom("user5"))
}

func TestSensitiveFilter(t *testing.T) {
	err := common.InitSensitiveFilter()
	if err != nil {
		t.Log(err)
		return
	}
	word := common.ReplaceWord("傻逼吧卧槽啊你妈的")
	assert.True(t, word == "sad")

}
