import Link from 'next/link';

export default function NotFound() {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen text-center p-4">
      <h1 className="text-4xl font-bold mb-4">Sayfa Bulunamadı</h1>
      <p className="text-lg mb-6">
        Aradığınız sayfa mevcut değil veya taşınmış olabilir.
      </p>
      <Link
        href="/"
        className="px-4 py-2 bg-fd-foreground text-black rounded-lg hover:bg-fd-primary transition-colors"
      >
        Ana Sayfaya Dön
      </Link>
    </div>
  );
} 