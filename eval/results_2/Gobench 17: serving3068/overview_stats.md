# overview Stats

## Program
| Info | Value |
| - | - |
| Number of go files | 1 |
| Number of lines | 77 |
| Number of non-empty lines | 63 |


## Trace
| Info | Value |
| - | - |
| Number of routines | 19 |
| Number of spawns | 5 |
| Number of atomics | 5 |
| Number of atomic operations | 26 |
| Number of channels | 1 |
| Number of channel operations | 12 |
| Number of selects | 0 |
| Number of select cases | 0 |
| Number of select channel operations | 0 |
| Number of select default operations | 0 |
| Number of mutexes | 1 |
| Number of mutex operations | 2 |
| Number of wait groups | 2 |
| Number of wait group operations | 17 |
| Number of cond vars | 0 |
| Number of cond var operations | 0 |
| Number of once | 1| 
| Number of once operations | 1 |


## Times
| Info | Value |
| - | - |
| Time for run without ADVOCATE | 1.003157 s |
| Time for run with ADVOCATE | 1.023428 s |
| Overhead of ADVOCATE | 2.020721 % |
| Replay without changes | 1.139767 s |
| Overhead of Replay | 13.618008 % s |
| Analysis | 0.011026 s |


## Results
==================== Summary ====================

-------------------- Critical -------------------
1 Possible send on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/serving3068/serving3068.go:49@94
	send : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/serving3068/serving3068.go:43@85
-------------------- Warning --------------------
2 Possible receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/serving3068/serving3068.go:49@94
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/serving3068/serving3068.go:29@80
3 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/serving3068/serving3068.go:49@94
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/serving3068/serving3068.go:29@91
