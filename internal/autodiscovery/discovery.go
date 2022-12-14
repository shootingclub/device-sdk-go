// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package autodiscovery

import (
	"sync"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"

	"git.dev.tengwanweigu.com/common/device-sdk-go.git/pkg/models"
)

type discoveryLocker struct {
	busy bool
	mux  sync.Mutex
}

var locker discoveryLocker

func DiscoveryWrapper(discovery models.ProtocolDiscovery, lc logger.LoggingClient) {
	locker.mux.Lock()
	if locker.busy {
		lc.Info("another device discovery process is currently running")
		locker.mux.Unlock()
		return
	}
	locker.busy = true
	locker.mux.Unlock()

	lc.Debug("protocol discovery triggered")
	discovery.Discover()

	// ReleaseLock
	locker.mux.Lock()
	locker.busy = false
	locker.mux.Unlock()
}
