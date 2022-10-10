// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2017-2018 Canonical Ltd
// Copyright (C) 2018-2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

// This package provides a simple example of a device service.
package main

import (
	"git.dev.tengwanweigu.com/common/device-sdk-go.git"
	"git.dev.tengwanweigu.com/common/device-sdk-go.git/example/driver"
	"git.dev.tengwanweigu.com/common/device-sdk-go.git/pkg/startup"
)

const (
	serviceName string = "device-simple"
)

func main() {
	sd := driver.SimpleDriver{}
	startup.Bootstrap(serviceName, device.Version, &sd)
}
