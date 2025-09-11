package config

type Config struct {
	Context    []string `json:"context"`
	AddIndexes bool     `json:"indexes"`     // always take the value from CLI even if not set
	SearchTerm string   `json:"search_term"` // due to mandatory flag req on CLI, not currently be used
	FilePath   string   `json:"file_path"`
}
