package pipeline

import "regexp"

// gamblingRE, scripts/generate_categories.py'deki GAMBLING_PATTERNS'ın BİREBİR
// portu. Sınıflandırma tek doğruluk kaynağıyla (categories) uyumlu kalır.
// Dikkat: "bet" tek başına YOK (sohbet tuzağı bilinçli); her desen ayırt edici olmalı.
var gamblingRE = regexp.MustCompile(`(?i)(bahis|bahsegel|casino|casib|kumar|iddaa|nesine|` +
	`bets10|bet10|xbet|betnano|sekabet|superbahis|tipobet|betboo|marsbahis|betpas|piabet|` +
	`pinbahis|betturkey|hovarda|matadorbet|betgaranti|betorspor|jojobet|holiganbet|betist|` +
	`betvole|supertotobet|tolobet|imajbet|meritking|dumanbet|slot|poker|rulet|baccarat|` +
	`gambling|mobilbahis|restbet|betcup|betebet|betasus|grandpashabet|pusulabet|betwoon|` +
	`betwinner|pashagaming|winxbet|betgit|elexbet|betpark|cratosslot|discountcasino|betlike|` +
	`maltcasino|ngsbahis|favorisen|tempobet|youwin|betsat|artemisbet|kralbet|bahigo|redwin|` +
	`sahabet|betpuan|jetbahis|asyabahis|perabet|dinamobet|betmarino|betorder)`)

// IsGambling, domain'in bahis desenine uyup uymadığını döndürür.
func IsGambling(domain string) bool { return gamblingRE.MatchString(domain) }

// phishRE, banka/kurum adı taşıyan (muhtemel phishing) domainleri yakalar. AYIRT
// EDİCİ markalar — genel kelime yok (false-positive'i azaltmak için). Gerçek kurumlar
// allowlist.txt'te sert-ret edilir; burada yalnız sahte varyantlar phishing sınıfına girer.
var phishRE = regexp.MustCompile(`(?i)(garantibbva|ziraatbank|isbank|akbank|yapikredi|halkbank|` +
	`vakifbank|denizbank|qnbfinans|enpara|papara|ininal|turkiyefinans|kuveytturk|edevlet)`)

// IsPhishing, domain'in banka/kurum desenine uyup uymadığını döndürür.
func IsPhishing(domain string) bool { return phishRE.MatchString(domain) }

// Classify, adayı "gambling", "phishing" veya "ads" olarak sınıflandırır.
func Classify(domain string) string {
	switch {
	case IsGambling(domain):
		return "gambling"
	case IsPhishing(domain):
		return "phishing"
	default:
		return "ads"
	}
}
