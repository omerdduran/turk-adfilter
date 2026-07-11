import './global.css';
import { RootProvider } from 'fumadocs-ui/provider';
import { Inter } from 'next/font/google';
import type { ReactNode } from 'react';
import type { Metadata } from 'next';

const inter = Inter({
  subsets: ['latin'],
});

const SITE_URL = 'https://reklamsiz-turkiye.com';
const TITLE = 'Turk-AdFilter — Türkçe Reklam & İzleyici Engelleme Listesi';
const DESCRIPTION =
  'Türkiye merkezli reklam, izleyici ve zararlı içerik sağlayıcılarını engelleyen topluluk tabanlı, açık kaynak filtre listesi. AdGuard, uBlock Origin, Pi-hole ve NextDNS ile uyumlu.';

export const metadata: Metadata = {
  title: {
    default: TITLE,
    template: '%s | Turk-AdFilter',
  },
  description: DESCRIPTION,
  keywords: [
    'reklam engelleme',
    'türkçe reklam engelleme listesi',
    'türkçe pihole listesi',
    'adblock',
    'türkiye',
    'ublock origin',
    'adguard',
    'pi-hole',
    'nextdns',
    'bahis engelleme',
    'phishing koruması',
  ],
  authors: [{ name: 'Ömer Duran', url: 'https://github.com/omerdduran' }],
  creator: 'Ömer Duran',
  publisher: 'Turk-AdFilter',
  metadataBase: new URL(SITE_URL),
  alternates: { canonical: '/' },
  robots: { index: true, follow: true },
  openGraph: {
    type: 'website',
    locale: 'tr_TR',
    url: SITE_URL,
    siteName: 'Turk-AdFilter',
    title: TITLE,
    description: DESCRIPTION,
    images: [{ url: '/assets/logo.png', alt: 'Turk-AdFilter' }],
  },
  twitter: {
    card: 'summary_large_image',
    title: TITLE,
    description: DESCRIPTION,
    images: ['/assets/logo.png'],
  },
};

const jsonLd = {
  '@context': 'https://schema.org',
  '@graph': [
    {
      '@type': 'Organization',
      '@id': `${SITE_URL}/#organization`,
      name: 'Turk-AdFilter',
      url: SITE_URL,
      logo: `${SITE_URL}/assets/logo.png`,
      sameAs: [
        'https://github.com/omerdduran/turk-adfilter',
        'https://codeberg.org/omerdduran/turk-adfilter',
      ],
    },
    {
      '@type': 'WebSite',
      '@id': `${SITE_URL}/#website`,
      url: SITE_URL,
      name: 'Turk-AdFilter',
      description: DESCRIPTION,
      inLanguage: 'tr-TR',
      publisher: { '@id': `${SITE_URL}/#organization` },
    },
    {
      '@type': 'SoftwareApplication',
      name: 'Turk-AdFilter',
      applicationCategory: 'SecurityApplication',
      operatingSystem: 'AdGuard, uBlock Origin, Pi-hole, AdGuard Home, NextDNS',
      description: DESCRIPTION,
      url: SITE_URL,
      inLanguage: 'tr-TR',
      isAccessibleForFree: true,
      license: 'https://www.gnu.org/licenses/gpl-3.0.html',
      offers: { '@type': 'Offer', price: '0', priceCurrency: 'USD' },
    },
  ],
};

export default function Layout({ children }: { children: ReactNode }) {
  return (
    <html lang="tr" className={inter.className} suppressHydrationWarning>
      <body className="flex flex-col min-h-screen">
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }}
        />
        <RootProvider>{children}</RootProvider>
      </body>
    </html>
  );
}
