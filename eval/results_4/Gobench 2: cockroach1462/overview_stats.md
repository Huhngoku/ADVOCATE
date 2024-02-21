# overview Stats

## Program
| Info | Value |
| - | - |
| Number of go files | 1 |
| Number of lines | 139 |
| Number of non-empty lines | 123 |


## Trace
| Info | Value |
| - | - |
| Number of routines | 18 |
| Number of spawns | 4 |
| Number of atomics | 2 |
| Number of atomic operations | 13 |
| Number of channels | 2 |
| Number of channel operations | 2 |
| Number of selects | 2 |
| Number of select cases | 4 |
| Number of select channel operations | 4 |
| Number of select default operations | 0 |
| Number of mutexes | 1 |
| Number of mutex operations | 4 |
| Number of wait groups | 1 |
| Number of wait group operations | 5 |
| Number of cond vars | 0 |
| Number of cond var operations | 0 |
| Number of once | 0| 
| Number of once operations | 0 |


## Times
| Info | Value |
| - | - |
| Time for run without ADVOCATE | 0.001356 s |
| Time for run with ADVOCATE | 0.006469 s |
| Overhead of ADVOCATE | 377.064897 % |
| Replay without changes | 0.050840 s |
| Overhead of Replay | 3649.262537 % s |
| Analysis | 0.010501 s |


## Results
==================== Summary ====================

-------------------- Critical -------------------
1 Possible negative waitgroup counter:
	done: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/cockroach1462/cockroach1462.go:55@47
	done/add: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/cockroach1462/cockroach1462.go:55@44
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/cockroach1462/cockroach1462.go:35@28
-------------------- Warning --------------------
2 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/cockroach1462/cockroach1462.go:77@36
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/cockroach1462/cockroach1462.go:122@39
