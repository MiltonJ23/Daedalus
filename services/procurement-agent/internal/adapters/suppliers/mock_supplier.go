package suppliers

import (
	"context"
	"hash/fnv"
	"math"

	"github.com/Daedalus/procurement-agent/internal/core/domain"
	"github.com/Daedalus/procurement-agent/internal/core/ports"
)

// MockSupplier — deterministic in-memory supplier used as a stand-in for real
// catalogs (Alibaba, IndustryStock, DirectIndustry…) until external adapters
// land. It returns reproducible offers based on a hash of the query so the
// service can be exercised end-to-end in dev and tests.
//
// FR-PROC-02 demands at least 3 supplier sources, so wire three distinct
// MockSupplier instances (or a mix of real + mock) in main.go.
type MockSupplier struct {
	name        string
	country     string
	rating      float64
	priceFactor float64
	leadOffset  int
	usdToXAF    float64
}

// NewMockSupplier builds a MockSupplier whose offers tilt toward the given
// rating / price / lead-time profile.
func NewMockSupplier(name, country string, rating, priceFactor float64, leadOffset int, usdToXAF float64) *MockSupplier {
	if usdToXAF == 0 {
		usdToXAF = 600 // sensible XAF/USD fallback
	}
	return &MockSupplier{
		name:        name,
		country:     country,
		rating:      rating,
		priceFactor: priceFactor,
		leadOffset:  leadOffset,
		usdToXAF:    usdToXAF,
	}
}

func (m *MockSupplier) Name() string { return m.name }

func (m *MockSupplier) Search(_ context.Context, q ports.SupplierQuery) ([]domain.EquipmentResult, error) {
	const offersPerSupplier = 3
	out := make([]domain.EquipmentResult, 0, offersPerSupplier)

	seed := hashString(q.Query + "|" + q.Category)

	for i := 0; i < offersPerSupplier; i++ {
		variant := float64(i + 1)
		basePrice := math.Round((5000+float64(seed%4000)+variant*1500)*m.priceFactor*100) / 100
		if q.MaxBudgetUSD > 0 && basePrice > q.MaxBudgetUSD {
			continue
		}

		r := domain.NewEquipmentResult("")
		r.Name = q.Query + " (variant " + intToStr(i+1) + ")"
		r.Model = m.name + "-" + intToStr(int(seed%9000)+1000+i)
		r.Supplier = m.name
		r.SupplierRating = m.rating
		r.PriceUSD = basePrice
		r.PriceXAF = math.Round(basePrice*m.usdToXAF*100) / 100
		r.LeadTimeDays = 14 + m.leadOffset + i*5
		r.SpecMatch = clamp01(0.6 + variant*0.1)
		r.Specifications = map[string]interface{}{
			"category": q.Category,
			"variant":  i + 1,
			"warranty": "12 months",
		}
		r.Dimensions = domain.Dimensions{
			WidthM:  2.0 + variant*0.3,
			DepthM:  1.5 + variant*0.2,
			HeightM: 1.8 + variant*0.1,
		}
		r.PowerKW = 5.0 + variant*1.5
		r.Country = m.country
		out = append(out, r)
	}
	return out, nil
}

// StaticConverter — fixed-rate USD→XAF converter for offline development.
type StaticConverter struct {
	Rate float64
}

func NewStaticConverter(rate float64) *StaticConverter {
	if rate <= 0 {
		rate = 600
	}
	return &StaticConverter{Rate: rate}
}

func (c *StaticConverter) USDToXAF(_ context.Context, amountUSD float64) (float64, error) {
	return math.Round(amountUSD*c.Rate*100) / 100, nil
}

// ── helpers ─────────────────────────────────────────────────────────

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func hashString(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}

func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [12]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
