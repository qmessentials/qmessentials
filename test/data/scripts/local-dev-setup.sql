
with prod_lookup as ( select id, product_name from products ),
     test_lookup as ( select id, test_name from tests where test_name in ('thickness', 'warp_count', 'fill_volume') ),
     config_matrix(prod_name, test_name, modifiers, unit, dec_places, min_v, max_v, critical)
         as ( values ('test product 1', 'thickness', array ['top', 'left'], 'micrometers', 2, 0.0, 12.5, true),
                     ('test product 1', 'warp count', null, 'pieces', 0, 0, 2, false)
            )
insert into product_test_configurations
    (product_id, test_id, specific_modifiers, unit, decimal_places, min_value, max_value, is_critical)
select
    p.id,
    t.id,
    m.modifiers,
    m.unit,
    m.dec_places,
    m.min_v,
    m.max_v,
    m.critical
from config_matrix m
join prod_lookup p on p.product_name = m.prod_name
join test_lookup t on t.test_name = m.test_name;