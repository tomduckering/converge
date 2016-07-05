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

package load

import (
	"fmt"

	"github.com/asteris-llc/converge/graph"
	"github.com/asteris-llc/converge/parse"
	"github.com/asteris-llc/converge/resource"
	"github.com/asteris-llc/converge/resource/shell"
	"github.com/hashicorp/hcl"
)

// SetResources loads the resources for each graph node
func SetResources(g *graph.Graph) (*graph.Graph, error) {
	return g.Transform(func(id string, v interface{}, edges []string) (interface{}, []string, error) {
		if id == "root" { // root
			return v, edges, nil
		}

		node, ok := v.(*parse.Node)
		if !ok {
			return v, edges, fmt.Errorf("SetResources can only be used on Graphs of *parse.Node. I got %T", v)
		}

		var dest resource.Resource
		switch node.Kind() {
		case "task":
			dest = new(shell.Preparer)

		default:
			return v, edges, fmt.Errorf("%q is not a valid resource type in %q", node.Kind(), node)
		}

		err := hcl.DecodeObject(dest, node.ObjectItem)
		if err != nil {
			return v, edges, err
		}

		return dest, edges, nil
	})
}