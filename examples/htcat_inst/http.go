package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ErikKassubek/deadlockDetectorGo/src/dedego"
)

const (
	_        = iota
	kB int64 = 1 << (10 * iota)
	mB
	gB
	tB
	pB
	eB
)

type HtCat struct {
	io.WriterTo
	d  defrag
	u  *url.URL
	cl *http.Client

	httpFragGenMu dedego.Mutex
	hfg           httpFragGen
}

type HttpStatusError struct {
	error
	Status string
}

func (cat *HtCat) startup(parallelism int) {
	req := http.Request{
		Method:     "GET",
		URL:        cat.u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Body:       nil,
		Host:       cat.u.Host,
	}

	resp, err := cat.cl.Do(&req)
	if err != nil {
		func() {
			DedegoRoutineIndex := dedego.SpawnPre()
			go func() {
				dedego.SpawnPost(DedegoRoutineIndex)
				{
					cat.d.cancel(err)
				}
			}()
		}()
		return
	}

	if resp.StatusCode != 200 {
		err = HttpStatusError{
			error: fmt.Errorf(
				"Expected HTTP Status 200, received: %q",
				resp.Status),
			Status: resp.Status}
		func() {
			DedegoRoutineIndex := dedego.SpawnPre()
			go func() {
				dedego.SpawnPost(DedegoRoutineIndex)
				{
					cat.d.cancel(err)
				}
			}()
		}()
		return
	}

	l := resp.Header.Get("Content-Length")

	noParallel := func(wtc writerToCloser) {
		f := cat.d.nextFragment()
		cat.d.setLast(cat.d.lastAllocated())
		f.contents = wtc
		cat.d.register(f)
	}

	if l == "" {
		func() {
			DedegoRoutineIndex := dedego.SpawnPre()
			go func() {
				dedego.SpawnPost(DedegoRoutineIndex)
				{

					noParallel(struct {
						io.WriterTo
						io.Closer
					}{
						WriterTo: bufio.NewReader(resp.Body),
						Closer:   resp.Body,
					})
				}
			}()
		}()
		return
	}

	length, err := strconv.ParseInt(l, 10, 64)
	if err != nil {
		func() {
			DedegoRoutineIndex := dedego.SpawnPre()
			go func() {
				dedego.SpawnPost(DedegoRoutineIndex)
				{

					cat.d.cancel(err)
				}
			}()
		}()
		return
	}

	cat.hfg.totalSize = length
	cat.hfg.targetFragSize = 1 + ((length - 1) / int64(parallelism))
	if cat.hfg.targetFragSize > 20*mB {
		cat.hfg.targetFragSize = 20 * mB
	}

	if cat.hfg.targetFragSize < 1*mB {
		cat.hfg.curPos = cat.hfg.totalSize
		er := newEagerReader(resp.Body, cat.hfg.totalSize)
		func() {
			DedegoRoutineIndex := dedego.SpawnPre()
			go func() {
				dedego.SpawnPost(DedegoRoutineIndex)
				{
					noParallel(er)
				}
			}()
		}()
		func() {
			DedegoRoutineIndex := dedego.SpawnPre()
			go func() {
				dedego.SpawnPost(DedegoRoutineIndex)
				{
					er.WaitClosed()
				}
			}()
		}()
		return
	}

	hf := cat.nextFragment()
	func() {
		DedegoRoutineIndex := dedego.SpawnPre()
		go func() {
			dedego.SpawnPost(DedegoRoutineIndex)
			{
				er := newEagerReader(
					struct {
						io.Reader
						io.Closer
					}{
						Reader: io.LimitReader(resp.Body, hf.size),
						Closer: resp.Body,
					},
					hf.size)

				hf.fragment.contents = er
				cat.d.register(hf.fragment)
				er.WaitClosed()

				cat.get()
			}
		}()
	}()

}

func New(client *http.Client, u *url.URL, parallelism int) *HtCat {
	cat := HtCat{
		u:  u,
		cl: client,
	}

	cat.d.initDefrag()
	cat.WriterTo = &cat.d
	cat.startup(parallelism)

	if cat.hfg.curPos == cat.hfg.totalSize {
		return &cat
	}

	for i := 1; i < parallelism; i += 1 {
		func() {
			DedegoRoutineIndex := dedego.SpawnPre()
			go func() {
				dedego.SpawnPost(DedegoRoutineIndex)
				{
					cat.get()
				}
			}()
		}()
	}

	return &cat
}

func (cat *HtCat) nextFragment() *httpFrag {
	cat.httpFragGenMu.Lock()
	defer cat.httpFragGenMu.Unlock()

	var hf *httpFrag

	if cat.hfg.hasNext() {
		f := cat.d.nextFragment()
		hf = cat.hfg.nextFragment(f)
	} else {
		cat.d.setLast(cat.d.lastAllocated())
	}

	return hf
}

func (cat *HtCat) get() {
	for {
		hf := cat.nextFragment()
		if hf == nil {
			return
		}

		req := http.Request{
			Method:     "GET",
			URL:        cat.u,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     hf.header,
			Body:       nil,
			Host:       cat.u.Host,
		}

		resp, err := cat.cl.Do(&req)
		if err != nil {
			cat.d.cancel(err)
			return
		}

		if !(resp.StatusCode == 206 || resp.StatusCode == 200) {
			err = HttpStatusError{
				error: fmt.Errorf("Expected HTTP Status "+
					"206 or 200, received: %q",
					resp.Status),
				Status: resp.Status}
			func() {
				DedegoRoutineIndex := dedego.SpawnPre()
				go func() {
					dedego.SpawnPost(DedegoRoutineIndex)
					{
						cat.d.cancel(err)
					}
				}()
			}()
			return
		}

		er := newEagerReader(resp.Body, hf.size)
		hf.fragment.contents = er
		cat.d.register(hf.fragment)
		er.WaitClosed()
	}
}
