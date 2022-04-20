CREATE TABLE IF NOT EXISTS blocks (
  block_number bigint,
  hash varchar(255),
  proposer_address varchar(255),
  created_date timestamp without time zone,
  PRIMARY KEY (block_number)
);

CREATE TABLE IF NOT EXISTS txs (
  hash varchar(255),
  block_number bigint,
  status varchar(50),
  fee bigint,
  fee_currency varchar(50),
  fee_payer_address varchar(50),
  created_date timestamp without time zone,
  PRIMARY KEY (hash),
  CONSTRAINT fk_block_number
    FOREIGN KEY(block_number) 
	  REFERENCES blocks(block_number)
	  ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS transfers (
  id varchar(255),
  tx_hash varchar(255),
  from_address varchar(255),
  to_address varchar(255),
  amount bigint,
  currency varchar(50),
  PRIMARY KEY (id),
  CONSTRAINT fk_tx_hash
    FOREIGN KEY(tx_hash) 
	  REFERENCES txs(hash)
	  ON DELETE CASCADE
);

CREATE OR REPLACE VIEW transfer_view AS SELECT
tr.id as id,
b.block_number as block_number, tr.tx_hash as tx_hash, t.status as status, t.fee as fee,
tr.from_address as from_address, tr.to_address as to_address, tr.amount as amount, 
tr.currency as currency, t.created_date as created_date
 from transfers as tr 
left join txs as t on t.hash = tr.tx_hash
left join blocks as b on b.block_number = t.block_number;