import { readFile, writeFile } from 'fs/promises';
import { connect } from 'http2';
import { URL } from 'url';

const README_PATH = 'README.md';
const START = '<!-- ARTICLES:START -->';
const END = '<!-- ARTICLES:END -->';

const SECRET_KEY = process.env.SECRET_KEY;
const USER_AGENT = process.env.USER_AGENT;
const BASE_URL = 'https://www.codegeekery.com/blog';
const BASE_POST_URL = 'https://www.codegeekery.com/posts/';
const TIMESTAMP_REGEX = /<!-- Última actualización: .*?-->/;

if (!SECRET_KEY) {
  console.error('SECRET_KEY no está definido en las variables de entorno.');
  process.exit(1);
}

async function fetchArticles() {
  return new Promise((resolve, reject) => {
    const url = new URL('https://www.codegeekery.com/api/latest');
    const client = connect(url.origin);

    const req = client.request({
      ':method': 'GET',
      ':path': url.pathname,
      'X-CODEGEEKERY': SECRET_KEY,
      'User-Agent': USER_AGENT,
    });

    let data = '';
    req.setEncoding('utf8');

    req.on('data', chunk => {
      data += chunk;
    });

    req.on('end', () => {
      client.close();
      try {
        const json = JSON.parse(data);
        resolve(json.slice(0, 3));
      } catch (err) {
        reject(new Error('Error parseando JSON: ' + err.message));
      }
    });

    req.on('error', err => {
      client.close();
      reject(err);
    });

    req.end();
  });
}

async function updateReadme(articles) {
  const content = await readFile(README_PATH, 'utf-8');

  const rows = [];
  const columnsPerRow = 3;

  for (let i = 0; i < articles.length; i += columnsPerRow) {
    const rowArticles = articles.slice(i, i + columnsPerRow);

    const imageRow = rowArticles.map(post => {
      const imageUrl = `${post.mainImage.asset.url}?w=200&h=200`;
      const link = `${BASE_POST_URL}${post.slug.current}`;
      return `[![${post.title}](${imageUrl})](${link})`;
    }).join(' | ');

    const titleRow = rowArticles.map(post => {
      const link = `${BASE_POST_URL}${post.slug.current}`;
      return `**[${post.title}](${link})**`;
    }).join(' | ');

    const separator = rowArticles.map(() => '---').join(' | ');

    rows.push(imageRow, separator, titleRow);
  }

  const tableMarkdown = rows.join('\n') + `\n\n[➡️ More blog posts](${BASE_URL})`;

  let updated = content.replace(
    new RegExp(`${START}[\\s\\S]*?${END}`),
    `${START}\n${tableMarkdown}\n${END}`
  );

  // 🔁 Añadir o actualizar el timestamp
  const timestamp = new Date().toISOString();
  const TIMESTAMP_REGEX = /<!-- Última actualización: .*?-->/;
  const newTimestampLine = `<!-- Última actualización: ${timestamp} -->`;

  if (TIMESTAMP_REGEX.test(updated)) {
    updated = updated.replace(TIMESTAMP_REGEX, newTimestampLine);
  } else {
    updated += `\n${newTimestampLine}\n`;
  }

  await writeFile(README_PATH, updated);
}

try {
  const articles = await fetchArticles();
  await updateReadme(articles);
} catch (err) {
  console.error('Error actualizando README:', err);
  process.exit(1);
}
