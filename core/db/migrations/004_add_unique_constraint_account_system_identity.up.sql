ALTER TABLE account ADD CONSTRAINT unique_account_system_id_identity_id UNIQUE (identity_id, system_id);
