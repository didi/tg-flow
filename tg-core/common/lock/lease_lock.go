package lock

type  LeaseLock interface {

	RegisterLock(key string, value string)
	
	TryLock(key string)
	
}