// admin/js/resources.js
// Definisi tiga jenis konten yang bisa dikelola: Project, Pengalaman, Tools.
//
// Setiap resource menjelaskan:
//   fields   : isian di editor (nama = nama kolom JSON di API)
//   row()    : tampilan satu baris di daftar
//   preview(): pratinjau langsung di editor, meniru tampilan di website
// Editor dan daftar (app.js) membaca definisi ini, jadi menambah jenis konten baru
// cukup menambah satu objek di sini.
import { assetSrc } from "./api.js";
import { el, ACCENTS, ACCENT_BG, TONES, TONE_BG, TONE_FG } from "./ui.js";

export const CATEGORY_LABEL = { broadcast: "Media & Broadcast", tech: "Tech & Dev", operations: "Operations" };
export const GROUP_LABEL = { frontend: "Frontend", backend: "Backend & Data", design: "Design & Media", dev: "Dev Environment" };

// ---------- Pembantu ----------

/** Gambar dari path/URL; bila kosong atau gagal dimuat, diganti teks cadangan. */
function picture(src, fallbackText, className = "") {
  const img = document.createElement("img");
  img.alt = "";
  img.loading = "lazy";
  const holder = el("span", "pic-empty", fallbackText);
  if (!src) return holder;
  img.className = className;
  img.src = assetSrc(src);
  img.addEventListener("error", () => img.replaceWith(el("span", "pic-empty", "Gambar gagal")));
  return img;
}

/** Teks dengan **tebal** -> node aman (tanpa innerHTML). */
function richText(text) {
  const frag = document.createDocumentFragment();
  text.split("**").forEach((part, i) => frag.append(i % 2 ? el("strong", "", part) : document.createTextNode(part)));
  return frag;
}

const splitList = (value, sep) => (value || "").split(sep).map((s) => s.trim()).filter(Boolean);

// ---------- 1. Projects ----------
const projects = {
  key: "projects",
  tab: "Projects",
  title: "Featured Projects",
  addLabel: "+ Project baru",
  singular: "project",
  note: "Tampil di carousel Featured Projects.",
  fields: [
    { name: "title", label: "Judul", type: "text", required: true, maxlength: 160 },
    { name: "stack", label: "Teknologi", hint: "pisahkan dengan koma", type: "text", placeholder: "Next.js, Go, PostgreSQL" },
    { name: "description", label: "Deskripsi", type: "textarea", rows: 4 },
    { name: "link", label: "Tautan project", type: "url", placeholder: "https://…", half: true },
    { name: "image", label: "Gambar", hint: "path atau URL", type: "text", placeholder: "/Projects/nama.png", half: true },
    { name: "category", label: "Kategori", type: "text", placeholder: "Web Development", half: true },
    { name: "sort_order", label: "Urutan", hint: "kecil = awal", type: "number", half: true },
    { name: "accent", label: "Warna label", type: "swatch", options: ACCENTS, fallback: "blue" },
    { name: "image_fit", label: "Tampilan gambar", type: "segmented", options: [["cover", "Cover"], ["contain", "Contain"]], fallback: "cover" },
  ],
  itemName: (item) => item.title,
  summary: (items) => `${items.length} project • tampil di website sesuai urutan`,
  emptyText: "Belum ada project. Tekan “Project baru” untuk menambah.",
  row(item) {
    const accent = ACCENTS.includes(item.accent) ? item.accent : "blue";
    const thumb = el("div", "thumb");
    thumb.style.background = ACCENT_BG[accent];
    thumb.dataset.fit = item.image_fit === "contain" ? "contain" : "cover";
    thumb.append(picture(item.image, "Tanpa gambar"));

    const meta = [el("span", `tag tag-${accent}`, item.category || "Tanpa kategori")];
    if (item.stack) meta.push(el("span", "stack-text", item.stack));
    return { thumb, title: item.title, meta };
  },
  preview(d, root) {
    const accent = ACCENTS.includes(d.accent) ? d.accent : "blue";
    const card = el("div", "pcard");
    card.dataset.fit = d.image_fit === "contain" ? "contain" : "cover";
    const media = el("div", "pcard-media");
    media.style.background = ACCENT_BG[accent];
    media.append(picture(d.image.trim(), "Tanpa gambar"));
    const caption = el("div", "pcard-caption");
    caption.append(el("span", "", d.category.trim() || "Kategori"), el("strong", "", d.title.trim() || "Judul project"));
    card.append(media, caption);
    root.replaceChildren(card);
  },
  previewNote: "Tampilan akhir di website mengikuti pilihan ini.",
};

// ---------- 2. Pengalaman ----------
const experiences = {
  key: "experiences",
  tab: "Pengalaman",
  title: "Career Experience",
  addLabel: "+ Pengalaman baru",
  singular: "pengalaman",
  note: "Kartu yang bisa dibalik di section Career Experience.",
  fields: [
    { name: "company", label: "Perusahaan", type: "text", required: true, maxlength: 120, half: true },
    { name: "place", label: "Lokasi", type: "text", placeholder: "Surabaya", half: true },
    { name: "role", label: "Peran (judul di depan kartu)", type: "text", placeholder: "MCR Operator & Motion Graphic Designer" },
    { name: "full_role", label: "Peran lengkap (di belakang kartu)", type: "text" },
    { name: "period", label: "Periode", type: "text", placeholder: "Agustus - November 2023", half: true },
    { name: "sort_order", label: "Urutan", hint: "kecil = awal", type: "number", half: true },
    { name: "summary", label: "Ringkasan", hint: "**kata** untuk huruf tebal", type: "textarea", rows: 3 },
    { name: "skills", label: "Skills", hint: "pisahkan dengan koma", type: "text", placeholder: "MCR Operation, Video Editing" },
    { name: "points", label: "Daftar tugas", hint: "satu tugas per baris", type: "textarea", rows: 5 },
    { name: "logo", label: "Logo", hint: "path atau URL", type: "text", placeholder: "/logos/nama.png", half: true },
    { name: "initials", label: "Inisial", hint: "bila logo kosong", type: "text", maxlength: 3, half: true },
    { name: "tone", label: "Warna inisial", type: "swatch", options: TONES, fallback: "blue" },
    { name: "categories", label: "Kategori filter", hint: "pilih satu atau lebih", type: "multi", options: Object.entries(CATEGORY_LABEL) },
  ],
  itemName: (item) => item.company,
  summary: (items) => `${items.length} pengalaman • kartu muncul sesuai urutan`,
  emptyText: "Belum ada pengalaman. Tekan “Pengalaman baru” untuk menambah.",
  row(item) {
    const tone = TONES.includes(item.tone) ? item.tone : "blue";
    const thumb = el("div", "thumb thumb-logo");
    if (item.logo) {
      thumb.append(picture(item.logo, item.initials || "—"));
    } else {
      const initials = el("span", "initials", item.initials || "—");
      initials.style.background = TONE_BG[tone];
      initials.style.color = TONE_FG[tone];
      thumb.append(initials);
    }
    const meta = [];
    if (item.period) meta.push(el("span", "stack-text", item.period));
    for (const c of splitList(item.categories, ",")) meta.push(el("span", "tag tag-slate", CATEGORY_LABEL[c] || c));
    return { thumb, title: item.company, subtitle: item.role, meta };
  },
  preview(d, root) {
    const tone = TONES.includes(d.tone) ? d.tone : "blue";
    const card = el("div", "xp-preview");

    const head = el("div", "xp-head");
    const logo = el("span", "xp-logo");
    if (d.logo.trim()) {
      logo.append(picture(d.logo.trim(), d.initials.trim() || "—"));
    } else {
      logo.textContent = d.initials.trim() || "—";
      logo.style.background = TONE_BG[tone];
      logo.style.color = TONE_FG[tone];
    }
    head.append(logo, el("strong", "", d.company.trim() || "Perusahaan"));

    const meta = [d.place.trim(), d.period.trim()].filter(Boolean).join("  •  ");
    const summary = el("p", "xp-summary");
    summary.append(richText(d.summary.trim() || "Ringkasan pengalaman"));

    const chips = el("div", "xp-chips");
    for (const skill of splitList(d.skills, ",").slice(0, 6)) chips.append(el("span", "", skill));

    card.append(head, el("h4", "", d.role.trim() || "Peran"), el("p", "xp-meta", meta || "Lokasi • Periode"), summary, chips);
    root.replaceChildren(card);
  },
  previewNote: "Sisi depan kartu. Sisi belakang berisi peran lengkap dan daftar tugas.",
};

// ---------- 3. Tools ----------
const tools = {
  key: "tools",
  tab: "Tools",
  title: "Tools & Capabilities",
  addLabel: "+ Tool baru",
  singular: "tool",
  note: "Logo teknologi di section Tools & Capabilities.",
  layout: "grid",
  fields: [
    { name: "name", label: "Nama", type: "text", required: true, maxlength: 60, placeholder: "Next.js" },
    { name: "icon", label: "Logo", hint: "path atau URL", type: "text", placeholder: "/tools/nextjs.svg" },
    { name: "sort_order", label: "Urutan dalam kelompok", hint: "kecil = awal", type: "number" },
    { name: "group", label: "Kelompok", type: "segmented", options: Object.entries(GROUP_LABEL), fallback: "frontend" },
  ],
  itemName: (item) => item.name,
  summary: (items) => `${items.length} tool • dikelompokkan di website`,
  emptyText: "Belum ada tool. Tekan “Tool baru” untuk menambah.",
  row(item) {
    const thumb = el("div", "thumb thumb-tool");
    thumb.append(picture(item.icon, "—"));
    const group = ["frontend", "backend", "design", "dev"].includes(item.group) ? item.group : "frontend";
    return { thumb, title: item.name, meta: [el("span", "tag tag-slate", GROUP_LABEL[group]), el("span", "order-chip", `urutan ${item.sort_order ?? 0}`)] };
  },
  preview(d, root) {
    const card = el("div", "tool-preview");
    const dot = el("span", "tool-preview-dot");
    dot.append(picture(d.icon.trim(), "—"));
    card.append(dot, el("strong", "", d.name.trim() || "Nama tool"), el("span", "tag tag-slate", GROUP_LABEL[d.group] || GROUP_LABEL.frontend));
    root.replaceChildren(card);
  },
  previewNote: "Bulatan logo seperti di website. Nama muncul saat kursor diarahkan.",
};

export const RESOURCES = [projects, experiences, tools];
export const byKey = (key) => RESOURCES.find((r) => r.key === key);
