"use client";

import { useState, useEffect, useCallback } from "react";
import type { Project } from "@/types/project";
import { api } from "@/lib/api";

export function useProjects(initialStatus: string = "active") {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [status, setStatus] = useState(initialStatus);

  const fetchProjects = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await api.projects.list(status);
      setProjects(data);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Échec du chargement des projets");
    } finally {
      setLoading(false);
    }
  }, [status]);

  useEffect(() => {
    fetchProjects();
  }, [fetchProjects]);

  return { projects, loading, error, status, setStatus, refresh: fetchProjects };
}
