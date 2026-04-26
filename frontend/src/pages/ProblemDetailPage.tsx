import { FormEvent, useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { api } from "../lib/api";
import { useAuth } from "../contexts/AuthContext";
import type { ProblemStats } from "../types";

export function ProblemDetailPage() {
  const { id = "" } = useParams();
  const { token } = useAuth();
  const [description, setDescription] = useState("");
  const [stats, setStats] = useState<ProblemStats | null>(null);
  const [file, setFile] = useState<File | null>(null);
  const [operatorId, setOperatorId] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");

  useEffect(() => {
    api.getProblem(id).then((res) => setDescription(res.description)).catch((err) => setError((err as Error).message));
    api.getProblemStats(id).then(setStats).catch(() => null);
  }, [id]);

  const submit = async (e: FormEvent) => {
    e.preventDefault();
    if (!file) return;
    setError("");
    setMessage("");
    setSubmitting(true);
    try {
      const res = await api.submit(id, file);
      setOperatorId(res.operatorId);
      setMessage("提交成功，請到我的提交頁查看結果。");
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <main className="container">
      <section className="section-head">
        <div>
          <h2>題目 {id}</h2>
          <p className="muted">查看題目內容、統計資訊並提交 zip 程式碼</p>
        </div>
      </section>
      <pre className="desc">{description || "載入中..."}</pre>
      {stats && (
        <div className="metrics">
          <div className="metric-card">
            <span>提交數</span>
            <strong>{stats.total_submissions}</strong>
          </div>
          <div className="metric-card">
            <span>AC 數</span>
            <strong>{stats.ac_count}</strong>
          </div>
          <div className="metric-card">
            <span>通過率</span>
            <strong>{stats.acceptance_rate}%</strong>
          </div>
        </div>
      )}

      {token ? (
        <form onSubmit={submit} className="form card">
          <h3>提交 zip 程式碼</h3>
          <input type="file" accept=".zip" onChange={(e) => setFile(e.target.files?.[0] || null)} />
          <button type="submit" disabled={!file || submitting}>
            {submitting ? "提交中..." : "送出評測"}
          </button>
          {operatorId && <p>operatorId: {operatorId}</p>}
        </form>
      ) : (
        <div className="panel">請先登入才可提交。</div>
      )}

      {message && <p className="ok">{message}</p>}
      {error && <p className="error">{error}</p>}
    </main>
  );
}
