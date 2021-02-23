// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"io"
	"os"
	"os/exec"

	"sigs.k8s.io/kustomize/kyaml/fn/runtime/backend"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type Exec struct {
	// Path is the path to the executable to run
	Path string `json:"path,omitempty" yaml:"path,omitempty"`
	// Args are the arguments to the executable
	Args []string `json:"args,omitempty" yaml:"args,omitempty"`
}

type Spec struct {
	Exec *Exec `json:"exec,omitempty" yaml:"exec,omitempty"`
}

type Filter struct {
	backend.FunctionFilter

	Spec *Spec
}

func (c *Filter) Filter(nodes []*yaml.RNode) ([]*yaml.RNode, error) {
	c.FunctionFilter.Run = c.Run
	return c.FunctionFilter.Filter(nodes)
}

func (c *Filter) Run(reader io.Reader, writer io.Writer) error {
	cmd := exec.Command(c.Spec.Exec.Path, c.Spec.Exec.Args...)
	cmd.Stdin = reader
	cmd.Stdout = writer
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type FunctionBackendProvider struct {
	// no params so far
}

func (c FunctionBackendProvider) UnmarshalAnnotation(a []byte) (backend.FunctionBackendSpec, error) {
	var s Spec
	if err := yaml.Unmarshal(a, &s); err != nil || s.Exec == nil {
		return nil, err
	}
	return &s, nil
}

func (s *Spec) NewFilter(api *yaml.RNode, globalScope bool, resultsFile string, deferFailure bool) (kio.Filter, error) {
	var f Filter

	f.Spec = s
	f.FunctionConfig = api
	f.GlobalScope = globalScope
	f.ResultsFile = resultsFile
	f.DeferFailure = deferFailure
	return &f, nil
}
