import { FormEvent, useState } from "react";
import { api } from "../lib/api";

export function AdminPage() {
  const [id, setId] = useState("");
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [timeLimit, setTimeLimit] = useState(1000);
  const [memoryLimit, setMemoryLimit] = useState(256);
  const [isVisible, setIsVisible] = useState(true);
  const [zip, setZip] = useState<File | null>(null);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  const createOrUpdate = async (e: FormEvent) => {
    e.preventDefault();
    setError("");
    setMessage("");
    try {
      await api.createProblem({
        id,
        title,
        description,
        time_limit: timeLimit,
        memory_limit: memoryLimit,
        is_visible: isVisible
      });
      if (zip) await api.uploadTestData(id, zip);
      setMessage("題目與測資已更新。");
    } catch (err) {
      setError((err as Error).message);
    }
  };

  const deleteProblem = async () => {
    if (!id) return;
    setError("");
    setMessage("");
    try {
      await api.deleteProblem(id);
      setMessage("題目已刪除。");
    } catch (err) {
      setError((err as Error).message);
    }
  };

  return (
    <main className="container">
      <section className="section-head">
        <div>
          <h2>管理員 - 題目管理</h2>
          <p className="muted">建立題目、更新設定、上傳測資與刪除題目</p>
        </div>
      </section>

      <form onSubmit={createOrUpdate} className="form card">
        <div className="grid-2">
          <input placeholder="題目 ID (例如 p1003)" value={id} onChange={(e) => setId(e.target.value)} />
          <input placeholder="標題" value={title} onChange={(e) => setTitle(e.target.value)} />
        </div>
        <textarea placeholder="題目敘述" value={description} onChange={(e) => setDescription(e.target.value)} rows={8} />
        <div className="grid-2">
          <input
            type="number"
            value={timeLimit}
            onChange={(e) => setTimeLimit(Number(e.target.value))}
            placeholder="time limit (ms)"
          />
          <input
            type="number"
            value={memoryLimit}
            onChange={(e) => setMemoryLimit(Number(e.target.value))}
            placeholder="memory limit (MB)"
          />
        </div>
        <label className="inline-label">
          <input type="checkbox" checked={isVisible} onChange={(e) => setIsVisible(e.target.checked)} />
          設為公開題目
        </label>
        <input type="file" accept=".zip" onChange={(e) => setZip(e.target.files?.[0] || null)} />
        <div className="row">
          <button type="submit">建立/更新題目</button>
          <button type="button" className="danger" onClick={deleteProblem} disabled={!id}>
            刪除題目（依 ID）
          </button>
        </div>
      </form>
      {message && <p className="ok">{message}</p>}
      {error && <p className="error">{error}</p>}
    </main>
  );
}
