package config

import (
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "")
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if c.Repo != defaultRepo {
		t.Errorf("Repo=%q", c.Repo)
	}
	if !c.DryRun {
		t.Error("token yokken DryRun true olmalı")
	}
	if c.ConfidenceMin != 70 {
		t.Errorf("ConfidenceMin=%d", c.ConfidenceMin)
	}
	if len(c.DNSServers) < 1 {
		t.Error("varsayılan DNS sunucusu yok")
	}
}

func TestLoadOverrides(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "tok")
	t.Setenv("HUNTER_CONFIDENCE_MIN", "85")
	t.Setenv("HUNTER_INTERVAL", "6h")
	t.Setenv("HUNTER_SOURCES", "mirror, crtsh")
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if c.DryRun {
		t.Error("token varken DryRun false olmalı")
	}
	if c.ConfidenceMin != 85 {
		t.Errorf("ConfidenceMin=%d", c.ConfidenceMin)
	}
	if c.Interval != 6*time.Hour {
		t.Errorf("Interval=%v", c.Interval)
	}
	if len(c.Sources) != 2 || c.Sources[0] != "mirror" || c.Sources[1] != "crtsh" {
		t.Errorf("Sources=%v", c.Sources)
	}
}

func TestValidateBadConfidence(t *testing.T) {
	t.Setenv("HUNTER_CONFIDENCE_MIN", "150")
	if _, err := Load(); err == nil {
		t.Error("geçersiz confidence hata vermeli")
	}
}

func TestRepoSplit(t *testing.T) {
	c := &Config{Repo: "owner/name"}
	if c.RepoOwner() != "owner" || c.RepoName() != "name" {
		t.Errorf("RepoOwner=%q RepoName=%q", c.RepoOwner(), c.RepoName())
	}
}
