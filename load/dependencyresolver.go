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
	"context"
	"fmt"
	"io/ioutil"
	"text/template"

	"github.com/asteris-llc/converge/graph"
	"github.com/asteris-llc/converge/helpers/logging"
	"github.com/asteris-llc/converge/parse"
	"github.com/asteris-llc/converge/render/extensions"
	"github.com/asteris-llc/converge/render/preprocessor"
)

type dependencyGenerator func(node *parse.Node) ([]string, error)

// ResolveDependencies examines the strings and depdendencies at each vertex of
// the graph and creates edges to fit them
func ResolveDependencies(ctx context.Context, g *graph.Graph) (*graph.Graph, error) {
	logger := logging.GetLogger(ctx).WithField("function", "ResolveDependencies")
	logger.Info("resolving dependencies")

	return g.Transform(ctx, func(id string, out *graph.Graph) error {
		if id == "root" { // skip root
			return nil
		}

		node, ok := out.Get(id).(*parse.Node)
		if !ok {
			return fmt.Errorf("ResolveDependencies can only be used on Graphs of *parse.Node. I got %T", out.Get(id))
		}

		depGenerators := []dependencyGenerator{
			getDepends,
			getParams,
			func(node *parse.Node) (out []string, err error) {
				return getXrefs(g, node)
			},
		}

		// we have dependencies from various sources, but they're always IDs, so we
		// can connect them pretty easily
		for _, source := range depGenerators {
			deps, err := source(node)
			if err != nil {
				return err
			}
			for _, dep := range deps {

				out.Connect(id, graph.SiblingID(id, dep))
			}
		}
		return nil
	})
}

func getDepends(node *parse.Node) ([]string, error) {
	deps, err := node.GetStringSlice("depends")

	switch err {
	case parse.ErrNotFound:
		return []string{}, nil

	case nil:
		return deps, nil

	default:
		return nil, err
	}
}

func getParams(node *parse.Node) (out []string, err error) {
	var strings []string
	strings, err = node.GetStrings()
	if err != nil {
		return nil, err
	}

	type stub struct{}
	language := extensions.DefaultLanguage()
	language.On("param", extensions.RememberCalls(&out, 0))
	for _, s := range strings {
		useless := stub{}
		tmpl, tmplErr := template.New("DependencyTemplate").Funcs(language.Funcs).Parse(s)
		if tmplErr != nil {
			return out, tmplErr
		}
		tmpl.Execute(ioutil.Discard, &useless)
	}
	for idx, val := range out {
		out[idx] = "param." + val
	}
	return out, err
}

func getXrefs(g *graph.Graph, node *parse.Node) (out []string, err error) {
	var strings []string
	var calls []string
	nodeRefs := make(map[string]struct{})
	strings, err = node.GetStrings()
	if err != nil {
		return nil, err
	}
	language := extensions.DefaultLanguage()
	language.On(extensions.RefFuncName, extensions.RememberCalls(&calls, 0))
	for _, s := range strings {
		tmpl, tmplErr := template.New("DependencyTemplate").Funcs(language.Funcs).Parse(s)
		if tmplErr != nil {
			return out, tmplErr
		}
		tmpl.Execute(ioutil.Discard, &struct{}{})
	}
	for _, call := range calls {
		vertex, _, found := preprocessor.VertexSplit(g, "root/"+call)
		if !found {
			return []string{}, fmt.Errorf("unresolvable call to %s", call)
		}
		vertex = vertex[len("root/"):]
		if _, ok := nodeRefs[vertex]; !ok {
			nodeRefs[vertex] = struct{}{}
			out = append(out, vertex)
		}
	}
	return out, err
}
