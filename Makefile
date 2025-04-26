test:
	@echo '--- Filter Lint ---'
	python3 scripts/filter_lint.py
	@echo '\n--- Domain Check ---'
	python3 scripts/domain_check.py
	@echo '\n--- Filter Stats ---'
	python3 scripts/filter_stats.py

dedup:
	python3 scripts/remove_duplicates.py 