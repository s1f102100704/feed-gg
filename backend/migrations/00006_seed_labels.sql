-- +goose Up
INSERT INTO label (name) VALUES
  ('寄り遅め'),
  ('単独行動'),
  ('無理な仕掛け'),
  ('ファーム優先'),
  ('オブジェクト後回し'),
  ('アグレッシブ'),
  ('オブジェクト優先'),
  ('連携重視')
ON CONFLICT (name) DO NOTHING;

-- +goose Down
DELETE FROM label
WHERE name IN (
  '寄り遅め',
  '単独行動',
  '無理な仕掛け',
  'ファーム優先',
  'オブジェクト後回し',
  'アグレッシブ',
  'オブジェクト優先',
  '連携重視'
);
