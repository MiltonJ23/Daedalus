"use client";

import { useEffect, useRef, useState, useCallback } from "react";
import type { ProjectUpdate } from "@/types/project";
import { api } from "@/lib/api";

interface UseAutoSaveOptions {
  projectId: string;
  data: ProjectUpdate;
  intervalMs?: number;
  enabled?: boolean;
}

export function useAutoSave({
  projectId,
  data,
  intervalMs = 60_000,
  enabled = true,
}: UseAutoSaveOptions) {
  const [lastSavedAt, setLastSavedAt] = useState<Date | null>(null);
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const previousData = useRef<string>("");
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);

  const save = useCallback(async () => {
    const serialized = JSON.stringify(data);
    if (serialized === previousData.current) return;

    setIsSaving(true);
    setError(null);
    try {
      const result = await api.projects.autosave(projectId, data);
      previousData.current = serialized;
      setLastSavedAt(new Date(result.updated_at));
    } catch (e) {
      setError(e instanceof Error ? e.message : "Échec de la sauvegarde automatique");
    } finally {
      setIsSaving(false);
    }
  }, [projectId, data]);

  useEffect(() => {
    if (!enabled) return;

    timerRef.current = setInterval(save, intervalMs);
    return () => {
      if (timerRef.current) clearInterval(timerRef.current);
    };
  }, [save, intervalMs, enabled]);

  const formattedTime = lastSavedAt
    ? lastSavedAt.toLocaleTimeString("fr-FR", { hour: "2-digit", minute: "2-digit" })
    : null;

  return { lastSavedAt, formattedTime, isSaving, error, saveNow: save };
}
