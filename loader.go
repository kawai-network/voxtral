package voxtral

import (
	"fmt"
	"os"
	"runtime"

	"github.com/ebitengine/purego"
)

type LibFuncs struct {
	FuncPtr any
	Name    string
}

// Init loads the shared library and registers the function pointers.
// If libPath is empty, it attempts to find the library in the current directory
// based on the OS (libgovoxtral.dylib or libgovoxtral.so).
func Init(libPath string) error {
	if libPath == "" {
		if runtime.GOOS == "darwin" {
			libPath = "./libgovoxtral.dylib"
		} else {
			libPath = "./libgovoxtral.so"
		}
	}

	// Check if file exists
	if _, err := os.Stat(libPath); os.IsNotExist(err) {
		return fmt.Errorf("library not found at %s", libPath)
	}

	gosd, err := purego.Dlopen(libPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return fmt.Errorf("failed to load library: %w", err)
	}

	libFuncs := []LibFuncs{
		{&CppLoadModel, "load_model"},
		{&CppTranscribe, "transcribe"},
		{&CppFreeResult, "free_result"},
	}

	for _, lf := range libFuncs {
		purego.RegisterLibFunc(lf.FuncPtr, gosd, lf.Name)
	}

	return nil
}
