// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"runtime"
	"sync"
	"time"

	"github.com/ErikKassubek/deadlockDetectorGo/src/dedego"
	"github.com/klauspost/compress/flate"
)

const (
	defaultBlockSize = 1 << 20
	tailSize         = 16384
	defaultBlocks    = 4
)

const (
	NoCompression       = flate.NoCompression
	BestSpeed           = flate.BestSpeed
	BestCompression     = flate.BestCompression
	DefaultCompression  = flate.DefaultCompression
	ConstantCompression = flate.ConstantCompression
	HuffmanOnly         = flate.HuffmanOnly
)

type Writer struct {
	Header
	w             io.Writer
	level         int
	wroteHeader   bool
	blockSize     int
	blocks        int
	currentBuffer []byte
	prevTail      []byte
	digest        hash.Hash32
	size          int
	closed        bool
	buf           [10]byte
	errMu         dedego.RWMutex
	err           error
	pushedErr     dedego.Chan[struct{}]
	results       dedego.Chan[result]
	dictFlatePool sync.Pool
	dstPool       sync.Pool
	wg            sync.WaitGroup
}

type result struct {
	result        dedego.Chan[[]byte]
	notifyWritten dedego.Chan[struct{}]
}

func (z *Writer) SetConcurrency(blockSize, blocks int) error {

	if blockSize <= tailSize {
		return fmt.Errorf("gzip: block size cannot be less than or equal to %d", tailSize)
	}
	if blocks <= 0 {
		return errors.New("gzip: blocks cannot be zero or less")
	}
	if blockSize == z.blockSize && blocks == z.blocks {
		return nil
	}
	z.blockSize = blockSize
	z.results = dedego.NewChan[result](int(blocks))
	z.blocks = blocks
	z.dstPool.New = func() interface{} { return make([]byte, 0, blockSize+(blockSize)>>4) }
	return nil
}

func NewWriter(w io.Writer) *Writer {

	z, _ := NewWriterLevel(w, DefaultCompression)
	return z
}

func NewWriterLevel(w io.Writer, level int) (*Writer, error) {

	if level < ConstantCompression || level > BestCompression {
		return nil, fmt.Errorf("gzip: invalid compression level: %d", level)
	}
	z := new(Writer)

	z.SetConcurrency(defaultBlockSize, runtime.GOMAXPROCS(0))
	z.init(w, level)
	return z, nil
}

func (z *Writer) pushError(err error) {

	z.errMu.Lock()
	if z.err != nil {
		z.errMu.Unlock()
		return
	}
	z.err = err
	z.pushedErr.Close()

	z.errMu.Unlock()

}

func (z *Writer) init(w io.Writer, level int) {

	z.wg.Wait()
	digest := z.digest
	if digest != nil {
		digest.Reset()
	} else {
		digest = crc32.NewIEEE()
	}
	z.Header = Header{OS: 255}
	z.w = w
	z.level = level
	z.digest = digest
	z.pushedErr = dedego.NewChan[struct{}](int(0))
	z.results = dedego.NewChan[result](int(0))
	z.err = nil
	z.closed = false
	z.Comment = ""
	z.Extra = nil
	z.ModTime = time.Time{}
	z.wroteHeader = false
	z.currentBuffer = nil
	z.buf = [10]byte{}
	z.prevTail = nil
	z.size = 0
	if z.dictFlatePool.New == nil {
		z.dictFlatePool.New = func() interface{} {
			f, _ := flate.NewWriterDict(w, level, nil)
			return f
		}
	}

}

func (z *Writer) Reset(w io.Writer) {

	if !z.results.IsNil() && !z.closed {
		z.results.Close()
	}
	z.SetConcurrency(defaultBlockSize, runtime.GOMAXPROCS(0))
	z.init(w, z.level)
}

func put2(p []byte, v uint16) {

	p[0] = uint8(v >> 0)
	p[1] = uint8(v >> 8)
}

func put4(p []byte, v uint32) {

	p[0] = uint8(v >> 0)
	p[1] = uint8(v >> 8)
	p[2] = uint8(v >> 16)
	p[3] = uint8(v >> 24)
}

func (z *Writer) writeBytes(b []byte) error {

	if len(b) > 0xffff {

		return errors.New("gzip.Write: Extra data is too large")
	}
	put2(z.buf[0:2], uint16(len(b)))
	_, err := z.w.Write(z.buf[0:2])
	if err != nil {
		return err
	}
	_, err = z.w.Write(b)
	return err
}

func (z *Writer) writeString(s string) (err error) {

	needconv := false
	for _, v := range s {
		if v == 0 || v > 0xff {
			return errors.New("gzip.Write: non-Latin-1 header string")
		}
		if v > 0x7f {
			needconv = true
		}
	}
	if needconv {
		b := make([]byte, 0, len(s))
		for _, v := range s {
			b = append(b, byte(v))
		}
		_, err = z.w.Write(b)
	} else {
		_, err = io.WriteString(z.w, s)
	}
	if err != nil {

		return err
	}

	z.buf[0] = 0
	_, err = z.w.Write(z.buf[0:1])

	return err
}

func (z *Writer) compressCurrent(flush bool) {

	c := z.currentBuffer
	if len(c) > z.blockSize {
		panic("len(z.currentBuffer) > z.blockSize (most likely due to concurrent Write race)")
	}

	r := result{}
	r.result = dedego.NewChan[[]byte](int(1))
	r.notifyWritten = dedego.NewChan[struct{}](int(0))
	{
		dedego.PreSelect(false, z.results.GetIdPre(false), z.pushedErr.GetIdPre(true))
		selectCaseDedego_5 := dedego.BuildMessage(r)
		switch dedegoFetchOrder[3] {
		case 0:
			select {

			case z.results.GetChan() <- selectCaseDedego_5:
				z.results.Post(false, selectCaseDedego_5)
			case <-time.After(1 * time.Second):
				select {
				case z.results.GetChan() <- selectCaseDedego_5:
					z.results.Post(false, selectCaseDedego_5)
				case selectCaseDedego_6 := <-z.pushedErr.GetChan():
					z.pushedErr.Post(true, selectCaseDedego_6)
					return
				}
			}
		case 1:
			select {
			case selectCaseDedego_6 := <-z.pushedErr.GetChan():
				z.pushedErr.Post(true, selectCaseDedego_6)
				return
			case <-time.After(1 * time.Second):
				select {
				case z.results.GetChan() <- selectCaseDedego_5:
					z.results.Post(false, selectCaseDedego_5)
				case selectCaseDedego_6 := <-z.pushedErr.GetChan():
					z.pushedErr.Post(true, selectCaseDedego_6)
					return
				}
			}
		default:
			select {
			case z.results.GetChan() <- selectCaseDedego_5:
				z.results.Post(false, selectCaseDedego_5)
			case selectCaseDedego_6 := <-z.pushedErr.GetChan():
				z.pushedErr.Post(true, selectCaseDedego_6)
				return
			}
		}
	}

	z.wg.Add(1)
	tail := z.prevTail
	if len(c) > tailSize {
		buf := z.dstPool.Get().([]byte)

		buf = append(buf[:0], c[len(c)-tailSize:]...)
		z.prevTail = buf
	} else {
		z.prevTail = nil
	}
	func() {
		DedegoRoutineIndex := dedego.SpawnPre()
		go func() {
			dedego.SpawnPost(DedegoRoutineIndex)
			{
				z.compressBlock(c, tail, r, z.closed)
			}
		}()
	}()

	z.currentBuffer = z.dstPool.Get().([]byte)
	z.currentBuffer = z.currentBuffer[:0]

	if flush {
		r.notifyWritten.Receive()
	}

}

func (z *Writer) checkError() error {

	z.errMu.RLock()
	err := z.err
	z.errMu.RUnlock()
	return err
}

func (z *Writer) Write(p []byte) (int, error) {

	if err := z.checkError(); err != nil {
		return 0, err
	}

	if !z.wroteHeader {
		z.wroteHeader = true
		z.buf[0] = gzipID1
		z.buf[1] = gzipID2
		z.buf[2] = gzipDeflate
		z.buf[3] = 0
		if z.Extra != nil {
			z.buf[3] |= 0x04
		}
		if z.Name != "" {
			z.buf[3] |= 0x08
		}
		if z.Comment != "" {
			z.buf[3] |= 0x10
		}
		put4(z.buf[4:8], uint32(z.ModTime.Unix()))
		if z.level == BestCompression {
			z.buf[8] = 2
		} else if z.level == BestSpeed {
			z.buf[8] = 4
		} else {
			z.buf[8] = 0
		}
		z.buf[9] = z.OS
		var n int
		var err error
		n, err = z.w.Write(z.buf[0:10])
		if err != nil {
			z.pushError(err)

			return n, err
		}
		if z.Extra != nil {
			err = z.writeBytes(z.Extra)
			if err != nil {
				z.pushError(err)

				return n, err
			}
		}
		if z.Name != "" {
			err = z.writeString(z.Name)
			if err != nil {
				z.pushError(err)

				return n, err
			}
		}
		if z.Comment != "" {
			err = z.writeString(z.Comment)
			if err != nil {
				z.pushError(err)

				return n, err
			}
		}
		func() {
			DedegoRoutineIndex := dedego.SpawnPre()
			go func() {
				dedego.SpawnPost(DedegoRoutineIndex)
				{
					listen := z.results
					var failed bool
					for {
						r, ok := listen.ReceiveOk()

						if !ok {
							return
						}
						if failed {
							r.notifyWritten.Close()

							continue
						}
						buf := r.result.Receive()

						n, err := z.w.Write(buf)
						if err != nil {
							z.pushError(err)
							r.notifyWritten.Close()

							failed = true
							continue
						}
						if n != len(buf) {
							z.pushError(fmt.Errorf("gzip: short write %d should be %d", n, len(buf)))
							failed = true
							r.notifyWritten.Close()

							continue
						}
						z.dstPool.Put(buf)
						r.notifyWritten.Close()

					}
				}
			}()
		}()

		z.currentBuffer = z.dstPool.Get().([]byte)
		z.currentBuffer = z.currentBuffer[:0]
	}
	q := p
	for len(q) > 0 {
		length := len(q)
		if length+len(z.currentBuffer) > z.blockSize {
			length = z.blockSize - len(z.currentBuffer)
		}
		z.digest.Write(q[:length])
		z.currentBuffer = append(z.currentBuffer, q[:length]...)
		if len(z.currentBuffer) > z.blockSize {
			panic("z.currentBuffer too large (most likely due to concurrent Write race)")
		}
		if len(z.currentBuffer) == z.blockSize {
			z.compressCurrent(false)
			if err := z.checkError(); err != nil {
				return len(p) - len(q), err
			}
		}
		z.size += length
		q = q[length:]
	}

	return len(p), z.checkError()
}

func (z *Writer) compressBlock(p, prevTail []byte, r result, closed bool) {

	defer func() {
		r.result.Close()

		z.wg.Done()
	}()
	buf := z.dstPool.Get().([]byte)
	dest := bytes.NewBuffer(buf[:0])

	compressor := z.dictFlatePool.Get().(*flate.Writer)
	compressor.ResetDict(dest, prevTail)
	compressor.Write(p)
	z.dstPool.Put(p)

	err := compressor.Flush()
	if err != nil {
		z.pushError(err)

		return
	}
	if closed {
		err = compressor.Close()
		if err != nil {
			z.pushError(err)

			return
		}
	}
	z.dictFlatePool.Put(compressor)

	if prevTail != nil {
		z.dstPool.Put(prevTail)
	}

	buf = dest.Bytes()
	r.result.Send(buf)

}

func (z *Writer) Flush() error {

	if err := z.checkError(); err != nil {

		return err
	}
	if z.closed {

		return nil
	}
	if !z.wroteHeader {
		_, err := z.Write(nil)
		if err != nil {

			return err
		}
	}

	z.compressCurrent(true)

	return z.checkError()
}

func (z *Writer) UncompressedSize() int {

	return z.size
}

func (z *Writer) Close() error {

	if err := z.checkError(); err != nil {

		return err
	}
	if z.closed {

		return nil
	}

	z.closed = true
	if !z.wroteHeader {
		z.Write(nil)
		if err := z.checkError(); err != nil {

			return err
		}
	}
	z.compressCurrent(true)
	if err := z.checkError(); err != nil {

		return err
	}
	z.results.Close()

	put4(z.buf[0:4], z.digest.Sum32())
	put4(z.buf[4:8], uint32(z.size))
	_, err := z.w.Write(z.buf[0:8])
	if err != nil {
		z.pushError(err)

		return err
	}

	return nil
}
