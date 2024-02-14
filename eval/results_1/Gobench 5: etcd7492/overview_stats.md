# overview Stats

## Program
| Info | Value |
| - | - |
| Number of go files | 1 |
| Number of lines | 157 |
| Number of non-empty lines | 121 |


## Trace
| Info | Value |
| - | - |
| Number of routines | 20 |
| Number of spawns | 6 |
| Number of atomics | 3 |
| Number of atomic operations | 40 |
| Number of channels | 3 |
| Number of channel operations | 7 |
| Number of selects | 20 |
| Number of select cases | 60 |
| Number of select channel operations | 60 |
| Number of select default operations | 0 |
| Number of mutexes | 2 |
| Number of mutex operations | 28 |
| Number of wait groups | 1 |
| Number of wait group operations | 5 |
| Number of cond vars | 0 |
| Number of cond var operations | 0 |
| Number of once | 0| 
| Number of once operations | 0 |


## Times
| Info | Value |
| - | - |
| Time for run without ADVOCATE | 0.002573 s |
| Time for run with ADVOCATE | 0.014169 s |
| Overhead of ADVOCATE | 450.680140 % |
| Analysis | 0.038049 s |


## Results
==================== Summary ====================

-------------------- Critical -------------------
1 Potential mixed deadlock:
	locks: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/etcd7492/etcd7492.go:88@70
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/etcd7492/etcd7492.go:94@67
	send/close-recv: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/etcd7492/etcd7492.go:68@105
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/etcd7492/etcd7492.go:51@108
2 Potential mixed deadlock:
	locks: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/sync/rwmutex.go:270@72
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/sync/rwmutex.go:270@69
	send/close-recv: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/etcd7492/etcd7492.go:68@105
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/etcd7492/etcd7492.go:51@108
3 Potential mixed deadlock:
	locks: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/sync/rwmutex.go:270@90
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/sync/rwmutex.go:270@118
	send/close-recv: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/etcd7492/etcd7492.go:68@157
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/etcd7492/etcd7492.go:51@161
4 Potential mixed deadlock:
	locks: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/etcd7492/etcd7492.go:88@83
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/etcd7492/etcd7492.go:94@116
	send/close-recv: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/etcd7492/etcd7492.go:68@157
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/etcd7492/etcd7492.go:51@161
