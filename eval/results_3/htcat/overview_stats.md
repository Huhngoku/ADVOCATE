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
| Number of atomics | 368 |
| Number of atomic operations | 24994 |
| Number of channels | 74 |
| Number of channel operations | 108 |
| Number of selects | 81 |
| Number of select cases | 173 |
| Number of select channel operations | 133 |
| Number of select default operations | 20 |
| Number of mutexes | 85 |
| Number of mutex operations | 10748 |
| Number of wait groups | 3 |
| Number of wait group operations | 22 |
| Number of cond vars | 5 |
| Number of cond var operations | 1338 |
| Number of once | 20| 
| Number of once operations | 428 |


## Times
Invalid time file
13.243596,0.232368,## Results
==================== Summary ====================

-------------------- Critical -------------------
1 Potential leak on conditional variable:
	conditional: /home/erikkassubek/go/pkg/mod/github.com/htcat/htcat@v1.0.2/eager_reader.go:64@17806
	
2 Potential leak on mutex:
	mutex: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/sync/cond.go:88@18118
	
3 Potential leak without possible partner:
	channel: /home/erikkassubek/go/pkg/mod/github.com/htcat/htcat@v1.0.2/eager_reader.go:105@4447
	partner: -
4 Potential leak without possible partner:
	channel: /home/erikkassubek/go/pkg/mod/github.com/htcat/htcat@v1.0.2/defrag.go:173@5720
	partner: -
5 Potential leak without possible partner:
	channel: /home/erikkassubek/go/pkg/mod/github.com/htcat/htcat@v1.0.2/defrag.go:173@5847
	partner: -
6 Potential leak without possible partner:
	channel: /home/erikkassubek/go/pkg/mod/github.com/htcat/htcat@v1.0.2/defrag.go:173@6644
	partner: -
7 Potential leak without possible partner:
	channel: /home/erikkassubek/go/pkg/mod/github.com/htcat/htcat@v1.0.2/defrag.go:173@6991
	partner: -
8 Possible negative waitgroup counter:
	done: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@454
	done/add: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@4799
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@4840
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@491
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4861
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4866
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4886
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4904
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@505
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:651@4616
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:651@4621
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:322@4473
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:322@4507
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:322@4591
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:322@4656
9 Possible negative waitgroup counter:
	done: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@491
	done/add: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@454
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@4799
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@4840
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4861
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4866
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4886
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4904
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@505
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:651@4616
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:651@4621
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:322@4473
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:322@4507
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:322@4591
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:322@4656
10 Possible negative waitgroup counter:
	done: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@4799
	done/add: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@454
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@4840
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@491
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4861
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4866
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4886
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4904
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@505
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:322@4591
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:322@4656
11 Possible negative waitgroup counter:
	done: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@4840
	done/add: 
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@454
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@4799
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/dnsclient_unix.go:655@491
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4861
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4866
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4886
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@4904
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:358@505
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:322@4591
		/home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/lookup.go:322@4656
-------------------- Warning --------------------
12 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/crypto/internal/randutil/randutil.go:28@559
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/crypto/internal/randutil/randutil.go:31@565
13 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1242@4221
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1389@104
14 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/crypto/internal/randutil/randutil.go:28@559
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/crypto/internal/randutil/randutil.go:31@5021
15 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/crypto/internal/randutil/randutil.go:28@559
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/crypto/internal/randutil/randutil.go:31@5108
16 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/crypto/internal/randutil/randutil.go:28@559
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/crypto/internal/randutil/randutil.go:31@5174
17 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/crypto/internal/randutil/randutil.go:28@559
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/crypto/internal/randutil/randutil.go:31@5186
18 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1242@5414
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1389@4501
19 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1242@5565
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1389@4581
20 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1242@6307
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1389@4648
21 Found receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1242@6453
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/net/http/transport.go:1389@4446
22 Found receive on closed channel:
	close: 
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/advocate/advocate.go:110@18163
23 Found receive on closed channel:
	close: 
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/advocate/advocate.go:110@18183
24 Found receive on closed channel:
	close: 
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/advocate/advocate.go:110@18203
25 Found receive on closed channel:
	close: 
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/advocate/advocate.go:110@18223
26 Found receive on closed channel:
	close: 
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/advocate/advocate.go:110@18243
27 Found receive on closed channel:
	close: 
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/advocate/advocate.go:110@18263
28 Found receive on closed channel:
	close: 
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/advocate/advocate.go:110@18283
29 Found receive on closed channel:
	close: 
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/advocate/advocate.go:110@18303
30 Found receive on closed channel:
	close: 
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/advocate/advocate.go:110@18323
31 Found receive on closed channel:
	close: 
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/advocate/advocate.go:110@18343
32 Found receive on closed channel:
	close: 
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/advocate/advocate.go:110@18363
33 Found receive on closed channel:
	close: 
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/advocate/advocate.go:110@18383
34 Found receive on closed channel:
	close: 
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/advocate/advocate.go:110@18403
35 Found receive on closed channel:
	close: 
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/go-patch/src/advocate/advocate.go:110@18423
