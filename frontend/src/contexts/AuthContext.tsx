import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { api } from "../lib/api";
import { clearAuth, getRole, getToken, saveAuth } from "../lib/auth";
import type { User } from "../types";

interface AuthContextValue {
  token: string | null;
  role: string | null;
  user: User | null;
  loading: boolean;
  login: (username: string, password: string) => Promise<void>;
  register: (username: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [token, setToken] = useState<string | null>(getToken());
  const [role, setRole] = useState<string | null>(getRole());
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const bootstrap = async () => {
      if (!token) {
        setLoading(false);
        return;
      }
      try {
        const me = await api.me();
        setUser(me);
      } catch {
        clearAuth();
        setToken(null);
        setRole(null);
      } finally {
        setLoading(false);
      }
    };
    void bootstrap();
  }, [token]);

  const value = useMemo<AuthContextValue>(
    () => ({
      token,
      role,
      user,
      loading,
      login: async (username, password) => {
        const res = await api.login(username, password);
        saveAuth(res.token, res.role);
        setToken(res.token);
        setRole(res.role);
        const me = await api.me();
        setUser(me);
      },
      register: async (username, password) => {
        await api.register(username, password);
      },
      logout: async () => {
        try {
          if (token) await api.logout();
        } finally {
          clearAuth();
          setToken(null);
          setRole(null);
          setUser(null);
        }
      }
    }),
    [token, role, user, loading]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
}
