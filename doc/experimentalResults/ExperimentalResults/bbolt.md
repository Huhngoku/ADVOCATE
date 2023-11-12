# bbolt

## Program 
| Info | Value |
| - | - |
| Number of go files | 73|
| Number of lines of code |19550|
## Trace 
| Info | Value |
| - | - |
| Number of routines | 16|
| Number of atomic variables | 119|
| Number of atomic operations | 447937|
| Number of channels | 5|
| Number of channel operations | 8|
| Number of mutexes | 15|
| Number of mutex operations | 110|
| Number of once variables | 4|
| Number of once operations | 5|
| Number of selects | 3|
| Number of select cases | 6|
| Number of executed select channel operations | 3|
| Number of executed select default cases | 0|
| Number of waitgroups | 0|
| Number of waitgroup operations | 0|
## Runtime 
| Info | Value |
| - | - |
| Runtime without modifications | 1.022|
| Runtime with modified runtime | 1.03|
| Runtime with modified runtime and trace creation | 2.268|
| Overhead of modified runtime [s] | 0.008|
| Overhead of modified runtime [\%] | 0.783|
| Overhead of modified runtime and trace creation [s] | 1.246|
| Overhead of modified runtime and trace creation [\%] | 121.918|
| Runtime for analysis [s] | 1.16|
## Found Results
==================== Summary ====================\
\
-------------------- Warning --------------------\
Found concurrent Send on same channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;send: /home/erikkassubek/Uni/HiWi/CoBuFi-Go/go-patch/src/runtime/mgcscavenge.go:652\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;send : /home/erikkassubek/Uni/HiWi/CoBuFi-Go/go-patch/src/runtime/mgcsweep.go:279\
Receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/Other/examples/bbolt/cmd/bbolt/main.go:1214\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/Other/examples/bbolt/cmd/bbolt/main.go:1420\
Receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/Other/examples/bbolt/cmd/bbolt/main.go:1347\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/Other/examples/bbolt/cmd/bbolt/main.go:1420\
