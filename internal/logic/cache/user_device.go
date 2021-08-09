package cache

import (
	"im/internal/logic/model"
	"im/pkg/db"
	"im/pkg/gerrors"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	UserDeviceKey    = "user_device:"
	UserDeviceExpire = 2 * time.Hour
)

type userDeviceCache struct{}

var UserDeviceCache = new(userDeviceCache)

// Get 获取指定用户的所有在线设备
func (c *userDeviceCache) Get(userId int64) ([]model.Device, error) {
	var devices []model.Device
	err := RedisUtil.Get(UserDeviceKey+strconv.FormatInt(userId, 10), &devices)
	if err != nil && err != redis.Nil {
		return nil, gerrors.WrapError(err)
	}

	if err == redis.Nil {
		return nil, nil
	}
	return devices, nil
}

// Set 将指定用户的所有在线设备存入缓存
func (c *userDeviceCache) Set(userId int64, devices []model.Device) error {
	err := RedisUtil.Set(UserDeviceKey+strconv.FormatInt(userId, 10), devices, UserDeviceExpire)
	return gerrors.WrapError(err)
}

// Del 删除用户的在线设备列表
func (c *userDeviceCache) Del(userIds ...int64) error {
	var ids = make([]string, len(userIds))
	for i := range userIds {
		ids[i] = UserDeviceKey + strconv.FormatInt(userIds[i], 10)
	}

	_, err := db.RedisCli.Del(ids...).Result()
	return gerrors.WrapError(err)
}
