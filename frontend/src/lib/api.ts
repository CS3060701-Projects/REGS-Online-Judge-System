import { getToken } from "./auth";
import type { ExampleCase, Problem, ProblemStats, Submission, User } from "../types";

const API_BASE = import.meta.env.VITE_API_BASE_URL || "http://localhost:8081/api";

async function request<T>(path: string, init: RequestInit = {}, auth = false): Promise<T> {
  const headers = new Headers(init.headers || {});
  if (!headers.has("Content-Type") && !(init.body instanceof FormData)) {
    headers.set("Content-Type", "application/json");
  }
  if (auth) {
    const token = getToken();
    if (token) headers.set("Authorization", `Bearer ${token}`);
  }

  const res = await fetch(`${API_BASE}${path}`, { ...init, headers });
  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.error || `Request failed: ${res.status}`);
  }
  return res.json() as Promise<T>;
}

export const api = {
  register: (username: string, password: string) =>
    request<{ message: string; user_id: number; role: string }>("/users/register", {
      method: "POST",
      body: JSON.stringify({ username, password })
    }),
  login: (username: string, password: string) =>
    request<{ token: string; role: string; message: string }>("/users/login", {
      method: "POST",
      body: JSON.stringify({ username, password })
    }),
  logout: () => request<{ message: string }>("/users/logout", { method: "POST" }, true),
  me: async () => {
    const raw = await request<{
      id?: number;
      username?: string;
      role?: string;
      created_at?: string;
      ID?: number;
      Username?: string;
      Role?: string;
      CreatedAt?: string;
    }>("/users/me", {}, true);

    return {
      id: raw.id ?? raw.ID ?? 0,
      username: raw.username ?? raw.Username ?? "",
      role: (raw.role ?? raw.Role ?? "User") as User["role"],
      created_at: raw.created_at ?? raw.CreatedAt ?? ""
    } satisfies User;
  },

  getProblems: (page = 1, limit = 20) =>
    request<{ total: number; page: number; limit: number; data: Problem[] }>(`/problems?page=${page}&limit=${limit}`),
  getProblem: (id: string) => request<{ description: string }>(`/problems/${id}`),
  getProblemExamples: (id: string) => request<{ examples: ExampleCase[] }>(`/problems/${id}/examples`),
  getProblemStats: (id: string) => request<ProblemStats>(`/stats/problems/${id}`),

  submit: async (problemId: string, file: File) => {
    const formData = new FormData();
    formData.append("problem_id", problemId);
    formData.append("file", file);
    return request<{ message: string; operatorId: string; userId: number }>(
      "/submissions",
      { method: "POST", body: formData },
      true
    );
  },
  mySubmissions: () => request<Submission[]>("/submissions", {}, true),
  submissionStatus: (operatorId: string) =>
    request<Submission>(`/submissions/${operatorId}`, {}, true),

  createProblem: (payload: {
    id: string;
    title: string;
    description: string;
    time_limit: number;
    memory_limit: number;
    is_visible: boolean;
  }) => request<{ message: string }>("/problems", { method: "PUT", body: JSON.stringify(payload) }, true),
  deleteProblem: (id: string) => request<{ message: string }>(`/problems/${id}`, { method: "DELETE" }, true),
  uploadTestData: async (id: string, file: File) => {
    const formData = new FormData();
    formData.append("file", file);
    return request<{ message: string }>(`/problems/${id}/testdata`, { method: "POST", body: formData }, true);
  }
};
