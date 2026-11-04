// admin/js/ui.js
// Pembantu tampilan kecil: membuat elemen, notifikasi, dan konstanta warna.
// Semua teks dari database dimasukkan lewat textContent (bukan innerHTML),
// sehingga isi data tidak bisa menyisipkan HTML/skrip (XSS).

export const $ = (id) => document.getElementById(id);

export function el(tag, className, text) {
  const node = document.createElement(tag);
  if (className) node.className = className;
  if (text !== undefined) node.textContent = text;
  return node;
}

export function toast(message, type = "ok") {
  const node = el("div", `toast ${type}`, message);
  $("toasts").append(node);
  setTimeout(() => node.remove(), 3800);
}

// Warna label project; dipakai thumbnail, pratinjau, dan pilihan warna.
export const ACCENTS = ["blue", "cyan", "violet", "orange", "slate"];
export const ACCENT_BG = { blue: "#e0e1ff", cyan: "#15132e", violet: "#ddf8fa", orange: "#ffeadf", slate: "#0b0f17" };

// Warna inisial pengalaman.
export const TONES = ["blue", "cyan", "orange", "violet"];
export const TONE_BG = { blue: "#e0e1ff", cyan: "#d8f8fa", orange: "#ffeadf", violet: "#eee9ff" };
export const TONE_FG = { blue: "#3525cd", cyan: "#087b8a", orange: "#9a3d12", violet: "#6338c9" };
