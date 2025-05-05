import * as fs from 'fs/promises';
import * as http2 from 'node:http2';


const README_PATH = 'README.md';
const START = '<!-- ARTICLES:START -->';
const END = '<!-- ARTICLES:END -->';

const SECRET_KEY = process.env.SECRET_KEY;
const BASE_URL = 'https://www.codegeekery.com/blog';
const BASE_POST_URL = 'https://www.codegeekery.com/posts/';

if (!SECRET_KEY) {
  console.error('SECRET_KEY no está definido en las variables de entorno.');
  process.exit(1);
}

function fetchArticles() {
  return new Promise((resolve, reject) => {
    const client = http2.connect('https://www.codegeekery.com');

    client.on('error', reject);

    const req = client.request({
      ':method': 'GET',
      ':path': '/api/latest',
      'X-CODEGEEKERY': SECRET_KEY,
    });

    let data = '';

    req.setEncoding('utf8');
    req.on('data', chunk => data += chunk);
    req.on('end', () => {
      try {
        const articles = JSON.parse(data);
        client.close();
        resolve(articles.slice(0, 3));
      } catch (err) {
        reject(new Error('Error al parsear JSON:' + err + '\nContenido:\n' + data.slice(0, 200)));
      }
    });

    req.end();
  });
}

async function updateReadme(articles) {
  const content = await fs.readFile(README_PATH, 'utf-8');

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

  const updated = content.replace(
    new RegExp(`${START}[\\s\\S]*?${END}`),
    `${START}\n${tableMarkdown}\n${END}`
  );

  await fs.writeFile(README_PATH, updated);
}

(async () => {
  try {
    const articles = await fetchArticles();
    await updateReadme(articles);
  } catch (err) {
    console.error('Error actualizando README:', err);
    process.exit(1);
  }
})();
