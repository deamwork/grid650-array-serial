package main

import (
	"github.com/XSAM/go-hybrid/metadata"

	"github.com/deamwork/grid650-array-serial/cmd/runtime"
)

func main() {
	metadata.SetAppName("grid650-array-serial")

	runtime.Start()
}
