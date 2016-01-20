CREATE TABLE IF NOT EXISTS groups (
    group_id serial primary key,
    phase_name text not null
);

CREATE TABLE IF NOT EXISTS users (
    user_id serial primary key,
    user_balance decimal not null,
    name text not null,
    email text not null,
    password text not null,
    date_created timestamp with time zone not null,
    start_date timestamp with time zone not null,
    max_co integer not null,
    group_id integer REFERENCES groups,
    end_date timestamp with time zone not null,
    last_login timestamp with time zone not null
);

CREATE TABLE IF NOT EXISTS vouchers (
    voucher_id serial primary key,
    amount decimal not null,
    user_id int REFERENCES users,
    balance_before decimal not null,
    balance_after decimal not null,
    notes text
);

CREATE TABLE IF NOT EXISTS submissions (
    submission_id serial primary key,
    user_id int REFERENCES users,
    date_submitted timestamp with time zone not null,
    co_level integer not null,
    cigs decimal not null,
    validated boolean not null,
    video_url text not null,
    notes text
);
