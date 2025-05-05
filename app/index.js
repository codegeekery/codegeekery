import * as fs from 'fs/promises';
import got from 'got';


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

// async function fetchArticles() {
//   const res = await got('https://www.codegeekery.com/api/latest', {
//     headers: {
//       'X-CODEGEEKERY': "DECRDp4424bzqF27IBJFB3F460Nth39mzSDD8iAkQEYjqIBdolFl52lQMB4y62E1NsfvZiLf2FkI7CB7B41FD29F",
//     },
//     http2: true, // 👈 habilita HTTP/2 si el servidor lo permite
//     responseType: 'json',
//   });

//   return res.body.slice(0, 3);
// }

const WEBHOOK_URL = 'https://webhook.site/538f72a5-8094-4fbd-87a9-ca082401fad1';

async function fetchArticles() {
  try {
    const res = await got('https://www.codegeekery.com/api/latest', {
      headers: {
        'X-CODEGEEKERY': "DECRDp4424bzqF27IBJFB3F460Nth39mzSDD8iAkQEYjqIBdolFl52lQMB4y62E1NsfvZiLf2FkI7CB7B41FD29F",
      },
      http2: true,
      responseType: 'json',
    });

    return res.body.slice(0, 3);
  } catch (error) {
    // Enviar error al webhook
    await got.post(WEBHOOK_URL, {
      json: {
        message: 'Error fetching articles',
        error: {
          name: error.name,
          message: error.message,
          stack: error.stack,
        },
        timestamp: new Date().toISOString(),
      },
    }).catch(console.error); // Evita que un segundo error bloquee la ejecución

    // Relanzar el error si lo necesitas manejar más arriba
    throw error;
  }
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
