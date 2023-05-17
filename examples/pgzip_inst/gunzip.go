// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"errors"
	"hash"
	"hash/crc32"
	"io"
	"sync"
	"time"

	"github.com/ErikKassubek/deadlockDetectorGo/src/dedego"
	"github.com/klauspost/compress/flate"
)

const (
	gzipID1     = 0x1f
	gzipID2     = 0x8b
	gzipDeflate = 8
	flagText    = 1 << 0
	flagHdrCrc  = 1 << 1
	flagExtra   = 1 << 2
	flagName    = 1 << 3
	flagComment = 1 << 4
)

func makeReader(r io.Reader) flate.Reader {

	if rr, ok := r.(flate.Reader); ok {
		return rr
	}
	return bufio.NewReader(r)
}

var (
	ErrChecksum = errors.New("gzip: invalid checksum")

	ErrHeader = errors.New("gzip: invalid header")
)

type Header struct {
	Comment string
	Extra   []byte
	ModTime time.Time
	Name    string
	OS      byte
}

type Reader struct {
	Header
	r            flate.Reader
	decompressor io.ReadCloser
	digest       hash.Hash32
	size         uint32
	flg          byte
	buf          [512]byte
	err          error
	closeErr     dedego.Chan[error]
	multistream  bool

	readAhead   dedego.Chan[read]
	roff        int
	current     []byte
	closeReader dedego.Chan[struct{}]
	lastBlock   bool
	blockSize   int
	blocks      int

	activeRA bool
	mu       dedego.Mutex

	blockPool dedego.Chan[[]byte]
}

type read struct {
	b   []byte
	err error
}

func NewReader(r io.Reader) (*Reader, error) {

	z := new(Reader)
	z.blocks = defaultBlocks
	z.blockSize = defaultBlockSize
	z.r = makeReader(r)
	z.digest = crc32.NewIEEE()
	z.multistream = true
	z.blockPool = dedego.NewChan[[]byte](z.blocks)
	z.readAhead.SetNil()
	z.closeErr.SetNil()
	z.closeReader.SetNil()
	for i := 0; i < z.blocks; i++ {
		z.blockPool.Send(make([]byte, z.blockSize))
	}
	if err := z.readHeader(true); err != nil {

		return nil, err
	}

	return z, nil
}

func NewReaderN(r io.Reader, blockSize, blocks int) (*Reader, error) {

	z := new(Reader)
	z.blocks = blocks
	z.blockSize = blockSize
	z.r = makeReader(r)
	z.digest = crc32.NewIEEE()
	z.multistream = true

	if z.blocks <= 0 {
		z.blocks = defaultBlocks
	}
	if z.blockSize <= 512 {
		z.blockSize = defaultBlockSize
	}
	z.blockPool = dedego.NewChan[[]byte](int(0))
	for i := 0; i < z.blocks; i++ {
		z.blockPool.Send(make([]byte, z.blockSize))
	}
	if err := z.readHeader(true); err != nil {

		return nil, err
	}

	return z, nil
}

func (z *Reader) Reset(r io.Reader) error {

	z.killReadAhead()
	z.r = makeReader(r)
	z.digest = crc32.NewIEEE()
	z.size = 0
	z.err = nil
	z.multistream = true

	if z.blocks <= 0 {
		z.blocks = defaultBlocks
	}
	if z.blockSize <= 512 {
		z.blockSize = defaultBlockSize
	}

	if z.blockPool.GetChan() == nil {
		z.blockPool = dedego.NewChan[[]byte](int(0))
		for i := 0; i < z.blocks; i++ {
			z.blockPool.Send(make([]byte, z.blockSize))
		}
	}

	return z.readHeader(true)
}

func (z *Reader) Multistream(ok bool) {

	z.multistream = ok
}

func get4(p []byte) uint32 {

	return uint32(p[0]) | uint32(p[1])<<8 | uint32(p[2])<<16 | uint32(p[3])<<24
}

func (z *Reader) readString() (string, error) {

	var err error
	needconv := false
	for i := 0; ; i++ {
		if i >= len(z.buf) {
			return "", ErrHeader
		}
		z.buf[i], err = z.r.ReadByte()
		if err != nil {
			return "", err
		}
		if z.buf[i] > 0x7f {
			needconv = true
		}
		if z.buf[i] == 0 {
			if needconv {
				s := make([]rune, 0, i)
				for _, v := range z.buf[0:i] {
					s = append(s, rune(v))
				}

				return string(s), nil
			}

			return string(z.buf[0:i]), nil
		}
	}
}

func (z *Reader) read2() (uint32, error) {

	_, err := io.ReadFull(z.r, z.buf[0:2])
	if err != nil {

		return 0, err
	}

	return uint32(z.buf[0]) | uint32(z.buf[1])<<8, nil
}

func (z *Reader) readHeader(save bool) error {

	z.killReadAhead()

	_, err := io.ReadFull(z.r, z.buf[0:10])
	if err != nil {
		return err
	}
	if z.buf[0] != gzipID1 || z.buf[1] != gzipID2 || z.buf[2] != gzipDeflate {
		return ErrHeader
	}
	z.flg = z.buf[3]
	if save {
		z.ModTime = time.Unix(int64(get4(z.buf[4:8])), 0)
		z.OS = z.buf[9]
	}
	z.digest.Reset()
	z.digest.Write(z.buf[0:10])

	if z.flg&flagExtra != 0 {
		n, err := z.read2()
		if err != nil {
			return err
		}
		data := make([]byte, n)
		if _, err = io.ReadFull(z.r, data); err != nil {
			return err
		}
		if save {
			z.Extra = data
		}
	}

	var s string
	if z.flg&flagName != 0 {
		if s, err = z.readString(); err != nil {
			return err
		}
		if save {
			z.Name = s
		}
	}

	if z.flg&flagComment != 0 {
		if s, err = z.readString(); err != nil {
			return err
		}
		if save {
			z.Comment = s
		}
	}

	if z.flg&flagHdrCrc != 0 {
		n, err := z.read2()
		if err != nil {
			return err
		}
		sum := z.digest.Sum32() & 0xFFFF
		if n != sum {
			return ErrHeader
		}
	}

	z.digest.Reset()
	z.decompressor = flate.NewReader(z.r)
	z.doReadAhead()

	return nil
}

func (z *Reader) killReadAhead() error {

	z.mu.Lock()
	defer z.mu.Unlock()
	if z.activeRA {
		if !z.closeReader.IsNil() {
			z.closeReader.Close()
		}

		e, ok := z.closeErr.ReceiveOk()

		z.activeRA = false

		select_break := false
		i := 0
		for {
			i++

			if select_break {
				break
			}
			select {
			case b := <-z.readAhead.GetChan():
				z.readAhead.Post(true, b)
				blk := b.GetInfo()
				if blk.b != nil {
					z.blockPool.Send(blk.b)
				}
			default:
			}
			if z.readAhead.IsClosed() {
				select_break = true
			}

		}

		if cap(z.current) > 0 {
			z.blockPool.Send(z.current)

			z.current = nil
		}
		if !ok {
			return nil
		}

		return e
	}

	return nil
}

func (z *Reader) doReadAhead() {

	z.mu.Lock()
	defer z.mu.Unlock()
	z.activeRA = true

	if z.blocks <= 0 {
		z.blocks = defaultBlocks
	}
	if z.blockSize <= 512 {
		z.blockSize = defaultBlockSize
	}
	ra := dedego.NewChan[read](int(0))
	z.readAhead = ra
	closeReader := dedego.NewChan[struct{}](int(0))
	z.closeReader = closeReader
	z.lastBlock = false
	closeErr := dedego.NewChan[error](int(1)) // BUG
	z.closeErr = closeErr
	z.size = 0
	z.roff = 0
	z.current = nil
	decomp := z.decompressor
	func() {

		DedegoRoutineIndex := dedego.SpawnPre()
		go func() {
			dedego.SpawnPost(DedegoRoutineIndex)
			{
				defer func() {
					closeErr.Send(decomp.Close())
					closeErr.Close()
					// ra.Close()

					z.readAhead.Close()

				}()

				digest := z.digest
				var wg sync.WaitGroup
				for {
					var buf []byte
					{
						dedego.PreSelect(false, z.blockPool.GetIdPre(true), closeReader.GetIdPre(true))
						switch dedegoFetchOrder[1] {
						case 0:
							select {

							case selectCaseDedego_1 := <-z.blockPool.GetChan():
								z.blockPool.Post(true, selectCaseDedego_1)
								buf = selectCaseDedego_1.GetInfo()
							case <-time.After(1 * time.Second):
								select {
								case selectCaseDedego_1 := <-z.blockPool.GetChan():
									z.blockPool.Post(true, selectCaseDedego_1)
									buf = selectCaseDedego_1.GetInfo()

								case selectCaseDedego_2 := <-closeReader.GetChan():
									closeReader.Post(true, selectCaseDedego_2)

									return
								}
							}
						case 1:
							select {
							case selectCaseDedego_2 := <-closeReader.GetChan():
								closeReader.Post(true, selectCaseDedego_2)

								return
							case <-time.After(1 * time.Second):
								select {
								case selectCaseDedego_1 := <-z.blockPool.GetChan():
									z.blockPool.Post(true, selectCaseDedego_1)

									buf = selectCaseDedego_1.GetInfo()
								case selectCaseDedego_2 := <-closeReader.GetChan():
									closeReader.Post(true, selectCaseDedego_2)
									return
								}
							}
						default:
							select {
							case selectCaseDedego_1 := <-z.blockPool.GetChan():
								z.blockPool.Post(true, selectCaseDedego_1)

								buf = selectCaseDedego_1.GetInfo()
							case selectCaseDedego_2 := <-closeReader.GetChan():

								closeReader.Post(true, selectCaseDedego_2)
								return
							}
						}
					}

					buf = buf[0:z.blockSize]

					n, err := io.ReadFull(decomp, buf)
					if err == io.ErrUnexpectedEOF {
						if n > 0 {
							err = nil
						} else {

							_, err = decomp.Read([]byte{})
							if err == io.EOF {
								err = nil
							}
						}
					}
					if n < len(buf) {
						buf = buf[0:n]
					}
					wg.Wait()
					wg.Add(1)
					func() {
						DedegoRoutineIndex := dedego.SpawnPre()
						go func() {
							dedego.SpawnPost(DedegoRoutineIndex)
							{
								digest.Write(buf)
								wg.Done()
							}
						}()
					}()

					z.size += uint32(n)

					if err != nil {
						wg.Wait()
					}
					{
						dedego.PreSelect(false, z.readAhead.GetIdPre(false), closeReader.GetIdPre(true))
						selectCaseDedego_3 := dedego.BuildMessage(read{b: buf, err: err}) // Had Bug
						switch dedegoFetchOrder[2] {
						case 0:
							select {

							case z.readAhead.GetChan() <- selectCaseDedego_3:
								z.readAhead.Post(false, selectCaseDedego_3)
							case <-time.After(1 * time.Second):
								select {
								case z.readAhead.GetChan() <- selectCaseDedego_3:
									z.readAhead.Post(false, selectCaseDedego_3)
								case selectCaseDedego_4 := <-closeReader.GetChan():
									closeReader.Post(true, selectCaseDedego_4)
									z.blockPool.Send(buf)

									return
								}
							}
						case 1:
							select {
							case selectCaseDedego_4 := <-closeReader.GetChan():
								closeReader.Post(true, selectCaseDedego_4)
								z.blockPool.Send(buf)

								return
							case <-time.After(1 * time.Second):
								select {
								case z.readAhead.GetChan() <- selectCaseDedego_3:
									z.readAhead.Post(false, selectCaseDedego_3)
								case selectCaseDedego_4 := <-closeReader.GetChan():
									closeReader.Post(true, selectCaseDedego_4)
									z.blockPool.Send(buf)

									return
								}
							}
						default:
							select {
							case z.readAhead.GetChan() <- selectCaseDedego_3:
								z.readAhead.Post(false, selectCaseDedego_3)
							case selectCaseDedego_4 := <-closeReader.GetChan():
								closeReader.Post(true, selectCaseDedego_4)
								z.blockPool.Send(buf)

								return
							}
						}
					}
					if err != nil {

						return
					}
				}
			}
		}()
	}()

}

func (z *Reader) Read(p []byte) (n int, err error) {

	if z.err != nil {
		return 0, z.err
	}
	if len(p) == 0 {

		return 0, nil
	}

	for {
		if len(z.current) == 0 && !z.lastBlock {

			read := z.readAhead.Receive()

			if read.err != nil {
				z.closeReader.SetNil()
				if read.err != io.EOF {
					z.err = read.err
					return
				}
				if read.err == io.EOF {
					z.lastBlock = true
					err = nil
				}
			}
			z.current = read.b
			z.roff = 0
		}

		avail := z.current[z.roff:]
		if len(p) >= len(avail) {

			n = copy(p, avail)
			z.blockPool.Send(z.current)

			z.current = nil
			if z.lastBlock {
				err = io.EOF
				break
			}
		} else {

			n = copy(p, avail)
			z.roff += n
		}

		return
	}

	if _, err := io.ReadFull(z.r, z.buf[0:8]); err != nil {
		z.err = err

		return 0, err
	}
	crc32, isize := get4(z.buf[0:4]), get4(z.buf[4:8])
	sum := z.digest.Sum32()
	if sum != crc32 || isize != z.size {
		z.err = ErrChecksum

		return 0, z.err
	}

	if !z.multistream {

		return 0, io.EOF
	}

	if err = z.readHeader(false); err != nil {
		z.err = err

		return
	}

	return z.Read(p)
}

func (z *Reader) WriteTo(w io.Writer) (n int64, err error) {

	total := int64(0)
	avail := z.current[z.roff:]
	if len(avail) != 0 {
		n, err := w.Write(avail)
		if n != len(avail) {
			return total, io.ErrShortWrite
		}
		total += int64(n)
		if err != nil {
			return total, err
		}
		z.blockPool.Send(z.current)

		z.current = nil
	}
	for {
		if z.err != nil {
			return total, z.err
		}

		for {
			read := z.readAhead.Receive()

			if read.err != nil {

				z.closeReader.SetNil()

				if read.err != io.EOF {
					z.err = read.err
					return total, z.err
				}
				if read.err == io.EOF {
					z.lastBlock = true
					err = nil
				}
			}

			n, err := w.Write(read.b)
			if n != len(read.b) {
				return total, io.ErrShortWrite
			}
			total += int64(n)
			if err != nil {
				return total, err
			}
			z.blockPool.Send(read.b)

			if z.lastBlock {
				break
			}
		}

		if _, err := io.ReadFull(z.r, z.buf[0:8]); err != nil {
			z.err = err
			return total, err
		}
		crc32, isize := get4(z.buf[0:4]), get4(z.buf[4:8])
		sum := z.digest.Sum32()
		if sum != crc32 || isize != z.size {
			z.err = ErrChecksum

			return total, z.err
		}

		if !z.multistream {

			return total, nil
		}

		err = z.readHeader(false)
		if err == io.EOF {

			return total, nil
		}
		if err != nil {
			z.err = err

			return total, err
		}
	}
}

func (z *Reader) Close() error {

	return z.killReadAhead()
}
