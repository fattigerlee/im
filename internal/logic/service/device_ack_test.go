package service

import (
	"context"
	"fmt"
	_ "im/internal/logic/dao"
	"im/pkg/db"
	"testing"
)

func init() {
	db.InitByTest()
}

func Test_deviceAckService_GetMaxByUserId(t *testing.T) {
	fmt.Println(DeviceAckService.Update(context.TODO(), 1, 2, 2))
}

func Test_deviceAckService_Update(t *testing.T) {
	fmt.Println(DeviceAckService.GetMaxByUserId(context.TODO(), 1))
}
