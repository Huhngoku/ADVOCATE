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
| Number of atomic operations | 4 |
| Number of channels | 1 |
| Number of channel operations | 2 |
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
1 Potential mixed deadlock:
	locks: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:691@19
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:697@26
	send/close-recv: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:693@25
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:698@29
