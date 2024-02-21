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
| Number of atomics | 0 |
| Number of atomic operations | 0 |
| Number of channels | 1 |
| Number of channel operations | 3 |
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
| Time for run without ADVOCATE | 0.104261 s |
| Time for run with ADVOCATE | 0.118889 s |
| Overhead of ADVOCATE | 14.030174 % |
| Replay without changes | 0.118925 s |
| Overhead of Replay | 14.064703 % s |
| Analysis | 0.047231 s |


## Results
==================== Summary ====================

-------------------- Critical -------------------
1 Potential leak with possible partner:
	channel: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:806@30
	partner: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:810@33
