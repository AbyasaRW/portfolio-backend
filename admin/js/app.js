// admin/js/app.js
// Pengendali utama Admin Panel: login, tab (Projects / Pengalaman / Tools), daftar,
// editor, dan hapus. Detail tiap jenis konten ada di resources.js, editor di editor.js.
import { request, session, loadSiteInfo, getSiteUrl } from "./api.js";
import { $, el, toast } from "./ui.js";
import { RESOURCES, byKey, GROUP_LABEL } from "./resources.js";
import { buildEditor, readValues, renderPreview } from "./editor.js";

const state = {
  active: RESOURCES[0], // resource yang sedang dibuka
  data: Object.fromEntries(RESOURCES.map((r) => [r.key, []])),
  editingId: null, // null = mode tambah
};

// ---------- Layar ----------
function showLogin() {
  $("loginView").hidden = false;
  $("dashView").hidden = true;
  $("logoutBtn").hidden = true;
  $("username").focus();
}

function showDashboard() {
  $("loginView").hidden = true;
  $("dashView").hidden = false;
  $("logoutBtn").hidden = false;
  selectTab(state.active.key);
}

function logout(expired = false) {
  session.clear();
  for (const r of RESOURCES) state.data[r.key] = [];
  showLogin();
  if (!expired) toast("Berhasil keluar.");
}
session.onExpired(() => {
  // Tutup dialog yang mungkin terbuka sebelum kembali ke login.
  for (const d of document.querySelectorAll("dialog[open]")) d.close();
  logout(true);
});

// ---------- Login ----------
$("loginForm").addEventListener("submit", async (event) => {
  event.preventDefault();
  const button = $("loginBtn");
  const error = $("loginError");
  error.hidden = true;
  button.disabled = true;
  button.textContent = "Memproses…";
  try {
    const data = await request("/login", { method: "POST", body: { username: $("username").value.trim(), password: $("password").value } });
    session.set(data.token);
    $("password").value = "";
    showDashboard();
  } catch (err) {
    error.textContent = err.message;
    error.hidden = false;
  } finally {
    button.disabled = false;
    button.textContent = "Masuk";
  }
});
$("logoutBtn").addEventListener("click", () => logout());

// ---------- Tab ----------
function buildTabs() {
  const bar = $("tabs");
  bar.replaceChildren();
  for (const r of RESOURCES) {
    const tab = el("button", "tab", r.tab);
    tab.type = "button";
    tab.id = `tab-${r.key}`;
    tab.setAttribute("role", "tab");
    tab.addEventListener("click", () => selectTab(r.key));
    bar.append(tab);
  }
}

function selectTab(key) {
  state.active = byKey(key);
  for (const tab of $("tabs").children) {
    const on = tab.id === `tab-${key}`;
    tab.classList.toggle("active", on);
    tab.setAttribute("aria-selected", String(on));
  }
  $("dashTitle").textContent = state.active.title;
  $("newBtn").textContent = state.active.addLabel;
  $("tabNote").textContent = state.active.note;
  load();
}

// ---------- Daftar ----------
async function load() {
  const resource = state.active;
  const list = $("list");
  list.className = "project-list";
  list.replaceChildren(el("li", "skeleton"), el("li", "skeleton"), el("li", "skeleton"));
  $("empty").hidden = true;
  $("summary").textContent = "Memuat…";
  try {
    const data = await request(`/${resource.key}`);
    state.data[resource.key] = Array.isArray(data) ? data : [];
    if (resource === state.active) render();
  } catch (err) {
    if (resource !== state.active) return;
    list.replaceChildren();
    $("summary").textContent = "Gagal memuat.";
    toast(err.message, "error");
  }
}

function renderRow(resource, item) {
  const parts = resource.row(item);
  const row = el("li", "project-row");

  const info = el("div", "project-info");
  info.append(el("h3", "", parts.title));
  if (parts.subtitle) info.append(el("p", "row-sub", parts.subtitle));
  const meta = el("div", "project-meta");
  meta.append(...parts.meta);
  if (resource.key !== "tools") meta.append(el("span", "order-chip", `urutan ${item.sort_order ?? 0}`));
  info.append(meta);

  const actions = el("div", "row-actions");
  const edit = el("button", "btn btn-soft", "Ubah");
  edit.type = "button";
  edit.addEventListener("click", () => openEditor(item));
  const del = el("button", "btn btn-soft", "Hapus");
  del.type = "button";
  del.addEventListener("click", () => askDelete(item));
  actions.append(edit, del);

  row.append(parts.thumb, info, actions);
  return row;
}

function render() {
  const resource = state.active;
  const items = state.data[resource.key];
  const list = $("list");
  list.replaceChildren();
  $("empty").textContent = resource.emptyText;
  $("empty").hidden = items.length > 0;
  $("summary").textContent = resource.summary(items);

  // Tools ditampilkan per kelompok supaya cocok dengan susunan di website.
  if (resource.key === "tools") {
    for (const [group, label] of Object.entries(GROUP_LABEL)) {
      const inGroup = items.filter((t) => (t.group || "frontend") === group);
      if (!inGroup.length) continue;
      list.append(el("li", "group-head", `${label} • ${inGroup.length}`));
      for (const item of inGroup) list.append(renderRow(resource, item));
    }
    return;
  }
  for (const item of items) list.append(renderRow(resource, item));
}

// ---------- Editor ----------
const editor = $("editor");
const form = $("editorForm");

function nextOrder(resource, group) {
  const pool = state.data[resource.key].filter((i) => resource.key !== "tools" || (i.group || "frontend") === group);
  return Math.max(0, ...pool.map((i) => i.sort_order || 0)) + 1;
}

function openEditor(item = null) {
  const resource = state.active;
  state.editingId = item ? item.ID : null;
  $("editorTitle").textContent = item ? `Ubah ${resource.singular}` : `${resource.singular[0].toUpperCase()}${resource.singular.slice(1)} baru`;
  $("saveBtn").textContent = item ? "Simpan perubahan" : `Tambah ${resource.singular}`;
  $("editorError").hidden = true;

  // Data baru: urutan diletakkan paling akhir.
  const values = item ?? { sort_order: nextOrder(resource, "frontend") };
  buildEditor(resource, values);
  renderPreview(resource);
  editor.showModal();
  form.querySelector("input,textarea")?.focus();
}

form.addEventListener("input", () => {
  $("editorError").hidden = true; // pesan kesalahan hilang begitu pengguna mulai memperbaiki
  renderPreview(state.active);
});
$("newBtn").addEventListener("click", () => openEditor());
for (const button of editor.querySelectorAll("[data-close]")) button.addEventListener("click", () => editor.close());
// Klik di luar kotak dialog (pada backdrop) menutupnya.
editor.addEventListener("click", (event) => { if (event.target === editor) editor.close(); });

form.addEventListener("submit", async (event) => {
  event.preventDefault();
  const resource = state.active;
  const error = $("editorError");
  const data = readValues(resource);
  const required = resource.fields.find((f) => f.required);

  if (required && !String(data[required.name]).trim()) {
    error.textContent = `${required.label} wajib diisi.`;
    error.hidden = false;
    form.elements[required.name].focus();
    return;
  }
  error.hidden = true;

  const button = $("saveBtn");
  const label = button.textContent;
  button.disabled = true;
  button.textContent = "Menyimpan…";
  try {
    const path = state.editingId ? `/${resource.key}/${state.editingId}` : `/${resource.key}`;
    await request(path, { method: state.editingId ? "PUT" : "POST", body: data, auth: true });
    editor.close();
    toast(state.editingId ? "Perubahan disimpan." : "Data ditambahkan.");
    await load();
  } catch (err) {
    error.textContent = err.message;
    error.hidden = false;
  } finally {
    button.disabled = false;
    button.textContent = label;
  }
});

// ---------- Hapus ----------
const confirmDialog = $("confirm");
function askDelete(item) {
  const resource = state.active;
  $("confirmText").textContent = `“${resource.itemName(item)}” akan dihapus dari website. Tindakan ini tidak bisa dibatalkan.`;
  confirmDialog.returnValue = "";
  confirmDialog.showModal();
  confirmDialog.addEventListener("close", async function onClose() {
    confirmDialog.removeEventListener("close", onClose);
    if (confirmDialog.returnValue !== "ok") return;
    try {
      await request(`/${resource.key}/${item.ID}`, { method: "DELETE", auth: true });
      toast("Data dihapus.");
      await load();
    } catch (err) {
      toast(err.message, "error");
    }
  });
}

// ---------- Mulai ----------
buildTabs();
loadSiteInfo().then(() => {
  $("siteLink").href = getSiteUrl();
  if (session.token) showDashboard();
  else showLogin();
});
