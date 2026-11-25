CREATE TABLE IF NOT EXISTS receitas (
    id                SERIAL PRIMARY KEY,
    titulo            VARCHAR(120) NOT NULL,
    categoria         VARCHAR(30)  NOT NULL,
    ingredientes      TEXT         NOT NULL,
    modo_preparo      TEXT         NOT NULL,
    tempo_preparo_min INTEGER      NOT NULL CHECK (tempo_preparo_min > 0),
    porcoes           INTEGER      NOT NULL CHECK (porcoes > 0),
    criada_em         TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS api_keys (
    id        SERIAL PRIMARY KEY,
    nome      VARCHAR(60) NOT NULL,
    key_hash  CHAR(64)    NOT NULL UNIQUE, -- SHA-256 em hexadecimal
    criada_em TIMESTAMPTZ NOT NULL DEFAULT now()
);
