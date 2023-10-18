package lock
/**
*	function:	租约锁
*	Author	:	dayunzhangyunfeng@didiglobal.com
*	Since	:	2019-11-25 19:50:00
*/
import (
	"fmt"
	"time"
	"github.com/didi/tg-flow/tg-core/common/redis"
	"github.com/didi/tg-flow/tg-core/common/utils"
	"github.com/didi/tg-flow/tg-core/consts"
	"context"
)


const (
	DEFAULT_LEASE_EXPIRE_TIME = 20
	DEFAULT_LEASE_PERIOD = 2
)

type RedisLeaseLock struct {
	LeaseLock
	
	LockMap map[string]*LockMessage
	
	LeaseExpireTime int64	// 过期时长（单位秒）
	
	LeasePeriod int64		// 间隔（单位秒）
}

func NewRedisLeaseLock(leasePeriod int64, leaseExpireTime int64) (*RedisLeaseLock, error) {
	rll := &RedisLeaseLock{}
	// 租约过期时长 至少 是续租间隔的4倍，否则抛出异常
	if (leaseExpireTime <= leasePeriod * 4) {
		return nil, fmt.Errorf("leaseExpireTime is too small, it must great than %v" , leasePeriod * 4)
	}
	
	rll.LeasePeriod = leasePeriod
	rll.LeaseExpireTime = leaseExpireTime
	rll.LockMap = make(map[string]*LockMessage)
	return rll, nil
}

func (this *RedisLeaseLock) RegisterLock(key, value string) error {
	if _, ok := this.LockMap[key]; ok {
		return fmt.Errorf("register lease lock error, key is exist! key=", key);
	}
	
	this.LockMap[key]= &LockMessage{
		Key			:	key,
		Value		:	value,
		ExpireTime	:	0,
		LockStatus 	:	0,
		ChangeTimes :	0}
	fmt.Println("register lease lock, key=" + key + ", value=" + value)

	// 注册时需要阻塞调用一次
	this.setLeaseLocks()
	
	go this.initTryLock();
	
	return nil
}

func (this *RedisLeaseLock) initTryLock() {
	for {
		time.Sleep(time.Duration(1) * time.Second)
		this.setLeaseLocks()
	}
}

func (this *RedisLeaseLock) setLeaseLocks() {
	for _, lockMessage :=range this.LockMap {
		this.setLeaseLock(lockMessage)
	}
}

func (this *RedisLeaseLock) setLeaseLock(lockMessage *LockMessage) {
	defer utils.Recover(context.TODO(), nil, consts.DLTagCronTask, "RedisLeaseLock.setLeaseLock")
	// 锁已过期，尝试nx操作获取锁
	if time.Now().Unix() >= lockMessage.GetExpireTime() {
		isOk, err := redis.Handler.SetNEx(context.TODO(), lockMessage.GetKey(), int(this.LeaseExpireTime), lockMessage.GetValue())
		if isOk == "OK" && err == nil {					
			lockMessage.SetExpireTime(time.Now().Unix() + this.LeaseExpireTime);
		} else if isOk == "" {
			// 如果NX不能获取锁，需要进一步确认此锁是否还属于自己
			// 如果属于自己，则取出当前锁的TTL赋给内存对象
			serverLockMessage, err := redis.Handler.GetString(context.TODO(), lockMessage.GetKey())
			if serverLockMessage != "" && err ==nil && serverLockMessage == lockMessage.GetValue() {							
				ttl, err := redis.Handler.TTL(context.TODO(), lockMessage.GetKey())
				if err == nil {
					ttl64 := int64(ttl)
					if ttl64 > this.LeasePeriod && err == nil {
						lockMessage.SetExpireTime(time.Now().Unix() + ttl64)
					}
				}
			}
		}
	} else {
		// 锁还未过期，尝试expire
		serverLockMessage, err := redis.Handler.GetString(context.TODO(), lockMessage.GetKey())
		if serverLockMessage == "" || serverLockMessage != lockMessage.GetValue() {
			return;
		}
		
		isOk, err := redis.Handler.Expire(context.TODO(), lockMessage.GetKey(), this.LeaseExpireTime)
		if isOk == 1 && err == nil {
			lockMessage.SetExpireTime(time.Now().Unix() + this.LeaseExpireTime);
		}
	}
}

/**
	获取任务对应的锁
	如果key不存在，抛异常
	如果获取成功，则返回true
	如果获取失败，则返回false
**/
func (this *RedisLeaseLock) TryLock(key string) bool {
	lockMessage, ok := this.LockMap[key]
	if !ok {
		panic("not registered key:"+ key)
	}
	
	if time.Now().Unix() < lockMessage.GetExpireTime() {
		lockMessage.SetNewStatus(int32(1))
		return true
	} else{
		lockMessage.SetNewStatus(int32(0))
		return false
	}
}
