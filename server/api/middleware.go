package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {

				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *Application) rateLimit(next http.Handler) http.Handler {

	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mux     sync.Mutex
		clients = make(map[string]*client)
	)

	// go routine for stale clients cleaning
	go func() {
		for {
			time.Sleep(time.Minute)
			mux.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mux.Unlock()
		}
	}()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if app.config.limiter.enabled {

			ip, _, err := net.SplitHostPort(r.RemoteAddr)

			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			mux.Lock()

			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
				}
			}

			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow() {
				app.rateLimitExceededResponse(w, r)
				return
			}

			mux.Unlock()
		}

		next.ServeHTTP(w, r)
	})
}

// global rate limiter
// func (app *Application) rateLimit(next http.Handler) http.Handler {

// 	limiter := rate.NewLimiter(2, 4)

// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if !limiter.Allow() {
// 			app.rateLimitExceededResponse(w, r)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }
