# overview Stats

## Program
| Info | Value |
| - | - |
| Number of go files | 1 |
| Number of lines | 94 |
| Number of non-empty lines | 81 |


## Trace
| Info | Value |
| - | - |
| Number of routines | 20 |
| Number of spawns | 6 |
| Number of atomics | 12 |
| Number of atomic operations | 31 |
| Number of channels | 4 |
| Number of channel operations | 4 |
| Number of selects | 0 |
| Number of select cases | 0 |
| Number of select channel operations | 0 |
| Number of select default operations | 0 |
| Number of mutexes | 3 |
| Number of mutex operations | 13 |
| Number of wait groups | 6 |
| Number of wait group operations | 10 |
| Number of cond vars | 0 |
| Number of cond var operations | 0 |
| Number of once | 0| 
| Number of once operations | 0 |


## Times
| Info | Value |
| - | - |
| Time for run without ADVOCATE | 0.003517 s |
| Time for run with ADVOCATE | 20.032289 s |
| Overhead of ADVOCATE | 569484.560705 % |
| Analysis | 0.043988 s |


## Results
==================== Summary ====================

-------------------- Critical -------------------
1 Potential leak on wait group:
	wait-group: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/cockroach1055/cockroach1055.go:45@96
	
2 Potential leak without possible partner:
	channel: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/cockroach1055/cockroach1055.go:77@88
	partner: -
3 Potential leak without possible partner:
	channel: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/cockroach1055/cockroach1055.go:77@106
	partner: -
4 Potential leak without possible partner:
	channel: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/cockroach1055/cockroach1055.go:93@39
	partner: -
5 Potential leak without possible partner:
	channel: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/cockroach1055/cockroach1055.go:77@85
	partner: -
