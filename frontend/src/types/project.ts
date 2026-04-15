export interface Project {
  id: string;
  name: string;
  industry_type: string;
  location: string;
  budget: number;
  floor_width: number;
  floor_depth: number;
  target_capacity?: string;
  status: "active" | "archived";
  version: number;
  is_archived: boolean;
  archived_at: string | null;
  created_at: string;
  updated_at: string;
}

export interface ProjectCreate {
  name: string;
  industry_type: string;
  location: string;
  budget: number;
  floor_width: number;
  floor_depth: number;
  target_capacity?: string;
}

export interface ProjectUpdate {
  name?: string;
  industry_type?: string;
  location?: string;
  budget?: number;
  floor_width?: number;
  floor_depth?: number;
  target_capacity?: string;
}

export interface AutoSaveResponse {
  message: string;
  version: number;
  updated_at: string;
}
