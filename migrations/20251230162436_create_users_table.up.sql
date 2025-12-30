CREATE TABLE users (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	username varchar(12) NOT NULL,
	password varchar NOT NULL,
	created_at timestamptz DEFAULT now()
)
