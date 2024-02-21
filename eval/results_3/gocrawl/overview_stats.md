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
| Number of atomics | 389 |
| Number of atomic operations | 4295 |
| Number of channels | 38 |
| Number of channel operations | 52 |
| Number of selects | 54 |
| Number of select cases | 111 |
| Number of select channel operations | 59 |
| Number of select default operations | 26 |
| Number of mutexes | 58 |
| Number of mutex operations | 435 |
| Number of wait groups | 4 |
| Number of wait group operations | 19 |
| Number of cond vars | 2 |
| Number of cond var operations | 11 |
| Number of once | 24| 
| Number of once operations | 574 |


## Times
| Info | Value |
| - | - |
| Time for run without ADVOCATE | 1.350567 s |
| Time for run with ADVOCATE | 1.380675 s |
| Overhead of ADVOCATE | 2.229286 % |
| Analysis | 0.057362 s |


## Results
==================== Summary ====================

-------------------- Critical -------------------
1 Possible negative waitgroup counter:
	done: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@878
	done/add: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@1631
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@1668
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@915
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@1682
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@929
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:651@1463
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:651@1466
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:322@1407
2 Possible negative waitgroup counter:
	done: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@915
	done/add: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@1631
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@1668
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@878
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@1682
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@929
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:651@1463
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:651@1466
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:322@1407
-------------------- Warning --------------------
3 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/fd_unix.go:105@954
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/fd_unix.go:118@953
4 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1242@1000
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1389@497
5 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/fd_unix.go:105@1707
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/fd_unix.go:118@1706
6 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/crypto/internal/randutil/randutil.go:28@1773
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/crypto/internal/randutil/randutil.go:31@1779
7 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/h2_bundle.go:931@5469
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/h2_bundle.go:904@5439
8 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1242@5474
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1389@1399
9 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/h2_bundle.go:9374@6045
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/h2_bundle.go:8312@5533
10 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/h2_bundle.go:9760@6250
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/h2_bundle.go:8477@5763
11 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/h2_bundle.go:8609@6323
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/h2_bundle.go:9619@6553
