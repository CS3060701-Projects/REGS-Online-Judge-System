export type Role = "Admin" | "User";

export interface User {
  id: number;
  username: string;
  role: Role;
  created_at: string;
}

export interface Problem {
  id: string;
  title: string;
  description?: string;
  time_limit: number;
  memory_limit: number;
  created_at?: string;
}

export interface Submission {
  operatorId?: string;
  operator_id?: string;
  problem_id?: string;
  status: string;
  run_time?: number;
  run_memory?: number;
  created_at?: string;
  Problem?: {
    id: string;
    title: string;
  };
}

export interface ProblemStats {
  total_submissions: number;
  ac_count: number;
  acceptance_rate: number;
  status_distribution: Record<string, number>;
}
