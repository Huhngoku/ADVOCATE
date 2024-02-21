# overview Stats

## Program
| Info | Value |
| - | - |
| Number of go files | 1 |
| Number of lines | 107 |
| Number of non-empty lines | 87 |


## Trace
| Info | Value |
| - | - |
| Number of routines | 17 |
| Number of spawns | 3 |
| Number of atomics | 1 |
| Number of atomic operations | 4 |
| Number of channels | 3 |
| Number of channel operations | 4 |
| Number of selects | 2 |
| Number of select cases | 4 |
| Number of select channel operations | 4 |
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
| Info | Value |
| - | - |
| Time for run without ADVOCATE | 1.005579 s |
| Time for run with ADVOCATE | 1.011280 s |
| Overhead of ADVOCATE | 0.566937 % |
| Replay without changes | 1.014757 s |
| Overhead of Replay | 0.912708 % s |
| Analysis | 0.009031 s |


## Results
==================== Summary ====================

-------------------- Critical -------------------
1 Potential mixed deadlock:
	locks: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/istio16224/istio16224.go:92@41
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/istio16224/istio16224.go:102@34
	send/close-recv: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/istio16224/istio16224.go:94@44
		/home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/istio16224/istio16224.go:104@40
2 Potential leak with possible partner:
	channel: 
	partner: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/istio16224/istio16224.go:106@51
