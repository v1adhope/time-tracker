create or replace function uuid6() returns uuid as $$
declare
begin
	return uuid6(clock_timestamp());
end $$ language plpgsql;

create or replace function uuid6(p_timestamp timestamp with time zone) returns uuid as $$
declare

	v_time double precision := null;

	v_gregorian_t bigint := null;
	v_clock_sequence_and_node bigint := null;

	v_gregorian_t_hex_a varchar := null;
	v_gregorian_t_hex_b varchar := null;
	v_clock_sequence_and_node_hex varchar := null;

	c_epoch double precision := 12219292800; -- RFC-9562 epoch: 1582-10-15
	c_100ns_factor double precision := 10^7; -- RFC-9562 precision: 100 ns

	c_version bigint := x'0000000000006000'::bigint; -- RFC-9562 version: b'0110...'
	c_variant bigint := x'8000000000000000'::bigint; -- RFC-9562 variant: b'10xx...'

begin

	v_time := extract(epoch from p_timestamp);

	v_gregorian_t := trunc((v_time + c_epoch) * c_100ns_factor);
	v_clock_sequence_and_node := trunc(random() * 2^30)::bigint << 32 | trunc(random() * 2^32)::bigint;

	v_gregorian_t_hex_a := lpad(to_hex((v_gregorian_t >> 12)), 12, '0');
	v_gregorian_t_hex_b := lpad(to_hex((v_gregorian_t & 4095) | c_version), 4, '0');
	v_clock_sequence_and_node_hex := lpad(to_hex(v_clock_sequence_and_node | c_variant), 16, '0');

	return (v_gregorian_t_hex_a || v_gregorian_t_hex_b  || v_clock_sequence_and_node_hex)::uuid;

end $$ language plpgsql;

create table if not EXISTS users (
  user_id uuid default uuid6(),
  surname varchar(255) not null,
  name varchar(255) not null,
  patronymic varchar(255) not null,
  address varchar(255) not null,
  passport_number varchar(255) not null unique,

  constraint pk_users_user_id primary key(user_id)
);

create index if not exists index_users_passport_number on users(passport_number);

create table if not exists tasks (
  task_id uuid default uuid6(),
  created_at timestamp not null,
  finished_at timestamp,
  user_id uuid,

  constraint pk_tasks_task_id primary key(task_id),
  constraint fk_tasks_users_user_id foreign key(user_id) references users(user_id) on delete cascade
);

create index if not exists index_tasks_user_id on tasks(user_id);
