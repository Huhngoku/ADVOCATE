# Changed files

The following is a list of all files in the Go runtime that have been 
added or changed.

Added files:

- src/runtime/dedego_routine.go
- src/runtime/dedego_trace.go
- src/runtime/dedego_util.go
- src/runtime/internal/atomic/dedego_atomic.go

Changed files (marked with DEDEGO-CHANGE):

- src/cmd/cgo/internal/test/testx.go
- src/runtime/proc.go
- src/runtime/runtime2.go
- src/runtime/chan.go
- src/runtime/select.go
- src/runtime/internal/atomic/doc.go
- src/runtime/internal/atomic/atomic_amd64.go
- src/runtime/internal/atomic/atomic_amd64.s
- src/sync/mutex.go
- src/sync/rwmutex.go
- src/sync/waitgroup.go
- src/sync/once.go

Disabled Tests (files contain disabled tests, marked with DEDEGO-REMOVE_TEST): 

- src/cmd/cgo/internal/test/cgo_test.go
- src/cmd/dist/test.go
- src/cmd/go/stript_test.go
- src/cmd/compile/internal/types2/sizeof_test.go
- src/context/x_test.go
- src/crypto/internal/nistec/nistec_test.go
- src/crypto/tls/tls_test.go
- src/go/build/deps_test.
- src/go/types/sizeof_test.go
- src/internal/intern/inter_test.go
- src/log/slog/text_handler_test.go
- src/net/netip/netip_test.go
- src/runtime/crash_cgo_test.go
- src/runtime/sizeof_test.go
- src/runtime/align_test.go
- src/runtime/metrics_test.go
- src/net/tcpsock_test.go
- src/reflect/all_test.go
- src/os/signal/signal_test.go
- src/go/internal/gccgoimporter/importer_test.go