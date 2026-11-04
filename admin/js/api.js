// admin/js/api.js
// Semua komunikasi dengan API backend + penyimpanan token login + alamat website.
const API = `${location.origin}/api`;
const TOKEN_KEY = "admin_token";

let token = localStorage.getItem(TOKEN_KEY);
let onUnauthorized = () => {};

export const session = {
  get token() { return token; },
  set(value) { token = value; localStorage.setItem(TOKEN_KEY, value); },
  clear() { token = null; localStorage.removeItem(TOKEN_KEY); },
  onExpired(handler) { onUnauthorized = handler; },
};

// ---- Alamat website portfolio ----
// Gambar & logo disimpan di folder public/ FRONTEND (path relatif seperti /Projects/x.png),
// jadi untuk menampilkannya di sini harus diarahkan ke alamat website.
const override = new URLSearchParams(location.search).get("site");
let siteUrl = (override || "http://localhost:3000").replace(/\/+$/, "");

export async function loadSiteInfo() {
  try {
    const data = await fetch(`${API}/site`).then((r) => r.json());
    if (data.site_url && !override) siteUrl = data.site_url;
  } catch { /* pakai alamat bawaan */ }
  return siteUrl;
}
export const getSiteUrl = () => siteUrl;

/** Path relatif -> alamat website; URL penuh dibiarkan. */
export function assetSrc(src) {
  if (!src) return "";
  return /^(https?:)?\/\//i.test(src) ? src : siteUrl + (src.startsWith("/") ? "" : "/") + src;
}

// ---- Permintaan ----
export async function request(path, { method = "GET", body, auth = false } = {}) {
  const headers = {};
  if (body) headers["Content-Type"] = "application/json";
  if (auth) headers.Authorization = `Bearer ${token}`;

  let res;
  try {
    res = await fetch(`${API}${path}`, { method, headers, body: body ? JSON.stringify(body) : undefined });
  } catch {
    throw new Error("Tidak bisa menghubungi server. Pastikan backend berjalan.");
  }

  // Token kedaluwarsa / tidak valid: kembali ke layar login.
  if (res.status === 401 && auth) {
    onUnauthorized();
    throw new Error("Sesi berakhir. Silakan masuk lagi.");
  }

  const data = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(data.error || `Permintaan gagal (${res.status})`);
  return data;
}
