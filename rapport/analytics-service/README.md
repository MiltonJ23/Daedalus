# Rapport — analytics-service

> SRS v2.0 §3.8 (ANA-001 → 004).

## Mission
Modélisation financière : breakdown coûts (équipement + shipping +
installation + contingency 10%) en XAF & USD, projection ROI 5 ans,
exports PDF/CSV (BOM + rapport projet complet).

## Endpoints
- `GET  /projects/:id/costs` — breakdown (ANA-001)
- `POST /projects/:id/roi`  — projection 5 ans (ANA-003)
- `GET  /projects/:id/bom.pdf` — export BOM (ANA-002, ANA-004)
- `GET  /projects/:id/bom.csv` — export CSV

## Modèle ROI
- CAPEX = total breakdown
- OPEX = paramètre utilisateur
- Cashflow annuel = volume × prix - OPEX
- NPV @ taux d'actualisation paramétrable
- IRR + payback période en mois

## Conversion devise (NFR-C01)
XAF primaire, USD calculé via taux de change cache 1 h (Exchange Rate API).
Watermark Daedalus sur exports si plan Free (cf. SRS §3.4.3).

## Métriques exposées (`/metrics`)
- `cost_calculations_total`
- `roi_projections_total`
- `exports_total{format,plan}`
- `fx_rate_xaf_usd` (gauge)
