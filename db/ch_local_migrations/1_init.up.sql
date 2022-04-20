CREATE TABLE IF NOT EXISTS blocks ON CLUSTER cluster_1 (
  block_number UInt64,
  hash String,
  proposer_address String,
  created_date DateTime('Europe/London')
) Engine = ReplicatedMergeTree('/clickhouse/tables/{layer}-{shard}/blocks', '{replica}')
ORDER BY block_number;

CREATE TABLE IF NOT EXISTS txs ON CLUSTER cluster_1 (
  hash String,
  block_number UInt64,
  status String,
  fee UInt64,
  fee_currency String,
  fee_payer_address String,
  created_date DateTime('Europe/London')
) Engine = ReplicatedMergeTree('/clickhouse/tables/{layer}-{shard}/txs', '{replica}')
ORDER BY hash;

CREATE TABLE IF NOT EXISTS transfers ON CLUSTER cluster_1 (
  id String,
  tx_hash String,
  from_address String,
  to_address String,
  amount UInt64,
  currency String
) Engine = ReplicatedMergeTree('/clickhouse/tables/{layer}-{shard}/transfers', '{replica}')
ORDER BY id;
