package models

type GrepFlags struct {
	FixedString  bool `json:"fixed_string"`
	PrintNumbers bool `json:"print_numbers"`
	IgnoreCase   bool `json:"ignore_case"`
	Invert       bool `json:"invert"`
	After        int  `json:"after"`
	Before       int  `json:"before"`
	CountOnly    bool `json:"count_only"`
}

type Task struct {
	ID              int       `json:"id"`
	Pattern         string    `json:"pattern"`
	Lines           []string  `json:"lines"`
	BeforeContext   []string  `json:"before_context"`
	AfterContext    []string  `json:"after_context"`
	StartLineNumber int       `json:"start_line_number"`
	Flags           GrepFlags `json:"flags"`
}

type Result struct {
	TaskID      int          `json:"task_id"`
	FoundBlocks []FoundBlock `json:"found_blocks"`
}

type FoundBlock struct {
	StartLineNumber int      `json:"start_line_number"`
	Lines           []string `json:"lines"`
}

type ParsedAddr struct {
	Raw    string // the original input
	Scheme string // http, https, ftp, etc.
	Host   string
	Port   string
}
