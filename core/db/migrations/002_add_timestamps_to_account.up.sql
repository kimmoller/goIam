ALTER TABLE account ADD COLUMN created_at timestamptz NOT NULL DEFAULT now();
ALTER TABLE account ADD COLUMN provisioned_at timestamptz;
ALTER TABLE account ADD COLUMN committed_at timestamptz;
