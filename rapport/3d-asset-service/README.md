# Rapport — 3d-asset-service

> SRS v2.0 §3.7 (3D-001 → 008) — **NOUVEAU v2.0** (mesh génératif).

## Mission
Pipeline de génération de mesh 3D **réaliste** pour équipements industriels à
partir des dimensions physiques + spécifications textuelles.

## Pipeline (SRS §3.7.1)
```
input: dimensions L×W×H + type/marque/modèle
         ↓
ResearchAgent (ADK) — images + fiches techniques
         ↓
┌────────────────────────────────────┐
│ Plan Premium  → Meshy.ai API       │  qualité max
│ Plan Standard → Shap-E (OpenAI)    │  open-source
│ Plan Basic    → géométrie procé.   │  fallback toujours OK
└────────────────────────────────────┘
         ↓
Post-processing (normalisation dims, PBR métal/acier)
         ↓
Stockage MinIO (GLB) + URL presignée
         ↓
Frontend useGLTF (React Three Fiber)
```

## Cache (3D-003)
Clé MinIO = `sha256(equipment_type + dims + brand + model + quality_level)`.
Une recherche identique ne relance pas la génération.

## Qualité par plan (3D-008)
| Plan        | Qualité            | Source               |
|-------------|--------------------|----------------------|
| Free / Starter | Procédural      | Box paramétré + détails |
| Business    | Standard           | Shap-E               |
| Enterprise  | Premium            | Meshy.ai             |

## Performance
- Scène R3F < 4 s (3D-001), FPS ≥ 30 (3D-006), LOD au-delà de 20 m.

## Métriques exposées (`/metrics`)
- `mesh_generation_total{quality,cache_hit}`
- `mesh_generation_duration_seconds{quality}`
- `meshy_api_errors_total`
- `minio_storage_bytes` (gauge)
