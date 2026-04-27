import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { api } from "../lib/api";
import type { Problem } from "../types";

export function ProblemsPage() {
  const [items, setItems] = useState<Problem[]>([]);
  const [keyword, setKeyword] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    api
      .getProblems(1, 50)
      .then((res) => setItems(res.data))
      .catch((err) => setError((err as Error).message))
      .finally(() => setLoading(false));
  }, []);

  const filtered = items.filter((p) => {
    const k = keyword.trim().toLowerCase();
    if (!k) return true;
    return p.id.toLowerCase().includes(k) || p.title.toLowerCase().includes(k);
  });

  return (
    <main className="container">
      <section className="section-head">
        <div>
          <h2>題目列表</h2>
          <p className="muted">可依題號或標題搜尋，點入查看題目與統計</p>
        </div>
        <input
          className="search"
          placeholder="搜尋題號或標題..."
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
        />
      </section>

      {error && <p className="error">{error}</p>}
      {loading ? (
        <div className="panel">載入題目中...</div>
      ) : filtered.length === 0 ? (
        <div className="panel">找不到符合條件的題目。</div>
      ) : (
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>標題</th>
                <th>時間限制(ms)</th>
                <th>記憶體限制(MB)</th>
              </tr>
            </thead>
            <tbody>
              {filtered.map((p) => (
                <tr key={p.id}>
                  <td>{p.id}</td>
                  <td>
                    <Link to={`/problems/${p.id}`}>{p.title}</Link>
                  </td>
                  <td>{p.time_limit}</td>
                  <td>{p.memory_limit}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </main>
  );
}
