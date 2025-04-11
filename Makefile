# (c) Cartesi and individual authors (see AUTHORS)
# SPDX-License-Identifier: Apache-2.0 (see LICENSE)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

ROLLUPS_CONTRACTS_ABI_BASEDIR := rollups-contracts/
CONTRACTS_VERSION := 2.0.0-rc.17
CONTRACTS_URL := https://github.com/cartesi/rollups-contracts/releases/download/
CONTRACTS_ARTIFACT := rollups-contracts-$(CONTRACTS_VERSION)-artifacts.tar.gz
CONTRACTS_SHA256 := f87096ea7e7a3ff38c2c3807ae3d1711f23b8ce0ff843423c5b73356ef8b8bc7

env:
	@echo export CARTESI_LOG_LEVEL="debug"
	@echo export CARTESI_BLOCKCHAIN_HTTP_ENDPOINT="http://localhost:8545"
	@echo export CARTESI_BLOCKCHAIN_WS_ENDPOINT="ws://localhost:8545"
	@echo export CARTESI_DATABASE_CONNECTION="postgres://postgres:password@localhost:5432/rollupsdb?sslmode=disable"
	@echo export PATH=\"$(CURDIR):$$PATH\"
	@echo export ESPRESSO_BASE_URL="http://localhost:24000"

migrate: ## Run migration on development database
	@echo "Running PostgreSQL migration"
	@go run dev/migrate/main.go

generate-db: ## Generate repository/db with Jet
	@echo "Generating internal/repository/db with jet"
	@rm -rf internal/repository/postgres/db
	@go run github.com/go-jet/jet/v2/cmd/jet -dsn=$$CARTESI_DATABASE_CONNECTION -schema=public -path=./internal/repository/postgres/db
	@rm -rf internal/repository/postgres/db/rollupsdb/public/model
	@go run github.com/go-jet/jet/v2/cmd/jet -dsn=$$CARTESI_DATABASE_CONNECTION -schema=espresso -path=./internal/repository/postgres/db
	@rm -rf internal/repository/postgres/db/rollupsdb/espresso/model

generate: $(ROLLUPS_CONTRACTS_ABI_BASEDIR)/.stamp ## Generate the file that are committed to the repo
	@echo "Generating Go files"
	@go generate ./internal/... ./pkg/...

$(ROLLUPS_CONTRACTS_ABI_BASEDIR)/.stamp:
	@echo "Downloading rollups-contracts artifacts"
	@mkdir -p $(ROLLUPS_CONTRACTS_ABI_BASEDIR)
	@curl -sSL $(CONTRACTS_URL)/v$(CONTRACTS_VERSION)/$(CONTRACTS_ARTIFACT) -o $(CONTRACTS_ARTIFACT)
	@echo "$(CONTRACTS_SHA256)  $(CONTRACTS_ARTIFACT)" | shasum -a 256 --check > /dev/null
	@tar -zxf $(CONTRACTS_ARTIFACT) -C $(ROLLUPS_CONTRACTS_ABI_BASEDIR)
	@touch $@
	@rm $(CONTRACTS_ARTIFACT)
