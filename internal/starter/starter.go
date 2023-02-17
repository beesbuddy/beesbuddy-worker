// Copyright 2017 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package shutdown provides a mechanism for registering handlers to be called
// on process shutdown.
package starter

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/beesbuddy/beesbuddy-worker/internal/log"
)

// GracePeriod specifies the maximum amount of time during which all shutdown
// handlers must complete before the process forcibly exits.
const GracePeriod = 1 * time.Minute

// Handle registers the onShutdown function to be run when the system is being
// shut down. On shutdown, registered functions are run in last-in-first-out
// order. Handle may be called concurrently.
func Handle(onShutdown func()) {
	starter.mu.Lock()
	defer starter.mu.Unlock()

	starter.sequence = append(starter.sequence, onShutdown)
}

// Stop calls all registered shutdown closures in last-in-first-out order and
// terminates the process with the given status code.
// It only executes once and guarantees termination within GracePeriod.
// Stop may be called concurrently.
func Stop(code int) {
	starter.once.Do(func() {
		log.Debug.Printf("shutdown: status code %d", code)

		// Ensure we terminate after a fixed amount of time.
		go func() {
			killSleep(GracePeriod)
			// Don't use log package here; it may have been flushed already.
			fmt.Fprintf(os.Stderr, "shutdown: %v elapsed since shutdown requested; exiting forcefully", GracePeriod)
			os.Exit(1)
		}()

		starter.mu.Lock() // No need to ever unlock.
		for i := len(starter.sequence) - 1; i >= 0; i-- {
			starter.sequence[i]()
		}

		os.Exit(code)
	})
}

// Testing hook.
var killSleep = time.Sleep

var starter struct {
	mu       sync.Mutex
	sequence []func()
	once     sync.Once
}

func Ignite(code int) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	sig := <-interrupt
	log.Info.Printf("shutdown: process received signal %v", sig)
	Stop(code)
}
