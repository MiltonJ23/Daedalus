# Rapport — procurement-agent

> SRS v2.0 §3.5 (PROC-001 → 006).

## Mission
Sourcing IA mondial d'équipements industriels. Parse une description NL en
spécifications structurées, recherche sur Alibaba / Made-in-China / TradeIndia
via Playwright, classe ≥ 5 résultats, applique les taux de change XAF/USD,
gère un cache Redis 24 h.

## Pipeline
```
NL prompt → LLM parse (PROC-001) → EquipmentSpec
                                       ↓
                        Cache Redis ? oui → return
                                       ↓ non
                        WebSearchTool (Playwright headless)
                                       ↓
                        RankingTool (prix, MOQ, lead time, rating)
                                       ↓
                        PricingTool (FX XAF/USD temps réel)
                                       ↓
                        Persist + cache 24 h + return
```

## API
- `POST /procurement/parse` — NL → `EquipmentSpec`
- `POST /procurement/search` — `EquipmentSpec` → résultats
- `POST /procurement/results/{id}/approve` — approval/reject

## Performances visées
- ≥ 5 résultats par requête (PROC-002)
- Cache hit < 50 ms ; recherche fraîche < 30 s (NFR-P01)

## Tests
- Tools atomiques testables (mock Playwright via `pytest-playwright`).
- Snapshots de réponses scrapées dans `tests/fixtures/`.

## Métriques exposées (`/metrics`)
- `procurement_searches_total{cache_hit}`
- `procurement_search_duration_seconds`
- `procurement_results_per_query` (histogram)
- `procurement_approvals_total{decision}`
