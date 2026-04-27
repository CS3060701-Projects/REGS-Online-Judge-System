import { Link, NavLink, Navigate, Route, Routes, useNavigate } from "react-router-dom";
import { useAuth } from "./contexts/AuthContext";
import { LoginPage } from "./pages/LoginPage";
import { RegisterPage } from "./pages/RegisterPage";
import { ProblemsPage } from "./pages/ProblemsPage";
import { ProblemDetailPage } from "./pages/ProblemDetailPage";
import { SubmissionsPage } from "./pages/SubmissionsPage";
import { AdminPage } from "./pages/AdminPage";

function Protected({ children }: { children: JSX.Element }) {
  const { token, loading } = useAuth();
  if (loading) return <p className="container">Loading...</p>;
  if (!token) return <Navigate to="/login" replace />;
  return children;
}

function AdminProtected({ children }: { children: JSX.Element }) {
  const { role } = useAuth();
  if (role !== "Admin") return <Navigate to="/" replace />;
  return children;
}

export default function App() {
  const { token, role, user, logout } = useAuth();
  const navigate = useNavigate();

  const onLogout = async () => {
    await logout();
    navigate("/login");
  };

  return (
    <div className="app-shell">
      <header className="topbar">
        <div>
          <h1>REGS Online Judge</h1>
          <p className="muted">題目管理、程式提交與評測查詢平台</p>
        </div>
        <div className="chip">{user ? `${user.username} (${user.role})` : "訪客模式"}</div>
      </header>

      <nav className="nav">
        <NavLink to="/" className={({ isActive }) => `nav-link${isActive ? " active" : ""}`}>
          題目
        </NavLink>
        {token && (
          <NavLink to="/submissions" className={({ isActive }) => `nav-link${isActive ? " active" : ""}`}>
            我的提交
          </NavLink>
        )}
        {role === "Admin" && (
          <NavLink to="/admin" className={({ isActive }) => `nav-link${isActive ? " active" : ""}`}>
            管理後台
          </NavLink>
        )}
        {!token && (
          <NavLink to="/login" className={({ isActive }) => `nav-link${isActive ? " active" : ""}`}>
            登入
          </NavLink>
        )}
        {!token && (
          <NavLink to="/register" className={({ isActive }) => `nav-link${isActive ? " active" : ""}`}>
            註冊
          </NavLink>
        )}
        {token && (
          <button className="ghost danger-text nav-link" onClick={onLogout}>
            登出
          </button>
        )}
      </nav>

      <Routes>
        <Route path="/" element={<ProblemsPage />} />
        <Route path="/problems/:id" element={<ProblemDetailPage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route
          path="/submissions"
          element={
            <Protected>
              <SubmissionsPage />
            </Protected>
          }
        />
        <Route
          path="/admin"
          element={
            <Protected>
              <AdminProtected>
                <AdminPage />
              </AdminProtected>
            </Protected>
          }
        />
      </Routes>
    </div>
  );
}
