package models

// GrepFlags represents the command-line flags for the grep command
type GrepFlags struct {
	FixedString  bool `json:"fixed_string"`
	PrintNumbers bool `json:"print_numbers"`
	IgnoreCase   bool `json:"ignore_case"`
	Invert       bool `json:"invert"`
	After        int  `json:"after"`
	Before       int  `json:"before"`
	CountOnly    bool `json:"count_only"`
}

// Task represents a grep task with its pattern, lines, and context
type Task struct {
	ID              int       `json:"id"`
	Pattern         string    `json:"pattern"`
	Lines           []string  `json:"lines"`
	BeforeContext   []string  `json:"before_context"`
	AfterContext    []string  `json:"after_context"`
	StartLineNumber int       `json:"start_line_number"`
	Flags           GrepFlags `json:"flags"`
}

// Result represents the result of a grep task
type Result struct {
	TaskID      int          `json:"task_id"`
	FoundBlocks []FoundBlock `json:"found_blocks"`
}

// FoundBlock represents a found block of lines
type FoundBlock struct {
	StartLineNumber int      `json:"start_line_number"`
	Lines           []string `json:"lines"`
}

// ParsedAddr represents a parsed address with its scheme, host, and port
type ParsedAddr struct {
	Raw    string // the original input
	Scheme string // http, https, ftp, etc.
	Host   string
	Port   string
}
