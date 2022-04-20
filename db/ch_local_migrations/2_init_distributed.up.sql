CREATE TABLE IF NOT EXISTS blocks_all ON CLUSTER cluster_1 as blocks ENGINE = Distributed(cluster_1, default, blocks, rand());
CREATE TABLE IF NOT EXISTS txs_all ON CLUSTER cluster_1 as txs ENGINE = Distributed(cluster_1, default, txs, rand());
CREATE TABLE IF NOT EXISTS transfers_all ON CLUSTER cluster_1 as transfers ENGINE = Distributed(cluster_1, default, transfers, rand());

CREATE OR REPLACE VIEW transfer_view ON CLUSTER cluster_1 AS SELECT 
tr.id as id,
b.block_number as block_number, tr.tx_hash as tx_hash, t.status as status, t.fee as fee,
tr.from_address as from_address, tr.to_address as to_address, tr.amount as amount, 
tr.currency as currency, t.created_date as created_date
from transfers_all as tr 
left join txs_all as t on t.hash = tr.tx_hash
left join blocks_all as b on b.block_number = t.block_number 
SETTINGS distributed_product_mode = 'allow';