# overview Stats

## Program
| Info | Value |
| - | - |
| Number of go files | 1 |
| Number of lines | 941 |
| Number of non-empty lines | 677 |


## Trace
| Info | Value |
| - | - |
| Number of routines | 18 |
| Number of spawns | 4 |
| Number of atomics | 1 |
| Number of atomic operations | 2 |
| Number of channels | 1 |
| Number of channel operations | 3 |
| Number of selects | 0 |
| Number of select cases | 0 |
| Number of select channel operations | 0 |
| Number of select default operations | 0 |
| Number of mutexes | 1 |
| Number of mutex operations | 2 |
| Number of wait groups | 0 |
| Number of wait group operations | 0 |
| Number of cond vars | 0 |
| Number of cond var operations | 0 |
| Number of once | 0| 
| Number of once operations | 0 |


## Times
| Info | Value |
| - | - |
| Time for run without ADVOCATE | 0.101638 s |
| Time for run with ADVOCATE | 0.107787 s |
| Overhead of ADVOCATE | 6.049903 % |
| Analysis | 0.019930 s |


## Results
==================== Summary ====================

-------------------- Critical -------------------
1 Possible send on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:365@38
	send : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:355@32
-------------------- Warning --------------------
2 Possible receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:365@38
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:361@28
