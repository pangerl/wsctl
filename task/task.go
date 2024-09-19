// Package task @Author lanpang
// @Date 2024/9/19 上午9:33:00
// @Desc
package task

import "fmt"

var Creators = map[string]Creator{}

type Creator func() Tasker

type Tasker interface {
	Check()
	Gather()
}

type Initializer interface {
	Init() error
}

func Add(name string, creator Creator) {
	Creators[name] = creator
}

func Get(name string) Tasker {
	fmt.Println("task name:", name)
	return Creators[name]()
}

func MayInit(t interface{}) error {
	if initializer, ok := t.(Initializer); ok {
		return initializer.Init()
	}
	return nil
}

func Do(name string) {
	tasker := Get(name)
	err := MayInit(tasker)
	if err != nil {
		return
	}
	// 采集数据
	tasker.Gather()
	// 检查数据
	tasker.Check()
}
