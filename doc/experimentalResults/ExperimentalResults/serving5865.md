# serving5865

## Program 
| Info | Value |
| - | - |
| Number of go files | 1|
| Number of lines of code |60|
## Trace 
| Info | Value |
| - | - |
| Number of routines | 7|
| Number of atomic variables | 2|
| Number of atomic operations | 4|
| Number of channels | 3|
| Number of channel operations | 7|
| Number of mutexes | 2|
| Number of mutex operations | 4|
| Number of once variables | 0|
| Number of once operations | 0|
| Number of selects | 0|
| Number of select cases | 0|
| Number of executed select channel operations | 0|
| Number of executed select default cases | 0|
| Number of waitgroups | 0|
| Number of waitgroup operations | 0|
## Runtime 
| Info | Value |
| - | - |
| Runtime without modifications | 2.003|
| Runtime with modified runtime | 2.006|
| Runtime with modified runtime and trace creation | 2.008|
| Overhead of modified runtime [s] | 0.003|
| Overhead of modified runtime [\%] | 0.15|
| Overhead of modified runtime and trace creation [s] | 0.005|
| Overhead of modified runtime and trace creation [\%] | 0.25|
| Runtime for analysis [s] | 0.04|
## Found Results
==================== Summary ====================\
\
-------------------- Critical -------------------\
Possible send on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/Other/examples/serving5865/main.go:17\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;send : /home/erikkassubek/Uni/HiWi/Other/examples/serving5865/main.go:29\
=================================================\
