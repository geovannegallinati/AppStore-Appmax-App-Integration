APP_PORT   ?= $(shell grep -E '^APP_PORT=' .env 2>/dev/null | cut -d= -f2 || echo 8080)
NGROK_URL  ?= $(shell grep -E '^NGROK_URL=' .env 2>/dev/null | cut -d= -f2 || echo "")
NGROK_AUTHTOKEN ?= $(shell grep -E '^NGROK_AUTHTOKEN=' .env 2>/dev/null | cut -d= -f2 || echo "")
BASE_URL    = http://localhost:$(APP_PORT)
HEALTH_URL  = $(BASE_URL)/health

COMPOSE = docker compose -f docker-compose.yml

.PHONY: install teardown env-init generate-key up down restart logs health validate rename-module migrate

install: env-init generate-key teardown up validate
	@echo "Stack is ready."

env-init:
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Created .env from .env.example"; \
	fi

generate-key:
	@key_val=$$(grep -E '^APP_KEY=' .env 2>/dev/null | cut -d= -f2 | tr -d '[:space:]'); \
	if [ -z "$$key_val" ] || [ $${#key_val} -ne 32 ]; then \
		new_key=$$(LC_ALL=C tr -dc 'A-Za-z0-9' < /dev/urandom | head -c 32); \
		sed -i.bak "s|^APP_KEY=.*|APP_KEY=$$new_key|" .env && rm -f .env.bak; \
		echo "Generated APP_KEY"; \
	else \
		echo "APP_KEY already set — skipping"; \
	fi

teardown:
	@echo "Removing containers and volumes..."
	$(COMPOSE) down -v --remove-orphans

up:
	$(COMPOSE) up -d --build

down:
	$(COMPOSE) down -v --remove-orphans

restart:
	$(COMPOSE) restart

logs:
	$(COMPOSE) logs -f

health:
	@for i in $$(seq 1 30); do \
		if curl -sf $(HEALTH_URL) > /dev/null 2>&1; then \
			exit 0; \
		fi; \
		sleep 3; \
	done; \
	echo "Healthcheck failed after 30 attempts: $(HEALTH_URL)"; exit 1

validate: health
	@echo "Validating endpoints..."
	@if [ -z "$(NGROK_AUTHTOKEN)" ]; then \
		echo "  [FAIL] NGROK_AUTHTOKEN is empty."; \
		echo "  Create a ngrok account at https://dashboard.ngrok.com/signup and set NGROK_AUTHTOKEN in .env"; \
		exit 1; \
	fi
	@status=$$(curl -o /dev/null -sw '%{http_code}' $(BASE_URL)/install/start 2>/dev/null); \
	if [ "$$status" -ge 500 ] || [ "$$status" -eq 000 ]; then \
		echo "  [FAIL] GET /install/start (HTTP $$status)"; exit 1; \
	fi
	@active_url=""; \
	ngrok_url="$(NGROK_URL)"; \
	configured_reachable=0; \
	if [ -n "$$ngrok_url" ]; then \
		case "$$ngrok_url" in http://*|https://*) ;; *) ngrok_url="https://$$ngrok_url" ;; esac; \
		for i in $$(seq 1 20); do \
			if curl -sf $$ngrok_url/health > /dev/null 2>&1; then \
				configured_reachable=1; break; \
			fi; \
			sleep 2; \
		done; \
	fi; \
	for i in $$(seq 1 40); do \
		active_url=$$($(COMPOSE) exec -T ngrok sh -lc "wget -qO- http://127.0.0.1:4040/api/tunnels 2>/dev/null" 2>/dev/null | grep -o '"public_url":"[^"]*"' | head -1 | cut -d'"' -f4); \
		if [ -n "$$active_url" ]; then \
			break; \
		fi; \
		sleep 2; \
	done; \
	if [ -z "$$active_url" ]; then \
		echo "  [FAIL] ngrok active tunnel URL not found in ngrok API"; exit 1; \
	fi; \
	front_status=$$(curl -o /dev/null -sw '%{http_code}' $$active_url/ 2>/dev/null); \
	if [ "$$front_status" -eq 000 ]; then \
		echo "  [FAIL] Frontend URL not reachable: $$active_url/"; exit 1; \
	fi; \
	health_status=$$(curl -o /dev/null -sw '%{http_code}' $$active_url/health 2>/dev/null); \
	if [ "$$health_status" -eq 000 ]; then \
		echo "  [FAIL] Health URL not reachable: $$active_url/health"; exit 1; \
	fi; \
	callback_status=$$(curl -o /dev/null -sw '%{http_code}' $$active_url/integrations/appmax/callback/install 2>/dev/null); \
	if [ "$$callback_status" -eq 000 ]; then \
		echo "  [FAIL] Callback URL not reachable: $$active_url/integrations/appmax/callback/install"; exit 1; \
	fi; \
	if [ "$$configured_reachable" -ne 0 ] && [ -n "$$ngrok_url" ] && [ "$$active_url" = "$$ngrok_url" ]; then :; fi; \
	echo "  Frontend URL: $$active_url/"; \
	echo "  Health URL: $$active_url/health"; \
	echo "  Callback URL: $$active_url/integrations/appmax/callback/install"
	@echo "All validations passed."

rename-module:
	@bash scripts/rename-module.sh $(NEW)

migrate:
	$(COMPOSE) exec app ./tmp/server artisan migrate
