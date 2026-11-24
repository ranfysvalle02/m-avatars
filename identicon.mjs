#!/usr/bin/env node
// m-avatars: one melody, many instruments. JavaScript performance of SPEC.md.
//
//   node identicon.mjs captain-pixel
//   node identicon.mjs --format svg captain-pixel > avatar.svg
//   node identicon.mjs --check golden.json

import { readFileSync } from "node:fs";
import { parseArgs } from "node:util";
import { fileURLToPath } from "node:url";

// =====================================================================
// LANE 1: THE MELODY
// Pure: data in, data out. No I/O, no clock, no platform opinions.
// =====================================================================

const ASCII_WHITESPACE = new Set([0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x20]);

export function compose(seed, size, density) {
  // Byte-level ASCII trim and fold, per SPEC.md. String#trim/toLowerCase are Unicode-aware and would diverge.
  const bytes = new TextEncoder().encode(seed);
  let start = 0;
  let end = bytes.length;
  while (start < end && ASCII_WHITESPACE.has(bytes[start])) start++;
  while (end > start && ASCII_WHITESPACE.has(bytes[end - 1])) end--;

  let state = 0x811c9dc5;
  for (let i = start; i < end; i++) {
    const b = bytes[i] >= 0x41 && bytes[i] <= 0x5a ? bytes[i] + 0x20 : bytes[i];
    state = Math.imul(state ^ b, 0x01000193) >>> 0;
  }
  if (state === 0) state = 0x6d2b79f5;

  const next = () => {
    state ^= state << 13;
    state ^= state >>> 17;
    state ^= state << 5;
    return (state >>>= 0);
  };

  const hue = next() % 360;
  const palette = [[0, 0, 94], [hue, 65, 56], [(hue + 145) % 360, 75, 60]];

  const grid = new Array(size * size).fill(0);
  const half = Math.floor((size + 1) / 2);
  for (let y = 0; y < size; y++) {
    for (let x = 0; x < half; x++) {
      if (next() % 100 < density) {
        const color = next() % 100 < 28 ? 2 : 1;
        grid[y * size + x] = color;
        grid[y * size + (size - 1 - x)] = color;
      }
    }
  }

  return { size, palette, grid };
}

// =====================================================================
// LANE 2: THE INSTRUMENTS
// Each one reads a sheet and makes no musical decisions of its own.
// =====================================================================

export function rows(sheet) {
  const out = [];
  for (let i = 0; i < sheet.grid.length; i += sheet.size) {
    out.push(sheet.grid.slice(i, i + sheet.size).join(""));
  }
  return out;
}

export function renderJSON(sheet) {
  return JSON.stringify({ size: sheet.size, palette: sheet.palette, grid: rows(sheet) });
}

export function renderText(sheet) {
  const glyphs = { 0: "  ", 1: "██", 2: "▒▒" };
  return rows(sheet).map((row) => [...row].map((c) => glyphs[c]).join("")).join("\n");
}

const STAMPS = {
  square: '<rect id="px" width="1" height="1"/>',
  rounded: '<rect id="px" x="0.05" y="0.05" width="0.9" height="0.9" rx="0.25"/>',
  circle: '<circle id="px" cx="0.5" cy="0.5" r="0.45"/>',
};

export function renderSVG(sheet, shape = "rounded") {
  const hsl = ([h, s, l]) => `hsl(${h}, ${s}%, ${l}%)`;
  const view = sheet.size + 1;
  const out = [
    `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 ${view} ${view}" width="420" height="420">`,
    `  <defs>${STAMPS[shape]}</defs>`,
    `  <rect width="100%" height="100%" fill="${hsl(sheet.palette[0])}"/>`,
  ];
  for (const color of [1, 2]) {
    out.push(`  <g fill="${hsl(sheet.palette[color])}">`);
    sheet.grid.forEach((cell, i) => {
      if (cell !== color) return;
      const x = (i % sheet.size) + 0.5;
      const y = Math.floor(i / sheet.size) + 0.5;
      out.push(`    <use href="#px" x="${x.toFixed(1)}" y="${y.toFixed(1)}"/>`);
    });
    out.push("  </g>");
  }
  out.push("</svg>");
  return out.join("\n");
}

// =====================================================================
// LANE 3: THE STAGE
// The only lane that touches the outside world: argv, files, stdout.
// =====================================================================

const HANDSHAKE = ["00100", "01110", "12121", "12021", "01210"];
const same = (a, b) => JSON.stringify(a) === JSON.stringify(b);

function check(path) {
  const golden = JSON.parse(readFileSync(path, "utf8"));
  let failures = 0;
  for (const v of golden.vectors) {
    const sheet = compose(v.seed, v.size, v.density);
    const ok = same(rows(sheet), v.grid) && same(sheet.palette, v.palette);
    if (!ok) failures++;
    console.log(`${ok ? "PASS" : "FAIL"}  ${JSON.stringify(v.seed)} size=${v.size} density=${v.density}`);
  }
  console.log(`js: ${golden.vectors.length - failures}/${golden.vectors.length} vectors match ${golden.spec}`);
  return failures ? 1 : 0;
}

function main(argv) {
  const usage = "usage: node identicon.mjs [--size N] [--density P] [--format text|svg|json] [--shape square|rounded|circle] SEED\n       node identicon.mjs --check golden.json";
  const fail = (msg) => {
    console.error(`error: ${msg}\n${usage}`);
    return 2;
  };

  let parsed;
  try {
    parsed = parseArgs({
      args: argv,
      allowPositionals: true,
      options: {
        size: { type: "string", default: "5" },
        density: { type: "string", default: "52" },
        format: { type: "string", default: "text" },
        shape: { type: "string", default: "rounded" },
        check: { type: "string" },
      },
    });
  } catch (err) {
    return fail(err.message);
  }
  const { values: opts, positionals } = parsed;

  if (!same(rows(compose("captain-pixel", 5, 52)), HANDSHAKE)) {
    console.error("handshake failed: this build does not play the m-avatars melody");
    return 2;
  }

  if (opts.check) return check(opts.check);

  const size = Number(opts.size);
  const density = Number(opts.density);
  if (positionals.length !== 1) return fail("exactly one SEED is required");
  if (!Number.isInteger(size) || size < 1 || size > 64) return fail("--size must be between 1 and 64");
  if (!Number.isInteger(density) || density < 0 || density > 100) return fail("--density must be between 0 and 100");
  if (!["text", "svg", "json"].includes(opts.format)) return fail("--format must be text, svg, or json");
  if (!(opts.shape in STAMPS)) return fail("--shape must be square, rounded, or circle");

  const sheet = compose(positionals[0], size, density);
  const render = { svg: () => renderSVG(sheet, opts.shape), json: () => renderJSON(sheet), text: () => renderText(sheet) };
  console.log(render[opts.format]());
  return 0;
}

if (process.argv[1] === fileURLToPath(import.meta.url)) {
  process.exitCode = main(process.argv.slice(2));
}
