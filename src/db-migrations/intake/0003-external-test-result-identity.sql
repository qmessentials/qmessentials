do $$
begin
    if exists (select 1 from test_results) then
        raise exception 'test_results must be empty before changing its external identity';
    end if;
end
$$;

alter table test_results
add column serial_number text not null references samples (serial_number);

alter table test_results
add column canonical_test_name text not null;

alter table test_results
drop column sample_id;

alter table test_results
drop column product_test_sequence;
