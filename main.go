package main

import (
	"fmt"
	"os"
	"time"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
)

const region = common.Region("cn-shanghai")

var ACCESS_KEY_ID = os.Getenv("ACCESS_KEY_ID")
var ACCESS_KEY_SECRET = os.Getenv("ACCESS_KEY_SECRET")
var client = ecs.NewClient(ACCESS_KEY_ID, ACCESS_KEY_SECRET)

func main() {
	if ACCESS_KEY_ID == "" || ACCESS_KEY_SECRET == "" {
		return
	}
	ch := make(chan int)

	r, _, _ := client.DescribeInstances(&ecs.DescribeInstancesArgs{RegionId: region})
	fmt.Println("Instance count:", len(r))

	for _, instance := range r {
		go reinstallInstance(ch, instance.InstanceId)
	}

	for i := 0; i < len(r); i++ {
		<-ch
	}

	fmt.Println("DONE!!!")
}

func reinstallInstance(ch chan int, instanceId string) {
	//spew.Dump(instanceId)

	client.StopInstance(instanceId, false)
	waitUntil(instanceId, "Stopped")

	diskId, _ := client.ReplaceSystemDisk(&ecs.ReplaceSystemDiskArgs{InstanceId: instanceId, ImageId: "ubuntu_16_0402_64_20G_alibase_20180326.vhd"})
	fmt.Printf("%s: diskId %s\n", instanceId, diskId)

	time.Sleep(5 * time.Second)

	client.StartInstance(instanceId)
	waitUntil(instanceId, "Running")

	ch <- 0
}

func waitUntil(instanceId string, state ecs.InstanceStatus) {
	r, _ := client.DescribeInstanceAttribute(instanceId)
	for r.Status != state {
		fmt.Printf("%s: sleep until %s\n", instanceId, state)
		time.Sleep(3 * time.Second)
		r, _ := client.DescribeInstanceAttribute(instanceId)
		if r.Status == state {
			break
		}
	}

	fmt.Printf("%s: %s\n", instanceId, state)
}
