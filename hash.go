package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

const (
	lotteryPrefix   = "lottery:"
	LotteryHashSize = 10000
	// 用户剩余抽奖次数 hash
	LotteryRemainUserId = lotteryPrefix + "remain:userid_%d"
)

// 用户剩余抽奖次数
func getUserRemainKey(userId int) string {
	return fmt.Sprintf(LotteryRemainUserId, int(math.Floor(float64(userId)/float64(LotteryHashSize))))
}

func GetUserRemain(userId int) (int, error) {
	hashKey := getUserRemainKey(userId)
	userKey := strconv.Itoa(userId)
	str, err := rdsClient.HGet(hashKey, userKey).Result()
	if err == redis.Nil {
		fmt.Printf("user_%d cache user remain not exist", userId)
		return 0, nil
	} else if err != nil {
		fmt.Printf("user_%d get cache user remain error=%v", err.Error())
		return 0, err
	}
	remain, err := strconv.Atoi(str)
	if err != nil {
		fmt.Printf("GetUserRemain strconv.Atoi error=%v", err.Error())
		return 0, err
	}
	fmt.Printf("GetUserRemain get user_%d remain=%d", userId, remain)
	return remain, nil
}

func ReloadUserRemain() {
	remainLessZeroUserIds := make([]int, 0)
	reloadLessZeroUserIds := make([]int, 0)

	kvm, err := rdsClient.HGetAll("lottery:remain:userid_0").Result()
	if err != nil {
		return
	}
	logrus.Debugf("kvm=%+v", kvm)

	for userIdStr, remainStr := range kvm {
		logrus.Debugf("userIdStr=%s remainStr=%s", userIdStr, remainStr)

		userId, _ := strconv.Atoi(userIdStr)
		remain, _ := strconv.Atoi(remainStr)

		logrus.Debugf("userId=%d remain=%d", userId, remain)
		if remain < 0 {
			remainLessZeroUserIds = append(remainLessZeroUserIds, userId)
			hashKey := getUserRemainKey(userId)
			userKey := strconv.Itoa(userId)
			err := rdsClient.HSet(hashKey, userKey, 0).Err()
			if err != nil {
				logrus.Errorf("rdsClient.HSet error=%v", err.Error())
			} else {
				reloadLessZeroUserIds = append(reloadLessZeroUserIds, userId)
			}
		}
	}

	logrus.Debugf("count=%d remainLessZeroUserIds=%+v", len(remainLessZeroUserIds), remainLessZeroUserIds)
	logrus.Debugf("count=%d reloadLessZeroUserIds=%+v", len(reloadLessZeroUserIds), reloadLessZeroUserIds)
}
