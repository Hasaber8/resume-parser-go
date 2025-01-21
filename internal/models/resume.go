package models

type Resume struct {
	Raw      map[string]string
	Sections map[string]Sections
	Metadata map[string]string
}

type Sections struct {
	Type    SectionType
	Content interface{}
}

type SectionType string

// most common resume sections
const (
	ContactSection  SectionType = "contact"
	TimelineSection SectionType = "timeline"
	ListSection     SectionType = "list"
	FreeformSection SectionType = "freeform"
)

// store generic contact info
type ContactContent struct {
	Name     string
	Email    []string
	Number   []string
	Location string
	Social   map[string]string
}

type TimelineContent struct {
	Entries []TimelineEntry
}

type TimelineEntry struct {
	Organization string
	Location     string
	Title        string
	StartDate    string
	EndDate      string
	Details      []string
	Metadata     map[string]string // for weird resume formats
}

// List section specific structures (for skills, etc.)
type ListContent struct {
	Categories []ListCategory // array of multiple comma seperated skills
}

type ListCategory struct {
	Name  string
	Items []string
}

// Freeform section for any unstructured content
type FreeformContent struct {
	Entries []FreeformEntry
}

type FreeformEntry struct {
	Heading string
	Content []string
}
