/**
 * Daedalus E2E Test Strategy
 * 
 * These tests verify the full stack: Frontend → API → Database
 * Run with: npx jest e2e-tests/ (or integrate with Jenkins pipeline)
 */

const API_URL = process.env.API_URL || "http://localhost:5000/api";

async function request(path, options = {}) {
  const res = await fetch(`${API_URL}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });
  return { status: res.status, body: await res.json().catch(() => ({})) };
}

// ── CRUD Flow ──────────────────────────────────

describe("Project Management E2E", () => {
  let projectId;

  test("POST /projects — Create project", async () => {
    const { status, body } = await request("/projects", {
      method: "POST",
      body: JSON.stringify({
        name: "E2E Test Factory",
        industry_type: "Agroalimentaire",
        location: "Douala, Cameroun",
        budget: 5000000,
        floor_width: 200,
        floor_depth: 150,
      }),
    });
    expect(status).toBe(201);
    expect(body.name).toBe("E2E Test Factory");
    expect(body.status).toBe("active");
    projectId = body.id;
  });

  test("GET /projects — List includes created project", async () => {
    const { status, body } = await request("/projects");
    expect(status).toBe(200);
    expect(body.some((p) => p.id === projectId)).toBe(true);
  });

  test("GET /projects/:id — Get single project", async () => {
    const { status, body } = await request(`/projects/${projectId}`);
    expect(status).toBe(200);
    expect(body.id).toBe(projectId);
    expect(body.floor_width).toBe(200);
  });

  test("PATCH /projects/:id/autosave — Auto-save updates version", async () => {
    const { status, body } = await request(`/projects/${projectId}/autosave`, {
      method: "PATCH",
      body: JSON.stringify({ budget: 6000000 }),
    });
    expect(status).toBe(200);
    expect(body.message).toBe("Auto-saved");
    expect(body.version).toBe(2);
  });

  test("PUT /projects/:id — Full update", async () => {
    const { status, body } = await request(`/projects/${projectId}`, {
      method: "PUT",
      body: JSON.stringify({ name: "E2E Updated Factory" }),
    });
    expect(status).toBe(200);
    expect(body.name).toBe("E2E Updated Factory");
    expect(body.version).toBe(3);
  });

  test("PATCH /projects/:id/archive — Archive project", async () => {
    const { status, body } = await request(`/projects/${projectId}/archive`, {
      method: "PATCH",
    });
    expect(status).toBe(200);
    expect(body.status).toBe("archived");
    expect(body.is_archived).toBe(true);
  });

  test("PATCH /projects/:id/archive?action=restore — Restore", async () => {
    const { status, body } = await request(
      `/projects/${projectId}/archive?action=restore`,
      { method: "PATCH" }
    );
    expect(status).toBe(200);
    expect(body.status).toBe("active");
  });

  test("DELETE /projects/:id — Requires confirmation", async () => {
    const { status } = await request(`/projects/${projectId}`, {
      method: "DELETE",
    });
    expect(status).toBe(400);
  });

  test("DELETE /projects/:id?confirm=true — Permanent delete", async () => {
    const { status } = await request(`/projects/${projectId}?confirm=true`, {
      method: "DELETE",
    });
    expect(status).toBe(200);

    const { status: getStatus } = await request(`/projects/${projectId}`);
    expect(getStatus).toBe(404);
  });
});

// ── Validation ─────────────────────────────────

describe("Validation E2E", () => {
  test("POST with missing fields → 422", async () => {
    const { status } = await request("/projects", {
      method: "POST",
      body: JSON.stringify({ name: "Incomplete" }),
    });
    expect(status).toBe(422);
  });

  test("POST with zero dimensions → 422", async () => {
    const { status } = await request("/projects", {
      method: "POST",
      body: JSON.stringify({
        name: "Bad",
        industry_type: "BTP",
        location: "Yaoundé",
        budget: 100,
        floor_width: 0,
        floor_depth: 50,
      }),
    });
    expect(status).toBe(422);
  });
});

// ── Health Check ───────────────────────────────

describe("System E2E", () => {
  test("GET /health — Returns healthy", async () => {
    const res = await fetch(`${API_URL.replace("/api", "")}/health`);
    const body = await res.json();
    expect(res.status).toBe(200);
    expect(body.status).toBe("healthy");
  });
});
