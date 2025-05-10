import './global.css';
import { RootProvider } from 'fumadocs-ui/provider';
import { Inter } from 'next/font/google';
import type { ReactNode } from 'react';
import type { Metadata } from 'next';

const inter = Inter({
  subsets: ['latin'],
});

export const metadata: Metadata = {
  title: "Turk-AdFilter",
  description: 'Türkiye merkezli reklam ve izleyici engelleme listesi',
  keywords: ['reklam engelleme', 'adblock', 'türkiye', 'ublock origin', 'adguard'],
  authors: [{ name: 'Ömer Duran' }],
  creator: 'Ömer Duran',
  publisher: 'Turk-AdFilter',
  metadataBase: new URL('https://reklamsiz-turkiye.com'),
  robots: {
    index: true,
    follow: true,
  },
};

export default function Layout({ children }: { children: ReactNode }) {
  return (
    <html lang="tr" className={inter.className} suppressHydrationWarning>
      <body className="flex flex-col min-h-screen">
        <RootProvider>{children}</RootProvider>
      </body>
    </html>
  );
}
