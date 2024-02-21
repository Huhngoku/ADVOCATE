# overview Stats

## Program
| Info | Value |
| - | - |
| Number of go files | 1 |
| Number of lines | 55 |
| Number of non-empty lines | 41 |


## Trace
| Info | Value |
| - | - |
| Number of routines | 17 |
| Number of spawns | 3 |
| Number of atomics | 2 |
| Number of atomic operations | 4 |
| Number of channels | 1 |
| Number of channel operations | 2 |
| Number of selects | 0 |
| Number of select cases | 0 |
| Number of select channel operations | 0 |
| Number of select default operations | 0 |
| Number of mutexes | 2 |
| Number of mutex operations | 4 |
| Number of wait groups | 0 |
| Number of wait group operations | 0 |
| Number of cond vars | 0 |
| Number of cond var operations | 0 |
| Number of once | 0| 
| Number of once operations | 0 |


## Times
| Info | Value |
| - | - |
| Time for run without ADVOCATE | 2.004197 s |
| Time for run with ADVOCATE | 2.022659 s |
| Overhead of ADVOCATE | 0.921167 % |
| Analysis | 0.010226 s |


## Results
==================== Summary ====================

-------------------- Critical -------------------
1 Possible send on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/serving5865/serving5865.go:15@44
	send : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/serving5865/serving5865.go:27@42
