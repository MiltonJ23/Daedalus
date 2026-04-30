package services

import (
	"math"
	"sort"

	"github.com/Daedalus/procurement-agent/internal/core/domain"
)

// Composite ranking weights (FR-PROC-03).
const (
	WeightPrice    = 0.40
	WeightSpec     = 0.35
	WeightSupplier = 0.15
	WeightLeadTime = 0.10
)

// rankResults applies the composite scoring algorithm, mutating Score in place
// and returning the slice sorted from best (highest score) to worst.
//
// Each component is normalised to [0,1] where 1 == best:
//   - price:        cheapest gets 1, most expensive gets 0 (linear)
//   - spec match:   already in [0,1] as supplied by the catalog adapter
//   - supplier:     rating in [0,5] → divided by 5
//   - lead time:    fastest gets 1, slowest gets 0 (linear)
func rankResults(results []domain.EquipmentResult) []domain.EquipmentResult {
	if len(results) == 0 {
		return results
	}

	minPrice, maxPrice := results[0].PriceUSD, results[0].PriceUSD
	minLead, maxLead := results[0].LeadTimeDays, results[0].LeadTimeDays
	for _, r := range results[1:] {
		if r.PriceUSD < minPrice {
			minPrice = r.PriceUSD
		}
		if r.PriceUSD > maxPrice {
			maxPrice = r.PriceUSD
		}
		if r.LeadTimeDays < minLead {
			minLead = r.LeadTimeDays
		}
		if r.LeadTimeDays > maxLead {
			maxLead = r.LeadTimeDays
		}
	}

	priceSpread := maxPrice - minPrice
	leadSpread := float64(maxLead - minLead)

	for i := range results {
		r := &results[i]

		priceN := 1.0
		if priceSpread > 0 {
			priceN = 1.0 - (r.PriceUSD-minPrice)/priceSpread
		}

		specN := clamp01(r.SpecMatch)
		supplierN := clamp01(r.SupplierRating / 5.0)

		leadN := 1.0
		if leadSpread > 0 {
			leadN = 1.0 - (float64(r.LeadTimeDays)-float64(minLead))/leadSpread
		}

		r.Score = round4(
			WeightPrice*priceN +
				WeightSpec*specN +
				WeightSupplier*supplierN +
				WeightLeadTime*leadN,
		)
	}

	sort.SliceStable(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	return results
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func round4(v float64) float64 {
	return math.Round(v*10000) / 10000
}
