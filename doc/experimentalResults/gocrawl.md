# gocrawl

## Program 
| Info | Value |
| - | - |
| Number of go files | 9|
| Number of lines of code |1497|
## Trace 
| Info | Value |
| - | - |
| Number of routines | 32|
| Number of atomic variables | 454|
| Number of atomic operations | 5122|
| Number of channels | 46|
| Number of channel operations | 70|
| Number of mutexes | 64|
| Number of mutex operations | 715|
| Number of once variables | 26|
| Number of once operations | 670|
| Number of selects | 32|
| Number of select cases | 68|
| Number of executed select channel operations | 30|
| Number of executed select default cases | 2|
| Number of waitgroups | 4|
| Number of waitgroup operations | 19|
## Runtime 
| Info | Value |
| - | - |
| Runtime without modifications | 2.328|
| Runtime with modified runtime | 2.301|
| Runtime with modified runtime and trace creation | 2.361|
| Overhead of modified runtime [s] | 0|
| Overhead of modified runtime [\%] | 0|
| Overhead of modified runtime and trace creation [s] | 0.033|
| Overhead of modified runtime and trace creation [\%] | 1.418|
| Runtime for analysis [s] | 0.6565|
## Found Results
==================== Summary ====================\
\
-------------------- Warning --------------------\
Receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/CoBuFi-Go/go-patch/src/net/fd_unix.go:105\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/CoBuFi-Go/go-patch/src/net/fd_unix.go:118\
Receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/CoBuFi-Go/go-patch/src/net/http/transport.go:1242\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/CoBuFi-Go/go-patch/src/net/http/transport.go:1389\
Receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/CoBuFi-Go/go-patch/src/crypto/internal/randutil/randutil.go:28\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/CoBuFi-Go/go-patch/src/crypto/internal/randutil/randutil.go:31\
Receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/CoBuFi-Go/go-patch/src/net/http/h2_bundle.go:931\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/CoBuFi-Go/go-patch/src/net/http/h2_bundle.go:904\
Receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/CoBuFi-Go/go-patch/src/net/http/h2_bundle.go:9374\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/CoBuFi-Go/go-patch/src/net/http/h2_bundle.go:8312\
Receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/CoBuFi-Go/go-patch/src/net/http/h2_bundle.go:9760\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/CoBuFi-Go/go-patch/src/net/http/h2_bundle.go:8477\
Receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/CoBuFi-Go/go-patch/src/net/http/h2_bundle.go:8609\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/CoBuFi-Go/go-patch/src/net/http/h2_bundle.go:9619\
Receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/Other/examples/gocrawl/crawler.go:307\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/Other/examples/gocrawl/worker.go:66\
\
=================================================\
Total runtime: 39.916288ms\
=================================================\
