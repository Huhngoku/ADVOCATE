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
| Number of routines | 9 |
| Number of spawns | 3 |
| Number of atomics | 1 |
| Number of atomic operations | 3 |
| Number of channels | 1 |
| Number of channel operations | 1 |
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
No time file provided


## Results
==================== Summary ====================

-------------------- Critical -------------------
1 Potential leak on mutex:
	mutex: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:739@23
	
2 Potential leak without possible partner:
	channel: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:745@22
	partner: -
