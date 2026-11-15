//go:build gts_compact_node_views

package parsercorephase0

// NativeNodeViewsEnabled reports whether this build publishes experimental
// native node views. Qualification switch: ordinary parsing uses the complete
// public adapter, including indexed queries and the exclusive mutation
// boundary.
const NativeNodeViewsEnabled = true
