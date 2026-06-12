alter table product_test_configurations
add column product_test_sequence integer;

with sequenced_rows as (
    select id, row_number() over (partition by product_id order by created_at, id) as seq
    from product_test_configurations
)
update product_test_configurations
set product_test_sequence = sequenced_rows.seq
from sequenced_rows
where product_test_configurations.id = sequenced_rows.id;

alter table product_test_configurations
alter column product_test_sequence set not null;

alter table product_test_configurations
add constraint product_test_sequence_unique unique (product_id, product_test_sequence);


