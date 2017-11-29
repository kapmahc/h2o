CREATE TABLE users (
  id                 BIGSERIAL PRIMARY KEY,
  name               VARCHAR(32)                 NOT NULL,
  email              VARCHAR(255)                NOT NULL,
  uid                VARCHAR(36)                 NOT NULL,
  password           bytea,
  provider_id        VARCHAR(255)                NOT NULL,
  provider_type      VARCHAR(32)                 NOT NULL,
  logo               VARCHAR(255),
  sign_in_count      INT                         NOT NULL DEFAULT 0,
  current_sign_in_at TIMESTAMP WITHOUT TIME ZONE,
  current_sign_in_ip INET,
  last_sign_in_at    TIMESTAMP WITHOUT TIME ZONE,
  last_sign_in_ip    INET,
  confirmed_at       TIMESTAMP WITHOUT TIME ZONE,
  locked_at          TIMESTAMP WITHOUT TIME ZONE,
  created_at         TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_at         TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE UNIQUE INDEX idx_users_uid
  ON users (uid);
CREATE UNIQUE INDEX idx_users_email
  ON users (email);
CREATE UNIQUE INDEX idx_users_provider_id_type
  ON users (provider_id, provider_type);
CREATE INDEX idx_users_name
  ON users (name);
CREATE INDEX idx_users_provider_id
  ON users (provider_id);
CREATE INDEX idx_users_provider_type
  ON users (provider_type);


CREATE TABLE logs (
  id         BIGSERIAL PRIMARY KEY,
  user_id    BIGINT                      NOT NULL REFERENCES users,
  ip         INET                        NOT NULL,
  message    VARCHAR(255)                NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE roles (
  id            BIGSERIAL PRIMARY KEY,
  name          VARCHAR(32)                 NOT NULL,
  resource_id   BIGINT,
  resource_type VARCHAR(255),
  created_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE UNIQUE INDEX idx_roles_name_resource_type_id
  ON roles (name, resource_type, resource_id);
CREATE INDEX idx_roles_name
  ON roles (name);
CREATE INDEX idx_roles_resource_type
  ON roles (resource_type);

CREATE TABLE policies (
  id         BIGSERIAL PRIMARY KEY,
  user_id    BIGINT                      NOT NULL REFERENCES users,
  role_id    BIGINT                      NOT NULL REFERENCES roles,
  start_up   DATE                        NOT NULL DEFAULT current_date,
  shut_down  DATE                        NOT NULL DEFAULT '2017-09-18',
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE UNIQUE INDEX idx_policies
  ON policies (user_id, role_id);


CREATE TABLE votes (
  id            BIGSERIAL PRIMARY KEY,
  resource_type VARCHAR(255)                NOT NULL,
  resource_id   BIGINT                      NOT NULL,
  point         INT                         NOT NULL DEFAULT 0,
  created_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE UNIQUE INDEX idx_votes_resources
  ON votes (resource_type, resource_id);
CREATE INDEX idx_votes_resource_type
  ON votes (resource_type);


CREATE TABLE attachments (
  id            BIGSERIAL PRIMARY KEY,
  title         VARCHAR(255)                NOT NULL,
  url           VARCHAR(255)                NOT NULL,
  length        INT                         NOT NULL,
  media_type    VARCHAR(32)                 NOT NULL,
  resource_type VARCHAR(255)                NOT NULL,
  resource_id   BIGINT                      NOT NULL,
  user_id       BIGINT                      NOT NULL REFERENCES users,
  created_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE UNIQUE INDEX idx_attachments_url
  ON attachments (url);
CREATE INDEX idx_attachments_title
  ON attachments (title);
CREATE INDEX idx_attachments_resource_type
  ON attachments (resource_type);
CREATE INDEX idx_attachments_media_type
  ON attachments (media_type);

CREATE TABLE leave_words (
  id         BIGSERIAL PRIMARY KEY,
  body       TEXT                        NOT NULL,
  type       VARCHAR(8)                  NOT NULL DEFAULT 'markdown',
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE links (
  id BIGSERIAL PRIMARY KEY,
  href VARCHAR(255) NOT NULL,
  label VARCHAR(255) NOT NULL,
  loc VARCHAR(16) NOT NULL,
  lang VARCHAR(8) NOT NULL,
  sort_order INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE INDEX idx_links_loc_lang ON links (loc, lang);
CREATE INDEX idx_links_lang ON links (lang);

CREATE TABLE cards (
  id BIGSERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  summary VARCHAR(2048) NOT NULL,
  type VARCHAR(8) NOT NULL DEFAULT 'markdown',
  action VARCHAR(32) NOT NULL,
  href VARCHAR(255) NOT NULL,
  logo VARCHAR(255) NOT NULL,
  loc VARCHAR(16) NOT NULL,
  lang VARCHAR(8) NOT NULL,
  sort_order INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE INDEX idx_cards_loc_lang ON cards (loc, lang);
CREATE INDEX idx_cards_lang ON cards (lang);

CREATE TABLE friend_links (
  id BIGSERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  home VARCHAR(255) NOT NULL,
  logo VARCHAR(255) NOT NULL,
  sort_order INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
