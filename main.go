package main

import (
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
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
	spew.Dump(r)
	spew.Dump(len(r))

	for _, instance := range r {
		go reinstallInstance(ch, instance.InstanceId)
	}

	for i := 0; i < len(r); i++ {
		spew.Dump(<-ch)
	}
}

func reinstallInstance(ch chan int, instanceId string) {
	client.StopInstance(instanceId, false)
	spew.Dump(instanceId)
	r, _ := client.DescribeInstanceAttribute(instanceId)
	for r.Status != "Stopped" {
		spew.Dump("sleep")
		time.Sleep(3 * time.Second)
		r, _ := client.DescribeInstanceAttribute(instanceId)
		if r.Status == "Stopped" {
			break
		}
	}
	spew.Dump("Stopped")

	diskId, _ := client.ReplaceSystemDisk(&ecs.ReplaceSystemDiskArgs{InstanceId: instanceId, ImageId: "ubuntu_16_0402_64_40G_alibase_20170711.vhd"})
	spew.Dump(diskId)

	time.Sleep(5 * time.Second)

	client.StartInstance(instanceId)

	ch <- 0
}
