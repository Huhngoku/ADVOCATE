# overview Stats

## Program
| Info | Value |
| - | - |
| Number of go files | 1 |
| Number of lines | 941 |
| Number of non-empty lines | 56 |


## Trace
| Info | Value |
| - | - |
| Number of routines | 11 |
| Number of spawns | 5 |
| Number of atomics | 1 |
| Number of atomic operations | 4 |
| Number of channels | 1 |
| Number of channel operations | 4 |
| Number of selects | 0 |
| Number of select cases | 0 |
| Number of select channel operations | 0 |
| Number of select default operations | 0 |
| Number of mutexes | 1 |
| Number of mutex operations | 4 |
| Number of wait groups | 0 |
| Number of wait group operations | 0 |
| Number of cond vars | 0 |
| Number of cond var operations | 0 |
| Number of once | 0| 
| Number of once operations | 0 |


## Times
No time file provided


## Results
==================== Summary ====================

-------------------- Critical -------------------
1 Found concurrent Recv on same channel:
	recv: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:640@35
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:645@24
