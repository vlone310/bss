CREATE TABLE accounts (
  id bigserial PRIMARY KEY,
  owner varchar NOT NULL,
  balance bigint NOT NULL,
  currency varchar NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE entries (
  id bigserial PRIMARY KEY,
  account_id bigint,
  amount_cents bigint NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE transfers (
  id bigserial PRIMARY KEY,
  from_account_id bigint NOT NULL,
  to_account_id bigint NOT NULL,
  amount_cents bigint NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE countries (
  code int PRIMARY KEY,
  name varchar NOT NULL,
  continent_name varchar NOT NULL
);

CREATE INDEX ON accounts (owner);

CREATE INDEX ON entries (account_id);

CREATE INDEX ON transfers (from_account_id);

CREATE INDEX ON transfers (to_account_id);

CREATE INDEX ON transfers (from_account_id, to_account_id);

COMMENT ON COLUMN entries.amount_cents IS 'can be negative and positive';

COMMENT ON COLUMN transfers.amount_cents IS 'must be positive';

ALTER TABLE entries ADD FOREIGN KEY (account_id) REFERENCES accounts (id);

ALTER TABLE transfers ADD FOREIGN KEY (from_account_id) REFERENCES accounts (id);

ALTER TABLE transfers ADD FOREIGN KEY (to_account_id) REFERENCES accounts (id);
