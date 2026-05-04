.PHONY: setup install test test-verbose lint format clean

PYTHON = python3
VENV = venv
ACTIVATE = . $(VENV)/bin/activate

setup:
	$(PYTHON) -m venv $(VENV)
	$(ACTIVATE) && pip install -r requirements.txt

install: setup
	$(ACTIVATE) && pip install -r requirements.txt

test: install
	$(ACTIVATE) && pytest tests/

test-verbose: install
	$(ACTIVATE) && pytest -v tests/

lint: install
	$(ACTIVATE) && flake8 src/ tests/

format: install
	$(ACTIVATE) && black src/ tests/

clean:
	rm -rf $(VENV) __pycache__ .pytest_cache src/__pycache__ tests/__pycache__ src/*.pyc tests/*.pyc
