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

package resource

// Tasker is a struct that is or contains an embedded resource.Task
type Tasker interface {
	GetTask() (Task, bool)
}

// Task does checking as Monitor does, but it can also make changes to make the
// checks pass.
type Task interface {
	Check(Renderer) (TaskStatus, error)
	Apply(Renderer) (TaskStatus, error)
}

// Resource adds metadata about the executed tasks
type Resource interface {
	Prepare(Renderer) (Task, error)
}

// Renderer is passed to resources
type Renderer interface {
	Value() (value string, present bool)
	Render(name, content string) (string, error)
	RequiredRender(name, content string) (string, error)
	RenderBool(name, content string) (bool, error)
	RenderStringSlice(name string, content []string) ([]string, error)
	RenderStringMapToStringSlice(name string, content map[string]string, toString func(string, string) string) ([]string, error)
}

// TaskWrapper provides an implementation of render.Tasker for tasks
type TaskWrapper struct {
	Task
}

// GetTask provides Tasker.GetTask ontop of a task
func (t *TaskWrapper) GetTask() (Task, bool) {
	return t.Task, true
}

// WrapTask creates a new TaskWrapper
func WrapTask(t Task) *TaskWrapper {
	return &TaskWrapper{t}
}

// ResolveTask unwraps Tasker layers until it finds an underlying Task or fails
func ResolveTask(w interface{}) (Task, bool) {
	if w == nil {
		return nil, false
	}
	if tasker, ok := w.(Tasker); ok {
		taskerTask, found := tasker.GetTask()
		if !found {
			return nil, false
		}
		return ResolveTask(taskerTask)
	} else if task, ok := w.(Task); ok {
		return task, true
	}
	return nil, false
}
