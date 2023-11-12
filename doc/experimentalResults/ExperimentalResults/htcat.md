# htcat

## Program 
| Info | Value |
| - | - |
| Number of go files | 7|
| Number of lines of code |721|
## Trace 
| Info | Value |
| - | - |
| Number of routines | 50|
| Number of atomic variables | 390|
| Number of atomic operations | 27783|
| Number of channels | 77|
| Number of channel operations | 122|
| Number of mutexes | 90|
| Number of mutex operations | 10416|
| Number of once variables | 21|
| Number of once operations | 432|
| Number of selects | 52|
| Number of select cases | 115|
| Number of executed select channel operations | 52|
| Number of executed select default cases | 0|
| Number of waitgroups | 3|
| Number of waitgroup operations | 22|
## Runtime 
| Info | Value |
| - | - |
| Runtime without modifications | 2.342|
| Runtime with modified runtime | 2.344|
| Runtime with modified runtime and trace creation | 2.351|
| Overhead of modified runtime [s] | 0.002|
| Overhead of modified runtime [\%] | 0.085|
| Overhead of modified runtime and trace creation [s] | 0.009|
| Overhead of modified runtime and trace creation [\%] | 0.384|
| Runtime for analysis [s] | 0.382|
## Found Results
==================== Summary ====================\
\
-------------------- Warning --------------------\
Receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/CoBuFi-Go/go-patch/src/crypto/internal/randutil/randutil.go:28\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/CoBuFi-Go/go-patch/src/crypto/internal/randutil/randutil.go:31\
Receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/CoBuFi-Go/go-patch/src/net/http/transport.go:1242\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/CoBuFi-Go/go-patch/src/net/http/transport.go:1389\
Receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/CoBuFi-Go/go-patch/src/net/http/transport.go:2257\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/CoBuFi-Go/go-patch/src/net/http/transport.go:2198\
Receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/CoBuFi-Go/go-patch/src/net/http/transport.go:2744\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/CoBuFi-Go/go-patch/src/net/http/transport.go:2417\
