alter table test_results
add column part_number text not null;

alter table test_results
add column product_test_sequence int not null;

alter table test_results
drop column product_test_configuration_id;

