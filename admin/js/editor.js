// admin/js/editor.js
// Editor generik: membuat isian dari definisi resource (resources.js), membaca/menulis
// nilainya, dan memperbarui pratinjau saat mengetik. Dipakai untuk Project, Pengalaman, dan Tools.
import { $, el } from "./ui.js";

/** Membuat satu isian sesuai tipenya. */
function buildField(field) {
  const id = `f_${field.name}`;
  const label = el("label", `field${field.half ? " half" : ""}`);
  if (field.type === "swatch" || field.type === "segmented" || field.type === "multi") {
    // Pilihan ganda dibungkus fieldset supaya terbaca benar oleh pembaca layar.
    const set = el("fieldset", `field${field.half ? " half" : ""}`);
    const legend = el("legend", "", field.label);
    if (field.hint) legend.append(el("small", "", field.hint));
    set.append(legend);

    const group = el("div", field.type === "swatch" ? "swatches" : field.type === "multi" ? "multi" : "segmented");
    const options = field.options.map((o) => (Array.isArray(o) ? o : [o, o]));
    for (const [value, text] of options) {
      const item = el("label");
      const input = document.createElement("input");
      input.type = field.type === "multi" ? "checkbox" : "radio";
      input.name = field.name;
      input.value = value;
      if (field.type === "swatch") item.append(input, el("span", `tag tag-${value}`, text));
      else item.append(input, el("span", "", text));
      group.append(item);
    }
    set.append(group);
    return set;
  }

  label.append(el("span", "", field.label + (field.required ? " *" : "")));
  if (field.hint) label.firstChild.append(el("small", "", field.hint));

  let input;
  if (field.type === "textarea") {
    input = document.createElement("textarea");
    input.rows = field.rows || 4;
  } else {
    input = document.createElement("input");
    input.type = field.type === "number" ? "number" : field.type === "url" ? "url" : "text";
    if (field.type === "number") { input.min = "0"; input.step = "1"; }
  }
  input.id = id;
  input.name = field.name;
  if (field.placeholder) input.placeholder = field.placeholder;
  if (field.maxlength) input.maxLength = field.maxlength;
  if (field.required) input.required = true;
  label.append(input);
  return label;
}

/** Mengisi #editorFields dengan isian resource dan memasang nilai awal. */
export function buildEditor(resource, values) {
  const root = $("editorFields");
  root.replaceChildren(...resource.fields.map(buildField));
  writeValues(resource, values);
}

export function writeValues(resource, values) {
  const form = $("editorForm");
  for (const f of resource.fields) {
    const value = values[f.name] ?? "";
    if (f.type === "swatch" || f.type === "segmented") {
      const options = f.options.map((o) => (Array.isArray(o) ? o[0] : o));
      form.elements[f.name].value = options.includes(value) ? value : f.fallback ?? options[0];
    } else if (f.type === "multi") {
      const chosen = String(value).split(",").map((s) => s.trim());
      for (const input of form.querySelectorAll(`input[name="${f.name}"]`)) input.checked = chosen.includes(input.value);
    } else {
      form.elements[f.name].value = value;
    }
  }
}

/** Membaca semua isian menjadi objek siap dikirim ke API. */
export function readValues(resource) {
  const form = $("editorForm");
  const out = {};
  for (const f of resource.fields) {
    if (f.type === "swatch" || f.type === "segmented") out[f.name] = form.elements[f.name].value;
    else if (f.type === "multi") out[f.name] = [...form.querySelectorAll(`input[name="${f.name}"]:checked`)].map((i) => i.value).join(",");
    else if (f.type === "number") out[f.name] = parseInt(form.elements[f.name].value, 10) || 0;
    else out[f.name] = form.elements[f.name].value;
  }
  return out;
}

export function renderPreview(resource) {
  resource.preview(readValues(resource), $("previewRoot"));
  $("previewNote").textContent = `${resource.previewNote} Perubahan muncul di website dalam ±30 detik.`;
}
