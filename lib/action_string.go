// Code generated by "stringer -type=Action"; DO NOT EDIT.

package lib

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[RestoreObject-0]
	_ = x[RestoreMetadata-1]
	_ = x[Delete-2]
}

const _Action_name = "RestoreObjectRestoreMetadataDelete"

var _Action_index = [...]uint8{0, 13, 28, 34}

func (i Action) String() string {
	if i < 0 || i >= Action(len(_Action_index)-1) {
		return "Action(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Action_name[_Action_index[i]:_Action_index[i+1]]
}
