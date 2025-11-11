.PHONY: help install install-dev test lint format clean run

help:
	@echo "OSE CLI Development Commands"
	@echo ""
	@echo "  make install       Install production dependencies"
	@echo "  make install-dev   Install development dependencies"
	@echo "  make test          Run tests with coverage"
	@echo "  make lint          Run linters (flake8, mypy)"
	@echo "  make format        Format code with black"
	@echo "  make clean         Remove build artifacts"
	@echo "  make run           Run CLI in development mode"
	@echo ""

install:
	pip install -e .

install-dev:
	pip install -e ".[dev]"
	pip install -r requirements-dev.txt

test:
	pytest tests/ -v --cov=cli --cov-report=html --cov-report=term

lint:
	flake8 cli/ --max-line-length=100
	mypy cli/

format:
	black cli/ tests/

clean:
	rm -rf build/
	rm -rf dist/
	rm -rf *.egg-info
	rm -rf .pytest_cache/
	rm -rf htmlcov/
	rm -rf .coverage
	find . -type d -name __pycache__ -exec rm -rf {} +
	find . -type f -name "*.pyc" -delete

run:
	python -m cli.main
