import Link from "next/link";

export default function Footer() {
  return (
    <footer className="py-12 px-4">
      <div className="max-w-6xl mx-auto">
        <div className="flex flex-col md:flex-row justify-between items-center">
          <div className="flex items-center mb-6 md:mb-0">
            <img
              src="/assets/logo.png"
              alt="Turk-AdFilter Logo"
              width={40}
              height={40}
              className="rounded-full shadow-sm"
            />
            <span className="ml-2 font-semibold text-gray-800 dark:text-gray-200">
              Turk-AdFilter
            </span>
          </div>

          <div className="flex gap-6">
            <Link
              href="/docs"
              className="text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200"
            >
              Dokümantasyon
            </Link>
            <Link
              href="/about"
              className="text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200"
            >
              Hakkımızda
            </Link>
            <a
              href="https://github.com/omerdduran/turk-adfilter"
              target="_blank"
              rel="noopener noreferrer"
              className="text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200"
            >
              GitHub
            </a>
          </div>
        </div>

        <div className="mt-8 pt-8 border-t border-gray-200 dark:border-gray-800 text-center text-sm text-gray-500 dark:text-gray-500">
          <p>
            © {new Date().getFullYear()} Turk-AdFilter. Tüm hakları saklıdır.
          </p>
          <p className="mt-2">
            Turk-AdFilter bir açık kaynak projesidir ve GPL-3.0 lisansı altında
            dağıtılmaktadır.
          </p>
        </div>
      </div>
    </footer>
  );
}
