import type { MetadataRoute } from 'next';
import { source } from '@/lib/source';

const SITE_URL = 'https://reklamsiz-turkiye.com';

export default function sitemap(): MetadataRoute.Sitemap {
  const now = new Date();

  const staticRoutes: MetadataRoute.Sitemap = [
    { url: `${SITE_URL}`, lastModified: now, changeFrequency: 'weekly', priority: 1 },
    { url: `${SITE_URL}/docs`, lastModified: now, changeFrequency: 'weekly', priority: 0.8 },
    { url: `${SITE_URL}/about`, lastModified: now, changeFrequency: 'monthly', priority: 0.6 },
    { url: `${SITE_URL}/releases`, lastModified: now, changeFrequency: 'weekly', priority: 0.7 },
    { url: `${SITE_URL}/feedback`, lastModified: now, changeFrequency: 'monthly', priority: 0.6 },
    { url: `${SITE_URL}/backups`, lastModified: now, changeFrequency: 'weekly', priority: 0.5 },
  ];

  const docsRoutes: MetadataRoute.Sitemap = source.getPages().map((page) => ({
    url: `${SITE_URL}${page.url}`,
    lastModified: now,
    changeFrequency: 'monthly',
    priority: 0.5,
  }));

  return [...staticRoutes, ...docsRoutes];
}
