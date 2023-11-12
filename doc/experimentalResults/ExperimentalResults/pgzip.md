# pgzip

## Program 
| Info | Value |
| - | - |
| Number of go files | 3|
| Number of lines of code |1224|
## Trace 
| Info | Value |
| - | - |
| Number of routines | 1556|
| Number of atomic variables | 12|
| Number of atomic operations | 375745|
| Number of channels | 5|
| Number of channel operations | 3096|
| Number of mutexes | 4|
| Number of mutex operations | 12|
| Number of once variables | 3|
| Number of once operations | 3|
| Number of selects | 3082|
| Number of select cases | 6164|
| Number of executed select channel operations | 3082|
| Number of executed select default cases | 0|
| Number of waitgroups | 1|
| Number of waitgroup operations | 4624|
## Runtime 
| Info | Value |
| - | - |
| Runtime without modifications | 7.191|
| Runtime with modified runtime | 7.286|
| Runtime with modified runtime and trace creation | 8.52|
| Overhead of modified runtime [s] | 0.095|
| Overhead of modified runtime [\%] | 1.321|
| Overhead of modified runtime and trace creation [s] | 1.329|
| Overhead of modified runtime and trace creation [\%] | 18.481|
| Runtime for analysis [s] | 157.52|
## Found Results
==================== Summary ====================\
\
-------------------- Warning --------------------\
Receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/Other/examples/pgzip/gunzip.go:395\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/Other/examples/pgzip/gunzip.go:345\
