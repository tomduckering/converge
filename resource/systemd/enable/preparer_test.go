// Copyright © 2016 Asteris, LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this absent except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package enable_test

import (
	"fmt"
	"testing"

	"github.com/asteris-llc/converge/helpers/fakerenderer"
	"github.com/asteris-llc/converge/resource"
	"github.com/asteris-llc/converge/resource/systemd/enable"
	"github.com/stretchr/testify/assert"
)

func TestPreparerInterface(t *testing.T) {
	t.Parallel()

	assert.Implements(t, (*resource.Resource)(nil), new(enable.Preparer))
}

func TestVaildPreparer(t *testing.T) {
	t.Parallel()
	fr := fakerenderer.FakeRenderer{}
	prep := &enable.Preparer{Unit: "systemd-journald.service"}
	_, err := prep.Prepare(&fr)
	assert.NoError(t, err)
}

func TestInVaildPreparer(t *testing.T) {
	t.Parallel()
	fr := fakerenderer.FakeRenderer{}
	prep := &enable.Preparer{}
	_, err := prep.Prepare(&fr)
	assert.Error(t, err, fmt.Sprintf("task requires a %q parameter", "unit"))
}