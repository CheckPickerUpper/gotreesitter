package recoveryturn

// MissingVersionTurnsForIncludedRanges narrows an existing certification only
// when the grammar's external scanner records an included-range mode
// exclusion. An absent hook preserves the caller's existing policy; this
// function cannot grant certification when certified is false.
func MissingVersionTurnsForIncludedRanges(certified bool, scanner any) bool {
	if !certified {
		return false
	}
	exclusion, ok := scanner.(interface {
		ExcludesMissingVersionTurnsForIncludedRanges() bool
	})
	return !ok || !exclusion.ExcludesMissingVersionTurnsForIncludedRanges()
}
