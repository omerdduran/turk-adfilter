import Link from 'next/link';

export default function HomePage() {
  return (
    <main className="flex flex-1 flex-col items-center justify-center text-center py-16 px-4">
      <img
        src="/assets/logo.png"
        alt="Turk-AdFilter Logo"
        width={120}
        height={120}
        className="mb-6 mx-auto rounded-full shadow-lg"
      />
      <h1 className="mb-4 text-4xl font-extrabold text-fd-foreground">🇹🇷 Turk-AdFilter</h1>
      <p className="mb-6 max-w-xl text-lg text-fd-muted-foreground mx-auto">
        Türkiye merkezli reklam, izleyici ve zararlı içerik sağlayıcılarını engelleyen topluluk tabanlı bir filtre listesi. NextDNS, AdGuard, uBlock Origin ve benzeri servislerle uyumlu şekilde çalışır.
      </p>
      <div className="flex flex-col sm:flex-row gap-4 justify-center mb-8">
        <Link
          href="/docs"
          className="bg-fd-foreground text-black px-6 py-3 rounded-lg font-semibold shadow hover:bg-fd-primary transition"
        >
          Dokümantasyonu Gör
        </Link>
        <a
          href="https://raw.githubusercontent.com/omerdduran/turk-adfilter/main/turk-adfilter.txt"
          target="_blank"
          rel="noopener noreferrer"
          className="bg-fd-muted-foreground text-black px-6 py-3 rounded-lg font-semibold shadow hover:bg-fd-primary/80 transition"
        >
          Filtre Listesini Görüntüle
        </a>
      </div>
      <p className="text-fd-muted-foreground text-sm mt-4">
        <strong>İnterneti temiz tut!</strong>
      </p>
    </main>
  );
}
