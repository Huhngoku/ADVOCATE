# overview Stats

## Program
| Info | Value |
| - | - |
| Number of go files | 73 |
| Number of lines | 19558 |
| Number of non-empty lines | 14671 |


## Trace
| Info | Value |
| - | - |
| Number of routines | 17 |
| Number of spawns | 0 |
| Number of atomics | 0 |
| Number of atomic operations | 0 |
| Number of channels | 0 |
| Number of channel operations | 0 |
| Number of selects | 2 |
| Number of select cases | 4 |
| Number of select channel operations | 4 |
| Number of select default operations | 0 |
| Number of mutexes | 0 |
| Number of mutex operations | 0 |
| Number of wait groups | 0 |
| Number of wait group operations | 0 |
| Number of cond vars | 0 |
| Number of cond var operations | 0 |
| Number of once | 0| 
| Number of once operations | 0 |


## Times
| Info | Value |
| - | - |
| Time for run without ADVOCATE | 1.020283 s |
| Time for run with ADVOCATE | 2.430540 s |
| Overhead of ADVOCATE | 138.222140 % |
| Analysis | 1.321094 s |


## Results
==================== Summary ====================

-------------------- Warning --------------------
1 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/Other/examples/bbolt/cmd/bbolt/main.go:1222@3124
	recv : /home/erikkassubek/Uni/HiWi/Other/examples/bbolt/cmd/bbolt/main.go:1428@148
2 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/Other/examples/bbolt/cmd/bbolt/main.go:1355@669538
	recv : /home/erikkassubek/Uni/HiWi/Other/examples/bbolt/cmd/bbolt/main.go:1428@3215
