// scripts/gen-seed.mjs
// Membuat backend/seed_content.json dari data frontend (profile.ts, experience.ts, tools.ts),
// supaya data awal di database SAMA PERSIS dengan data cadangan di website.
//
// Jalankan dari folder backend:   node --no-warnings scripts/gen-seed.mjs
// (Node 22+ bisa membaca file .ts langsung.)
import { writeFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import path from "node:path";

const dir = path.dirname(fileURLToPath(import.meta.url));
const data = (name) => import(new URL(`../../frontend/data/${name}.ts`, import.meta.url).href);

const profile = (await data("profile")).default;
const { experienceMeta } = await data("experience");
const toolStacks = (await data("tools")).default;

const experiences = profile.experience.map((item, index) => {
  const meta = experienceMeta[index];
  return {
    company: meta.company,
    role: meta.role,
    full_role: item.role,
    place: meta.place,
    period: item.period,
    categories: meta.categories.join(","),
    tone: meta.tone,
    initials: meta.initials,
    logo: meta.logo,
    summary: meta.summary,
    skills: meta.skills.join(", "),
    points: item.points.join("\n"),
    sort_order: index + 1,
  };
});

const tools = Object.entries(toolStacks).flatMap(([group, list]) =>
  list.map((tool, index) => ({ group, name: tool.name, icon: tool.src, sort_order: index + 1 })),
);

const out = path.join(dir, "..", "seed_content.json");
writeFileSync(out, JSON.stringify({ experiences, tools }, null, 2) + "\n", "utf8");
console.log(`seed_content.json: ${experiences.length} experiences, ${tools.length} tools`);
