import * as fs from 'node:fs/promises';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
// 📌 Ruta de los archivos
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename); // Esto apunta a /dist
// post-latest.json está en la misma carpeta (/dist)
const POSTS_FILE = path.resolve(__dirname, 'post-latest.json');
// README.md está un nivel arriba (en la raíz)
const README_FILE = path.resolve(__dirname, '../README.md');
// 📌 Marcadores en el README
const START_MARKER = '<!-- ARTICLES:START -->';
const END_MARKER = '<!-- ARTICLES:END -->';
const main = async () => {
    try {
        console.log('📥 Leyendo post-latest.json...');
        const rawData = await fs.readFile(POSTS_FILE, 'utf8');
        const articles = JSON.parse(rawData); // Usamos el tipo Article
        if (!Array.isArray(articles) || articles.length === 0) {
            console.warn('⚠️ No se encontraron artículos en post-latest.json.');
            return;
        }
        console.log('📄 Leyendo README.md...');
        const markdown = await fs.readFile(README_FILE, 'utf8');
        // Generar nuevo contenido para el README
        const newContent = generateArticlesContent(articles);
        const updatedMarkdown = replaceContentBetweenMarkers(markdown, START_MARKER, END_MARKER, newContent);
        console.log('💾 Guardando cambios en README.md...');
        await fs.writeFile(README_FILE, updatedMarkdown, 'utf8');
        console.log('✅ README.md actualizado con éxito.');
    }
    catch (error) {
        console.error('❌ Error durante la actualización del README:', error);
    }
};
// 📌 Genera el contenido en Markdown con los artículos
// const generateArticlesContent = (articles: Article[]): string => {
//   return articles
//     .map(article => `- [${article.title}](https://www.codegeekery.com/posts/${article.slug.current})\n  ![Image](${article.mainImage.asset.url})`)
//     .join('\n\n');
// };
export const generateArticlesContent = (articles) => {
    let markdown = `# Latest Blog Posts 📝\n\n`;
    articles?.forEach((article) => {
        markdown += `- [${article.title}](https://www.codegeekery.com/posts/${article.slug.current})\n`;
    });
    markdown += `\n[➡️ More blog posts](https://www.codegeekery.com/blog)`;
    return markdown;
};
// 📌 Reemplaza el contenido dentro de los marcadores en el README
const replaceContentBetweenMarkers = (markdown, startMarker, endMarker, newContent) => {
    const regex = new RegExp(`(${startMarker})([\\s\\S]*?)(${endMarker})`, 'g');
    return markdown.replace(regex, `$1\n${newContent}\n$3`);
};
main();
