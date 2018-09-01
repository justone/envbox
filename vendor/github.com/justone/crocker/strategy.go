package crocker

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Stategy is an interface for finding the right credential helper.
type Strategy interface {
	Helper() (string, error)
}

// ProgPrefix is prefix that is prepended to all command names before checking
// for existence.
var ProgPrefix = "docker-credential-"

// progWith prefixes with the ProgPrefix
func progWith(suffix string) string {
	return fmt.Sprintf("%s%s", ProgPrefix, suffix)
}

// binExists checks if a binary exists somewhere
func binExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// StockStrategy is a strategy for finding helpers.
type StockStrategy struct{}

// Helper in StockStrategy checks for the existence of the appropriate platform
// helper.  It only checks for the stock, upstream, helpers.
func (s StockStrategy) Helper() (string, error) {
	var prog string

	if runtime.GOOS == "linux" && binExists(progWith("secretservice")) {
		prog = progWith("secretservice")
	} else if runtime.GOOS == "darwin" && binExists(progWith("osxkeychain")) {
		prog = progWith("osxkeychain")
	} else if runtime.GOOS == "windows" && binExists(progWith("wincred")) {
		prog = progWith("wincred")
	}

	if len(prog) == 0 {
		return "", fmt.Errorf("unable to find helper")
	}

	return prog, nil
}

// MemThenStockStrategy is a strategy for finding helpers.
type MemThenStockStrategy struct{}

// Helper in MemThenStockStrategy checks for the cachemem (in memory store)
// helper and then falls back to the StockStrategy.
func (s MemThenStockStrategy) Helper() (string, error) {
	if binExists(progWith("cachemem")) {
		return progWith("cachemem"), nil
	}

	ss := StockStrategy{}
	return ss.Helper()
}
