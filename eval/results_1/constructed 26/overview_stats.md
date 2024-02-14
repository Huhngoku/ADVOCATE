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
| Number of routines | 19 |
| Number of spawns | 5 |
| Number of atomics | 1 |
| Number of atomic operations | 4 |
| Number of channels | 0 |
| Number of channel operations | 0 |
| Number of selects | 0 |
| Number of select cases | 0 |
| Number of select channel operations | 0 |
| Number of select default operations | 0 |
| Number of mutexes | 0 |
| Number of mutex operations | 0 |
| Number of wait groups | 1 |
| Number of wait group operations | 4 |
| Number of cond vars | 0 |
| Number of cond var operations | 0 |
| Number of once | 0| 
| Number of once operations | 0 |


## Times
| Info | Value |
| - | - |
| Time for run without ADVOCATE | 0.504233 s |
| Time for run with ADVOCATE | 0.508737 s |
| Overhead of ADVOCATE | 0.893238 % |
| Analysis | 0.009999 s |


## Results
==================== Summary ====================

-------------------- Critical -------------------
1 Possible negative waitgroup counter:
	done: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:531@34
	done/add: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:535@36
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:522@30
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:526@32
2 Possible negative waitgroup counter:
	done: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:535@36
	done/add: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:531@34
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:522@30
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:526@32
