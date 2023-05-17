package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"path"

	"github.com/ErikKassubek/deadlockDetectorGo/src/dedego"
	"github.com/PuerkitoBio/goquery"
	"github.com/andybalholm/cascadia"
	"github.com/temoto/robotstxt"
	"golang.org/x/net/html"
)

type worker struct {
	host  string
	index int

	push    dedego.Chan[*workerResponse]
	pop     popChannel
	stop    dedego.Chan[struct{}]
	enqueue dedego.Chan[interface{}]
	wg      *sync.WaitGroup

	robotsGroup *robotstxt.Group

	logFunc func(LogFlags, string, ...interface{})

	wait           dedego.Chan[time.Time]
	lastFetch      *FetchInfo
	lastCrawlDelay time.Duration
	opts           *Options
}

func (w *worker) run() {
	defer func() {
		w.logFunc(LogInfo, "worker done.")
		w.wg.Done()
	}()

	for {
		var idleChan <-chan time.Time

		w.logFunc(LogInfo, "waiting for pop...")

		if w.opts.WorkerIdleTTL > 0 {
			idleChan = time.After(w.opts.WorkerIdleTTL)
		}
		{
			dedego.PreSelect(false, w.stop.GetIdPre(true), w.pop.GetIdPre(true))
			switch dedegoFetchOrder[4] {
			case 0:
				select {

				case selectCaseDedego_7 := <-w.stop.GetChan():
					w.stop.Post(true, selectCaseDedego_7)
					w.logFunc(LogInfo, "stop signal received.")
					return
				case <-time.After(2 * time.Second):
					select {
					case selectCaseDedego_7 := <-w.stop.GetChan():
						w.stop.Post(true, selectCaseDedego_7)
						w.logFunc(LogInfo, "stop signal received.")
						return

					case <-idleChan:
						w.logFunc(LogInfo, "idle timeout received.")
						w.sendResponse(nil, false, nil, true)
						return

					case selectCaseDedego_8 := <-w.pop.GetChan():
						w.pop.Post(true, selectCaseDedego_8)
						batch := selectCaseDedego_8.GetInfo()

						for _, ctx := range batch {
							w.logFunc(LogInfo, "popped: %s", ctx.url)

							if ctx.IsRobotsURL() {
								w.requestRobotsTxt(ctx)
							} else if w.isAllowedPerRobotsPolicies(ctx.url) {
								w.requestURL(ctx, ctx.HeadBeforeGet)
							} else {

								w.opts.Extender.Disallowed(ctx)
								w.sendResponse(ctx, false, nil, false)
							}
							{
								dedego.PreSelect(true, w.stop.GetIdPre(true))
								switch dedegoFetchOrder[3] {
								case 0:
									select {

									case selectCaseDedego_6 := <-w.stop.GetChan():
										w.stop.Post(true, selectCaseDedego_6)
										w.logFunc(LogInfo, "stop signal received.")
										return
									case <-time.After(2 * time.Second):
										select {
										case selectCaseDedego_6 := <-w.stop.GetChan():
											w.stop.Post(true, selectCaseDedego_6)
											w.logFunc(LogInfo, "stop signal received.")
											return
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
										case selectCaseDedego_6 := <-w.stop.GetChan():
											w.stop.Post(true, selectCaseDedego_6)
											w.logFunc(LogInfo, "stop signal received.")
											return
										default:
											dedego.PostDefault()

										}
									}
								default:
									select {
									case selectCaseDedego_6 := <-w.stop.GetChan():
										w.stop.Post(true, selectCaseDedego_6)
										w.logFunc(LogInfo, "stop signal received.")
										return
									default:
										dedego.PostDefault()

									}
								}
							}
						}
					}
				}
			case 1:
				select {
				case <-idleChan:
					w.logFunc(LogInfo, "idle timeout received.")
					w.sendResponse(nil, false, nil, true)
					return
				case <-time.After(2 * time.Second):
					select {
					case selectCaseDedego_7 := <-w.stop.GetChan():
						w.stop.Post(true, selectCaseDedego_7)
						w.logFunc(LogInfo, "stop signal received.")
						return

					case <-idleChan:
						w.logFunc(LogInfo, "idle timeout received.")
						w.sendResponse(nil, false, nil, true)
						return

					case selectCaseDedego_8 := <-w.pop.GetChan():
						w.pop.Post(true, selectCaseDedego_8)
						batch := selectCaseDedego_8.GetInfo()

						for _, ctx := range batch {
							w.logFunc(LogInfo, "popped: %s", ctx.url)

							if ctx.IsRobotsURL() {
								w.requestRobotsTxt(ctx)
							} else if w.isAllowedPerRobotsPolicies(ctx.url) {
								w.requestURL(ctx, ctx.HeadBeforeGet)
							} else {

								w.opts.Extender.Disallowed(ctx)
								w.sendResponse(ctx, false, nil, false)
							}
							{
								dedego.PreSelect(true, w.stop.GetIdPre(true))
								switch dedegoFetchOrder[3] {
								case 0:
									select {

									case selectCaseDedego_6 := <-w.stop.GetChan():
										w.stop.Post(true, selectCaseDedego_6)
										w.logFunc(LogInfo, "stop signal received.")
										return
									case <-time.After(2 * time.Second):
										select {
										case selectCaseDedego_6 := <-w.stop.GetChan():
											w.stop.Post(true, selectCaseDedego_6)
											w.logFunc(LogInfo, "stop signal received.")
											return
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
										case selectCaseDedego_6 := <-w.stop.GetChan():
											w.stop.Post(true, selectCaseDedego_6)
											w.logFunc(LogInfo, "stop signal received.")
											return
										default:
											dedego.PostDefault()

										}
									}
								default:
									select {
									case selectCaseDedego_6 := <-w.stop.GetChan():
										w.stop.Post(true, selectCaseDedego_6)
										w.logFunc(LogInfo, "stop signal received.")
										return
									default:
										dedego.PostDefault()

									}
								}
							}
						}
					}
				}
			case 2:
				select {
				case selectCaseDedego_8 := <-w.pop.GetChan():
					w.pop.Post(true, selectCaseDedego_8)
					batch := selectCaseDedego_8.GetInfo()

					for _, ctx := range batch {
						w.logFunc(LogInfo, "popped: %s", ctx.url)

						if ctx.IsRobotsURL() {
							w.requestRobotsTxt(ctx)
						} else if w.isAllowedPerRobotsPolicies(ctx.url) {
							w.requestURL(ctx, ctx.HeadBeforeGet)
						} else {

							w.opts.Extender.Disallowed(ctx)
							w.sendResponse(ctx, false, nil, false)
						}
						{
							dedego.PreSelect(true, w.stop.GetIdPre(true))
							switch dedegoFetchOrder[3] {
							case 0:
								select {

								case selectCaseDedego_6 := <-w.stop.GetChan():
									w.stop.Post(true, selectCaseDedego_6)
									w.logFunc(LogInfo, "stop signal received.")
									return
								case <-time.After(2 * time.Second):
									select {
									case selectCaseDedego_6 := <-w.stop.GetChan():
										w.stop.Post(true, selectCaseDedego_6)
										w.logFunc(LogInfo, "stop signal received.")
										return
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
									case selectCaseDedego_6 := <-w.stop.GetChan():
										w.stop.Post(true, selectCaseDedego_6)
										w.logFunc(LogInfo, "stop signal received.")
										return
									default:
										dedego.PostDefault()

									}
								}
							default:
								select {
								case selectCaseDedego_6 := <-w.stop.GetChan():
									w.stop.Post(true, selectCaseDedego_6)
									w.logFunc(LogInfo, "stop signal received.")
									return
								default:
									dedego.PostDefault()

								}
							}
						}
					}
				case <-time.After(2 * time.Second):
					select {
					case selectCaseDedego_7 := <-w.stop.GetChan():
						w.stop.Post(true, selectCaseDedego_7)
						w.logFunc(LogInfo, "stop signal received.")
						return

					case <-idleChan:
						w.logFunc(LogInfo, "idle timeout received.")
						w.sendResponse(nil, false, nil, true)
						return

					case selectCaseDedego_8 := <-w.pop.GetChan():
						w.pop.Post(true, selectCaseDedego_8)
						batch := selectCaseDedego_8.GetInfo()

						for _, ctx := range batch {
							w.logFunc(LogInfo, "popped: %s", ctx.url)

							if ctx.IsRobotsURL() {
								w.requestRobotsTxt(ctx)
							} else if w.isAllowedPerRobotsPolicies(ctx.url) {
								w.requestURL(ctx, ctx.HeadBeforeGet)
							} else {

								w.opts.Extender.Disallowed(ctx)
								w.sendResponse(ctx, false, nil, false)
							}
							{
								dedego.PreSelect(true, w.stop.GetIdPre(true))
								switch dedegoFetchOrder[3] {
								case 0:
									select {

									case selectCaseDedego_6 := <-w.stop.GetChan():
										w.stop.Post(true, selectCaseDedego_6)
										w.logFunc(LogInfo, "stop signal received.")
										return
									case <-time.After(2 * time.Second):
										select {
										case selectCaseDedego_6 := <-w.stop.GetChan():
											w.stop.Post(true, selectCaseDedego_6)
											w.logFunc(LogInfo, "stop signal received.")
											return
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
										case selectCaseDedego_6 := <-w.stop.GetChan():
											w.stop.Post(true, selectCaseDedego_6)
											w.logFunc(LogInfo, "stop signal received.")
											return
										default:
											dedego.PostDefault()

										}
									}
								default:
									select {
									case selectCaseDedego_6 := <-w.stop.GetChan():
										w.stop.Post(true, selectCaseDedego_6)
										w.logFunc(LogInfo, "stop signal received.")
										return
									default:
										dedego.PostDefault()

									}
								}
							}
						}
					}
				}
			default:
				select {
				case selectCaseDedego_7 := <-w.stop.GetChan():
					w.stop.Post(true, selectCaseDedego_7)
					w.logFunc(LogInfo, "stop signal received.")
					return

				case <-idleChan:
					w.logFunc(LogInfo, "idle timeout received.")
					w.sendResponse(nil, false, nil, true)
					return

				case selectCaseDedego_8 := <-w.pop.GetChan():
					w.pop.Post(true, selectCaseDedego_8)
					batch := selectCaseDedego_8.GetInfo()

					for _, ctx := range batch {
						w.logFunc(LogInfo, "popped: %s", ctx.url)

						if ctx.IsRobotsURL() {
							w.requestRobotsTxt(ctx)
						} else if w.isAllowedPerRobotsPolicies(ctx.url) {
							w.requestURL(ctx, ctx.HeadBeforeGet)
						} else {

							w.opts.Extender.Disallowed(ctx)
							w.sendResponse(ctx, false, nil, false)
						}
						{
							dedego.PreSelect(true, w.stop.GetIdPre(true))
							switch dedegoFetchOrder[3] {
							case 0:
								select {

								case selectCaseDedego_6 := <-w.stop.GetChan():
									w.stop.Post(true, selectCaseDedego_6)
									w.logFunc(LogInfo, "stop signal received.")
									return
								case <-time.After(2 * time.Second):
									select {
									case selectCaseDedego_6 := <-w.stop.GetChan():
										w.stop.Post(true, selectCaseDedego_6)
										w.logFunc(LogInfo, "stop signal received.")
										return
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
									case selectCaseDedego_6 := <-w.stop.GetChan():
										w.stop.Post(true, selectCaseDedego_6)
										w.logFunc(LogInfo, "stop signal received.")
										return
									default:
										dedego.PostDefault()

									}
								}
							default:
								select {
								case selectCaseDedego_6 := <-w.stop.GetChan():
									w.stop.Post(true, selectCaseDedego_6)
									w.logFunc(LogInfo, "stop signal received.")
									return
								default:
									dedego.PostDefault()

								}
							}
						}
					}
				}
			}
		}
	}
}

func (w *worker) isAllowedPerRobotsPolicies(u *url.URL) bool {
	if w.robotsGroup != nil {

		ok := w.robotsGroup.Test(u.Path)
		if !ok {
			w.logFunc(LogIgnored, "ignored on robots.txt policy: %s", u.String())
		}
		return ok
	}

	return true
}

func (w *worker) requestURL(ctx *URLContext, headRequest bool) {
	if res, ok := w.fetchURL(ctx, w.opts.UserAgent, headRequest); ok {
		var harvested interface{}
		var visited bool

		defer res.Body.Close()

		if res.StatusCode >= 200 && res.StatusCode < 300 {

			harvested = w.visitURL(ctx, res)
			visited = true
		} else {

			w.opts.Extender.Error(newCrawlErrorMessage(ctx, res.Status, CekHttpStatusCode))
			w.logFunc(LogError, "ERROR status code for %s: %s", ctx.url, res.Status)
		}
		w.sendResponse(ctx, visited, harvested, false)
	}
}

func (w *worker) requestRobotsTxt(ctx *URLContext) {

	if robData, reqRob := w.opts.Extender.RequestRobots(ctx, w.opts.RobotUserAgent); !reqRob {
		w.logFunc(LogInfo, "using robots.txt from cache")
		w.robotsGroup = w.getRobotsTxtGroup(ctx, robData, nil)

	} else if res, ok := w.fetchURL(ctx, w.opts.UserAgent, false); ok {

		defer res.Body.Close()
		w.robotsGroup = w.getRobotsTxtGroup(ctx, nil, res)
	}
}

func (w *worker) getRobotsTxtGroup(ctx *URLContext, b []byte, res *http.Response) (g *robotstxt.Group) {
	var data *robotstxt.RobotsData
	var e error

	if res != nil {
		var buf bytes.Buffer
		io.Copy(&buf, res.Body)
		res.Body = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))
		data, e = robotstxt.FromResponse(res)

		res.Body = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))

		w.opts.Extender.FetchedRobots(ctx, res)
	} else {
		data, e = robotstxt.FromBytes(b)
	}

	if e != nil {
		w.opts.Extender.Error(newCrawlError(nil, e, CekParseRobots))
		w.logFunc(LogError, "ERROR parsing robots.txt for host %s: %s", w.host, e)
	} else {
		g = data.FindGroup(w.opts.RobotUserAgent)
	}
	return g
}

func (w *worker) setCrawlDelay() {
	var robDelay time.Duration

	if w.robotsGroup != nil {
		robDelay = w.robotsGroup.CrawlDelay
	}
	w.lastCrawlDelay = w.opts.Extender.ComputeDelay(w.host,
		&DelayInfo{
			w.opts.CrawlDelay,
			robDelay,
			w.lastCrawlDelay,
		},
		w.lastFetch)
	w.logFunc(LogInfo, "using crawl-delay: %v", w.lastCrawlDelay)
}

func (w *worker) fetchURL(ctx *URLContext, agent string, headRequest bool) (res *http.Response, ok bool) {
	var e error
	var silent bool

	for {

		w.logFunc(LogTrace, "waiting for crawl delay")
		if w.wait.GetChan() != nil {
			w.wait.Receive()

			w.wait.SetNil()
		}

		w.setCrawlDelay()

		now := time.Now()

		if res, e = w.opts.Extender.Fetch(ctx, agent, headRequest); e != nil {

			if ue, ok := e.(*url.Error); ok {

				if ue.Err == ErrEnqueueRedirect {

					silent = true

					if ur, e := ctx.url.Parse(ue.URL); e != nil {

						w.opts.Extender.Error(newCrawlError(nil, e, CekParseRedirectURL))
						w.logFunc(LogError, "ERROR parsing redirect URL %s: %s", ue.URL, e)
					} else {
						w.logFunc(LogTrace, "redirect to %s from %s, linked from %s", ur, ctx.URL(), ctx.SourceURL())

						rCtx := ctx.cloneForRedirect(ur, w.opts.URLNormalizationFlags)
						w.enqueue.Send(rCtx)

					}
				}
			}

			w.lastFetch = nil

			if !silent {

				w.opts.Extender.Error(newCrawlError(ctx, e, CekFetch))
				w.logFunc(LogError, "ERROR fetching %s: %s", ctx.url, e)
			}

			w.sendResponse(ctx, false, nil, false)
			return nil, false

		}

		fetchDuration := time.Now().Sub(now)

		go func() {
			a := <-time.After(w.lastCrawlDelay)
			w.wait.Send(a)
		}()

		w.lastFetch = &FetchInfo{
			ctx,
			fetchDuration,
			res.StatusCode,
			headRequest,
		}

		if headRequest {

			defer res.Body.Close()

			headRequest = false

			if !w.opts.Extender.RequestGet(ctx, res) {
				w.logFunc(LogIgnored, "ignored on HEAD filter policy: %s", ctx.url)
				w.sendResponse(ctx, false, nil, false)
				ok = false
				break
			}
		} else {
			ok = true
			break
		}
	}
	return
}

func (w *worker) sendResponse(ctx *URLContext, visited bool, harvested interface{}, idleDeath bool) {

	if ctx == nil || !isRobotsURL(ctx.url) {
		{
			dedego.PreSelect(true, w.stop.GetIdPre(true))
			switch dedegoFetchOrder[5] {
			case 0:
				select {

				case selectCaseDedego_9 := <-w.stop.GetChan():
					w.stop.Post(true, selectCaseDedego_9)
					w.logFunc(LogInfo, "ignoring send response, will stop.")
					return
				case <-time.After(2 * time.Second):
					select {
					case selectCaseDedego_9 := <-w.stop.GetChan():
						w.stop.Post(true, selectCaseDedego_9)
						w.logFunc(LogInfo, "ignoring send response, will stop.")
						return
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
					case selectCaseDedego_9 := <-w.stop.GetChan():
						w.stop.Post(true, selectCaseDedego_9)
						w.logFunc(LogInfo, "ignoring send response, will stop.")
						return
					default:
						dedego.PostDefault()

					}
				}
			default:
				select {
				case selectCaseDedego_9 := <-w.stop.GetChan():
					w.stop.Post(true, selectCaseDedego_9)
					w.logFunc(LogInfo, "ignoring send response, will stop.")
					return
				default:
					dedego.PostDefault()

				}
			}
		}

		res := &workerResponse{
			ctx,
			visited,
			harvested,
			w.host,
			idleDeath,
		}
		w.push.Send(res)

	}
}

func (w *worker) visitURL(ctx *URLContext, res *http.Response) interface{} {
	var doc *goquery.Document
	var harvested interface{}
	var doLinks bool

	if bd, e := ioutil.ReadAll(res.Body); e != nil {
		w.opts.Extender.Error(newCrawlError(ctx, e, CekReadBody))
		w.logFunc(LogError, "ERROR reading body %s: %s", ctx.url, e)
	} else {
		if node, e := html.Parse(bytes.NewBuffer(bd)); e != nil {
			w.opts.Extender.Error(newCrawlError(ctx, e, CekParseBody))
			w.logFunc(LogError, "ERROR parsing %s: %s", ctx.url, e)
		} else {
			doc = goquery.NewDocumentFromNode(node)
			doc.Url = res.Request.URL
		}

		res.Body = ioutil.NopCloser(bytes.NewBuffer(bd))
	}

	if harvested, doLinks = w.opts.Extender.Visit(ctx, res, doc); doLinks {

		if doc != nil {
			harvested = w.processLinks(doc)
		} else {
			w.opts.Extender.Error(newCrawlErrorMessage(ctx, "No goquery document to process links.", CekProcessLinks))
			w.logFunc(LogError, "ERROR processing links %s", ctx.url)
		}
	}

	w.opts.Extender.Visited(ctx, harvested)

	return harvested
}

func handleBaseTag(root *url.URL, baseHref string, aHref string) string {
	resolvedBase, err := root.Parse(baseHref)
	if err != nil {
		return ""
	}

	parsedURL, err := url.Parse(aHref)
	if err != nil {
		return ""
	}

	if parsedURL.Host == "" && !strings.HasPrefix(aHref, "/") {
		aHref = path.Join(resolvedBase.Path, aHref)
	}

	resolvedURL, err := resolvedBase.Parse(aHref)
	if err != nil {
		return ""
	}
	return resolvedURL.String()
}

var (
	aHrefMatcher    = cascadia.MustCompile("a[href]")
	baseHrefMatcher = cascadia.MustCompile("base[href]")
)

func (w *worker) processLinks(doc *goquery.Document) (result []*url.URL) {
	baseURL, _ := doc.FindMatcher(baseHrefMatcher).Attr("href")
	urls := doc.FindMatcher(aHrefMatcher).Map(func(_ int, s *goquery.Selection) string {
		val, _ := s.Attr("href")
		if baseURL != "" {
			val = handleBaseTag(doc.Url, baseURL, val)
		}
		return val
	})
	for _, s := range urls {

		if len(s) > 0 && !strings.HasPrefix(s, "#") {
			if parsed, e := url.Parse(s); e == nil {
				parsed = doc.Url.ResolveReference(parsed)
				result = append(result, parsed)
			} else {
				w.logFunc(LogIgnored, "ignore on unparsable policy %s: %s", s, e.Error())
			}
		}
	}
	return
}
