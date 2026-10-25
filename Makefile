.PHONY: all test test-unit lint dedup install-deps stats backup backup-force backup-info restore restore-list backup-cleanup

# Default target (optional)
all: test test-unit

test:
	@echo '--- Filter Lint ---'
	python3 scripts/filter_lint.py
	@echo '\n--- Domain Check (DNS, süre sınırlı) ---'
	python3 scripts/domain_check.py
	@echo '\n--- Filter Stats (çevrimdışı) ---'
	python3 scripts/filter_stats.py --no-resolve

test-unit:
	@echo '--- Unit Tests (generate_formats + invariants + freshness) ---'
	python3 -m unittest scripts.test_generate_formats
	python3 -m unittest scripts.test_filter_invariants
	python3 -m unittest scripts.test_refresh_candidates

lint:
	@echo '--- AGLint (AdGuard filter linter) ---'
	npx aglint

dedup:
	python3 scripts/remove_duplicates.py

# Target to install dependencies
install-deps:
	python3 -m pip install -r requirements.txt

# Target to generate the plot, depends on dependencies being installed
stats: install-deps
	python3 scripts/plot_filter_stats.py

# Backup related targets
backup:
	@echo '--- Filtre Yedekleme ---'
	python3 scripts/backup_filters.py

backup-force:
	@echo '--- Zorla Filtre Yedekleme ---'
	python3 scripts/backup_filters.py --force

backup-info:
	@echo '--- Yedek Bilgileri ---'
	python3 scripts/backup_filters.py --info

restore:
	@echo '--- İnteraktif Geri Yükleme ---'
	python3 scripts/restore_backup.py --interactive

restore-list:
	@echo '--- Mevcut Yedekler ---'
	python3 scripts/restore_backup.py --list

backup-cleanup:
	@echo '--- Eski Yedekleri Temizle (30+ gün) ---'
	python3 scripts/restore_backup.py --cleanup 30 