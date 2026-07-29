help:
	@echo "Makefile targets:"

setup:
	cd emulator-agent && pip install -r requirements.txt
	cd llm-orchestrator && pip install -r requirements.txt

test:
	cd emulator-agent && pytest tests/
	cd llm-orchestrator && pytest tests/

lint:
	cd emulator-agent && flake8 .
	cd llm-orchestrator && flake8 .

# Defaults match local insecure_dev_mode + deploy/devrealm (see docs/TESTING.md)
ASTRATE_URL ?= http://localhost:8080
REALM       ?= test
REALM_KEY   ?= ../../deploy/devrealm/realm_private.pem

interfaces-install:
	@for f in astrate-interfaces/*.json; do \
	  astartectl realm-management interfaces install "$$f" \
	    -u $(ASTRATE_URL) -k $(REALM_KEY) -r $(REALM); \
	done

run-emulator-stub:
	cd emulator-agent && python -m emulator_agent.main --stub --insecure \
	  --astrate-url $(ASTRATE_URL) --realm $(REALM) --device-id $(DEVICE_ID)

# ROM is never committed. Example:
#   make run-emulator-rom DEVICE_ID=… \
#     ROM="/Users/…/Pokemon - Red Version (UE)[!]/Pokemon Red.gb"
ROM ?=
run-emulator-rom:
	@test -n "$(DEVICE_ID)" || (echo "DEVICE_ID is required" && exit 1)
	@test -n "$(ROM)" || (echo "ROM path is required (do not commit the ROM)" && exit 1)
	cd emulator-agent && python -m emulator_agent.main \
	  --rom "$(ROM)" --insecure \
	  --astrate-url $(ASTRATE_URL) --realm $(REALM) --device-id $(DEVICE_ID)

run-emulator-mtls:
	cd emulator-agent && python -m emulator_agent.main --stub \
	  --astrate-url $(ASTRATE_URL) --realm $(REALM) --device-id $(DEVICE_ID) \
	  --cert $(CERT) --key $(KEY) --ca $(CA)

run-orchestrator:
	cd llm-orchestrator && python -m llm_orchestrator.main

clean:
	find . -type d -name "__pycache__" -exec rm -rf {} +
	find . -type d -name ".pytest_cache" -exec rm -rf {} +
