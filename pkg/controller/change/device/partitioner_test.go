// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package device

import (
	types "github.com/onosproject/onos-api/go/onos/config"
	devicechange "github.com/onosproject/onos-api/go/onos/config/change/device"
	"github.com/onosproject/onos-config/pkg/controller"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDevicePartitioner(t *testing.T) {
	partitioner := &Partitioner{}
	key, err := partitioner.Partition(types.ID(devicechange.NewID("change-1", "device-1", "1.0.0")))
	assert.NoError(t, err)
	assert.Equal(t, controller.PartitionKey("device-1"), key)
}
