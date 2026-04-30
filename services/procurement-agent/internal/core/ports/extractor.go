package ports

import "context"

// SpecExtractor parses a free-text equipment description (FR-PROC-01)
// into a structured specification map. Implementations can be:
//   - heuristic / regex-based (offline, deterministic — see adapters/extractor)
//   - LLM-backed (OpenAI, Anthropic, …) for richer NL understanding
type SpecExtractor interface {
	Extract(ctx context.Context, naturalLanguageQuery string) (map[string]interface{}, error)
}
