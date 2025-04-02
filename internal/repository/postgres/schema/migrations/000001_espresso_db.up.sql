-- (c) Cartesi and individual authors (see AUTHORS)
-- SPDX-License-Identifier: Apache-2.0 (see LICENSE)

CREATE SCHEMA espresso;

CREATE TABLE IF NOT EXISTS espresso.espresso_config
(
    "application_address" ethereum_address PRIMARY KEY,
	"starting_block" uint64 NOT NULL,
    "namespace" uint64 NOT NULL
);

CREATE TABLE IF NOT EXISTS espresso.espresso_nonce
(
    "sender_address" ethereum_address NOT NULL,
    "application_address" ethereum_address NOT NULL,
    "nonce" BIGINT NOT NULL,
    UNIQUE("sender_address", "application_address")
);

CREATE TABLE IF NOT EXISTS espresso.espresso_block
(
    "application_address" ethereum_address PRIMARY KEY,
	"last_processed_espresso_block" uint64 NOT NULL
);

CREATE TABLE IF NOT EXISTS espresso.input_index
(
    "application_address" ethereum_address PRIMARY KEY,
	"index" BIGINT NOT NULL
);
