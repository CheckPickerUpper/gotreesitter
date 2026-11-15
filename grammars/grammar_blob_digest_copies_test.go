package grammars

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestNoFileHoldsAGrammarBlobDigest refuses every tracked file that holds a
// copy of a shipped grammar blob's SHA-256. Code, tests, receipts and docs
// read a blob's digest from the blob they load; a copy goes stale the moment
// the blob is regenerated, and anything gated on the copy then changes
// behaviour without an error. The search covers the hex spelling and the Go
// byte-literal spelling of each digest.
func TestNoFileHoldsAGrammarBlobDigest(t *testing.T) {
	blobPaths, err := filepath.Glob(filepath.Join("grammar_blobs", "*.bin"))
	if err != nil {
		t.Fatal(err)
	}
	if len(blobPaths) == 0 {
		t.Fatal("no grammar blobs under grammar_blobs")
	}
	arguments := []string{"-C", "..", "grep", "-I", "-i", "-n", "-F"}
	for _, blobPath := range blobPaths {
		blob, err := os.ReadFile(blobPath)
		if err != nil {
			t.Fatal(err)
		}
		digest := sha256.Sum256(blob)
		arguments = append(arguments, "-e", hex.EncodeToString(digest[:]), "-e", goByteLiteralPrefix(digest[:8]))
	}
	arguments = append(arguments, "--", ".", ":(exclude)grammars/grammar_blobs")
	output, err := exec.Command("git", arguments...).Output()
	var exitError *exec.ExitError
	switch {
	case err == nil:
		t.Fatalf("these lines hold a copy of a grammar blob's SHA-256; read the digest from the loaded blob instead:\n%s", output)
	case errors.As(err, &exitError) && exitError.ExitCode() == 1:
		// git grep exits 1 when no line matches.
	default:
		t.Fatalf("git grep for grammar blob digests: %v\n%s", err, stderrOf(err))
	}
}

// goByteLiteralPrefix spells bytes the way gofmt writes a byte-array literal,
// such as "0x44, 0x56, 0xf9".
func goByteLiteralPrefix(bytes []byte) string {
	spelled := make([]string, len(bytes))
	for index, value := range bytes {
		spelled[index] = fmt.Sprintf("0x%02x", value)
	}
	return strings.Join(spelled, ", ")
}

// stderrOf returns the standard error a failed command wrote, when the error
// carries it.
func stderrOf(err error) string {
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		return string(exitError.Stderr)
	}
	return ""
}
