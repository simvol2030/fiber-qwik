import type { RequestHandler } from "@builder.io/qwik-city";

export const onGet: RequestHandler = async ({ send }) => {
  const baseUrl = "https://example.com";
  const pages = [
    { url: "/", priority: "1.0", changefreq: "daily" },
    { url: "/login", priority: "0.5", changefreq: "monthly" },
    { url: "/register", priority: "0.5", changefreq: "monthly" },
  ];

  const sitemap = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
${pages
  .map(
    (page) => `  <url>
    <loc>${baseUrl}${page.url}</loc>
    <changefreq>${page.changefreq}</changefreq>
    <priority>${page.priority}</priority>
  </url>`
  )
  .join("\n")}
</urlset>`;

  send(200, sitemap);
};
