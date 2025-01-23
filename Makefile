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

env:
	@echo export CGO_CFLAGS=\"$(CGO_CFLAGS)\"
	@echo export CGO_LDFLAGS=\"$(CGO_LDFLAGS)\"
	@echo export CARTESI_LOG_LEVEL="debug"
	@echo export CARTESI_BLOCKCHAIN_HTTP_ENDPOINT=""
	@echo export CARTESI_BLOCKCHAIN_WS_ENDPOINT=""
	@echo export CARTESI_BLOCKCHAIN_ID="11155111"
	@echo export CARTESI_CONTRACTS_INPUT_BOX_ADDRESS="0x593E5BCf894D6829Dd26D0810DA7F064406aebB6"
	@echo export CARTESI_CONTRACTS_INPUT_BOX_DEPLOYMENT_BLOCK_NUMBER="6850934"
	@echo export CARTESI_AUTH_MNEMONIC=\"test test test test test test test test test test test junk\"
	@echo export CARTESI_POSTGRES_ENDPOINT="postgres://postgres:password@localhost:5432/rollupsdb?sslmode=disable"
	@echo export CARTESI_TEST_POSTGRES_ENDPOINT="postgres://test_user:password@localhost:5432/test_rollupsdb?sslmode=disable"
	@echo export CARTESI_TEST_MACHINE_IMAGES_PATH=\"$(CARTESI_TEST_MACHINE_IMAGES_PATH)\"
	@echo export PATH=$(CURDIR):$$PATH
	@echo export ESPRESSO_BASE_URL="https://query.decaf.testnet.espresso.network"
	@echo export ESPRESSO_STARTING_BLOCK="1409980"
	@echo export ESPRESSO_NAMESPACE="55555"

migrate: ## Run migration on development database
	@echo "Running PostgreSQL migration"
	@go run dev/migrate/main.go

generate-db: ## Generate repository/db with Jet
	@echo "Generating internal/repository/db with jet"
	@rm -rf internal/repository/postgres/db
	@go run github.com/go-jet/jet/v2/cmd/jet -dsn=$$CARTESI_POSTGRES_ENDPOINT -schema=public -path=./internal/repository/postgres/db
	@rm -rf internal/repository/postgres/db/rollupsdb/public/model
	@go run github.com/go-jet/jet/v2/cmd/jet -dsn=$$CARTESI_POSTGRES_ENDPOINT -schema=espresso -path=./internal/repository/postgres/db
	@rm -rf internal/repository/postgres/db/rollupsdb/espresso/model
