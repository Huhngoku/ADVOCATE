# overview Stats

## Program
| Info | Value |
| - | - |
| Number of go files | 9 |
| Number of lines | 1505 |
| Number of non-empty lines | 958 |


## Trace
| Info | Value |
| - | - |
| Number of routines | 32 |
| Number of spawns | 18 |
| Number of atomics | 206 |
| Number of atomic operations | 1439 |
| Number of channels | 37 |
| Number of channel operations | 50 |
| Number of selects | 52 |
| Number of select cases | 107 |
| Number of select channel operations | 57 |
| Number of select default operations | 25 |
| Number of mutexes | 48 |
| Number of mutex operations | 395 |
| Number of wait groups | 4 |
| Number of wait group operations | 19 |
| Number of cond vars | 2 |
| Number of cond var operations | 11 |
| Number of once | 18| 
| Number of once operations | 234 |


## Times
| Info | Value |
| - | - |
| Time for run without ADVOCATE | 1.248361 s |
| Time for run with ADVOCATE | 1.309881 s |
| Overhead of ADVOCATE | 4.928062 % |
| Analysis | 0.062167 s |


## Results
==================== Summary ====================

-------------------- Critical -------------------
1 Potential mixed deadlock:
	locks: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/h2_bundle.go:9694@6073
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/h2_bundle.go:7903@5524
	send/close-recv: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/h2_bundle.go:9374@6059
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/h2_bundle.go:8312@5543
2 Possible negative waitgroup counter:
	done: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@881
	done/add: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@1643
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@1680
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@918
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@1694
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@932
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:651@1474
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:651@1477
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:322@1418
3 Possible negative waitgroup counter:
	done: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@918
	done/add: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@1643
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@1680
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@881
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@1694
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@932
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:651@1474
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:651@1477
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:322@1418
-------------------- Warning --------------------
4 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/fd_unix.go:105@957
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/fd_unix.go:118@956
5 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1242@1003
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1389@491
6 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/fd_unix.go:105@1719
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/fd_unix.go:118@1718
7 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/crypto/internal/randutil/randutil.go:28@1785
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/crypto/internal/randutil/randutil.go:31@1791
8 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/h2_bundle.go:931@5479
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/h2_bundle.go:904@5449
9 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1242@5484
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1389@1410
10 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/h2_bundle.go:9374@6059
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/h2_bundle.go:8312@5543
11 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/h2_bundle.go:9760@6355
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/h2_bundle.go:8477@5773
12 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/h2_bundle.go:8609@6423
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/h2_bundle.go:9619@6575
