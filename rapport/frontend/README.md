# Rapport — frontend

> SRS v2.0 §7.1 — Interface utilisateur (SPA Next.js + R3F).

## Vues principales (SRS §7.1)
1. Landing Page — pricing + CTA
2. Auth — login/register + OAuth Google/GitHub
3. Dashboard — projets, statut, **quotas restants du plan**
4. Pricing — tableau comparatif, intégration Paystack checkout
5. Project Config — floor plan, industrie, budget
6. Procurement Workspace — input NL, résultats IA, approve/reject
7. 3D Factory Viewer — R3F, orbit, walkthrough, raycasting, layer toggles
8. Cost Analysis — BOM, ROI chart, exports
9. Notifications — cloche, historique, préférences
10. Account / Billing — abonnement, paiements, membres workspace
11. Admin Panel — utilisateurs, agents config, monitoring

## Stack
- Next.js 14 (App Router) + React 18 + TypeScript 5
- Tailwind CSS 3 + framer-motion + lucide-react
- React Three Fiber + drei (3D viewer)
- Zustand (state) + Axios + React Query (HTTP)

## Build & qualité
```bash
cd frontend
npm install
npm run lint
npm run build   # next build → output .next standalone
```

Lighthouse score cible ≥ 75 (NFR-P02).

## Tests E2E
`e2e-tests/` — Playwright + Jest (login → projet → 3D → export).
