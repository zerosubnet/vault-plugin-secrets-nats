//go:build generate
// +build generate

package claims

// Generate deepcopy methodsets
//go:generate go run -tags generate sigs.k8s.io/controller-tools/cmd/controller-gen object:headerFile=../../hack/boilerplate.go.txt paths=./...

import (
	_ "sigs.k8s.io/controller-tools/cmd/controller-gen" //nolint:typecheck
)