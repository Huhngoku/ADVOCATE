# overview Stats

## Program
| Info | Value |
| - | - |
| Number of go files | 1 |
| Number of lines | 109 |
| Number of non-empty lines | 92 |


## Trace
| Info | Value |
| - | - |
| Number of routines | 18 |
| Number of spawns | 4 |
| Number of atomics | 0 |
| Number of atomic operations | 0 |
| Number of channels | 1 |
| Number of channel operations | 1 |
| Number of selects | 7 |
| Number of select cases | 14 |
| Number of select channel operations | 10 |
| Number of select default operations | 2 |
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
| Time for run without ADVOCATE | 1.015832 s |
| Time for run with ADVOCATE | 1.022709 s |
| Overhead of ADVOCATE | 0.676982 % |
| Analysis | 0.009633 s |


## Results
==================== Summary ====================

-------------------- Warning --------------------
1 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/grpc1687/grpc1687.go:40@44
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/grpc1687/grpc1687.go:49@43
