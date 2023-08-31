package common

import (
	"douyin/dal/model"
	"douyin/dal/mysql"
	"github.com/bits-and-blooms/bloom/v3"
	"go.uber.org/zap"
	"log"
)

var bloomUserFilter *bloom.BloomFilter
var bloomCommentFilter *bloom.BloomFilter
var bloomWorkCountFilter *bloom.BloomFilter
var bloomFavoriteUserIdFilter *bloom.BloomFilter
var bloomFavoriteVideoIdFilter *bloom.BloomFilter
var bloomRelationFollowIdFilter *bloom.BloomFilter
var bloomRelationFollowerIdFilter *bloom.BloomFilter

func InitUserBloomFilter() {
	// 初始化用户名布隆过滤器
	bloomUserFilter = bloom.NewWithEstimates(100000, 0.01) // 假设预期元素数量为 100000，误判率为 0.01
}

func InitCommentBloomFilter() {
	// 初始化评论布隆过滤器
	bloomCommentFilter = bloom.NewWithEstimates(100000, 0.01) // 假设预期元素数量为 100000，误判率为 0.01
}

func InitWorkCountFilter() {
	// 初始化作品数布隆过滤器
	bloomWorkCountFilter = bloom.NewWithEstimates(100000, 0.01) // 假设预期元素数量为 100000，误判率为 0.01
}

func InitFavoriteUserIdFilter() {
	// 初始化用户点赞数布隆过滤器
	bloomFavoriteUserIdFilter = bloom.NewWithEstimates(100000, 0.01) // 假设预期元素数量为 100000，误判率为 0.01
}

func InitFavoriteVideoIdFilter() {
	// 初始化视频点赞数布隆过滤器
	bloomFavoriteVideoIdFilter = bloom.NewWithEstimates(100000, 0.01) // 假设预期元素数量为 100000，误判率为 0.01
}

func InitRelationFollowIdFilter() {
	// 初始化关注布隆过滤器
	bloomRelationFollowIdFilter = bloom.NewWithEstimates(100000, 0.01) // 假设预期元素数量为 100000，误判率为 0.01
}

func InitRelationFollowerIdFilter() {
	// 初始化粉丝布隆过滤器
	bloomRelationFollowerIdFilter = bloom.NewWithEstimates(100000, 0.01) // 假设预期元素数量为 100000，误判率为 0.01
}

func AddToUserBloom(data string) {
	bloomUserFilter.Add([]byte(data))
}

func TestUserBloom(data string) bool {
	return bloomUserFilter.Test([]byte(data))
}

func LoadUsernamesToBloomFilter() {
	var usernames []string
	err := mysql.DB.Model(&model.User{}).Pluck("name", &usernames).Error
	if err != nil {
		log.Fatal("Failed to retrieve usernames from database:", err)
	}

	for _, username := range usernames {
		AddToUserBloom(username)
	}

	zap.L().Info("Loaded %d usernames to the bloom filter.\n", zap.Int("size", len(usernames)))
}

func AddToCommentBloom(data string) {
	bloomCommentFilter.Add([]byte(data))
}

func TestCommentBloom(data string) bool {
	return bloomCommentFilter.Test([]byte(data))
}

func LoadCommentVideoIdToBloomFilter() {
	var videoIdList []string
	mysql.DB.Model(&model.Comment{}).Distinct().Pluck("video_id", &videoIdList)
	for _, videoId := range videoIdList {
		AddToCommentBloom(videoId)
	}
	zap.L().Info("Loaded %d comments to the bloom filter.\n", zap.Int("size", len(videoIdList)))
}

func AddToWorkCountBloom(data string) {
	bloomWorkCountFilter.Add([]byte(data))
}

func TestWorkCountBloom(data string) bool {
	return bloomWorkCountFilter.Test([]byte(data))
}

func LoadWorkCountToBloomFilter() {
	var authorIdList []string
	mysql.DB.Model(&model.Video{}).Distinct().Pluck("author_id", &authorIdList)
	for _, authorId := range authorIdList {
		AddToWorkCountBloom(authorId)
	}
	zap.L().Info("Loaded %d authors from video to the bloom filter.\n", zap.Int("size", len(authorIdList)))
}

func AddToFavoriteUserIdBloom(data string) {
	bloomFavoriteUserIdFilter.Add([]byte(data))
}

func TestFavoriteUserIdBloom(data string) bool {
	return bloomFavoriteUserIdFilter.Test([]byte(data))
}

func LoadFavoriteUserIdToBloomFilter() {
	var userIdList []string
	mysql.DB.Model(&model.Favorite{}).Distinct().Pluck("user_id", &userIdList)
	for _, userId := range userIdList {
		AddToFavoriteUserIdBloom(userId)
	}
	zap.L().Info("Loaded %d user from favorite to the bloom filter.\n", zap.Int("size", len(userIdList)))
}

func AddToFavoriteVideoIdBloom(data string) {
	bloomFavoriteVideoIdFilter.Add([]byte(data))
}

func TestFavoriteVideoIdBloom(data string) bool {
	return bloomFavoriteVideoIdFilter.Test([]byte(data))
}

func LoadFavoriteVideoIdToBloomFilter() {
	var videoIdList []string
	mysql.DB.Model(&model.Favorite{}).Distinct().Pluck("video_id", &videoIdList)
	for _, videoId := range videoIdList {
		AddToFavoriteVideoIdBloom(videoId)
	}
	zap.L().Info("Loaded %d video from favorite to the bloom filter.\n", zap.Int("size", len(videoIdList)))
}

func AddToRelationFollowIdBloom(data string) {
	bloomRelationFollowIdFilter.Add([]byte(data))
}

func TestRelationFollowIdBloom(data string) bool {
	return bloomRelationFollowIdFilter.Test([]byte(data))
}

func LoadRelationFollowIdToBloomFilter() {
	var followIdList []string
	mysql.DB.Model(&model.UserFollow{}).Distinct().Pluck("user_id", &followIdList)
	for _, followId := range followIdList {
		AddToRelationFollowIdBloom(followId)
	}
	zap.L().Info("Loaded %d followId from follow to the bloom filter.\n", zap.Int("size", len(followIdList)))
}

func AddToRelationFollowerIdBloom(data string) {
	bloomRelationFollowerIdFilter.Add([]byte(data))
}

func TestRelationFollowerIdBloom(data string) bool {
	return bloomRelationFollowerIdFilter.Test([]byte(data))
}

func LoadRelationFollowerIdToBloomFilter() {
	var followerIdList []string
	mysql.DB.Model(&model.UserFollow{}).Distinct().Pluck("follow_id", &followerIdList)
	for _, followerId := range followerIdList {
		AddToRelationFollowerIdBloom(followerId)
	}
	zap.L().Info("Loaded %d follower from follow to the bloom filter.\n", zap.Int("size", len(followerIdList)))
}
