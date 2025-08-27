package models

// GrepFlags is the struct for the grep flags
type GrepFlags struct {
	FixedString  bool `json:"fixed_string"`
	PrintNumbers bool `json:"print_numbers"`
	IgnoreCase   bool `json:"ignore_case"`
	Invert       bool `json:"invert"`
	After        int  `json:"after"`
	Before       int  `json:"before"`
	CountOnly    bool `json:"count_only"`
}

// Request is the struct for the request
type Request struct {
	ID              int       `json:"id"`
	Pattern         string    `json:"pattern"`
	Lines           []string  `json:"lines"`
	BeforeContext   []string  `json:"before_context"`
	AfterContext    []string  `json:"after_context"`
	StartLineNumber int       `json:"start_line_number"`
	Flags           GrepFlags `json:"flags"`
}

// Response is the struct for the response
type Response struct {
	TaskID      int          `json:"task_id"`
	FoundBlocks []FoundBlock `json:"found_blocks"`
}

// FoundBlock is the struct for the found block
type FoundBlock struct {
	StartLineNumber int      `json:"start_line_number"`
	Lines           []string `json:"lines"`
}
