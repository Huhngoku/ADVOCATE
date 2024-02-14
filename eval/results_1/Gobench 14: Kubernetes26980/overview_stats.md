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
| Number of mutex operations | 10 |
| Number of wait groups | 0 |
| Number of wait group operations | 0 |
| Number of cond vars | 1 |
| Number of cond var operations | 1 |
| Number of once | 0| 
| Number of once operations | 0 |


## Times
| Info | Value |
| - | - |
| Time for run without ADVOCATE | 0.000992 s |
| Time for run with ADVOCATE | 0.018197 s |
| Overhead of ADVOCATE | 1734.375000 % |
| Analysis | 0.032710 s |


## Results
==================== Summary ====================

-------------------- Critical -------------------
1 Potential mixed deadlock:
	locks: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/sync/rwmutex.go:270@48
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/sync/rwmutex.go:270@30
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
3 Potential leak on mutex:
	mutex: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/kubernetes26980/kubernetes26980.go:23@58
	
4 Potential leak on mutex:
	mutex: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/sync/rwmutex.go:270@62
	
-------------------- Warning --------------------
5 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/kubernetes26980/kubernetes26980.go:58@53
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/kubernetes26980/kubernetes26980.go:60@46
