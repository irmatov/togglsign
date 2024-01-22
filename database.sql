CREATE SEQUENCE signatures_id_seq;

CREATE TABLE signatures (
    id INTEGER NOT NULL PRIMARY KEY DEFAULT nextval('signatures_id_seq'),
    email VARCHAR NOT NULL,
    signed_at TIMESTAMP WITH TIME ZONE NOT NULL,
    sig VARCHAR NOT NULL,
    UNIQUE (email, sig)
);

CREATE SEQUENCE responses_id_seq;

CREATE TABLE responses (
    id INTEGER NOT NULL PRIMARY KEY DEFAULT nextval('responses_id_seq'),
    sort_key INTEGER NOT NULL,
    question VARCHAR NOT NULL,
    answer VARCHAR NOT NULL,
    sig_id INTEGER NOT NULL REFERENCES signatures(id)
);
