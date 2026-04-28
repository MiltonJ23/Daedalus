# Rapport — orchestrator-agent

> SRS v2.0 §3.5, §5.2 — Agent racine ADK.

## Mission
Décompose les objectifs utilisateur en sous-tâches et délègue aux sub-agents
(`procurement-agent`, `layout-engine` agent, `3d-asset-service` agent).
Enforced les **quotas runs IA** par plan (PROC-006) avant chaque session.

## Architecture multi-agents (SRS §5.2.2)
```
OrchestratorAgent (root)
├── ProcurementAgent
│   ├── WebSearchTool   (Playwright / SerpAPI)
│   ├── PricingTool     (exchange rates)
│   └── CacheTool       (Redis)
├── LayoutAgent
│   ├── SpatialOptimizerTool (NetworkX)
│   └── CollisionDetectionTool
└── AssetGenerationAgent
    ├── WebResearchTool
    ├── MeshyAITool
    ├── ShapETool       (fallback)
    └── MinIOStorageTool
```

## Pattern
- Sessions ADK sérialisées en **Redis** (`session:adk:<session_id>`).
- Communication inter-agents via **Redis Streams** (PROC-005) puis **RabbitMQ**
  pour les events sortants (`agent.step`, `agent.completed`).
- Streaming temps réel vers le frontend via **SSE** (`/agent/run`) ou WebSocket
  (`notification-service`).

## Quotas (PROC-006)
Avant `AgentSession.run()` :
```python
quota = await redis.get(f"plan:quota:{user_id}:runs")
plan_limit = PLAN_LIMITS[user_plan]
if quota and int(quota) >= plan_limit:
    raise QuotaExceeded()
```

## Métriques exposées (`/metrics`)
- `agent_sessions_total{status}`
- `agent_step_duration_seconds{agent,tool}`
- `agent_quota_exceeded_total{plan}`
- `llm_tokens_total{model,direction}`
