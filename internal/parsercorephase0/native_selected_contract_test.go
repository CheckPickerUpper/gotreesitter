package parsercorephase0

import (
	"encoding/hex"
	"testing"
)

func TestNativeSelectedSnapshotContracts(t *testing.T) {
	for _, tc := range []struct {
		digest    string
		byDefault bool
	}{
		{"4456f9b11706face73bd511d6e80bbd2f0d468055065296db6f5334fa101d379", true},
		{"214322327c9ee7dda29febf32e415f7ed9e00a90b7690d3c32522e393211b8da", true},
		{"df63fc35604c4e4e7a484abde9eb2110b61640045601c23991723f323a48310d", true},
		{"be6bd2f7c9d14160d8f74b29d63dc9a9294bc03ac1f1c4bde9522273f7c87b46", false},
	} {
		bytes, err := hex.DecodeString(tc.digest)
		if err != nil {
			t.Fatal(err)
		}
		identity := [32]byte(bytes)
		if !NativeSelectedSnapshotContract(identity) || NativeSelectedSnapshotDefaultContract(identity) != tc.byDefault {
			t.Fatalf("unexpected publisher contract for %x", identity)
		}
		for i := range identity {
			changed := identity
			changed[i] ^= 1
			if NativeSelectedSnapshotContract(changed) || NativeSelectedSnapshotDefaultContract(changed) {
				t.Fatalf("unqualified table change admitted: %x", changed)
			}
		}
	}
}
