import { Github, Users, Shield, Code, Heart } from 'lucide-react';
import Link from 'next/link';

export default function AboutPage() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-background to-muted/20">
      <main className="max-w-4xl mx-auto px-4 py-16">
        <div className="text-center mb-16">
          <h1 className="text-4xl font-bold mb-4 bg-gradient-to-r from-red-600 to-red-700 dark:from-red-500 dark:to-red-400 bg-clip-text text-transparent">
            Turk-AdFilter Hakkında
          </h1>
          <p className="text-lg text-muted-foreground">
            Türk internet kullanıcıları için özel olarak tasarlanmış reklam engelleme projesi
          </p>
        </div>

        <div className="space-y-12">
          {/* Proje Hikayesi */}
          <section>
            <h2 className="text-2xl font-semibold mb-4 text-foreground">Proje Hikayesi</h2>
            <div className="prose dark:prose-invert max-w-none">
              <p>
                Turk-AdFilter, Türkiye'deki internet kullanıcılarının karşılaştığı reklam ve izleme sorunlarına çözüm olarak 2025 yılında başlatıldı. 
                Proje, Türk web sitelerindeki reklamları ve izleyicileri hedefleyen, topluluk tarafından desteklenen açık kaynaklı bir filtre listesi olarak geliştirildi.
              </p>
              <p>
                Amacımız, Türk internet kullanıcılarına daha temiz, daha hızlı ve daha güvenli bir internet deneyimi sunmaktır. 
                Bu hedefe ulaşmak için sürekli olarak filtre listemizi güncelliyor ve geliştiriyoruz.
              </p>
            </div>
          </section>

          {/* Misyon ve Vizyon */}
          <section>
            <h2 className="text-2xl font-semibold mb-4 text-foreground">Misyon ve Vizyon</h2>
            <div className="grid md:grid-cols-2 gap-6">
              <div className="bg-card p-6 rounded-lg shadow-sm border">
                <h3 className="text-xl font-semibold mb-3 text-foreground">Misyonumuz</h3>
                <p className="text-muted-foreground">
                  Türk internet kullanıcılarına reklamsız, izleyicisiz ve daha güvenli bir internet deneyimi sunmak.
                </p>
              </div>
              <div className="bg-card p-6 rounded-lg shadow-sm border">
                <h3 className="text-xl font-semibold mb-3 text-foreground">Vizyonumuz</h3>
                <p className="text-muted-foreground">
                  Türkiye'nin en kapsamlı ve güncel reklam engelleme listesi olmak ve topluluk katkılarıyla sürekli gelişmek.
                </p>
              </div>
            </div>
          </section>

          {/* Değerlerimiz */}
          <section>
            <h2 className="text-2xl font-semibold mb-4 text-foreground">Değerlerimiz</h2>
            <div className="grid md:grid-cols-3 gap-6">
              <div className="bg-card p-6 rounded-lg shadow-sm border">
                <div className="bg-red-100 dark:bg-red-900/30 p-3 rounded-full w-12 h-12 flex items-center justify-center mb-4">
                  <Shield className="text-red-600 dark:text-red-400" />
                </div>
                <h3 className="text-lg font-semibold mb-2 text-foreground">Güvenlik</h3>
                <p className="text-muted-foreground">
                  Kullanıcılarımızın güvenliği ve gizliliği her zaman önceliğimizdir.
                </p>
              </div>
              <div className="bg-card p-6 rounded-lg shadow-sm border">
                <div className="bg-red-100 dark:bg-red-900/30 p-3 rounded-full w-12 h-12 flex items-center justify-center mb-4">
                  <Users className="text-red-600 dark:text-red-400" />
                </div>
                <h3 className="text-lg font-semibold mb-2 text-foreground">Topluluk</h3>
                <p className="text-muted-foreground">
                  Topluluk katkıları ve geri bildirimleri projemizin temelidir.
                </p>
              </div>
              <div className="bg-card p-6 rounded-lg shadow-sm border">
                <div className="bg-red-100 dark:bg-red-900/30 p-3 rounded-full w-12 h-12 flex items-center justify-center mb-4">
                  <Code className="text-red-600 dark:text-red-400" />
                </div>
                <h3 className="text-lg font-semibold mb-2 text-foreground">Açıklık</h3>
                <p className="text-muted-foreground">
                  Açık kaynak felsefesiyle şeffaf ve katılımcı bir geliştirme süreci.
                </p>
              </div>
            </div>
          </section>

          {/* Katkıda Bulunma */}
          <section>
            <h2 className="text-2xl font-semibold mb-4 text-foreground">Katkıda Bulunma</h2>
            <div className="bg-card p-6 rounded-lg shadow-sm border">
              <p className="text-muted-foreground mb-4">
                Turk-AdFilter, topluluk katkılarıyla büyüyen bir projedir. Projeye katkıda bulunmak için:
              </p>
              <ul className="list-disc list-inside space-y-2 text-muted-foreground mb-6">
                <li>GitHub üzerinden issue açabilirsiniz</li>
                <li>Pull request gönderebilirsiniz</li>
                <li>Yeni domain önerilerinde bulunabilirsiniz</li>
                <li>Dokümantasyona katkı sağlayabilirsiniz</li>
              </ul>
              <a
                href="https://github.com/omerdduran/turk-adfilter"
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center text-red-600 dark:text-red-400 hover:underline"
              >
                <Github className="mr-2" />
                GitHub'da Katkıda Bulun
              </a>
            </div>
          </section>

          {/* İletişim */}
          <section>
            <h2 className="text-2xl font-semibold mb-4 text-foreground">İletişim</h2>
            <div className="bg-card p-6 rounded-lg shadow-sm border">
              <p className="text-muted-foreground mb-4">
                Sorularınız, önerileriniz veya geri bildirimleriniz için:
              </p>
              <ul className="space-y-2 text-muted-foreground">
                <li>GitHub Issues üzerinden iletişime geçebilirsiniz</li>
                <li>Discord sunucumuza katılabilirsiniz</li>
                <li>E-posta ile bize ulaşabilirsiniz</li>
              </ul>
            </div>
          </section>
        </div>

        {/* Footer */}
        <footer className="mt-16 pt-8 border-t text-center text-sm text-muted-foreground">
          <p>© {new Date().getFullYear()} Turk-AdFilter. Tüm hakları saklıdır.</p>
          <p className="mt-2">
            Made with <Heart className="inline text-red-500" size={14} /> by the Turk-AdFilter Community
          </p>
        </footer>
      </main>
    </div>
  );
} 