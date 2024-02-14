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
| Number of routines | 17 |
| Number of spawns | 3 |
| Number of atomics | 0 |
| Number of atomic operations | 0 |
| Number of channels | 1 |
| Number of channel operations | 2 |
| Number of selects | 0 |
| Number of select cases | 0 |
| Number of select channel operations | 0 |
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
| Time for run without ADVOCATE | 0.201463 s |
| Time for run with ADVOCATE | 0.206858 s |
| Overhead of ADVOCATE | 2.677911 % |
| Analysis | 0.045834 s |


## Results
==================== Summary ====================

-------------------- Critical -------------------
1 Possible send on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:441@29
	send : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:433@27
