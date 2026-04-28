# Rapport — layout-engine

> SRS v2.0 §3.6 (LAY-001 → 006).

## Mission
Optimisation spatiale des équipements industriels sur le floor plan.
Algorithme NetworkX : graphe de flux matière, minimisation des distances,
clearances de sécurité 1.2 m enforcés, détection de collisions, génération
d'un manifeste JSON consommable par le frontend R3F.

## Contraintes (LAY-002)
- Clearance 1.2 m autour de chaque équipement → **bloquant** si violé.
- Pas de chevauchement de bounding box.
- Respect de l'ordre du flux production passé en entrée.

## Algorithme
1. Construction du graphe `G = (équipements, edges flux)`.
2. Placement initial : Kamada-Kawai (NetworkX) projeté sur la grille floor.
3. Recuit simulé : minimisation `Σ distance(flux) + λ × violations`.
4. Validation : aucun overlap, clearances OK → manifeste accepté.
5. Plan Business+ : génération de 3 variantes (LAY-004).

## Manifeste de sortie (LAY-006)
```json
{
  "id": "uuid",
  "score": 142.7,
  "placements": [
    { "equipment_id": "uuid", "position": [x,y,z], "rotation_deg": 90 }
  ]
}
```

## Performance
< 60 s pour ≤ 30 équipements (LAY-001) — sinon retour HTTP 202 + job async.

## Métriques exposées (`/metrics`)
- `layout_generation_duration_seconds`
- `layout_collisions_total`
- `layout_clearance_violations_total`
- `layout_variants_generated_total{plan}`
