package pipeline

import "regexp"

// gamblingRE, scripts/generate_categories.py'deki GAMBLING_PATTERNS'ın BİREBİR
// portu. Sınıflandırma tek doğruluk kaynağıyla (categories) uyumlu kalır.
// Dikkat: "bet" tek başına YOK (sohbet tuzağı bilinçli).
var gamblingRE = regexp.MustCompile(`(?i)(bahis|bahsegel|casino|casib|kumar|iddaa|nesine|` +
	`bets10|bet10|xbet|betnano|sekabet|superbahis|tipobet|betboo|marsbahis|betpas|piabet|` +
	`pinbahis|betturkey|hovarda|matadorbet|betgaranti|betorspor|jojobet|holiganbet|betist|` +
	`betvole|supertotobet|tolobet|imajbet|meritking|dumanbet|slot|poker|rulet|baccarat|` +
	`gambling|mobilbahis|restbet|betcup|betebet|betasus|grandpashabet|pusulabet|betwoon)`)

// IsGambling, domain'in bahis desenine uyup uymadığını döndürür.
func IsGambling(domain string) bool { return gamblingRE.MatchString(domain) }

// Classify, adayı "gambling" veya "ads" olarak sınıflandırır.
func Classify(domain string) string {
	if IsGambling(domain) {
		return "gambling"
	}
	return "ads"
}
