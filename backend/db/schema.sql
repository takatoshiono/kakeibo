--  保有金融機関
CREATE TABLE sources (
    id varchar(255) not null,
    name varchar(255) not null,
    display_order int not null,
    created_at datetime not null,
    updated_at datetime not null,
    primary key (id)
);
CREATE UNIQUE INDEX idx_name_on_sources on sources(name);

-- 項目
CREATE TABLE categories (
    id varchar(255) not null,
    name varchar(255) not null,
    level int not null,
    display_order int not null,
    parent_id varchar(255),
    created_at datetime not null,
    updated_at datetime not null,
    primary key (id)
);
CREATE UNIQUE INDEX idx_name_level_on_categories on categories(name, level);
CREATE INDEX idx_parent_id on categories(parent_id);

/*
    マネーフォワードからダウンロードした入出金履歴データ
    計算対象かつ振替ではないデータだけ保存する
    IDはそのまま主キーに使う
*/
CREATE TABLE money_forward_records (
    id varchar(255) not null, -- ID
    recorded_on date not null, -- 日付. 計上日
    title varchar(255) not null, -- 内容
    amount int not null, -- 金額
    source_id varchar(255) not null, -- 保有金融機関id
    category_id varchar(255) not null, -- 中項目id
    memo varchar(255) not null, -- メモ
    created_at datetime not null,
    updated_at datetime not null,
    primary key (id)
);

-- マネーフォワードの入出金履歴データファイル
CREATE TABLE money_forward_files (
    id varchar(255) not null,
    name text not null,
    created_at datetime not null,
    imported_at datetime,
    primary key (id)
);
