# overview Stats

## Program
| Info | Value |
| - | - |
| Number of go files | 1 |
| Number of lines | 942 |
| Number of non-empty lines | 678 |


## Trace
| Info | Value |
| - | - |
| Number of routines | 19 |
| Number of spawns | 5 |
| Number of atomics | 1 |
| Number of atomic operations | 3 |
| Number of channels | 0 |
| Number of channel operations | 0 |
| Number of selects | 0 |
| Number of select cases | 0 |
| Number of select channel operations | 0 |
| Number of select default operations | 0 |
| Number of mutexes | 0 |
| Number of mutex operations | 0 |
| Number of wait groups | 1 |
| Number of wait group operations | 3 |
| Number of cond vars | 0 |
| Number of cond var operations | 0 |
| Number of once | 0| 
| Number of once operations | 0 |


## Times
| Info | Value |
| - | - |
| Time for run without ADVOCATE | 0.501595 s |
| Time for run with ADVOCATE | 0.504913 s |
| Overhead of ADVOCATE | 0.661490 % |
| Replay without changes | 0.504883 s |
| Overhead of Replay | 0.655509 % s |
| Analysis | 0.009847 s |


## Results
==================== Summary ====================

-------------------- Critical -------------------
1 Possible negative waitgroup counter:
	done: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:536@34
	done/add: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:523@31
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:527@32
