/*
Copyright 2024 - 2025 Zen HuiFer

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package initialize

import (
	"context"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	LockTime       = 4 * time.Second
	RsDistlockNs   = "tdln:"
	ReleaseLockLua = `
        if redis.call('get',KEYS[1])==ARGV[1] then
            return redis.call('del', KEYS[1])
        else
            return 0
        end
    `
)

type RedisDistLock struct {
	id          string
	lockName    string
	redisClient *redis.Client
	m           sync.Mutex
}

func NewRedisDistLock(redisClient *redis.Client, lockName string) *RedisDistLock {
	return &RedisDistLock{
		lockName:    lockName,
		redisClient: redisClient,
	}
}

func (lock *RedisDistLock) Lock() {
	for !lock.TryLock() {
		time.Sleep(5 * time.Second)
	}
}

func (lock *RedisDistLock) TryLock() bool {
	if lock.id != "" {
		// 处于加锁中
		return false
	}
	lock.m.Lock()
	defer lock.m.Unlock()
	if lock.id != "" {
		// 处于加锁中
		return false
	}
	ctx := context.Background()
	id := uuid.New().String()
	reply := lock.redisClient.SetNX(ctx, RsDistlockNs+lock.lockName, id, LockTime)
	if reply.Err() == nil && reply.Val() {
		lock.id = id
		return true
	}

	return false
}

func (lock *RedisDistLock) Unlock() {
	if lock.id == "" {
		// 未加锁
		panic("解锁失败，因为未加锁")
	}
	lock.m.Lock()
	defer lock.m.Unlock()
	if lock.id == "" {
		// 未加锁
		panic("解锁失败，因为未加锁")
	}
	ctx := context.Background()
	reply := lock.redisClient.Eval(ctx, ReleaseLockLua, []string{RsDistlockNs + lock.lockName}, lock.id)
	if reply.Err() != nil {
		panic("释放锁失败！")
	} else {
		lock.id = ""
	}
}
