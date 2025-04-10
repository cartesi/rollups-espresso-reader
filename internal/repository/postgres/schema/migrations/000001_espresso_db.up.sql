-- (c) Cartesi and individual authors (see AUTHORS)
-- SPDX-License-Identifier: Apache-2.0 (see LICENSE)

CREATE SCHEMA espresso;

CREATE TABLE IF NOT EXISTS espresso.app_info
(
    "application_address" ethereum_address PRIMARY KEY,
	"starting_block" uint64 NOT NULL,
    "namespace" uint64 NOT NULL,
    "last_processed_espresso_block" uint64,
    "index" BIGINT
);

CREATE TABLE IF NOT EXISTS espresso.espresso_nonce
(
    "sender_address" ethereum_address NOT NULL,
    "application_address" ethereum_address NOT NULL,
    "nonce" BIGINT NOT NULL,
    UNIQUE("sender_address", "application_address")
);
