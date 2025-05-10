import Link from 'next/link';
import { Shield, Download, Users, BookOpen, Code, Check, ExternalLink } from 'lucide-react';

export default function HomePage() {
  return (
    <div className="flex flex-col min-h-screen bg-gradient-to-b ">
      {/* Hero Section */}
      <header className="relative overflow-hidden py-20 md:py-32 px-4 min-h-screen">
        <div className="max-w-6xl mx-auto flex flex-col items-center">
          <div className="flex items-center mb-6">
            <img
              src="/assets/logo.png"
              alt="Turk-AdFilter Logo"
              width={90}
              height={90}
              className="rounded-full"
            />
            <h1 className="ml-4 text-4xl md:text-5xl font-extrabold bg-clip-text text-transparent bg-gradient-to-r from-red-600 to-red-700 dark:from-red-500 dark:to-red-400">
              Turk-AdFilter
            </h1>
          </div>
          
          <p className="text-xl md:text-2xl font-medium text-center max-w-2xl mb-8 text-gray-700 dark:text-gray-300">
            Türkiye merkezli reklam, izleyici ve zararlı içerik engelleyici
          </p>
          
          <p className="mb-10 max-w-2xl text-center text-gray-600 dark:text-gray-400">
            Topluluk tabanlı filtre listemiz, Pi-hole, AdGuard, uBlock Origin ve benzeri servislerle
            uyumlu şekilde çalışarak Türk web sitelerindeki reklamları ve izleyicileri engeller.
          </p>
          
          <div className="flex flex-col sm:flex-row gap-4 mb-12">
            <a
              href="https://raw.githubusercontent.com/omerdduran/turk-adfilter/main/turk-adfilter.txt"
              target="_blank"
              rel="noopener noreferrer"
              className="bg-red-600 hover:bg-red-700 text-white px-8 py-4 rounded-lg font-semibold shadow-lg flex items-center justify-center transition-all"
            >
              <Download size={20} className="mr-2" />
              Filtre Listesini İndir
            </a>
            <Link
              href="/docs"
              className="bg-white dark:bg-[#191919] text-gray-900 dark:text-white border border-gray-300 dark:border-gray-700 px-8 py-4 rounded-lg font-semibold shadow hover:shadow-md flex items-center justify-center transition-all"
            >
              <BookOpen size={20} className="mr-2" />
              Dokümantasyonu Gör
            </Link>
          </div>
          
          <div className="flex flex-wrap justify-center gap-4 text-sm text-gray-500 dark:text-gray-400">
            <span className="flex items-center"><Check size={16} className="mr-1 text-green-500" /> %100 Açık Kaynak</span>
            <span className="flex items-center"><Check size={16} className="mr-1 text-green-500" /> Düzenli Güncelleme</span>
            <span className="flex items-center"><Check size={16} className="mr-1 text-green-500" /> Topluluk Destekli</span>
          </div>
        </div>
      </header>

      {/* Features */}
      <section className="py-16 px-4 bg-white dark:bg-[#191919]">
        <div className="max-w-6xl mx-auto">
          <h2 className="text-3xl font-bold text-center mb-16 text-gray-800 dark:text-gray-200">
            Neden Turk-AdFilter?
          </h2>
          
          <div className="grid md:grid-cols-3 gap-8">
            <div className="bg-gray-50 dark:bg-[#121212] p-6 rounded-xl shadow-sm">
              <div className="bg-red-100 dark:bg-red-900/30 p-3 rounded-full w-12 h-12 flex items-center justify-center mb-4">
                <Shield size={24} className="text-red-600 dark:text-red-400" />
              </div>
              <h3 className="text-xl font-bold mb-2 text-gray-800 dark:text-gray-200">Türkiye Odaklı Koruma</h3>
              <p className="text-gray-600 dark:text-gray-400">
                Özellikle Türk web sitelerindeki rahatsız edici reklamları ve izleyicileri engeller.
              </p>
            </div>
            
            <div className="bg-gray-50 dark:bg-[#121212] p-6 rounded-xl shadow-sm">
              <div className="bg-red-100 dark:bg-red-900/30 p-3 rounded-full w-12 h-12 flex items-center justify-center mb-4">
                <Users size={24} className="text-red-600 dark:text-red-400" />
              </div>
              <h3 className="text-xl font-bold mb-2 text-gray-800 dark:text-gray-200">Topluluk Gücü</h3>
              <p className="text-gray-600 dark:text-gray-400">
                Filtremiz sürekli gelişiyor ve topluluğun katkılarıyla büyüyor.
              </p>
            </div>
            
            <div className="bg-gray-50 dark:bg-[#121212] p-6 rounded-xl shadow-sm">
              <div className="bg-red-100 dark:bg-red-900/30 p-3 rounded-full w-12 h-12 flex items-center justify-center mb-4">
                <Code size={24} className="text-red-600 dark:text-red-400" />
              </div>
              <h3 className="text-xl font-bold mb-2 text-gray-800 dark:text-gray-200">Geniş Uyumluluk</h3>
              <p className="text-gray-600 dark:text-gray-400">
                Pi-hole, AdGuard, uBlock Origin ve diğer popüler araçlarla uyumlu.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* How to Use */}
      <section className="py-16 px-4 ">
        <div className="max-w-6xl mx-auto">
          <h2 className="text-3xl font-bold text-center mb-16 text-gray-800 dark:text-gray-200">
            Nasıl Kullanılır
          </h2>
          
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
            <div className="bg-white dark:bg-[#191919] p-6 rounded-xl shadow">
              <h3 className="text-xl font-bold mb-4 text-gray-800 dark:text-gray-200">uBlock Origin ile</h3>
              <ol className="list-decimal list-inside space-y-2 text-gray-600 dark:text-gray-400">
                <li>uBlock Origin ayarlarını açın</li>
                <li>Filtre listeleri sekmesine gidin</li>
                <li>"İçe aktar" butonuna tıklayın</li>
                <li>Aşağıdaki URL'yi yapıştırın:</li>
              </ol>
              <div className="mt-4 p-3 bg-gray-100 dark:bg-[#121212] rounded text-sm font-mono text-gray-800 dark:text-gray-300 break-all">
                https://raw.githubusercontent.com/omerdduran/turk-adfilter/main/turk-adfilter.txt
              </div>
            </div>
            
            <div className="bg-white dark:bg-[#191919] p-6 rounded-xl shadow">
              <h3 className="text-xl font-bold mb-4 text-gray-800 dark:text-gray-200">AdGuard ile</h3>
              <ol className="list-decimal list-inside space-y-2 text-gray-600 dark:text-gray-400">
                <li>AdGuard ayarlarını açın</li>
                <li>Filtreler bölümüne gidin</li>
                <li>"Özel filtre ekle" seçeneğini bulun</li>
                <li>Aşağıdaki URL'yi ekleyin:</li>
              </ol>
              <div className="mt-4 p-3 bg-gray-100 dark:bg-[#121212] rounded text-sm font-mono text-gray-800 dark:text-gray-300 break-all">
                https://raw.githubusercontent.com/omerdduran/turk-adfilter/main/turk-adfilter.txt
              </div>
            </div>
            
            <div className="bg-white dark:bg-[#191919] p-6 rounded-xl shadow">
              <h3 className="text-xl font-bold mb-4 text-gray-800 dark:text-gray-200">Pi-hole ile</h3>
              <ol className="list-decimal list-inside space-y-2 text-gray-600 dark:text-gray-400">
                <li>Pi-hole yönetim panelinize giriş yapın</li>
                <li>"Group Management" &gt; "Adlists" sekmesine gidin</li>
                <li>"Add a new adlist" kısmına aşağıdaki URL'yi ekleyin:</li>
                <li>"Add" butonuna tıklayın ve ardından "Update Gravity" ile güncellemeleri uygulayın</li>
              </ol>
              <div className="mt-4 p-3 bg-gray-100 dark:bg-[#121212] rounded text-sm font-mono text-gray-800 dark:text-gray-300 break-all">
                https://raw.githubusercontent.com/omerdduran/turk-adfilter/main/turk-adfilter.txt
              </div>
            </div>
          </div>
          
          <div className="mt-10 text-center">
            <Link
              href="/docs/kurulum"
              className="inline-flex items-center text-red-600 dark:text-red-400 font-medium hover:underline"
            >
              Tüm kurulum talimatlarını görüntüle
              <ExternalLink size={16} className="ml-1" />
            </Link>
          </div>
        </div>
      </section>

      {/* Stats */}
      <section className="py-16 px-4 bg-white dark:bg-[#191919]">
        <div className="max-w-5xl mx-auto">
          <div className="grid md:grid-cols-3 gap-8 text-center">
            <div>
              <div className="text-4xl font-bold text-red-600 dark:text-red-500 mb-2">10,000+</div>
              <p className="text-gray-600 dark:text-gray-400">Engellenen Domain</p>
            </div>
            <div>
              <div className="text-4xl font-bold text-red-600 dark:text-red-500 mb-2">5,000+</div>
              <p className="text-gray-600 dark:text-gray-400">Aktif Kullanıcı</p>
            </div>
            <div>
              <div className="text-4xl font-bold text-red-600 dark:text-red-500 mb-2">%40</div>
              <p className="text-gray-600 dark:text-gray-400">Daha Hızlı Sayfa Yükleme</p>
            </div>
          </div>
        </div>
      </section>
      
      {/* CTA */}
      <section className="py-16 px-4 bg-gradient-to-r from-red-600 to-red-700 dark:from-red-800 dark:to-red-900 text-white">
        <div className="max-w-4xl mx-auto text-center">
          <h2 className="text-3xl font-bold mb-6">İnterneti Temiz Tut!</h2>
          <p className="mb-8 text-lg max-w-2xl mx-auto">
            Turk-AdFilter ile Türk web sitelerinde daha güvenli ve temiz bir internet deneyimi yaşayın.
            Şimdi filtre listemizi ekleyin ve internette gezinirken rahatsız edici reklamlardan, izleyicilerden kurtulun.
          </p>
          
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <a
              href="https://github.com/omerdduran/turk-adfilter"
              target="_blank"
              rel="noopener noreferrer"
              className="bg-white text-red-600 hover:bg-gray-100 px-6 py-3 rounded-lg font-semibold shadow flex items-center justify-center"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" className="mr-2" viewBox="0 0 16 16">
                <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.012 8.012 0 0 0 16 8c0-4.42-3.58-8-8-8z"/>
              </svg>
              GitHub'da Katkıda Bulun
            </a>
            
            <a
              href="https://raw.githubusercontent.com/omerdduran/turk-adfilter/main/turk-adfilter.txt"
              target="_blank"
              rel="noopener noreferrer"
              className="bg-red-700 hover:bg-red-800 dark:bg-red-900 dark:hover:bg-red-950 text-white px-6 py-3 rounded-lg font-semibold shadow flex items-center justify-center border border-red-500"
            >
              <Download size={20} className="mr-2" />
              Filtre Listesini İndir
            </a>
          </div>
        </div>
      </section>

      {/* Footer */}
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
              <span className="ml-2 font-semibold text-gray-800 dark:text-gray-200">Turk-AdFilter</span>
            </div>
            
            <div className="flex gap-6">
              <Link href="/docs" className="text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200">
                Dokümantasyon
              </Link>
              <Link href="/about" className="text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200">
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
            <p>© {new Date().getFullYear()} Turk-AdFilter. Tüm hakları saklıdır.</p>
            <p className="mt-2">Turk-AdFilter bir açık kaynak projesidir ve GPL-3.0 lisansı altında dağıtılmaktadır.</p>
          </div>
        </div>
      </footer>
    </div>
  );
}