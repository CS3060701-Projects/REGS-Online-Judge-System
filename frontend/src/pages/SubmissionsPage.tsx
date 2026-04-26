import { FormEvent, useState } from "react";
import { api } from "../lib/api";
import type { Submission } from "../types";

export function SubmissionsPage() {
  const [items, setItems] = useState<Submission[]>([]);
  const [operatorId, setOperatorId] = useState("");
  const [single, setSingle] = useState<Submission | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const getStatusClass = (status: string) => {
    if (status === "AC") return "status-pill status-ac";
    if (status === "WA") return "status-pill status-wa";
    if (status === "CE" || status === "SE") return "status-pill status-ce";
    if (status === "RE" || status === "TLE") return "status-pill status-re";
    if (status === "Pending" || status === "Compiling" || status === "Judging") return "status-pill status-pending";
    return "status-pill";
  };

  const loadMine = async () => {
    setError("");
    setLoading(true);
    try {
      const data = await api.mySubmissions();
      setItems(data);
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setLoading(false);
    }
  };

  const queryOne = async (e: FormEvent) => {
    e.preventDefault();
    setError("");
    try {
      const data = await api.submissionStatus(operatorId);
      setSingle(data);
    } catch (err) {
      setError((err as Error).message);
    }
  };

  return (
    <main className="container">
      <section className="section-head">
        <div>
          <h2>我的提交</h2>
          <p className="muted">可查詢全部提交，或用 operatorId 即時查看單筆評測狀態</p>
        </div>
        <button onClick={loadMine}>{loading ? "載入中..." : "重新整理提交"}</button>
      </section>

      <form onSubmit={queryOne} className="form inline">
        <input
          placeholder="輸入 operatorId 查詢"
          value={operatorId}
          onChange={(e) => setOperatorId(e.target.value)}
        />
        <button type="submit">查詢單筆</button>
      </form>

      {single && (
        <div className="panel">
          <strong>查詢結果:</strong> {single.operatorId} | {single.status} | {single.run_time ?? "-"}ms |{" "}
          {single.run_memory ?? "-"}KB
        </div>
      )}

      <div className="table-wrap">
        <table>
          <thead>
            <tr>
              <th>operator</th>
              <th>題目</th>
              <th>狀態</th>
              <th>時間</th>
              <th>記憶體</th>
            </tr>
          </thead>
          <tbody>
            {items.map((s) => (
              <tr key={s.operator_id || s.operatorId}>
                <td>{s.operator_id || s.operatorId}</td>
                <td>{s.Problem?.id ?? s.problem_id ?? "-"}</td>
                <td>
                  <span className={getStatusClass(s.status)}>{s.status}</span>
                </td>
                <td>{s.run_time ?? "-"}</td>
                <td>{s.run_memory ?? "-"}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {!loading && items.length === 0 && <div className="panel">目前沒有提交紀錄。</div>}

      {error && <p className="error">{error}</p>}
    </main>
  );
}
