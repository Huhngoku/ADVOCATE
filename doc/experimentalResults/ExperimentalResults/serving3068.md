# serving3068

## Program 
| Info | Value |
| - | - |
| Number of go files | 1|
| Number of lines of code |82|
## Trace 
| Info | Value |
| - | - |
| Number of routines | 9|
| Number of atomic variables | 5|
| Number of atomic operations | 27|
| Number of channels | 3|
| Number of channel operations | 17|
| Number of mutexes | 1|
| Number of mutex operations | 2|
| Number of once variables | 1|
| Number of once operations | 1|
| Number of selects | 0|
| Number of select cases | 0|
| Number of executed select channel operations | 0|
| Number of executed select default cases | 0|
| Number of waitgroups | 2|
| Number of waitgroup operations | 18|
## Runtime 
| Info | Value |
| - | - |
| Runtime without modifications | 1.008|
| Runtime with modified runtime | 1.008|
| Runtime with modified runtime and trace creation | 1.011|
| Overhead of modified runtime [s] | 0|
| Overhead of modified runtime [\%] | 0|
| Overhead of modified runtime and trace creation [s] | 0.003|
| Overhead of modified runtime and trace creation [\%] | 0.298|
| Runtime for analysis [s] | 0.05|
## Found Results
==================== Summary ====================\
\
-------------------- Critical -------------------\
Possible send on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/Other/examples/serving3068/main.go:51\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;send : /home/erikkassubek/Uni/HiWi/Other/examples/serving3068/main.go:45\
-------------------- Warning --------------------\
Possible receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/Other/examples/serving3068/main.go:51\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/Other/examples/serving3068/main.go:31\
Receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/Other/examples/serving3068/main.go:51\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/Other/examples/serving3068/main.go:31\
=================================================\
