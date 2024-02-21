# overview Stats

## Program
| Info | Value |
| - | - |
| Number of go files | 1 |
| Number of lines | 214 |
| Number of non-empty lines | 190 |


## Trace
| Info | Value |
| - | - |
| Number of routines | 19 |
| Number of spawns | 5 |
| Number of atomics | 4 |
| Number of atomic operations | 29 |
| Number of channels | 2 |
| Number of channel operations | 5 |
| Number of selects | 0 |
| Number of select cases | 0 |
| Number of select channel operations | 0 |
| Number of select default operations | 0 |
| Number of mutexes | 4 |
| Number of mutex operations | 28 |
| Number of wait groups | 0 |
| Number of wait group operations | 0 |
| Number of cond vars | 0 |
| Number of cond var operations | 0 |
| Number of once | 0| 
| Number of once operations | 0 |


## Times
| Info | Value |
| - | - |
| Time for run without ADVOCATE | 0.000958 s |
| Time for run with ADVOCATE | 0.004120 s |
| Overhead of ADVOCATE | 330.062630 % |
| Analysis | 0.009594 s |


## Results
==================== Summary ====================

-------------------- Warning --------------------
1 Possible receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/etcd7443/etcd7443.go:181@63
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/etcd7443/etcd7443.go:40@35
2 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/etcd7443/etcd7443.go:211@120
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/GoBench/etcd7443/etcd7443.go:213@34
