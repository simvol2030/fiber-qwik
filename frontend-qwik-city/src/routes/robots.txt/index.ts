import type { RequestHandler } from "@builder.io/qwik-city";

export const onGet: RequestHandler = async ({ send }) => {
  const body = `User-agent: *
Allow: /
Disallow: /admin/
Disallow: /dashboard/
Disallow: /login
Disallow: /register

Sitemap: https://example.com/sitemap.xml
`;

  send(200, body);
};
