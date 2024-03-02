// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"internal/race"
	// ADVOCATE-CHANGE-START
	"runtime"
	// ADVOCATE-CHANGE-END
	"sync/atomic"
	"unsafe"
)

// There is a modified copy of this file in runtime/rwmutex.go.
// If you make any changes here, see if you should make them there.

// A RWMutex is a reader/writer mutual exclusion lock.
// The lock can be held by an arbitrary number of readers or a single writer.
// The zero value for a RWMutex is an unlocked mutex.
//
// A RWMutex must not be copied after first use.
//
// If a goroutine holds a RWMutex for reading and another goroutine might
// call Lock, no goroutine should expect to be able to acquire a read lock
// until the initial read lock is released. In particular, this prohibits
// recursive read locking. This is to ensure that the lock eventually becomes
// available; a blocked Lock call excludes new readers from acquiring the
// lock.
//
// In the terminology of the Go memory model,
// the n'th call to Unlock “synchronizes before” the m'th call to Lock
// for any n < m, just as for Mutex.
// For any call to RLock, there exists an n such that
// the n'th call to Unlock “synchronizes before” that call to RLock,
// and the corresponding call to RUnlock “synchronizes before”
// the n+1'th call to Lock.
type RWMutex struct {
	w           Mutex        // held if there are pending writers
	writerSem   uint32       // semaphore for writers to wait for completing readers
	readerSem   uint32       // semaphore for readers to wait for completing writers
	readerCount atomic.Int32 // number of pending readers
	readerWait  atomic.Int32 // number of departing readers
	// ADVOCATE-CHANGE-START
	id uint64 // id for the mutex
	// ADVOCATE-CHANGE-END
}

const rwmutexMaxReaders = 1 << 30

// Happens-before relationships are indicated to the race detector via:
// - Unlock  -> Lock:  readerSem
// - Unlock  -> RLock: readerSem
// - RUnlock -> Lock:  writerSem
//GetAdvocateObjectID
// The methods below temporarily disable handling of race synchronization
// events in order to provide the more precise model above to the race
// detector.
//
// For example, atomic.AddInt32 in RLock should not appear to provide
// acquire-release semantics, which would incorrectly synchronize racing
// readers, thus potentially missing races.

// RLock locks rw for reading.
//
// It should not be used for recursive read locking; a blocked Lock
// call excludes new readers from acquiring the lock. See the
// documentation on the RWMutex type.
func (rw *RWMutex) RLock() {
	// ADVOCATE-CHANGE-START
	enabled, valid, replayElem := runtime.WaitForReplay(runtime.OperationRWMutexRLock, 2)
	if enabled && valid {
		if replayElem.Blocked {
			if rw.id == 0 {
				rw.id = runtime.GetAdvocateObjectID()
			}
			_ = runtime.AdvocateMutexLockPre(rw.id, true, true)
			runtime.BlockForever()
		}
	}

	// RWMutexe don't need to be initialized in default go code. Because
	// go does not have constructors, the only way to initialize a RWMutex
	// is directly in the lock function. If the id of the channel is the default
	// value, it is set to a new, unique object id
	if rw.id == 0 {
		rw.id = runtime.GetAdvocateObjectID()
	}

	// AdvocateMutexLockPre records, that a routine tries to lock a mutex.
	// AdvocatePost is called, if the mutex was locked successfully.
	// In this case, the Lock event in the trace is updated to include
	// this information. advocateIndex is used for AdvocatePost to find the
	// pre event.
	advocateIndex := runtime.AdvocateMutexLockPre(rw.id, true, true)
	defer runtime.AdvocatePost(advocateIndex)
	// ADVOCATE-CHANGE-END

	if race.Enabled {
		_ = rw.w.state
		race.Disable()
	}
	if rw.readerCount.Add(1) < 0 {
		// A writer is pending, wait for it.
		runtime_SemacquireRWMutexR(&rw.readerSem, false, 0)
	}
	if race.Enabled {
		race.Enable()
		race.Acquire(unsafe.Pointer(&rw.readerSem))
	}
}

// TryRLock tries to lock rw for reading and reports whether it succeeded.
//
// Note that while correct uses of TryRLock do exist, they are rare,
// and use of TryRLock is often a sign of a deeper problem
// in a particular use of mutexes.
func (rw *RWMutex) TryRLock() bool {
	// ADVOCATE-CHANGE-START
	enabled, valid, replayElem := runtime.WaitForReplay(runtime.OperationRWMutexTryRLock, 2)
	if enabled && valid {
		if replayElem.Blocked {
			if rw.id == 0 {
				rw.id = runtime.GetAdvocateObjectID()
			}
			_ = runtime.AdvocateMutexLockTry(rw.id, true, true)
			runtime.BlockForever()
		}
		// if !replayElem.Suc {
		// 	if rw.id == 0 {
		// 		rw.id = runtime.GetAdvocateObjectID()
		// 	}
		// 	advocateIndex := runtime.AdvocateMutexLockTry(rw.id, true, true)
		// 	runtime.AdvocatePostTry(advocateIndex, false)
		// 	return false
		// }
	}
	// RWMutexe don't need to be initialized in default go code. Because
	// go does not have constructors, the only way to initialize a RWMutex
	// is directly in the lock function. If the id of the channel is the default
	// value, it is set to a new, unique object id
	if rw.id == 0 {
		rw.id = runtime.GetAdvocateObjectID()
	}
	// AdvocateMutexLockPre records, that a routine tries to lock a mutex.
	// advocateIndex is used for AdvocatePostTry to find the pre event.
	advocateIndex := runtime.AdvocateMutexLockTry(rw.id, true, true)
	// ADVOCATE-CHANGE-END

	if race.Enabled {
		_ = rw.w.state
		race.Disable()
	}
	for {
		c := rw.readerCount.Load()
		if c < 0 {
			if race.Enabled {
				race.Enable()
			}
			// ADVOCATE-CHANGE-START
			// If the mutex was not locked successfully, AdvocatePostTry is called
			// to update the trace.
			runtime.AdvocatePostTry(advocateIndex, false)
			// ADVOCATE-CHANGE-END
			return false
		}
		if rw.readerCount.CompareAndSwap(c, c+1) {
			if race.Enabled {
				race.Enable()
				race.Acquire(unsafe.Pointer(&rw.readerSem))
			}
			// ADVOCATE-CHANGE-START
			// If the mutex was locked successfully, AdvocatePostTry is called
			// to update the trace.
			runtime.AdvocatePostTry(advocateIndex, true)
			// ADVOCATE-CHANGE-END
			return true
		}
	}
}

// RUnlock undoes a single RLock call;
// it does not affect other simultaneous readers.
// It is a run-time error if rw is not locked for reading
// on entry to RUnlock.
func (rw *RWMutex) RUnlock() {
	// ADVOCATE-CHANGE-START
	enabled, valid, replayElem := runtime.WaitForReplay(runtime.OperationRWMutexRUnlock, 2)
	if enabled && valid {
		if replayElem.Blocked {
			_ = runtime.AdvocateUnlockPre(rw.id, true, true)
			runtime.BlockForever()
		}
	}

	// AdvocateUnlockPre is used to record the unlocking of a mutex.
	// AdvocatePost records the successful unlocking of a mutex.
	advocateIndex := runtime.AdvocateUnlockPre(rw.id, true, true)
	defer runtime.AdvocatePost(advocateIndex)
	// ADVOCATE-CHANGE-END

	if race.Enabled {
		_ = rw.w.state
		race.ReleaseMerge(unsafe.Pointer(&rw.writerSem))
		race.Disable()
	}
	if r := rw.readerCount.Add(-1); r < 0 {
		// Outlined slow-path to allow the fast-path to be inlined
		rw.rUnlockSlow(r)
	}
	if race.Enabled {
		race.Enable()
	}
}

func (rw *RWMutex) rUnlockSlow(r int32) {
	if r+1 == 0 || r+1 == -rwmutexMaxReaders {
		race.Enable()
		fatal("sync: RUnlock of unlocked RWMutex")
	}
	// A writer is pending.
	if rw.readerWait.Add(-1) == 0 {
		// The last reader unblocks the writer.
		runtime_Semrelease(&rw.writerSem, false, 1)
	}
}

// Lock locks rw for writing.
// If the lock is already locked for reading or writing,
// Lock blocks until the lock is available.
func (rw *RWMutex) Lock() {
	// ADVOCATE-CHANGE-START
	enabled, valid, replayElem := runtime.WaitForReplay(runtime.OperationRWMutexLock, 2)
	if enabled && valid {
		if replayElem.Blocked {
			if rw.id == 0 {
				rw.id = runtime.GetAdvocateObjectID()
			}
			_ = runtime.AdvocateMutexLockPre(rw.id, true, false)
			runtime.BlockForever()
		}
	}
	// RWMutexe don't need to be initialized in default go code. Because
	// go does not have constructors, the only way to initialize a RWMutex
	// is directly in the lock function. If the id of the channel is the default
	// value, it is set to a new, unique object id
	if rw.id == 0 {
		rw.id = runtime.GetAdvocateObjectID()
	}

	// AdvocateMutexLockPre records, that a routine tries to lock a mutex.
	// AdvocatePost is called, if the mutex was locked successfully.
	// In this case, the Lock event in the trace is updated to include
	// this information. advocateIndex is used for AdvocatePost to find the
	// pre event.
	advocateIndex := runtime.AdvocateMutexLockPre(rw.id, true, false)
	defer runtime.AdvocatePost(advocateIndex)
	// ADVOCATE-CHANGE-END

	if race.Enabled {
		_ = rw.w.state
		race.Disable()
	}
	// First, resolve competition with other writers.
	rw.w.Lock()
	// Announce to readers there is a pending writer.
	r := rw.readerCount.Add(-rwmutexMaxReaders) + rwmutexMaxReaders
	// Wait for active readers.
	if r != 0 && rw.readerWait.Add(r) != 0 {
		runtime_SemacquireRWMutex(&rw.writerSem, false, 0)
	}
	if race.Enabled {
		race.Enable()
		race.Acquire(unsafe.Pointer(&rw.readerSem))
		race.Acquire(unsafe.Pointer(&rw.writerSem))
	}
}

// TryLock tries to lock rw for writing and reports whether it succeeded.
//
// Note that while correct uses of TryLock do exist, they are rare,
// and use of TryLock is often a sign of a deeper problem
// in a particular use of mutexes.
func (rw *RWMutex) TryLock() bool {
	// ADVOCATE-CHANGE-START
	enabled, valid, replayElem := runtime.WaitForReplay(runtime.OperationRWMutexTryLock, 2)
	if enabled && valid {
		if replayElem.Blocked {
			if rw.id == 0 {
				rw.id = runtime.GetAdvocateObjectID()
			}
			// AdvocateMutexLockPre records, that a routine tries to lock a mutex.
			// advocateIndex is used for AdvocatePostTry to find the pre event.
			_ = runtime.AdvocateMutexLockTry(rw.id, true, false)
			runtime.BlockForever()
		}
		// if !replayElem.Suc {
		// 	advocateIndex := runtime.AdvocateMutexLockTry(rw.id, true, false)
		// 	runtime.AdvocatePostTry(advocateIndex, false)
		// 	return false
		// }
	}
	// RWMutexe don't need to be initialized in default go code. Because
	// go does not have constructors, the only way to initialize a RWMutex
	// is directly in the lock function. If the id of the channel is the default
	// value, it is set to a new, unique object id
	if rw.id == 0 {
		rw.id = runtime.GetAdvocateObjectID()
	}
	// AdvocateMutexLockPre records, that a routine tries to lock a mutex.
	// advocateIndex is used for AdvocatePostTry to find the pre event.
	advocateIndex := runtime.AdvocateMutexLockTry(rw.id, true, false)
	// ADVOCATE-CHANGE-END
	if race.Enabled {
		_ = rw.w.state
		race.Disable()
	}
	if !rw.w.TryLock() {
		if race.Enabled {
			race.Enable()
		}
		// ADVOCATE-CHANGE-START
		// If the mutex was not locked successfully, AdvocatePostTry is called
		// to update the trace.
		runtime.AdvocatePostTry(advocateIndex, false)
		// ADVOCATE-CHANGE-END
		return false
	}
	if !rw.readerCount.CompareAndSwap(0, -rwmutexMaxReaders) {
		rw.w.Unlock()
		if race.Enabled {
			race.Enable()
		}
		// ADVOCATE-CHANGE-START
		// If the mutex was not locked successfully, AdvocatePostTry is called
		// to update the trace.
		runtime.AdvocatePostTry(advocateIndex, false)
		// ADVOCATE-CHANGE-END
		return false
	}
	if race.Enabled {
		race.Enable()
		race.Acquire(unsafe.Pointer(&rw.readerSem))
		race.Acquire(unsafe.Pointer(&rw.writerSem))
	}
	// ADVOCATE-CHANGE-START
	// If the mutex was locked successfully, AdvocatePostTry is called
	// to update the trace.
	runtime.AdvocatePostTry(advocateIndex, true)
	// ADVOCATE-CHANGE-END

	return true
}

// Unlock unlocks rw for writing. It is a run-time error if rw is
// not locked for writing on entry to Unlock.
//
// As with Mutexes, a locked RWMutex is not associated with a particular
// goroutine. One goroutine may RLock (Lock) a RWMutex and then
// arrange for another goroutine to RUnlock (Unlock) it.
func (rw *RWMutex) Unlock() {
	// ADVOCATE-CHANGE-START
	_, _, _ = runtime.WaitForReplay(runtime.OperationRWMutexUnlock, 2)
	// AdvocateUnlockPre is used to record the unlocking of a mutex.
	// AdvocatePost records the successful unlocking of a mutex.
	// For non rw mutexe, the unlock cannot fail. Therefore it is not
	// strictly necessary to record the post for the unlocking of a mutex.
	advocateIndex := runtime.AdvocateUnlockPre(rw.id, true, false)
	defer runtime.AdvocatePost(advocateIndex)
	// ADVOCATE-CHANGE-END

	if race.Enabled {
		_ = rw.w.state
		race.Release(unsafe.Pointer(&rw.readerSem))
		race.Disable()
	}

	// Announce to readers there is no active writer.
	r := rw.readerCount.Add(rwmutexMaxReaders)
	if r >= rwmutexMaxReaders {
		race.Enable()
		fatal("sync: Unlock of unlocked RWMutex")
	}
	// Unblock blocked readers, if any.
	for i := 0; i < int(r); i++ {
		runtime_Semrelease(&rw.readerSem, false, 0)
	}
	// Allow other writers to proceed.
	rw.w.Unlock()
	if race.Enabled {
		race.Enable()
	}

}

// RLocker returns a Locker interface that implements
// the Lock and Unlock methods by calling rw.RLock and rw.RUnlock.
func (rw *RWMutex) RLocker() Locker {
	return (*rlocker)(rw)
}

type rlocker RWMutex

func (r *rlocker) Lock()   { (*RWMutex)(r).RLock() }
func (r *rlocker) Unlock() { (*RWMutex)(r).RUnlock() }
