package lock
/**
*	function:	锁信息
*	Author	:	dayunzhangyunfeng@didiglobal.com
*	Since	:	2019-11-25 19:50:00
*/
import (
	"sync/atomic"
)

type LockMessage struct {
	
	Key	string
	
	Value string
	
	ExpireTime int64
	
	LockStatus int32
	
	ChangeTimes int64
}

func NewLockMessage(key, value string) *LockMessage {
	lockMessage := &LockMessage{}
	lockMessage.Key = key
	lockMessage.Value = value
	return lockMessage
}

func (this *LockMessage) GetKey() string {
	return this.Key;
}

func (this *LockMessage) GetValue() string {
	return this.Value;
}

func (this *LockMessage) SetNewStatus(newStatus int32) bool {
	return atomic.CompareAndSwapInt32(&this.LockStatus, 1-newStatus, newStatus)
}

/**
 * 锁过期时间
 * 如果当前节点持有锁，则expireTime和redis中的锁过期时间保持一致
 * 如果当前节点不持有锁，则expireTime是当前节点最后一次获取到锁的过期时间
 */
func (this *LockMessage) GetExpireTime() int64 {
	return atomic.LoadInt64(&this.ExpireTime)
}

func (this *LockMessage) SetExpireTime(expireTime int64) {
	atomic.StoreInt64(&this.ExpireTime, expireTime)
}
	