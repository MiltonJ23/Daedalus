import type { Project, ProjectCreate, ProjectUpdate, AutoSaveResponse } from "@/types/project";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api";

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${API_URL}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });
  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.error || body.errors || `HTTP ${res.status}`);
  }
  return res.json();
}

export const api = {
  projects: {
    list: (status: string = "active") =>
      request<Project[]>(`/projects?status=${status}`),

    get: (id: string) =>
      request<Project>(`/projects/${id}`),

    create: (data: ProjectCreate) =>
      request<Project>("/projects", {
        method: "POST",
        body: JSON.stringify(data),
      }),

    update: (id: string, data: ProjectUpdate) =>
      request<Project>(`/projects/${id}`, {
        method: "PUT",
        body: JSON.stringify(data),
      }),

    autosave: (id: string, data: ProjectUpdate) =>
      request<AutoSaveResponse>(`/projects/${id}/autosave`, {
        method: "PATCH",
        body: JSON.stringify(data),
      }),

    archive: (id: string) =>
      request<Project>(`/projects/${id}/archive`, { method: "PATCH" }),

    restore: (id: string) =>
      request<Project>(`/projects/${id}/archive?action=restore`, { method: "PATCH" }),

    delete: (id: string) =>
      request<{ message: string }>(`/projects/${id}?confirm=true`, { method: "DELETE" }),
  },
};
