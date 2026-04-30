// Package extractor implements SpecExtractor adapters (FR-PROC-01).
//
// HeuristicExtractor is a deterministic, offline implementation that pulls
// numeric quantities, units, materials and capacity hints out of a free-text
// equipment description using regular expressions. It serves as the default
// adapter so the service runs without any LLM API key, and as the fallback
// when an LLM call fails.
package extractor

import (
	"context"
	"regexp"
	"strconv"
	"strings"
)

type HeuristicExtractor struct{}

func NewHeuristicExtractor() *HeuristicExtractor { return &HeuristicExtractor{} }

var (
	reCapacity   = regexp.MustCompile(`(?i)(\d+(?:\.\d+)?)\s*(l|liters?|litres?|kg|tons?|tonnes?|t/h|kg/h|units?/h|pcs/h)`)
	rePower      = regexp.MustCompile(`(?i)(\d+(?:\.\d+)?)\s*(kw|kilowatts?|hp|horsepower|w\b)`)
	reVoltage    = regexp.MustCompile(`(?i)(\d{2,4})\s*(v|volts?)\b`)
	reDimensions = regexp.MustCompile(`(?i)(\d+(?:\.\d+)?)\s*(?:m|metres?|meters?)?\s*[x×]\s*(\d+(?:\.\d+)?)\s*(?:m|metres?|meters?)?(?:\s*[x×]\s*(\d+(?:\.\d+)?)\s*(?:m|metres?|meters?)?)?`)
	reAxes       = regexp.MustCompile(`(?i)(\d)\s*-?\s*axis`)
	reBudget     = regexp.MustCompile(`(?i)(?:under|below|less than|max(?:imum)?|<=?)\s*\$?\s*(\d{3,8})`)
)

var materialKeywords = []string{"stainless steel", "carbon steel", "aluminium", "aluminum", "cast iron", "plastic", "ceramic"}

var categoryKeywords = map[string]string{
	"cnc":          "machining",
	"milling":      "machining",
	"lathe":        "machining",
	"injection":    "molding",
	"mixer":        "food-processing",
	"oven":         "food-processing",
	"packaging":    "packaging",
	"conveyor":     "material-handling",
	"forklift":     "material-handling",
	"compressor":   "utilities",
	"generator":    "utilities",
	"3d printer":   "additive-manufacturing",
	"laser cutter": "cutting",
}

func (h *HeuristicExtractor) Extract(_ context.Context, query string) (map[string]interface{}, error) {
	q := strings.ToLower(strings.TrimSpace(query))
	out := map[string]interface{}{
		"raw_query": query,
	}

	if m := reCapacity.FindStringSubmatch(q); m != nil {
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			out["capacity_value"] = v
			out["capacity_unit"] = strings.TrimSpace(m[2])
		}
	}
	if m := rePower.FindStringSubmatch(q); m != nil {
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			unit := strings.ToLower(m[2])
			switch unit {
			case "hp", "horsepower":
				out["power_kw"] = round2(v * 0.7457)
			case "w":
				out["power_kw"] = round2(v / 1000)
			default:
				out["power_kw"] = v
			}
		}
	}
	if m := reVoltage.FindStringSubmatch(q); m != nil {
		if v, err := strconv.Atoi(m[1]); err == nil {
			out["voltage_v"] = v
		}
	}
	if m := reDimensions.FindStringSubmatch(q); m != nil {
		dims := map[string]float64{}
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			dims["width_m"] = v
		}
		if v, err := strconv.ParseFloat(m[2], 64); err == nil {
			dims["depth_m"] = v
		}
		if m[3] != "" {
			if v, err := strconv.ParseFloat(m[3], 64); err == nil {
				dims["height_m"] = v
			}
		}
		if len(dims) > 0 {
			out["dimensions"] = dims
		}
	}
	if m := reAxes.FindStringSubmatch(q); m != nil {
		if v, err := strconv.Atoi(m[1]); err == nil {
			out["axes"] = v
		}
	}
	if m := reBudget.FindStringSubmatch(q); m != nil {
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			out["max_budget_usd"] = v
		}
	}
	for _, mat := range materialKeywords {
		if strings.Contains(q, mat) {
			out["material"] = mat
			break
		}
	}
	for kw, cat := range categoryKeywords {
		if strings.Contains(q, kw) {
			out["category"] = cat
			out["equipment_type"] = kw
			break
		}
	}

	return out, nil
}

func round2(v float64) float64 {
	return float64(int(v*100+0.5)) / 100
}
