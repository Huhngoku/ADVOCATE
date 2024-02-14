# overview Stats

## Program
| Info | Value |
| - | - |
| Number of go files | 7 |
| Number of lines | 728 |
| Number of non-empty lines | 495 |


## Trace
| Info | Value |
| - | - |
| Number of routines | 51 |
| Number of spawns | 37 |
| Number of atomics | 190 |
| Number of atomic operations | 2244 |
| Number of channels | 68 |
| Number of channel operations | 97 |
| Number of selects | 83 |
| Number of select cases | 178 |
| Number of select channel operations | 138 |
| Number of select default operations | 20 |
| Number of mutexes | 76 |
| Number of mutex operations | 688 |
| Number of wait groups | 3 |
| Number of wait group operations | 22 |
| Number of cond vars | 4 |
| Number of cond var operations | 59 |
| Number of once | 17| 
| Number of once operations | 88 |


## Times
| Info | Value |
| - | - |
| Time for run without ADVOCATE | 7.453610 s |
| Time for run with ADVOCATE | 3.794589 s |
| Overhead of ADVOCATE | 0.000000 % |
| Analysis | 3.555098 s |


## Results
==================== Summary ====================

-------------------- Critical -------------------
1 Potential mixed deadlock:
	locks: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/internal/singleflight/singleflight.go:95@4841
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/internal/singleflight/singleflight.go:71@4485
	send/close-recv: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/internal/singleflight/singleflight.go:101@4846
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:332@4500
2 Potential mixed deadlock:
	locks: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:2725@42543
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:2568@4246
	send/close-recv: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:2257@42542
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:2198@42522
3 Potential leak with possible partner:
	channel: 
	partner: /home/erikkassubek/go/pkg/mod/github.com/htcat/htcat@v1.0.2/defrag.go:154@51102
4 Potential leak with possible partner:
	channel: 
	partner: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:2744@51090
5 Potential leak on mutex:
	mutex: /home/erikkassubek/go/pkg/mod/github.com/htcat/htcat@v1.0.2/http.go:173@49531
	
6 Potential leak without possible partner:
	channel: /home/erikkassubek/go/pkg/mod/github.com/htcat/htcat@v1.0.2/eager_reader.go:105@50775
	partner: -
7 Possible negative waitgroup counter:
	done: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@452
	done/add: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@4802
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@4839
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@489
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4859
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4863
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4873
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4897
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@503
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:651@4612
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:651@4623
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:322@4483
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:322@4495
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:322@4615
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:322@4665
8 Possible negative waitgroup counter:
	done: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@4802
	done/add: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@452
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@4839
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@489
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4859
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4863
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4873
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4897
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@503
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:322@4665
9 Possible negative waitgroup counter:
	done: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@489
	done/add: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@452
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@4802
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@4839
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4859
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4863
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4873
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4897
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@503
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:651@4612
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:651@4623
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:322@4483
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:322@4495
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:322@4615
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:322@4665
-------------------- Warning --------------------
10 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/crypto/internal/randutil/randutil.go:28@557
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/crypto/internal/randutil/randutil.go:31@563
11 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1242@4227
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1389@104
12 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/crypto/internal/randutil/randutil.go:28@557
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/crypto/internal/randutil/randutil.go:31@5126
13 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/crypto/internal/randutil/randutil.go:28@557
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/crypto/internal/randutil/randutil.go:31@5175
14 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/crypto/internal/randutil/randutil.go:28@557
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/crypto/internal/randutil/randutil.go:31@5192
15 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/crypto/internal/randutil/randutil.go:28@557
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/crypto/internal/randutil/randutil.go:31@5209
16 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1242@5699
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1389@4468
17 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1242@5884
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1389@4656
18 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1242@5967
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1389@4478
19 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1242@6145
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1389@4589
20 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:2257@42542
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:2198@42522
21 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:2744@42621
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:2417@4325
22 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:2744@50685
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:2417@6206
