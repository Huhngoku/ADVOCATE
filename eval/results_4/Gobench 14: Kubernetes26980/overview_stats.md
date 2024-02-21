# overview Stats

## Trace
| Info | Value |
| - | - |
| Number of routines | 18 |
| Number of spawns | 4 |
| Number of atomics | 3 |
| Number of atomic operations | 9 |
| Number of channels | 2 |
| Number of channel operations | 3 |
| Number of selects | 0 |
| Number of select cases | 0 |
| Number of select channel operations | 0 |
| Number of select default operations | 0 |
| Number of mutexes | 2 |
| Number of mutex operations | 8 |
| Number of wait groups | 0 |
| Number of wait group operations | 0 |
| Number of cond vars | 1 |
| Number of cond var operations | 1 |
| Number of once | 0| 
| Number of once operations | 0 |


## Times
| Info | Value |
| - | - |
| Time for run without ADVOCATE | 0.000909 s |
| Time for run with ADVOCATE | 0.004043 s |
| Overhead of ADVOCATE | 344.774477 % |
| Analysis | 0.009184 s |


## Results
==================== Summary ====================

-------------------- Critical -------------------
1 Potential mixed deadlock:
	locks: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/sync/rwmutex.go:266@48
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/sync/rwmutex.go:266@30
	send/close-recv: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/kubernetes26980/kubernetes26980.go:58@53
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/kubernetes26980/kubernetes26980.go:60@46
2 Potential mixed deadlock:
	locks: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/kubernetes26980/kubernetes26980.go:57@47
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/kubernetes26980/kubernetes26980.go:15@29
	send/close-recv: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/kubernetes26980/kubernetes26980.go:58@53
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/kubernetes26980/kubernetes26980.go:60@46
-------------------- Warning --------------------
3 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/kubernetes26980/kubernetes26980.go:58@53
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/kubernetes26980/kubernetes26980.go:60@46
