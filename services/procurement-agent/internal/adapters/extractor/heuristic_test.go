package extractor

import (
	"context"
	"testing"
)

func TestHeuristic_Extract(t *testing.T) {
	ex := NewHeuristicExtractor()
	spec, err := ex.Extract(context.Background(),
		"Need a 5-axis CNC milling machine 3.5 x 2.0 x 1.8 m, 15kW, 380V, stainless steel, under $50000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if spec["axes"] != 5 {
		t.Errorf("expected axes=5, got %v", spec["axes"])
	}
	if spec["power_kw"].(float64) != 15 {
		t.Errorf("expected power_kw=15, got %v", spec["power_kw"])
	}
	if spec["voltage_v"] != 380 {
		t.Errorf("expected voltage_v=380, got %v", spec["voltage_v"])
	}
	if spec["material"] != "stainless steel" {
		t.Errorf("expected stainless steel, got %v", spec["material"])
	}
	if spec["category"] != "machining" {
		t.Errorf("expected machining, got %v", spec["category"])
	}
	if spec["max_budget_usd"].(float64) != 50000 {
		t.Errorf("expected max_budget_usd=50000, got %v", spec["max_budget_usd"])
	}
	dims, ok := spec["dimensions"].(map[string]float64)
	if !ok {
		t.Fatalf("expected dimensions map, got %T", spec["dimensions"])
	}
	if dims["width_m"] != 3.5 || dims["depth_m"] != 2.0 || dims["height_m"] != 1.8 {
		t.Errorf("unexpected dims: %+v", dims)
	}
}

func TestHeuristic_HPConvertedToKW(t *testing.T) {
	ex := NewHeuristicExtractor()
	spec, _ := ex.Extract(context.Background(), "20 hp compressor")
	if v, ok := spec["power_kw"].(float64); !ok || v < 14 || v > 16 {
		t.Errorf("expected ~14.91 kW, got %v", spec["power_kw"])
	}
}
