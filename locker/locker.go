package locker

type Locker interface {
	Lock() (bool, error)
	IsLocked() (bool, error)
	Release() (bool, error)
}

type Factory interface {
	NewLocker(key string, waitMs int, timeoutMs int) Locker
}
