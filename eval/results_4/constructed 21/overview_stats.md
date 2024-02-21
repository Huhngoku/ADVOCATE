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
| Number of routines | 18 |
| Number of spawns | 4 |
| Number of atomics | 0 |
| Number of atomic operations | 0 |
| Number of channels | 1 |
| Number of channel operations | 4 |
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
| Time for run without ADVOCATE | 0.301506 s |
| Time for run with ADVOCATE | 0.310064 s |
| Overhead of ADVOCATE | 2.838418 % |
| Replay without changes | 0.309105 s |
| Overhead of Replay | 2.520348 % s |
| Analysis | 0.014377 s |


## Results
==================== Summary ====================

-------------------- Critical -------------------
1 Found concurrent Recv on same channel:
	recv: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:452@32
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:456@29
