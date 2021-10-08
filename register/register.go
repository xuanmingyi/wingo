package register

import (
	"fmt"
	"strings"
	"sync"
)

type Record interface{}

type Register struct {
	Driver  Driver
	Lock    sync.Mutex
	Records map[string]Record
}

func (r *Register) DefaultRecord(key string, record Record) {

}

func (r *Register) Open(path string) {
	var err error

	if strings.HasPrefix(path, "sqlite://") {
		r.Driver = new(SQliteDriver)
		r.Driver.Open(path)
	}

	// 缓存所有数据
	r.Records, err = r.Driver.Search("*")
	if err != nil {
		panic(err)
	}
}

func (r *Register) Close() {
	if r.Driver.IsValid() {
		r.Driver.Close()
	}
}

// 获取
func (r *Register) Get(key string) (record Record, err error) {
	r.Lock.Lock()
	defer r.Lock.Unlock()

	if record, found := r.Records[key]; found {
		return record, nil
	} else {
		return nil, fmt.Errorf("key not found: %s", key)
	}
}

// 添加
// 更新
func (r *Register) Set(key string, value Record) (err error) {
	r.Lock.Lock()
	defer r.Lock.Unlock()

	if _, found := r.Records[key]; found {
		// found
		err = r.Driver.Update(key, value)
		if err != nil {
			return err
		}
	} else {
		// not found
		err = r.Driver.Create(key, value)
		if err != nil {
			return err
		}
	}

	r.Records[key] = value
	return nil
}

// 删除
func (r *Register) Del(key string) {
	r.Lock.Lock()
	defer r.Lock.Unlock()

	if _, found := r.Records[key]; found {
		// found
		err := r.Driver.Delete(key)
		if err != nil {
			panic(err)
		}
	}
	delete(r.Records, key)
}
