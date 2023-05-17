package main

import (
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/ErikKassubek/deadlockDetectorGo/src/dedego"
)

type workerResponse struct {
	ctx           *URLContext
	visited       bool
	harvestedURLs interface{}
	host          string
	idleDeath     bool
}

type Crawler struct {
	Options *Options

	logFunc         func(LogFlags, string, ...interface{})
	push            dedego.Chan[*workerResponse]
	enqueue         dedego.Chan[interface{}]
	stop            dedego.Chan[struct{}]
	wg              *sync.WaitGroup
	pushPopRefCount int
	visits          int

	visited map[string]struct{}
	hosts   map[string]struct{}
	workers map[string]*worker
}

func NewCrawlerWithOptions(opts *Options) *Crawler {
	ret := new(Crawler)
	ret.Options = opts
	return ret
}

func NewCrawler(ext Extender) *Crawler {
	return NewCrawlerWithOptions(NewOptions(ext))
}

func (c *Crawler) Run(seeds interface{}) error {

	c.logFunc = getLogFunc(c.Options.Extender, c.Options.LogFlags, -1)

	seeds = c.Options.Extender.Start(seeds)
	ctxs := c.toURLContexts(seeds, nil)
	c.init(ctxs)

	c.enqueueUrls(ctxs)
	err := c.collectUrls()

	c.Options.Extender.End(err)
	return err
}

func (c *Crawler) init(ctxs []*URLContext) {

	c.hosts = make(map[string]struct{}, len(ctxs))
	for _, ctx := range ctxs {

		if _, ok := c.hosts[ctx.normalizedURL.Host]; !ok {
			c.hosts[ctx.normalizedURL.Host] = struct{}{}
		}
	}

	hostCount := len(c.hosts)
	l := len(ctxs)
	c.logFunc(LogTrace, "init() - seeds length: %d", l)
	c.logFunc(LogTrace, "init() - host count: %d", hostCount)
	c.logFunc(LogInfo, "robot user-agent: %s", c.Options.RobotUserAgent)

	c.wg = new(sync.WaitGroup)

	c.visited = make(map[string]struct{}, l)
	c.pushPopRefCount, c.visits = 0, 0

	c.stop = dedego.NewChan[struct{}](int(0))
	if c.Options.SameHostOnly {
		c.workers, c.push = make(map[string]*worker, hostCount),
			dedego.NewChan[*workerResponse](int(hostCount))
	} else {
		c.workers, c.push = make(map[string]*worker, c.Options.HostBufferFactor*hostCount),
			dedego.NewChan[*workerResponse](int(0))
	}

	c.enqueue = dedego.NewChan[interface{}](int(0))
	c.setExtenderEnqueueChan()
}

func (c *Crawler) setExtenderEnqueueChan() {
	defer func() {
		if err := recover(); err != nil {

			c.logFunc(LogError, "cannot set the enqueue channel: %s", err)
		}
	}()

	v := reflect.ValueOf(c.Options.Extender)
	el := v.Elem()
	if el.Kind() != reflect.Struct {
		c.logFunc(LogInfo, "extender is not a struct, cannot set the enqueue channel")
		return
	}
	ec := el.FieldByName("EnqueueChan")
	if !ec.IsValid() {
		c.logFunc(LogInfo, "extender.EnqueueChan does not exist, cannot set the enqueue channel")
		return
	}
	t := ec.Type()
	if t.Kind() != reflect.Chan || t.ChanDir() != reflect.SendDir {
		c.logFunc(LogInfo, "extender.EnqueueChan is not of type chan<-interface{}, cannot set the enqueue channel")
		return
	}
	tt := t.Elem()
	if tt.Kind() != reflect.Interface || tt.NumMethod() != 0 {
		c.logFunc(LogInfo, "extender.EnqueueChan is not of type chan<-interface{}, cannot set the enqueue channel")
		return
	}
	src := reflect.ValueOf(c.enqueue)
	ec.Set(src)
}

func (c *Crawler) launchWorker(ctx *URLContext) *worker {

	i := len(c.workers) + 1
	pop := newPopChannel()

	w := &worker{
		host:    ctx.normalizedURL.Host,
		index:   i,
		push:    c.push,
		pop:     pop,
		stop:    c.stop,
		enqueue: c.enqueue,
		wg:      c.wg,
		logFunc: getLogFunc(c.Options.Extender, c.Options.LogFlags, i),
		opts:    c.Options,
	}

	c.wg.Add(1)
	func() {
		DedegoRoutineIndex := dedego.SpawnPre()
		go func() {
			dedego.SpawnPost(DedegoRoutineIndex)
			{

				w.run()
			}
		}()
	}()
	c.logFunc(LogInfo, "worker %d launched for host %s", i, w.host)
	c.workers[w.host] = w

	return w
}

func (c *Crawler) isSameHost(ctx *URLContext) bool {

	if ctx.normalizedSourceURL != nil {
		return ctx.normalizedURL.Host == ctx.normalizedSourceURL.Host
	}

	_, ok := c.hosts[ctx.normalizedURL.Host]
	return ok
}

func (c *Crawler) enqueueUrls(ctxs []*URLContext) (cnt int) {
	for _, ctx := range ctxs {
		var isVisited, enqueue bool

		if ctx.IsRobotsURL() {
			continue
		}

		_, isVisited = c.visited[ctx.normalizedURL.String()]

		if enqueue = c.Options.Extender.Filter(ctx, isVisited); !enqueue {

			c.logFunc(LogIgnored, "ignore on filter policy: %s", ctx.normalizedURL)
			continue
		}

		if !ctx.normalizedURL.IsAbs() {

			c.logFunc(LogIgnored, "ignore on absolute policy: %s", ctx.normalizedURL)

		} else if !strings.HasPrefix(ctx.normalizedURL.Scheme, "http") {
			c.logFunc(LogIgnored, "ignore on scheme policy: %s", ctx.normalizedURL)

		} else if c.Options.SameHostOnly && !c.isSameHost(ctx) {

			c.logFunc(LogIgnored, "ignore on same host policy: %s", ctx.normalizedURL)

		} else {

			w, ok := c.workers[ctx.normalizedURL.Host]
			if !ok {

				w = c.launchWorker(ctx)

				if robCtx, e := ctx.getRobotsURLCtx(); e != nil {
					c.Options.Extender.Error(newCrawlError(ctx, e, CekParseRobots))
					c.logFunc(LogError, "ERROR parsing robots.txt from %s: %s", ctx.normalizedURL, e)
				} else {
					c.logFunc(LogEnqueued, "enqueue: %s", robCtx.url)
					c.Options.Extender.Enqueued(robCtx)
					w.pop.stack(robCtx)
				}
			}

			cnt++
			c.logFunc(LogEnqueued, "enqueue: %s", ctx.url)
			c.Options.Extender.Enqueued(ctx)
			w.pop.stack(ctx)
			c.pushPopRefCount++

			if !isVisited {

				c.visited[ctx.normalizedURL.String()] = struct{}{}
			}
		}
	}
	return
}

func (c *Crawler) collectUrls() error {
	defer func() {
		c.logFunc(LogInfo, "waiting for goroutines to complete...")
		c.wg.Wait()
		c.logFunc(LogInfo, "crawler done.")
	}()

	for {

		if c.pushPopRefCount == 0 && len(c.enqueue.GetChan()) == 0 {
			c.logFunc(LogInfo, "sending STOP signals...")
			c.stop.Close()

			return nil
		}
		{
			dedego.PreSelect(false, c.push.GetIdPre(true), c.enqueue.GetIdPre(true), c.stop.GetIdPre(true))
			switch dedegoFetchOrder[1] {
			case 0:
				select {

				case selectCaseDedego_1 := <-c.push.GetChan():
					c.push.Post(true, selectCaseDedego_1)
					res := selectCaseDedego_1.GetInfo()

					if res.visited {
						c.visits++
						if c.Options.MaxVisits > 0 && c.visits >= c.Options.MaxVisits {

							c.logFunc(LogInfo, "sending STOP signals...")
							c.stop.Close()

							return ErrMaxVisits
						}
					}
					if res.idleDeath {

						delete(c.workers, res.host)
						c.logFunc(LogInfo, "worker for host %s cleared on idle policy", res.host)
					} else {
						c.enqueueUrls(c.toURLContexts(res.harvestedURLs, res.ctx.url))
						c.pushPopRefCount--
					}
				case <-time.After(2 * time.Second):
					select {
					case selectCaseDedego_1 := <-c.push.GetChan():
						c.push.Post(true, selectCaseDedego_1)
						res := selectCaseDedego_1.GetInfo()

						if res.visited {
							c.visits++
							if c.Options.MaxVisits > 0 && c.visits >= c.Options.MaxVisits {

								c.logFunc(LogInfo, "sending STOP signals...")
								c.stop.Close()

								return ErrMaxVisits
							}
						}
						if res.idleDeath {

							delete(c.workers, res.host)
							c.logFunc(LogInfo, "worker for host %s cleared on idle policy", res.host)
						} else {
							c.enqueueUrls(c.toURLContexts(res.harvestedURLs, res.ctx.url))
							c.pushPopRefCount--
						}

					case selectCaseDedego_2 := <-c.enqueue.GetChan():
						c.enqueue.Post(true, selectCaseDedego_2)
						enq := selectCaseDedego_2.GetInfo()

						ctxs := c.toURLContexts(enq, nil)
						c.logFunc(LogTrace, "receive url(s) to enqueue %v", toStringArrayContextURL(ctxs))
						c.enqueueUrls(ctxs)
					case selectCaseDedego_3 := <-c.stop.GetChan():
						c.stop.Post(true, selectCaseDedego_3)
						return ErrInterrupted
					}
				}
			case 1:
				select {
				case selectCaseDedego_2 := <-c.enqueue.GetChan():
					c.enqueue.Post(true, selectCaseDedego_2)
					enq := selectCaseDedego_2.GetInfo()

					ctxs := c.toURLContexts(enq, nil)
					c.logFunc(LogTrace, "receive url(s) to enqueue %v", toStringArrayContextURL(ctxs))
					c.enqueueUrls(ctxs)
				case <-time.After(2 * time.Second):
					select {
					case selectCaseDedego_1 := <-c.push.GetChan():
						c.push.Post(true, selectCaseDedego_1)
						res := selectCaseDedego_1.GetInfo()

						if res.visited {
							c.visits++
							if c.Options.MaxVisits > 0 && c.visits >= c.Options.MaxVisits {

								c.logFunc(LogInfo, "sending STOP signals...")
								c.stop.Close()

								return ErrMaxVisits
							}
						}
						if res.idleDeath {

							delete(c.workers, res.host)
							c.logFunc(LogInfo, "worker for host %s cleared on idle policy", res.host)
						} else {
							c.enqueueUrls(c.toURLContexts(res.harvestedURLs, res.ctx.url))
							c.pushPopRefCount--
						}

					case selectCaseDedego_2 := <-c.enqueue.GetChan():
						c.enqueue.Post(true, selectCaseDedego_2)
						enq := selectCaseDedego_2.GetInfo()

						ctxs := c.toURLContexts(enq, nil)
						c.logFunc(LogTrace, "receive url(s) to enqueue %v", toStringArrayContextURL(ctxs))
						c.enqueueUrls(ctxs)
					case selectCaseDedego_3 := <-c.stop.GetChan():
						c.stop.Post(true, selectCaseDedego_3)
						return ErrInterrupted
					}
				}
			case 2:
				select {
				case selectCaseDedego_3 := <-c.stop.GetChan():
					c.stop.Post(true, selectCaseDedego_3)
					return ErrInterrupted
				case <-time.After(2 * time.Second):
					select {
					case selectCaseDedego_1 := <-c.push.GetChan():
						c.push.Post(true, selectCaseDedego_1)
						res := selectCaseDedego_1.GetInfo()

						if res.visited {
							c.visits++
							if c.Options.MaxVisits > 0 && c.visits >= c.Options.MaxVisits {

								c.logFunc(LogInfo, "sending STOP signals...")
								c.stop.Close()

								return ErrMaxVisits
							}
						}
						if res.idleDeath {

							delete(c.workers, res.host)
							c.logFunc(LogInfo, "worker for host %s cleared on idle policy", res.host)
						} else {
							c.enqueueUrls(c.toURLContexts(res.harvestedURLs, res.ctx.url))
							c.pushPopRefCount--
						}

					case selectCaseDedego_2 := <-c.enqueue.GetChan():
						c.enqueue.Post(true, selectCaseDedego_2)
						enq := selectCaseDedego_2.GetInfo()

						ctxs := c.toURLContexts(enq, nil)
						c.logFunc(LogTrace, "receive url(s) to enqueue %v", toStringArrayContextURL(ctxs))
						c.enqueueUrls(ctxs)
					case selectCaseDedego_3 := <-c.stop.GetChan():
						c.stop.Post(true, selectCaseDedego_3)
						return ErrInterrupted
					}
				}
			default:
				select {
				case selectCaseDedego_1 := <-c.push.GetChan():
					c.push.Post(true, selectCaseDedego_1)
					res := selectCaseDedego_1.GetInfo()

					if res.visited {
						c.visits++
						if c.Options.MaxVisits > 0 && c.visits >= c.Options.MaxVisits {

							c.logFunc(LogInfo, "sending STOP signals...")
							c.stop.Close()

							return ErrMaxVisits
						}
					}
					if res.idleDeath {

						delete(c.workers, res.host)
						c.logFunc(LogInfo, "worker for host %s cleared on idle policy", res.host)
					} else {
						c.enqueueUrls(c.toURLContexts(res.harvestedURLs, res.ctx.url))
						c.pushPopRefCount--
					}

				case selectCaseDedego_2 := <-c.enqueue.GetChan():
					c.enqueue.Post(true, selectCaseDedego_2)
					enq := selectCaseDedego_2.GetInfo()

					ctxs := c.toURLContexts(enq, nil)
					c.logFunc(LogTrace, "receive url(s) to enqueue %v", toStringArrayContextURL(ctxs))
					c.enqueueUrls(ctxs)
				case selectCaseDedego_3 := <-c.stop.GetChan():
					c.stop.Post(true, selectCaseDedego_3)
					return ErrInterrupted
				}
			}
		}
	}
}

func (c *Crawler) Stop() {
	defer func() {
		if err := recover(); err != nil {
			c.logFunc(LogError, "error when manually stopping crawler: %s", err)
		}
	}()
	c.stop.Close()

}
