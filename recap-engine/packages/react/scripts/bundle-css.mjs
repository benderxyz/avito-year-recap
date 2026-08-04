import { mkdir, readFile, writeFile } from 'node:fs/promises';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const pkgRoot = path.resolve(__dirname, '..');
const entry = path.join(pkgRoot, 'src/styles.css');
const outFile = path.join(pkgRoot, 'dist/styles.css');

const importRe = /@import\s+['"](.+?)['"]\s*;/g;

async function bundleCss(filePath, seen = new Set()) {
  const abs = path.resolve(filePath);
  if (seen.has(abs)) return '';
  seen.add(abs);

  const source = await readFile(abs, 'utf8');
  let result = '';
  let lastIndex = 0;

  for (const match of source.matchAll(importRe)) {
    const [statement, spec] = match;
    const index = match.index ?? 0;
    result += source.slice(lastIndex, index);
    lastIndex = index + statement.length;

    const imported = path.resolve(path.dirname(abs), spec);
    result += await bundleCss(imported, seen);
    result += '\n';
  }

  result += source.slice(lastIndex);
  return result;
}

const css = (await bundleCss(entry)).replace(/\n{3,}/g, '\n\n').trimStart();

await mkdir(path.dirname(outFile), { recursive: true });
await writeFile(outFile, `${css}\n`, 'utf8');
console.log(`bundled css → ${path.relative(pkgRoot, outFile)}`);
