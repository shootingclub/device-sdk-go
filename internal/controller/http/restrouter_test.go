//
// Copyright (c) 2019 Intel Corporation
// Copyright (C) 2020 IOTech Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package http

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"testing"

	bootstrapContainer "github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/shootingclub/device-sdk-go/internal/config"
	"github.com/shootingclub/device-sdk-go/internal/container"
)

func TestAddRoute(t *testing.T) {

	tests := []struct {
		Name          string
		Route         string
		ErrorExpected bool
	}{
		{"Success", "/api/v2/test", false},
		{"Reserved Route", common.ApiVersionRoute, true},
	}

	lc := logger.NewMockClient()
	dic := di.NewContainer(di.ServiceConstructorMap{
		bootstrapContainer.LoggingClientInterfaceName: func(get di.Get) interface{} {
			return lc
		},
		container.ConfigurationName: func(get di.Get) interface{} {
			return &config.ConfigurationStruct{}
		},
	})

	for _, test := range tests {
		r := mux.NewRouter()
		controller := NewRestController(r, dic, uuid.NewString())
		controller.InitRestRoutes()

		err := controller.AddRoute(test.Route, func(http.ResponseWriter, *http.Request) {}, http.MethodPost)
		if test.ErrorExpected {
			assert.Error(t, err, "Expected an error")
		} else {
			if !assert.NoError(t, err, "Unexpected an error") {
				t.Fatal()
			}

			err := controller.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
				path, err := route.GetPathTemplate()
				if err != nil {
					return err
				}

				// Have to skip all the reserved routes that have previously been added.
				if controller.reservedRoutes[path] {
					return nil
				}

				routeMethods, err := route.GetMethods()
				if err != nil {
					return err
				}

				assert.Equal(t, test.Route, path)
				assert.Equal(t, http.MethodPost, routeMethods[0], "Expected POST Method")
				return nil
			})

			assert.NoError(t, err, "Unexpected error examining route")
		}
	}
}

func TestInitRestRoutes(t *testing.T) {
	lc := logger.NewMockClient()
	dic := di.NewContainer(di.ServiceConstructorMap{
		bootstrapContainer.LoggingClientInterfaceName: func(get di.Get) interface{} {
			return lc
		},
		container.ConfigurationName: func(get di.Get) interface{} {
			return &config.ConfigurationStruct{}
		},
	})
	r := mux.NewRouter()
	controller := NewRestController(r, dic, uuid.NewString())
	controller.InitRestRoutes()

	err := controller.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return err
		}

		// Verify the route is reserved by attempting to add it as 'external' route.
		// If tests fails then the route was not added to the reserved list
		err = controller.AddRoute(path, func(http.ResponseWriter, *http.Request) {})
		assert.Error(t, err, path, fmt.Sprintf("Expected error for '%s'", path))
		return nil
	})

	assert.NoError(t, err, "Unexpected error examining route")
}
