# Katkıda Bulunma Rehberi

Turk-AdFilter; %100 açık kaynaklı, topluluk tabanlı bir projedir ve katkılarınıza açıktır. Bu rehber, projeye nasıl katkıda bulunabileceğinizi açıklar.

## Katkı Türleri

Projeye aşağıdaki şekillerde katkıda bulunabilirsiniz:

1. **Yeni domain ekleme**: Türkiye'de karşılaştığınız reklam ve izleme domainlerini ekleyebilirsiniz
2. **Hata düzeltme**: Mevcut listede gördüğünüz hataları düzeltebilirsiniz
3. **Yanlış pozitifleri bildirme**: Listedeki bir kuralın normal web sitelerini engellediğini fark ederseniz bildirebilirsiniz
4. **Belgelendirme**: Bu dokümantasyonu geliştirebilir veya çeviriler ekleyebilirsiniz

## GitHub Üzerinden Katkı

### 1. Issue (Sorun) Açma

  #### 🐞 Genel Hata Bildirimi
Diğer hata türleri için genel bir şablon. Teknik sorunları ve beklenmeyen hataları buradan bildirin.

  #### 🚫 False Positive
Yanlışlıkla engellenen siteleri buradan bildirin. Hangi kuralın ve sitenin etkilendiğini belirtin.

  #### ✨ Feature Request
Yeni filtre, kural veya geliştirme önerilerinizi paylaşın.

  #### 📢 Missed Ad
Engellenmeyen reklamları veya izleyicileri buradan bildirin. Site ve reklam detayını ekleyin.

  #### 📝 Blank Issue
Belirli bir şablona uymayan konular için boş şablon.

### 2. Pull Request (PR) Gönderme

Listeye doğrudan katkıda bulunmak için:

1. Projeyi fork edin (GitHub'da "Fork" butonuna tıklayarak)
2. Fork ettiğiniz repoyu yerel bilgisayarınıza klonlayın:
   ```
   git clone https://github.com/KULLANICI_ADINIZ/turk-adfilter.git
   ```
3. Değişikliklerinizi yapın (turk-adfilter.txt dosyasını düzenleyin)
4. Değişikliklerinizi commit edin:
   ```
   git commit -m "Yeni domainler eklendi: example.com, ads.example.com"
   ```
5. Değişiklikleri GitHub'a push edin:
   ```
   git push origin main
   ```
6. GitHub'da ana repoya PR gönderin

## Domain Ekleme Kuralları

Yeni domainler eklerken aşağıdaki kurallara uyun:

1. **Sadece Türkiye merkezli** reklam, izleyici ve zararlı içerik sağlayıcıları ekleyin
2. Domainleri alfabetik sıraya göre ekleyin
3. Her domain için kanıt sunun (ekran görüntüsü, açıklama vb.)
4. Domain eklerken tam sözdizimini kullanın:
   ```
   ||example.com^
   ```

## Yanlış Pozitifleri Bildirme

Eğer liste normal web sitelerini engelliyor ve sorunlara neden oluyorsa:

1. [GitHub Issues](https://github.com/omerdduran/turk-adfilter/issues) sayfasında yeni bir issue açın
2. "Yanlış Pozitif Bildirimi" şablonunu seçin
3. Hangi kuralın soruna neden olduğunu, hangi sitede sorun yaşadığınızı ve nasıl tekrarlanabileceğini detaylı olarak belirtin

## Teşekkürler

Katkıda bulunan tüm katılımcılara teşekkür ederiz! Sizin sayenizde daha iyi bir Türkçe web deneyimi sunabiliyoruz.

## İletişim

- GitHub: [@omerdduran/turk-adfilter](https://github.com/omerdduran/turk-adfilter)
- Matrix: [#reklamsiz-turkiye:matrix.org](https://matrix.to/#/#reklamsiz-turkiye:matrix.org)
