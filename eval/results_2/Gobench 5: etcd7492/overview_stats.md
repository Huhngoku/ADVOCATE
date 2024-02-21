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
| Number of atomics | 4 |
| Number of atomic operations | 39 |
| Number of channels | 3 |
| Number of channel operations | 7 |
| Number of selects | 31 |
| Number of select cases | 93 |
| Number of select channel operations | 93 |
| Number of select default operations | 0 |
| Number of mutexes | 2 |
| Number of mutex operations | 24 |
| Number of wait groups | 1 |
| Number of wait group operations | 5 |
| Number of cond vars | 0 |
| Number of cond var operations | 0 |
| Number of once | 0| 
| Number of once operations | 0 |


## Times
Invalid time file
0.004964,0.013228,## Results
==================== Summary ====================

-------------------- Critical -------------------
1 Potential mixed deadlock:
	locks: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/sync/rwmutex.go:266@137
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/sync/rwmutex.go:266@131
	send/close-recv: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/etcd7492/etcd7492.go:68@156
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/etcd7492/etcd7492.go:51@155
2 Potential mixed deadlock:
	locks: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/etcd7492/etcd7492.go:88@135
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/etcd7492/etcd7492.go:94@129
	send/close-recv: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/etcd7492/etcd7492.go:68@156
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/etcd7492/etcd7492.go:51@155
3 Potential leak without possible partner:
	channel: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/etcd7492/etcd7492.go:61@234
	partner: -
