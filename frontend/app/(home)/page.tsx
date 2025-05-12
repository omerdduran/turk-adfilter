"use client";

import Link from 'next/link';
import { Shield, Download, Users, BookOpen, Code, Check, ExternalLink, Copy, ChevronDown } from 'lucide-react';
import React, { useState, useEffect } from 'react';

type SourceType = 'github' | 'codeberg';

type Metrics = {
  totalDomains: number;
  filterListDomains: number;
  hostsListDomains: number;
  lastUpdated: string;
};

export default function HomePage() {
  const [sharedListSource, setSharedListSource] = useState<SourceType>('github');
  const [metrics, setMetrics] = useState<Metrics | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [activeDropdown, setActiveDropdown] = useState<string | null>(null);

  const filterListLinks = {
    github: "https://raw.githubusercontent.com/omerdduran/turk-adfilter/main/turk-adfilter.txt",
    codeberg: "https://codeberg.org/omerdduran/turk-adfilter/raw/branch/main/turk-adfilter.txt"
  };
  const hostsListLinks = {
    github: "https://raw.githubusercontent.com/omerdduran/turk-adfilter/refs/heads/main/hosts.txt",
    codeberg: "https://codeberg.org/omerdduran/turk-adfilter/raw/branch/main/hosts.txt"
  };

  useEffect(() => {
    const fetchMetrics = async () => {
      try {
        const response = await fetch('/api/metrics');
        const data = await response.json();
        setMetrics(data);
      } catch (error) {
        console.error('Error fetching metrics:', error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchMetrics();
    // Refresh metrics every 15 minutes (matching the server-side cache)
    const interval = setInterval(fetchMetrics, 15 * 60 * 1000);
    return () => clearInterval(interval);
  }, []);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      const target = event.target as HTMLElement;
      if (!target.closest('.dropdown-container')) {
        setActiveDropdown(null);
      }
    };

    document.addEventListener('click', handleClickOutside);
    return () => document.removeEventListener('click', handleClickOutside);
  }, []);

  const toggleDropdown = (dropdownName: string) => {
    setActiveDropdown(activeDropdown === dropdownName ? null : dropdownName);
  };

  return (
    <div className="flex flex-col min-h-screen bg-gradient-to-b ">
      {/* Hero Section */}
      <header className="relative overflow-hidden py-20 md:py-32 px-4 min-h-screen">
        <div className="max-w-6xl mx-auto flex flex-col items-center">
          <div className="flex flex-col items-center mb-8">
            <img
              src="/assets/logo.png"
              alt="Turk-AdFilter Logo"
              width={90}
              height={90}
              className="rounded-full mb-4"
            />
            <h1 className="text-4xl md:text-5xl font-extrabold bg-clip-text text-transparent bg-gradient-to-r from-red-600 to-red-700 dark:from-red-500 dark:to-red-400">
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
          
          <div className="flex flex-wrap justify-center gap-4 text-sm text-gray-500 dark:text-gray-400 mb-12">
            <span className="flex items-center"><Check size={16} className="mr-1 text-green-500" /> %100 Açık Kaynak</span>
            <span className="flex items-center"><Check size={16} className="mr-1 text-green-500" /> Düzenli Güncelleme</span>
            <span className="flex items-center"><Check size={16} className="mr-1 text-green-500" /> Topluluk Destekli</span>
          </div>

          <div className="flex flex-col sm:flex-row gap-4 mb-16 justify-center items-center">
            {/* Dropdown Button for Filter List */}
            <div className="relative dropdown-container">
              <button
                type="button"
                onClick={() => toggleDropdown('browser')}
                className="bg-red-600 hover:bg-red-700 text-white px-8 py-4 rounded-lg font-semibold shadow-lg flex items-center justify-center transition-all w-full sm:w-auto"
              >
                <Download size={20} className="mr-2" />
                Tarayıcılar için Reklam Filtresi
                <ChevronDown size={20} className="ml-2" />
              </button>
              <div className={`absolute ${activeDropdown === 'browser' ? 'opacity-100 scale-100' : 'opacity-0 scale-95 pointer-events-none'} transform transition-all duration-200 ease-out top-full mt-px w-full bg-red-600 dark:bg-red-700 shadow-lg rounded-md py-1 z-20`}>
                <a
                  href={filterListLinks.github}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="block px-4 py-2 text-sm text-white hover:bg-red-700 dark:hover:bg-red-800"
                  onClick={() => setActiveDropdown(null)}
                >
                  GitHub
                </a>
                <a
                  href={filterListLinks.codeberg}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="block px-4 py-2 text-sm text-white hover:bg-red-700 dark:hover:bg-red-800"
                  onClick={() => setActiveDropdown(null)}
                >
                  Codeberg
                </a>
              </div>
            </div>
            <Link
              href="/docs"
              className="bg-white dark:bg-[#191919] text-gray-900 dark:text-white border border-gray-300 dark:border-gray-700 px-8 py-4 rounded-lg font-semibold shadow hover:shadow-md flex items-center justify-center transition-all"
            >
              <BookOpen size={20} className="mr-2" />
              Dokümantasyonu Gör
            </Link>
          </div>
          
          <div className="grid grid-cols-2 md:grid-cols-4 gap-8 mb-12 max-w-4xl mx-auto w-full text-center">
            <div>
              <div className="text-3xl md:text-4xl font-bold text-red-600 dark:text-red-500 mb-2">
                {isLoading ? (
                  <div className="animate-pulse bg-gray-200 dark:bg-gray-700 h-10 w-28 mx-auto rounded"></div>
                ) : (
                  metrics?.totalDomains.toLocaleString() ?? '0'
                )}
              </div>
              <p className="text-gray-600 dark:text-gray-400">Toplam Domain</p>
            </div>
            <div>
              <div className="text-3xl md:text-4xl font-bold text-red-600 dark:text-red-500 mb-2">
                {isLoading ? (
                  <div className="animate-pulse bg-gray-200 dark:bg-gray-700 h-10 w-28 mx-auto rounded"></div>
                ) : (
                  metrics?.filterListDomains.toLocaleString() ?? '0'
                )}
              </div>
              <p className="text-gray-600 dark:text-gray-400">Reklam Filtresi</p>
            </div>
            <div>
              <div className="text-3xl md:text-4xl font-bold text-red-600 dark:text-red-500 mb-2">
                {isLoading ? (
                  <div className="animate-pulse bg-gray-200 dark:bg-gray-700 h-10 w-28 mx-auto rounded"></div>
                ) : (
                  metrics?.hostsListDomains.toLocaleString() ?? '0'
                )}
              </div>
              <p className="text-gray-600 dark:text-gray-400">DNS Filtresi</p>
            </div>
            <div>
              <div className="text-3xl md:text-4xl font-bold text-red-600 dark:text-red-500 mb-2">%40</div>
              <p className="text-gray-600 dark:text-gray-400">Hız Artışı</p>
            </div>
          </div>
          
          {metrics?.lastUpdated && (
            <p className="text-sm text-gray-500 dark:text-gray-500 mb-12 text-center">
              Son güncelleme: {metrics.lastUpdated} <span className="ml-1">(GMT+3)</span>
            </p>
          )}
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
          
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6 min-w-0">
            <div className="bg-white dark:bg-[#191919] p-6 rounded-xl shadow min-w-0">
              <h3 className="text-xl font-bold mb-4 text-gray-800 dark:text-gray-200">uBlock Origin ile</h3>
              <ol className="list-decimal list-inside space-y-2 text-gray-600 dark:text-gray-400">
                <li>uBlock Origin ayarlarını açın</li>
                <li>Filtre listeleri sekmesine gidin</li>
                <li>"İçe aktar" butonuna tıklayın</li>
                <li>Aşağıdaki URL'lerden birini yapıştırın:</li>
              </ol>
              <div className="flex items-center gap-2 text-sm mt-3 mb-2">
                <span className="text-gray-600 dark:text-gray-400">Kaynak:</span>
                <button
                  onClick={() => setSharedListSource('github')}
                  className={`px-2 py-1 rounded text-xs font-medium ${sharedListSource === 'github' ? 'bg-red-600 text-white' : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-600'}`}
                >
                  GitHub
                </button>
                <button
                  onClick={() => setSharedListSource('codeberg')}
                  className={`px-2 py-1 rounded text-xs font-medium ${sharedListSource === 'codeberg' ? 'bg-[#2D81C9] text-white' : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-600'}`}
                >
                  Codeberg
                </button>
              </div>
              <CopyableLink link={filterListLinks[sharedListSource]} />
            </div>
            
            <div className="bg-white dark:bg-[#191919] p-6 rounded-xl shadow min-w-0">
              <h3 className="text-xl font-bold mb-4 text-gray-800 dark:text-gray-200">AdGuard ile</h3>
              <ol className="list-decimal list-inside space-y-2 text-gray-600 dark:text-gray-400">
                <li>AdGuard ayarlarını açın</li>
                <li>Filtreler bölümüne gidin</li>
                <li>"Özel filtre ekle" seçeneğini bulun</li>
                <li>Aşağıdaki URL'lerden birini ekleyin:</li>
              </ol>
              <div className="flex items-center gap-2 text-sm mt-3 mb-2">
                <span className="text-gray-600 dark:text-gray-400">Kaynak:</span>
                <button
                  onClick={() => setSharedListSource('github')}
                  className={`px-2 py-1 rounded text-xs font-medium ${sharedListSource === 'github' ? 'bg-red-600 text-white' : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-600'}`}
                >
                  GitHub
                </button>
                <button
                  onClick={() => setSharedListSource('codeberg')}
                  className={`px-2 py-1 rounded text-xs font-medium ${sharedListSource === 'codeberg' ? 'bg-[#2D81C9] text-white' : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-600'}`}
                >
                  Codeberg
                </button>
              </div>
              <CopyableLink link={filterListLinks[sharedListSource]} />
            </div>
            
            <div className="bg-white dark:bg-[#191919] p-6 rounded-xl shadow min-w-0">
              <h3 className="text-xl font-bold mb-4 text-gray-800 dark:text-gray-200">Pi-hole ile</h3>
              <ol className="list-decimal list-inside space-y-2 text-gray-600 dark:text-gray-400">
                <li>Pi-hole yönetim panelinize giriş yapın</li>
                <li>"Group Management" &gt; "Adlists" sekmesine gidin</li>
                <li>"Add a new adlist" kısmına aşağıdaki URL\'lerden birini ekleyin:</li>
                <li>"Add" butonuna tıklayın ve ardından "Update Gravity" ile güncellemeleri uygulayın</li>
              </ol>
              <div className="flex items-center gap-2 text-sm mt-3 mb-2">
                <span className="text-gray-600 dark:text-gray-400">Kaynak:</span>
                <button
                  onClick={() => setSharedListSource('github')}
                  className={`px-2 py-1 rounded text-xs font-medium ${sharedListSource === 'github' ? 'bg-red-600 text-white' : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-600'}`}
                >
                  GitHub
                </button>
                <button
                  onClick={() => setSharedListSource('codeberg')}
                  className={`px-2 py-1 rounded text-xs font-medium ${sharedListSource === 'codeberg' ? 'bg-[#2D81C9] text-white' : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-600'}`}
                >
                  Codeberg
                </button>
              </div>
              <CopyableLink link={hostsListLinks[sharedListSource]} />
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
      
      {/* CTA */}
      <section className="py-16 px-4 bg-gradient-to-r from-red-600 to-red-700 dark:from-red-800 dark:to-red-900 text-white">
        <div className="max-w-4xl mx-auto text-center">
          <h2 className="text-3xl font-bold mb-6">İnterneti Temiz Tut!</h2>
          <p className="mb-8 text-lg max-w-2xl mx-auto">
            Turk-AdFilter ile Türk web sitelerinde daha güvenli ve temiz bir internet deneyimi yaşayın.
            Şimdi filtre listemizi ekleyin ve internette gezinirken rahatsız edici reklamlardan, izleyicilerden kurtulun.
          </p>
          
          <div className="flex flex-col sm:flex-row gap-4 justify-center flex-wrap items-start">
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
            
            {/* Dropdown Button for Hosts List */}
            <div className="relative dropdown-container">
              <button
                type="button"
                onClick={() => toggleDropdown('dns')}
                className="bg-red-700 hover:bg-red-800 dark:bg-red-900 dark:hover:bg-red-950 text-white px-6 py-3 rounded-lg font-semibold shadow flex items-center justify-center border border-red-500 w-full sm:w-auto"
              >
                <Download size={20} className="mr-2" />
                DNS Filtresi
                <ChevronDown size={20} className="ml-2" />
              </button>
              <div className={`absolute ${activeDropdown === 'dns' ? 'opacity-100 scale-100' : 'opacity-0 scale-95 pointer-events-none'} transform transition-all duration-200 ease-out top-full mt-px w-full bg-red-700 dark:bg-red-800 shadow-lg rounded-md py-1 z-20`}>
                <a
                  href={hostsListLinks.github}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="block px-4 py-2 text-sm text-white hover:bg-red-800 dark:hover:bg-red-900"
                  onClick={() => setActiveDropdown(null)}
                >
                  GitHub
                </a>
                <a
                  href={hostsListLinks.codeberg}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="block px-4 py-2 text-sm text-white hover:bg-red-800 dark:hover:bg-red-900"
                  onClick={() => setActiveDropdown(null)}
                >
                  Codeberg
                </a>
              </div>
            </div>

            {/* Dropdown Button for RAW Filter List */}
            <div className="relative dropdown-container">
              <button
                type="button"
                onClick={() => toggleDropdown('raw')}
                className="bg-red-700 hover:bg-red-800 dark:bg-red-900 dark:hover:bg-red-950 text-white px-6 py-3 rounded-lg font-semibold shadow flex items-center justify-center border border-red-500 w-full sm:w-auto"
              >
                <Download size={20} className="mr-2" />
                Eklentiler için Filtre
                <ChevronDown size={20} className="ml-2" />
              </button>
              <div className={`absolute ${activeDropdown === 'raw' ? 'opacity-100 scale-100' : 'opacity-0 scale-95 pointer-events-none'} transform transition-all duration-200 ease-out top-full mt-px w-full bg-red-700 dark:bg-red-800 shadow-lg rounded-md py-1 z-20`}>
                <a
                  href={filterListLinks.github}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="block px-4 py-2 text-sm text-white hover:bg-red-800 dark:hover:bg-red-900"
                  onClick={() => setActiveDropdown(null)}
                >
                  GitHub
                </a>
                <a
                  href={filterListLinks.codeberg}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="block px-4 py-2 text-sm text-white hover:bg-red-800 dark:hover:bg-red-900"
                  onClick={() => setActiveDropdown(null)}
                >
                  Codeberg
                </a>
              </div>
            </div>
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

function CopyableLink({ link }: { link: string }) {
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(link);
      setCopied(true);
      setTimeout(() => setCopied(false), 1500);
    } catch (err) {
      setCopied(false);
    }
  };

  return (
    <div className="flex items-center gap-2 mt-4 min-w-0">
      <div
        className="flex-1 min-w-0 max-w-full p-2 md:p-1 bg-gray-100 dark:bg-[#121212] rounded text-sm md:text-xs font-mono text-gray-800 dark:text-gray-300 break-all overflow-x-auto whitespace-nowrap select-all truncate custom-scrollbar-x"
        tabIndex={0}
        title={link}
      >
        {link}
      </div>
      <button
        onClick={handleCopy}
        className="ml-1 p-2 bg-transparent hover:bg-gray-200 dark:hover:bg-[#222] rounded transition focus:outline-none focus:ring-2 focus:ring-red-400"
        type="button"
        aria-label="Kopyala"
      >
        <span className="sr-only">Kopyala</span>
        <Copy size={18} className={copied ? 'text-green-500' : 'text-gray-500'} />
        <span className="absolute opacity-0 pointer-events-none group-hover:opacity-100 group-focus:opacity-100 bg-gray-800 text-white text-xs rounded px-2 py-1 -mt-10 left-1/2 -translate-x-1/2 transition-opacity duration-200">
          {copied ? 'Kopyalandı!' : 'Kopyala'}
        </span>
      </button>
    </div>
  );
}