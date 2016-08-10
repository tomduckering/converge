// Copyright © 2016 Asteris, LLC
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

package disable_test

import (
	"testing"

	"github.com/asteris-llc/converge/helpers"
	"github.com/asteris-llc/converge/resource"
	"github.com/asteris-llc/converge/resource/file/absent"
	"github.com/asteris-llc/converge/resource/systemd/disable"
	"github.com/asteris-llc/converge/resource/systemd/enable"
	"github.com/stretchr/testify/assert"
)

func TestTemplateInterface(t *testing.T) {
	t.Parallel()

	assert.Implements(t, (*resource.Task)(nil), new(absent.Absent))
}

func TestCheck(t *testing.T) {
	defer helpers.HideLogs(t)()

	tasks := []resource.Task{
		&disable.Disable{Unit: "systemd-journald.service"},
	}

	checks := []helpers.CheckValidator{
		helpers.CheckValidatorCreator("", true, ""),
	}
	helpers.TaskCheckValidator(tasks, checks, t)
}

func TestApply(t *testing.T) {
	revert := &enable.Enable{Unit: "systemd-journald.service"}
	defer revert.Apply()
	tasks := []resource.Task{
		&disable.Disable{Unit: "systemd-journald.service"},
	}

	errs := []string{
		"",
	}
	helpers.TaskApplyValidator(tasks, errs, t)
}
