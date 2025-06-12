.PHONY: test dedup plot-stats install-deps backup restore backup-info backup-cleanup

# Default target (optional)
all: test

test:
	@echo '--- Filter Lint ---'
	python3 scripts/filter_lint.py
	@echo '\n--- Domain Check ---'
	python3 scripts/domain_check.py
	@echo '\n--- Filter Stats ---'
	python3 scripts/filter_stats.py

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