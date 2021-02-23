package backend

import (
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type FunctionBackendProvider interface {
	UnmarshalAnnotation(a []byte) (FunctionBackendSpec, error)
}

type FunctionBackendSpec interface {
	NewFilter(
		api *yaml.RNode,
		globalScope bool,
		resultsFile string,
		deferFailure bool) (kio.Filter, error)
}
