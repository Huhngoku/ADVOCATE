package main

import (
	"advocate"
	"flag"
	"gobench/cockroach1055"
	"gobench/cockroach1462"
	"gobench/etcd6873"
	"gobench/etcd7443"
	"gobench/etcd7492"
	"gobench/etcd7902"
	"gobench/grpc1353"
	"gobench/grpc1460"
	"gobench/grpc1687"
	"gobench/istio16224"
	"gobench/kubernetes10182"
	"gobench/kubernetes1321"
	"gobench/kubernetes26980"
	"gobench/kubernetes6632"
	"gobench/moby28462"
	"gobench/serving2137"
	"gobench/serving3068"
	"gobench/serving5865"
	"os"
	"time"
)

func main() {
	if true {
		// init tracing
		advocate.InitTracing(0)
		defer advocate.Finish()
	} else {
		// init replay
		advocate.EnableReplay()
		defer advocate.WaitForReplayFinish()
	}

	list := flag.Bool("l", false, "List tests. Do not run any test.")
	testCase := flag.Int("c", 0, "Test to run. If not set, all are run.")
	// replay := flag.Bool("r", false, "Replay")
	timeout := flag.Int("t", 0, "Timeout")
	flag.Parse()

	const n = 18
	testNames := [n]string{
		"Test 01: Cockroach1055",
		"Test 02: Cockroach1462",
		"Test 03: Etcd6873",
		"Test 04: Etcd7443",
		"Test 05: Etcd7492",
		"Test 06: Etcd7902",
		"Test 07: Grpc1353",
		"Test 08: Grpc1460",
		"Test 09: Grpc1687",
		"Test 10: Istio16224",
		"Test 11: Kubernetes1321",
		"Test 12: Kubernetes6632",
		"Test 13: Kubernetes10182",
		"Test 14: Kubernetes26980",
		"Test 15: Moby28462",
		"Test 16: Serving2137",
		"Test 17: Serving3068",
		"Test 18: Serving5865",
	}

	testFuncs := [n]func(){
		cockroach1055.Cockroach1055,
		cockroach1462.Cockroach1462,
		etcd6873.Etcd6873,
		etcd7443.Etcd7443,
		etcd7492.Etcd7492,
		etcd7902.Etcd7902,
		grpc1353.Grpc1353,
		grpc1460.Grpc1460,
		grpc1687.Grpc1687,
		istio16224.Istio16224,
		kubernetes1321.Kubernetes1321,
		kubernetes6632.Kubernetes6632,
		kubernetes10182.Kubernetes10182,
		kubernetes26980.Kubernetes26980,
		moby28462.Moby28462,
		serving2137.Serving2137,
		serving3068.Serving3068,
		serving5865.Serving5865,
	}

	if list != nil && *list {
		for i := 0; i < n; i++ {
			println(testNames[i])
		}
		return
	}

	// cancel test if time has run out
	go func() {
		if timeout != nil && *timeout != 0 {
			time.Sleep(time.Duration(*timeout) * time.Second)
			// advocate.Finish()
			os.Exit(42)
		}
	}()

	if testCase != nil && *testCase != 0 {
		if *testCase > n {
			println("Invalid test case")
			return
		}

		println(testNames[*testCase-1])
		testFuncs[*testCase-1]()
	} else {
		for i := 0; i < n; i++ {
			println(testNames[i])
			testFuncs[i]()
			println("Done: ", i+1, " of ", n)
			time.Sleep(5 * time.Second)
		}
	}
}
