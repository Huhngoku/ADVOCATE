package main

import (
	"io"
	"sync"

	"github.com/ErikKassubek/deadlockDetectorGo/src/dedego"
)

type eagerReader struct {
	closeNotify dedego.Chan[struct{}]
	rc          io.ReadCloser

	buf   []byte
	more  *sync.Cond
	begin int
	end   int

	lastErr error
}

func newEagerReader(r io.ReadCloser, bufSz int64) *eagerReader {
	er := eagerReader{
		closeNotify: dedego.NewChan[struct{}](int(0)),
		rc:          r,
		buf:         make([]byte, bufSz, bufSz),
	}

	er.more = sync.NewCond(new(sync.Mutex))
	func() {
		DedegoRoutineIndex := dedego.SpawnPre()
		go func() {
			dedego.SpawnPost(DedegoRoutineIndex)
			{

				er.buffer()
			}
		}()
	}()

	return &er
}

func (er *eagerReader) buffer() {
	for er.lastErr == nil && er.end != len(er.buf) {
		var n int

		er.more.L.Lock()
		n, er.lastErr = er.rc.Read(er.buf[er.end:])
		er.end += n

		er.more.Broadcast()
		er.more.L.Unlock()
	}
}

func (er *eagerReader) writeOnce(dst io.Writer) (int64, error) {

	er.more.L.Lock()
	defer er.more.L.Unlock()

	for er.begin == er.end {
		if er.lastErr != nil {
			return 0, er.lastErr
		}

		if er.begin == len(er.buf) {
			return 0, io.EOF
		}

		er.more.Wait()
	}

	n, err := dst.Write(er.buf[er.begin:er.end])
	er.begin += n
	return int64(n), err
}

func (er *eagerReader) WriteTo(dst io.Writer) (int64, error) {
	var written int64

	for {
		n, err := er.writeOnce(dst)
		written += n
		switch err {
		case io.EOF:

			return 0, nil
		case nil:

			continue
		default:

			return written, err
		}
	}
}

func (er *eagerReader) Close() error {
	err := er.rc.Close()
	er.closeNotify.Send(struct{}{})

	return err
}

func (er *eagerReader) WaitClosed() {
	er.closeNotify.Receive()

}
