insert into products (part_number, product_name)
values ('KE-TST-TC00001', 'Toaster Case 1'),
('KE-TST-HE00001', 'Toaster Heating Element 1'),
('KE-TST-LF00001', 'Toaster Spring Lift 1');

-- =============================================================================
-- Toaster Component Test Configurations — Demo / Testing Data
-- Products: Toaster Case, Heating Element, Spring Lift
-- =============================================================================
-- All product and test lookups use domain keys; no hardcoded IDs.
--
-- Each product_test_configurations INSERT uses INSERT ... SELECT with a CTE
-- that resolves part_number -> product_id and canonical_test_name -> test_id.
--
-- Each modifier combination INSERT resolves the parent configuration by joining
-- product_test_configurations back through product and canonical test keys.
-- =============================================================================


-- =============================================================================
-- TOASTER CASE (KE-TST-TC00001)
-- =============================================================================

-- Overall width (left-to-right external dimension)
with p as (
    select id from products
    where part_number = 'KE-TST-TC00001'
),

t as (
    select id from tests
    where canonical_test_name = 'width'
)

insert into product_test_configurations
(
    product_id, test_id, specific_modifiers, unit, decimal_places,
    min_value, max_value, is_critical
)
select
    p.id,
    t.id,
    null,
    'millimeters',
    1,
    280.0,
    310.0,
    true
from p, t;

-- Overall height
with p as (
    select id from products
    where part_number = 'KE-TST-TC00001'
),

t as (
    select id from tests
    where canonical_test_name = 'height'
)

insert into product_test_configurations
(
    product_id, test_id, specific_modifiers, unit, decimal_places,
    min_value, max_value, is_critical
)
select
    p.id,
    t.id,
    null,
    'millimeters',
    1,
    195.0,
    215.0,
    true
from p, t;

-- Wall thickness — top/bottom/left/right sides
with p as (
    select id from products
    where part_number = 'KE-TST-TC00001'
),

t as (
    select id from tests
    where canonical_test_name = 'thickness'
)

insert into product_test_configurations
(
    product_id, test_id, specific_modifiers, unit, decimal_places,
    min_value, max_value, is_critical
)
select
    p.id,
    t.id,
    array['top', 'bottom', 'left', 'right'],
    'millimeters',
    2,
    2.00,
    4.00,
    false
from p, t;

-- Gross weight (assembled case, no internals)
with p as (
    select id from products
    where part_number = 'KE-TST-TC00001'
),

t as (
    select id from tests
    where canonical_test_name = 'gross_weight'
)

insert into product_test_configurations
(
    product_id, test_id, specific_modifiers, unit, decimal_places,
    min_value, max_value, is_critical
)
select
    p.id,
    t.id,
    null,
    'grams',
    0,
    380.0,
    460.0,
    false
from p, t;

-- Surface roughness Ra — external visible surfaces only
with p as (
    select id from products
    where part_number = 'KE-TST-TC00001'
),

t as (
    select id from tests
    where canonical_test_name = 'surface_roughness_ra'
)

insert into product_test_configurations
(
    product_id, test_id, specific_modifiers, unit, decimal_places,
    min_value, max_value, is_critical
)
select
    p.id,
    t.id,
    array['outside'],
    'micrometers ra',
    2,
    0.40,
    1.60,
    false
from p, t;

-- Top-surface deflection under a 5 kg point load (quality/durability check)
with p as (
    select id from products
    where part_number = 'KE-TST-TC00001'
),

t as (
    select id from tests
    where canonical_test_name = 'deflection'
)

insert into product_test_configurations
(
    product_id, test_id, specific_modifiers, unit, decimal_places,
    min_value, max_value, is_critical
)
select
    p.id,
    t.id,
    array['top'],
    'millimeters',
    2,
    null,
    2.00,
    true
from p, t;

-- External surface temperature at operating limits (safety-critical)
with p as (
    select id from products
    where part_number = 'KE-TST-TC00001'
),

t as (
    select id from tests
    where canonical_test_name = 'operating_temperature'
)

insert into product_test_configurations
(
    product_id, test_id, specific_modifiers, unit, decimal_places,
    min_value, max_value, is_critical
)
select
    p.id,
    t.id,
    array['top', 'left', 'right'],
    'degrees Celsius',
    1,
    null,
    75.0,
    true
from p, t;

-- Crumb-tray insertion force (ease-of-use spec)
with p as (
    select id from products
    where part_number = 'KE-TST-TC00001'
),

t as (
    select id from tests
    where canonical_test_name = 'insertion_force'
)

insert into product_test_configurations
(
    product_id, test_id, specific_modifiers, unit, decimal_places,
    min_value, max_value, is_critical
)
select
    p.id,
    t.id,
    null,
    'newtons',
    1,
    2.0,
    12.0,
    false
from p, t;

-- =============================================================================
-- TOASTER HEATING ELEMENT (KE-TST-HE00001)
-- =============================================================================

-- Element length (active coil span)
with p as (
    select id from products
    where part_number = 'KE-TST-HE00001'
),

t as (
    select id from tests
    where canonical_test_name = 'length'
)

insert into product_test_configurations
(
    product_id, test_id, specific_modifiers, unit, decimal_places,
    min_value, max_value, is_critical
)
select
    p.id,
    t.id,
    null,
    'millimeters',
    1,
    245.0,
    265.0,
    true
from p, t;

-- Wire diameter (determines resistance and heat output)
with p as (
    select id from products
    where part_number = 'KE-TST-HE00001'
),

t as (
    select id from tests
    where canonical_test_name = 'diameter'
)

insert into product_test_configurations
(
    product_id, test_id, specific_modifiers, unit, decimal_places,
    min_value, max_value, is_critical
)
select
    p.id,
    t.id,
    null,
    'millimeters',
    3,
    0.280,
    0.340,
    true
from p, t;

-- Mass variance (lot-to-lot consistency check)
with p as (
    select id from products
    where part_number = 'KE-TST-HE00001'
),

t as (
    select id from tests
    where canonical_test_name = 'mass_variance'
)

insert into product_test_configurations
(
    product_id, test_id, specific_modifiers, unit, decimal_places,
    min_value, max_value, is_critical
)
select
    p.id,
    t.id,
    null,
    'grams',
    2,
    null,
    0.50,
    false
from p, t;

-- Tensile strength (must survive thermal cycling stress)
with p as (
    select id from products
    where part_number = 'KE-TST-HE00001'
),

t as (
    select id from tests
    where canonical_test_name = 'tensile_strength'
)

insert into product_test_configurations
(
    product_id, test_id, specific_modifiers, unit, decimal_places,
    min_value, max_value, is_critical
)
select
    p.id,
    t.id,
    null,
    'newtons',
    1,
    45.0,
    null,
    true
from p, t;

-- Power consumption at rated voltage
with p as (
    select id from products
    where part_number = 'KE-TST-HE00001'
),

t as (
    select id from tests
    where canonical_test_name = 'power_consumption'
)

insert into product_test_configurations
(
    product_id, test_id, specific_modifiers, unit, decimal_places,
    min_value, max_value, is_critical
)
select
    p.id,
    t.id,
    null,
    'watts',
    0,
    720.0,
    800.0,
    true
from p, t;

-- Peak element surface temperature (safety and materials limit)
with p as (
    select id from products
    where part_number = 'KE-TST-HE00001'
),

t as (
    select id from tests
    where canonical_test_name = 'operating_temperature'
)

insert into product_test_configurations
(
    product_id, test_id, specific_modifiers, unit, decimal_places,
    min_value, max_value, is_critical
)
select
    p.id,
    t.id,
    null,
    'degrees Celsius',
    0,
    null,
    760.0,
    true
from p, t;

-- =============================================================================
-- TOASTER SPRING LIFT (KE-TST-LF00001)
-- =============================================================================

-- Free (uncompressed) and fully compressed length
with p as (
    select id from products
    where part_number = 'KE-TST-LF00001'
),

t as (
    select id from tests
    where canonical_test_name = 'length'
)

insert into product_test_configurations
(
    product_id, test_id, specific_modifiers, unit, decimal_places,
    min_value, max_value, is_critical
)
select
    p.id,
    t.id,
    array['initial', 'residual'],
    'millimeters',
    1,
    null,
    null,
    false
from p, t;

-- Travel distance (stroke from fully compressed to fully extended)
with p as (
    select id from products
    where part_number = 'KE-TST-LF00001'
),

t as (
    select id from tests
    where canonical_test_name = 'travel_distance'
)

insert into product_test_configurations
(
    product_id, test_id, specific_modifiers, unit, decimal_places,
    min_value, max_value, is_critical
)
select
    p.id,
    t.id,
    null,
    'millimeters',
    1,
    62.0,
    78.0,
    true
from p, t;

-- Push-down (insertion) force — how hard the user presses the lever
with p as (
    select id from products
    where part_number = 'KE-TST-LF00001'
),

t as (
    select id from tests
    where canonical_test_name = 'insertion_force'
)

insert into product_test_configurations
(
    product_id, test_id, specific_modifiers, unit, decimal_places,
    min_value, max_value, is_critical
)
select
    p.id,
    t.id,
    array['push'],
    'newtons',
    1,
    18.0,
    35.0,
    false
from p, t;

-- Spring return (extraction) force — must eject toast reliably
with p as (
    select id from products
    where part_number = 'KE-TST-LF00001'
),

t as (
    select id from tests
    where canonical_test_name = 'extraction_force'
)

insert into product_test_configurations
(
    product_id, test_id, specific_modifiers, unit, decimal_places,
    min_value, max_value, is_critical
)
select
    p.id,
    t.id,
    array['pull'],
    'newtons',
    1,
    14.0,
    30.0,
    true
from p, t;

-- Deflection under maximum rated load (structural integrity)
with p as (
    select id from products
    where part_number = 'KE-TST-LF00001'
),

t as (
    select id from tests
    where canonical_test_name = 'deflection'
)

insert into product_test_configurations
(
    product_id, test_id, specific_modifiers, unit, decimal_places,
    min_value, max_value, is_critical
)
select
    p.id,
    t.id,
    null,
    'millimeters',
    2,
    null,
    1.50,
    false
from p, t;

-- Pivot fastener torque — left and right pivot points
with p as (
    select id from products
    where part_number = 'KE-TST-LF00001'
),

t as (
    select id from tests
    where canonical_test_name = 'fastener_torque'
)

insert into product_test_configurations
(
    product_id, test_id, specific_modifiers, unit, decimal_places,
    min_value, max_value, is_critical
)
select
    p.id,
    t.id,
    array['left', 'right'],
    'newton-meters',
    2,
    0.80,
    1.20,
    true
from p, t;

-- Carriage cycle time (lever release to toast-at-top)
with p as (
    select id from products
    where part_number = 'KE-TST-LF00001'
),

t as (
    select id from tests
    where canonical_test_name = 'cycle_time'
)

insert into product_test_configurations
(
    product_id, test_id, specific_modifiers, unit, decimal_places,
    min_value, max_value, is_critical
)
select
    p.id,
    t.id,
    null,
    'seconds',
    2,
    0.30,
    0.90,
    false
from p, t;

-- Gross weight of the assembled lift mechanism
with p as (
    select id from products
    where part_number = 'KE-TST-LF00001'
),

t as (
    select id from tests
    where canonical_test_name = 'gross_weight'
)

insert into product_test_configurations
(
    product_id, test_id, specific_modifiers, unit, decimal_places,
    min_value, max_value, is_critical
)
select
    p.id,
    t.id,
    null,
    'grams',
    1,
    55.0,
    80.0,
    false
from p, t;

insert into samples (serial_number, part_number, status) values (
    'KETST-000001-001', 'KE-TST-TC00001', 'TESTING'
);
