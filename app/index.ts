import fs from 'fs/promises';
import http2 from 'http2';

const API_URL = 'https://www.codegeekery.com/api/latest';
const README_PATH = 'README.md';
const START = '<!-- ARTICLES:START -->';
const END = '<!-- ARTICLES:END -->';

const SECRET_KEY = process.env.SECRET_KEY!;
const BASE_URL = 'https://www.codegeekery.com/blog';
const BASE_POST_URL = 'https://www.codegeekery.com/posts/';

interface Article {
  _id: string;
  title: string;
  slug: {
    current: string;
    _type: string;
  };
  mainImage: {
    asset: {
      url: string;
    };
  };
}

function fetchArticles(): Promise<Article[]> {
  return new Promise((resolve, reject) => {
    const client = http2.connect('https://www.codegeekery.com');

    client.on('error', reject);

    const req = client.request({
      ':method': 'GET',
      ':path': '/api/latest',
      'X-CODEGEEKERY': "DECRDp4424bzqF27IBJFB3F460Nth39mzSDD8iAkQEYjqIBdolFl52lQMB4y62E1NsfvZiLf2FkI7CB7B41FD29F",
    });

    let data = '';

    req.setEncoding('utf8');
    req.on('data', chunk => data += chunk);
    req.on('end', () => {
      try {
        const articles: Article[] = JSON.parse(data);
        client.close();
        resolve(articles.slice(0, 3));
      } catch (err) {
        reject(new Error('Error al parsear JSON: ' + err + '\nContenido:\n' + data.slice(0, 200)));
      }
    });

    req.end();
  });
}

async function updateReadme(articles: Article[]) {
  const content = await fs.readFile(README_PATH, 'utf-8');

  const articleEntries = articles.map(post => {
    const baseImageUrl = post.mainImage.asset.url;
    const imageUrl = `${baseImageUrl}?w=200&h=200`; // Agregado w=200&h=200
    const link = `${BASE_POST_URL}${post.slug.current}`;
    return `[![${post.title}](${imageUrl})](${link})\n**[${post.title}](${link})**`;
  }).join('\n\n');

  const moreLink = `[➡️ More blog posts](${BASE_URL})`;

  const newSection = `${articleEntries}\n\n${moreLink}`;

  const updated = content.replace(
    new RegExp(`${START}[\\s\\S]*?${END}`),
    `${START}\n${newSection}\n${END}`
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
