CREATE TABLE vpn_users (
  id         BIGSERIAL PRIMARY KEY,
  full_name  VARCHAR(255) NOT NULL,
  email      VARCHAR(255) NOT NULL,
  password   VARCHAR(255) NOT NULL,
  details    TEXT NOT NULL,
  online     BOOLEAN                     NOT NULL DEFAULT FALSE,
  enable     BOOLEAN                     NOT NULL DEFAULT FALSE,
  start_up   DATE                        NOT NULL DEFAULT '2016-12-13',
  shut_down  DATE                        NOT NULL DEFAULT current_date,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE UNIQUE INDEX idx_vpn_users_email ON vpn_users (email);
CREATE INDEX idx_vpn_users_full_name ON vpn_users (full_name);

CREATE TABLE vpn_logs (
  id           BIGSERIAL PRIMARY KEY,
  user_id      BIGINT NOT NULL REFERENCES vpn_users,
  trusted_ip   INET,
  trusted_port SMALLINT,
  remote_ip    INET,
  remote_port  SMALLINT,
  start_up     TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  shut_down    TIMESTAMP WITHOUT TIME ZONE,
  received     FLOAT                       NOT NULL DEFAULT '0.0',
  send         FLOAT                       NOT NULL DEFAULT '0.0'
);
