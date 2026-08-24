package col

import "errors"

var (
	ErrBlockNotFound  = errors.New("column block not found")
	ErrBlockSealed    = errors.New("column block is sealed")
	ErrBlockReclaimed = errors.New("column block has been reclaimed")
	ErrSegmentOpen    = errors.New("segment file could not be opened")
	ErrSegmentMissing = errors.New("segment file is missing")
	ErrEmptyColumn    = errors.New("column has no values")
)
