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
| Number of spawns | 3 |
| Number of atomics | 118 |
| Number of atomic operations | 721467 |
| Number of channels | 2 |
| Number of channel operations | 2 |
| Number of selects | 2 |
| Number of select cases | 4 |
| Number of select channel operations | 4 |
| Number of select default operations | 0 |
| Number of mutexes | 11 |
| Number of mutex operations | 98 |
| Number of wait groups | 0 |
| Number of wait group operations | 0 |
| Number of cond vars | 0 |
| Number of cond var operations | 0 |
| Number of once | 3| 
| Number of once operations | 4 |


## Times
| Info | Value |
| - | - |
| Time for run without ADVOCATE | 1.017436 s |
| Time for run with ADVOCATE | 2.413929 s |
| Overhead of ADVOCATE | 137.256103 % |
| Analysis | 1.300918 s |


## Results
==================== Summary ====================

-------------------- Warning --------------------
1 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/Other/examples/bbolt/cmd/bbolt/main.go:1222@3128
	recv : /home/erikkassubek/Uni/HiWi/Other/examples/bbolt/cmd/bbolt/main.go:1428@144
2 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/Other/examples/bbolt/cmd/bbolt/main.go:1355@721638
	recv : /home/erikkassubek/Uni/HiWi/Other/examples/bbolt/cmd/bbolt/main.go:1428@3225
