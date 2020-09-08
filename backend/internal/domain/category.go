package domain

// Category represents 項目 and its 階層
type Category struct {
	// ID
	ID string

	// Name
	Name string

	// 階層
	Level CategoryLevel

	// 表示順
	DisplayOrder int

	// 親CategoryID
	ParentID string
}

// CategoryLevel represents category hierarchy level.
type CategoryLevel int

const (
	// CategoryLevel1 is the top level category.
	CategoryLevel1 CategoryLevel = iota + 1
	// CategoryLevel2 is the second level category.
	CategoryLevel2
)
