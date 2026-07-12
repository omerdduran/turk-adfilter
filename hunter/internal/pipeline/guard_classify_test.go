package pipeline

import "testing"

func TestGuardLabelBoundary(t *testing.T) {
	g := NewGuard(nil, []string{"sozcu.com.tr"})
	cases := map[string]Verdict{
		"site.gov.tr":         Reject, // resmi suffix
		"ads.google.com":      Reject, // korunan eTLD+1 alt alanı
		"google.com":          Reject, // tam eşleşme
		"evilgoogle.com":      Allow,  // LABEL SINIRI: alt alan değil, önerilebilir
		"cdn.sozcu.com.tr":    Reject, // auto-seed (crawl sitesinin kendi domaini)
		"x.cloudfront.net":    Hold,   // paylaşımlı altyapı → beklet
		"amazonaws.com":       Hold,
		"membrana.media":      Allow, // gerçek aday
		"static.qovani.com":   Allow,
		// Phishing güvenliği: GERÇEK banka allowlist'te sert-ret, SAHTE varyant geçer.
		"garantibbva.com.tr":      Reject,
		"ziraatbank.com.tr":       Reject,
		"garantibbva-tr-giris.xyz": Allow,
	}
	for d, want := range cases {
		if got := g.Check(d); got != want {
			t.Errorf("Check(%q)=%v, beklenen %v", d, got, want)
		}
	}
}

func TestGuardExtra(t *testing.T) {
	g := NewGuard([]string{"benim-cdn.net"}, nil)
	if g.Check("a.benim-cdn.net") != Reject {
		t.Error("HUNTER_ALLOWLIST_EXTRA girişi reddedilmeli")
	}
}

func TestClassify(t *testing.T) {
	gambling := []string{
		"11587bets10.com", "casibom.com", "sekabet99.net", "mobilbahis1.com", "grandpashabet.co",
		// yeni eklenen markalar
		"betwinner12.com", "pashagaming.net", "betpark55.com", "perabet7.com", "sahabet.vip", "kralbet.info",
	}
	for _, d := range gambling {
		if Classify(d) != "gambling" {
			t.Errorf("Classify(%q)=%s, gambling olmalı", d, Classify(d))
		}
	}
	// "sohbet" hâlâ ads olmalı — "bet" tek başına yok, perabet/sahabet gibi tam markalar sohbet'te geçmez.
	ads := []string{"sohbet.com", "sohbethane.net", "membrana.media", "static.qovani.com", "adcluster.com.tr"}
	for _, d := range ads {
		if Classify(d) != "ads" {
			t.Errorf("Classify(%q)=%s, ads olmalı ('sohbet'/'bet' tuzağı)", d, Classify(d))
		}
	}
}

func TestClassifyPhishing(t *testing.T) {
	phish := []string{"garantibbva-tr.xyz", "ziraatbank-giris.com", "isbank-mobil.net", "papara-tr.info", "edevlet-basvuru.xyz"}
	for _, d := range phish {
		if Classify(d) != "phishing" {
			t.Errorf("Classify(%q)=%s, phishing olmalı", d, Classify(d))
		}
	}
	// Bahis önceliği: hem bahis hem banka deseni varsa gambling kazanır.
	if Classify("sekabet-ziraatbank.com") != "gambling" {
		t.Error("bahis deseni phishing'ten önce gelmeli")
	}
}

func TestHasAdToken(t *testing.T) {
	yes := []string{"adcluster.com.tr", "analytics.x.com", "reklam.net", "tracker.io"}
	for _, d := range yes {
		if !HasAdToken(d) {
			t.Errorf("HasAdToken(%q)=false, true olmalı", d)
		}
	}
	if HasAdToken("sohbet.com") {
		t.Error("sohbet.com ad token içermemeli")
	}
}
