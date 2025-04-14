package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/intel-sandbox/kube-score/pkg/common"
)

func TestDBQuery(t *testing.T) {
	ctx := context.Background()
	fmt.Printf("failed to create client")
	dbconfig := common.RunConfigDB{}
	dbconfig.Redis.Address = "localhost:6379"
	dbconfig.Redis.DB = 0
	dbconfig.Redis.Password = ""

	client := RedisClient{}
	err := client.NewClient(ctx, dbconfig)
	if err != nil {
		fmt.Printf("failed to create client : %v\n", err)
		t.Failed()
	}
	key := "kube-score"
	value := "hello-world"

	if err := client.AddKey(ctx, key, value); err != nil {
		t.Fail()
	}

	retVal, err := client.GetVaule(ctx, key)
	if err != nil {
		t.Fail()
	}

	if retVal != value {
		t.Fail()
	}
}
