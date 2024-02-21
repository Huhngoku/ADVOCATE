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
| Number of routines | 17 |
| Number of spawns | 3 |
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
| Number of wait group operations | 2 |
| Number of cond vars | 0 |
| Number of cond var operations | 0 |
| Number of once | 0| 
| Number of once operations | 0 |


## Times
| Info | Value |
| - | - |
| Time for run without ADVOCATE | 0.104182 s |
| Time for run with ADVOCATE | 0.119072 s |
| Overhead of ADVOCATE | 14.292296 % |
| Replay without changes | 0.118555 s |
| Overhead of Replay | 13.796049 % s |
| Analysis | 0.040579 s |


## Results
==================== Summary ====================

-------------------- Critical -------------------
1 Potential leak on wait group:
	wait-group: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:825@29
	
