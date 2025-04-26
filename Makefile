.PHONY: test dedup plot-stats install-deps

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