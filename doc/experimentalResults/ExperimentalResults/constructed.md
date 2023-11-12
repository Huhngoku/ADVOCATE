# constructed

## Program 
| Info | Value |
| - | - |
| Number of go files | 0|
| Number of lines of code |0|
## Trace 
| Info | Value |
| - | - |
| Number of routines | 34|
| Number of atomic variables | 16|
| Number of atomic operations | 45|
| Number of channels | 24|
| Number of channel operations | 50|
| Number of mutexes | 10|
| Number of mutex operations | 30|
| Number of once variables | 4|
| Number of once operations | 7|
| Number of selects | 0|
| Number of select cases | 0|
| Number of executed select channel operations | 0|
| Number of executed select default cases | 0|
| Number of waitgroups | 2|
| Number of waitgroup operations | 6|
## Runtime 
| Info | Value |
| - | - |
| Runtime without modifications | 4.112|
| Runtime with modified runtime | 4.114|
| Runtime with modified runtime and trace creation | 4.123|
| Overhead of modified runtime [s] | 0.002|
| Overhead of modified runtime [\%] | 0.049|
| Overhead of modified runtime and trace creation [s] | 0.011|
| Overhead of modified runtime and trace creation [\%] | 0.268|
| Runtime for analysis [s] | 0.009|
## Found Results
==================== Summary ====================\
\
-------------------- Critical -------------------\
Possible send on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:129\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;send : /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:121\
Possible send on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:143\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;send : /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:138\
Possible send on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:180\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;send : /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:186\
Possible send on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:181\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;send : /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:196\
Possible send on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:313\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;send : /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:301\
Possible send on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:360\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;send : /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:350\
Possible send on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:436\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;send : /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:428\
-------------------- Warning --------------------\
Possible receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:129\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:125\
Possible receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:143\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:139\
Receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:168\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:165\
Possible receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:180\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:197\
Possible receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:181\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:191\
Receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:339\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:333\
Possible receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:360\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:356\
=================================================\