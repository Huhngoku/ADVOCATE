package main

import (
	"io"
	"sync/atomic"
	"time"

	"github.com/ErikKassubek/deadlockDetectorGo/src/dedego"
)

type writerToCloser interface {
	io.WriterTo
	io.Closer
}

type fragment struct {
	ord      int64
	contents writerToCloser
}

type defrag struct {
	lastWritten int64

	lastAlloc int64

	lastOrdinal       int64
	lastOrdinalNotify dedego.Chan[int64]

	future map[int64]*fragment

	registerNotify dedego.Chan[*fragment]

	cancellation error
	cancelNotify dedego.Chan[error]

	written int64

	done dedego.Chan[struct{}]
}

func newDefrag() *defrag {
	ret := defrag{}
	ret.initDefrag()

	return &ret
}

func (d *defrag) initDefrag() {
	d.future = make(map[int64]*fragment)
	d.registerNotify = dedego.NewChan[*fragment](int(0))
	d.cancelNotify = dedego.NewChan[error](int(0))
	d.lastOrdinalNotify = dedego.NewChan[int64](int(0))
	d.done = dedego.NewChan[struct{}](int(0))
}

func (d *defrag) nextFragment() *fragment {
	atomic.AddInt64(&d.lastAlloc, 1)
	f := fragment{ord: d.lastAlloc}

	return &f
}

func (d *defrag) cancel(err error) {
	d.cancelNotify.Send(err)

}

func (d *defrag) WriteTo(dst io.Writer) (written int64, err error) {
	defer d.done.Close()

	if d.cancellation != nil {
		return d.written, d.cancellation
	}

	for {

		if d.lastWritten >= d.lastOrdinal && d.lastOrdinal > 0 {
			break
		}
		{
			dedego.PreSelect(false, d.registerNotify.GetIdPre(true), d.cancelNotify.GetIdPre(true), d.lastOrdinalNotify.GetIdPre(true))
			switch dedegoFetchOrder[1] {
			case 0:
				select {

				case selectCaseDedego_1 := <-d.registerNotify.GetChan():
					d.registerNotify.Post(true, selectCaseDedego_1)
					frag := selectCaseDedego_1.GetInfo()

					next := d.lastWritten + 1
					if frag.ord == next {

						n, err := d.writeConsecutive(dst, frag)
						d.written += n
						if err != nil {
							return d.written, err
						}
					} else if frag.ord > next {

						d.future[frag.ord] = frag
					} else {
						return d.written, assertErrf(
							"Unexpected retrograde fragment %v, "+
								"expected at least %v",
							frag.ord, next)
					}
				case <-time.After(2 * time.Second):
					select {
					case selectCaseDedego_1 := <-d.registerNotify.GetChan():
						d.registerNotify.Post(true, selectCaseDedego_1)
						frag := selectCaseDedego_1.GetInfo()

						next := d.lastWritten + 1
						if frag.ord == next {

							n, err := d.writeConsecutive(dst, frag)
							d.written += n
							if err != nil {
								return d.written, err
							}
						} else if frag.ord > next {

							d.future[frag.ord] = frag
						} else {
							return d.written, assertErrf(
								"Unexpected retrograde fragment %v, "+
									"expected at least %v",
								frag.ord, next)
						}

					case selectCaseDedego_2 := <-d.cancelNotify.GetChan():
						d.cancelNotify.Post(true, selectCaseDedego_2)
						d.cancellation = selectCaseDedego_2.GetInfo()

						d.future = nil
						return d.written, d.cancellation

					case selectCaseDedego_3 := <-d.lastOrdinalNotify.GetChan():
						d.lastOrdinalNotify.Post(true, selectCaseDedego_3)
						d.lastOrdinal = selectCaseDedego_3.GetInfo()

						continue
					}
				}
			case 1:
				select {
				case selectCaseDedego_2 := <-d.cancelNotify.GetChan():
					d.cancelNotify.Post(true, selectCaseDedego_2)
					d.cancellation = selectCaseDedego_2.GetInfo()

					d.future = nil
					return d.written, d.cancellation
				case <-time.After(2 * time.Second):
					select {
					case selectCaseDedego_1 := <-d.registerNotify.GetChan():
						d.registerNotify.Post(true, selectCaseDedego_1)
						frag := selectCaseDedego_1.GetInfo()

						next := d.lastWritten + 1
						if frag.ord == next {

							n, err := d.writeConsecutive(dst, frag)
							d.written += n
							if err != nil {
								return d.written, err
							}
						} else if frag.ord > next {

							d.future[frag.ord] = frag
						} else {
							return d.written, assertErrf(
								"Unexpected retrograde fragment %v, "+
									"expected at least %v",
								frag.ord, next)
						}

					case selectCaseDedego_2 := <-d.cancelNotify.GetChan():
						d.cancelNotify.Post(true, selectCaseDedego_2)
						d.cancellation = selectCaseDedego_2.GetInfo()

						d.future = nil
						return d.written, d.cancellation

					case selectCaseDedego_3 := <-d.lastOrdinalNotify.GetChan():
						d.lastOrdinalNotify.Post(true, selectCaseDedego_3)
						d.lastOrdinal = selectCaseDedego_3.GetInfo()

						continue
					}
				}
			case 2:
				select {
				case selectCaseDedego_3 := <-d.lastOrdinalNotify.GetChan():
					d.lastOrdinalNotify.Post(true, selectCaseDedego_3)
					d.lastOrdinal = selectCaseDedego_3.GetInfo()

					continue
				case <-time.After(2 * time.Second):
					select {
					case selectCaseDedego_1 := <-d.registerNotify.GetChan():
						d.registerNotify.Post(true, selectCaseDedego_1)
						frag := selectCaseDedego_1.GetInfo()

						next := d.lastWritten + 1
						if frag.ord == next {

							n, err := d.writeConsecutive(dst, frag)
							d.written += n
							if err != nil {
								return d.written, err
							}
						} else if frag.ord > next {

							d.future[frag.ord] = frag
						} else {
							return d.written, assertErrf(
								"Unexpected retrograde fragment %v, "+
									"expected at least %v",
								frag.ord, next)
						}

					case selectCaseDedego_2 := <-d.cancelNotify.GetChan():
						d.cancelNotify.Post(true, selectCaseDedego_2)
						d.cancellation = selectCaseDedego_2.GetInfo()

						d.future = nil
						return d.written, d.cancellation

					case selectCaseDedego_3 := <-d.lastOrdinalNotify.GetChan():
						d.lastOrdinalNotify.Post(true, selectCaseDedego_3)
						d.lastOrdinal = selectCaseDedego_3.GetInfo()

						continue
					}
				}
			default:
				select {
				case selectCaseDedego_1 := <-d.registerNotify.GetChan():
					d.registerNotify.Post(true, selectCaseDedego_1)
					frag := selectCaseDedego_1.GetInfo()

					next := d.lastWritten + 1
					if frag.ord == next {

						n, err := d.writeConsecutive(dst, frag)
						d.written += n
						if err != nil {
							return d.written, err
						}
					} else if frag.ord > next {

						d.future[frag.ord] = frag
					} else {
						return d.written, assertErrf(
							"Unexpected retrograde fragment %v, "+
								"expected at least %v",
							frag.ord, next)
					}

				case selectCaseDedego_2 := <-d.cancelNotify.GetChan():
					d.cancelNotify.Post(true, selectCaseDedego_2)
					d.cancellation = selectCaseDedego_2.GetInfo()

					d.future = nil
					return d.written, d.cancellation

				case selectCaseDedego_3 := <-d.lastOrdinalNotify.GetChan():
					d.lastOrdinalNotify.Post(true, selectCaseDedego_3)
					d.lastOrdinal = selectCaseDedego_3.GetInfo()

					continue
				}
			}
		}
	}

	return d.written, nil
}

func (d *defrag) setLast(lastOrdinal int64) {
	{
		dedego.PreSelect(false, d.lastOrdinalNotify.GetIdPre(false), d.done.GetIdPre(true))
		selectCaseDedego_4 := dedego.BuildMessage(lastOrdinal)
		switch dedegoFetchOrder[2] {
		case 0:
			select {

			case d.lastOrdinalNotify.GetChan() <- selectCaseDedego_4:
				d.lastOrdinalNotify.Post(false, selectCaseDedego_4)
			case <-time.After(2 * time.Second):
				select {
				case d.lastOrdinalNotify.GetChan() <- selectCaseDedego_4:
					d.lastOrdinalNotify.Post(false, selectCaseDedego_4)
				case selectCaseDedego_5 := <-d.done.GetChan():
					d.done.Post(true, selectCaseDedego_5)
				}
			}
		case 1:
			select {
			case selectCaseDedego_5 := <-d.done.GetChan():
				d.done.Post(true, selectCaseDedego_5)
			case <-time.After(2 * time.Second):
				select {
				case d.lastOrdinalNotify.GetChan() <- selectCaseDedego_4:
					d.lastOrdinalNotify.Post(false, selectCaseDedego_4)
				case selectCaseDedego_5 := <-d.done.GetChan():
					d.done.Post(true, selectCaseDedego_5)
				}
			}
		default:
			select {
			case d.lastOrdinalNotify.GetChan() <- selectCaseDedego_4:
				d.lastOrdinalNotify.Post(false, selectCaseDedego_4)
			case selectCaseDedego_5 := <-d.done.GetChan():
				d.done.Post(true, selectCaseDedego_5)
			}
		}
	}
}

func (d *defrag) lastAllocated() int64 {
	return atomic.LoadInt64(&d.lastAlloc)
}

func (d *defrag) register(frag *fragment) {
	d.registerNotify.Send(frag)

}

func (d *defrag) writeConsecutive(dst io.Writer, start *fragment) (
	int64, error) {

	written, err := start.contents.WriteTo(dst)
	if err != nil {
		return int64(written), err
	}

	if err := start.contents.Close(); err != nil {
		return int64(written), err
	}

	d.lastWritten += 1

	for {
		{
			dedego.PreSelect(true, d.cancelNotify.GetIdPre(true))
			switch dedegoFetchOrder[3] {
			case 0:
				select {

				case selectCaseDedego_6 := <-d.cancelNotify.GetChan():
					d.cancelNotify.Post(true, selectCaseDedego_6)
					d.cancellation = selectCaseDedego_6.GetInfo()
					d.future = nil
					return 0, d.cancellation
				case <-time.After(2 * time.Second):
					select {
					case selectCaseDedego_6 := <-d.cancelNotify.GetChan():
						d.cancelNotify.Post(true, selectCaseDedego_6)
						d.cancellation = selectCaseDedego_6.GetInfo()
						d.future = nil
						return 0, d.cancellation
					default:
						dedego.PostDefault()
					}
				}
			case 1:
				select {
				default:
					dedego.PostDefault()
				case <-time.After(2 * time.Second):
					select {
					case selectCaseDedego_6 := <-d.cancelNotify.GetChan():
						d.cancelNotify.Post(true, selectCaseDedego_6)
						d.cancellation = selectCaseDedego_6.GetInfo()
						d.future = nil
						return 0, d.cancellation
					default:
						dedego.PostDefault()
					}
				}
			default:
				select {
				case selectCaseDedego_6 := <-d.cancelNotify.GetChan():
					d.cancelNotify.Post(true, selectCaseDedego_6)
					d.cancellation = selectCaseDedego_6.GetInfo()
					d.future = nil
					return 0, d.cancellation
				default:
					dedego.PostDefault()
				}
			}
		}

		next := d.lastWritten + 1
		if frag, ok := d.future[next]; ok {

			delete(d.future, next)
			n, err := frag.contents.WriteTo(dst)
			written += n
			defer frag.contents.Close()
			if err != nil {
				return int64(written), err
			}

			d.lastWritten = next
		} else {

			return int64(written), nil
		}
	}
}
