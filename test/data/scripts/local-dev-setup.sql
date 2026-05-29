insert into products (product_name)
values ('test product 1');

insert into test_unit_categories (category)
values
    -- Spatial & Geometry
    ('length'),
    ('volume'),
    ('angle'),

    -- Mass, Force, & Mechanics
    ('weight'),
    ('mass'),
    ('force'),
    ('torque'),
    ('pressure'),

    -- Thermodynamics & Fluids
    ('temperature'),
    ('density'),
    ('flow rate'),

    -- Time, Dynamics, & Rates
    ('duration'),
    ('frequency'),
    ('velocity'),

    -- Energy & Electricity
    ('power'),
    ('energy'),
    ('current'),
    ('voltage'),
    ('resistance'),

    -- Material Properties & Chemistry
    ('concentration'),
    ('surface roughness'),

    -- Data & Discrete Tracking
    ('data'),
    ('count');

insert into test_units (unit, unit_category, unit_abbreviation, measurement_system)
values
    -- Length
    ('nanometers', 'length', 'nm', 'metric'),
    ('micrometers', 'length', 'µm', 'metric'),
    ('millimeters', 'length', 'mm', 'metric'),
    ('centimeters', 'length', 'cm', 'metric'),
    ('meters', 'length', 'm', 'metric'),
    ('kilometers', 'length', 'km', 'metric'),
    ('inches', 'length', 'in', 'imperial'),
    ('feet', 'length', 'ft', 'imperial'),
    ('yards', 'length', 'yd', 'imperial'),
    ('miles', 'length', 'mi', 'imperial'),

    -- Volume
    ('milliliters', 'volume', 'mL', 'metric'),
    ('centiliters', 'volume', 'cL', 'metric'),
    ('deciliters', 'volume', 'dL', 'metric'),
    ('liters', 'volume', 'L', 'metric'),
    ('fluid ounces', 'volume', 'fl oz', 'imperial'),
    ('cups', 'volume', 'cup', 'imperial'),
    ('pints', 'volume', 'pt', 'imperial'),
    ('quarts', 'volume', 'qt', 'imperial'),
    ('gallons', 'volume', 'gal', 'imperial'),

    -- Weight
    ('ounces', 'weight', 'oz', 'imperial'),
    ('pounds', 'weight', 'lb', 'imperial'),

    -- Mass
    ('micrograms', 'mass', 'µg', 'metric'),
    ('milligrams', 'mass', 'mg', 'metric'),
    ('grams', 'mass', 'g', 'metric'),
    ('kilograms', 'mass', 'kg', 'metric'),

    -- Force
    ('newtons', 'force', 'N', 'metric'),

    -- Pressure
    ('pounds per square inch', 'pressure', 'psi', 'imperial'),
    ('kilopascals', 'pressure', 'kPa', 'metric'),

    --Torque
    ('millinewton-meters', 'torque', 'mNm', 'metric'),
    ('newton-meters', 'torque', 'Nm', 'metric'),
    ('kilonewton-meters', 'torque', 'kNm', 'metric'),
    ('ounce-inches', 'torque', 'oz-in', 'imperial'),
    ('pound-inches', 'torque', 'lb-in', 'imperial'),
    ('pound-feet', 'torque', 'lb-ft', 'imperial'),

    -- Thermodynamics & Fluids
    ('degrees Celsius', 'temperature', '°C', 'metric'),
    ('degrees Fahrenheit', 'temperature', '°F', 'imperial'),
    ('kelvin', 'temperature', 'K', 'metric'),
    ('grams per cubic centimeter', 'density', 'g/cm³', 'metric'),
    ('pounds per cubic foot', 'density', 'lb/ft³', 'imperial'),
    ('liters per minute', 'flow rate', 'L/min', 'metric'),
    ('gallons per minute', 'flow rate', 'gpm', 'imperial'),

    -- Time, Dynamics, & Rates
    ('milliseconds', 'duration', 'ms', 'metric'),
    ('seconds', 'duration', 's', 'metric'),
    ('minutes', 'duration', 'min', 'metric'),
    ('hours', 'duration', 'h', 'metric'),
    ('hertz', 'frequency', 'Hz', 'metric'),
    ('kilohertz', 'frequency', 'kHz', 'metric'),
    ('meters per second', 'velocity', 'm/s', 'metric'),
    ('feet per second', 'velocity', 'ft/s', 'imperial'),
    ('miles per hour', 'velocity', 'mph', 'imperial'),

    -- Energy & Electricity
    ('watts', 'power', 'W', 'metric'),
    ('kilowatts', 'power', 'kW', 'metric'),
    ('horsepower', 'power', 'hp', 'imperial'),
    ('joules', 'energy', 'J', 'metric'),
    ('kilowatt-hours', 'energy', 'kWh', 'metric'),
    ('British thermal units', 'energy', 'BTU', 'imperial'),
    ('volts', 'voltage', 'V', 'metric'),
    ('amperes', 'current', 'A', 'metric'),
    ('ohms', 'resistance', 'Ω', 'metric'),

    -- Material Properties & Chemistry
    ('parts per million', 'concentration', 'ppm', 'metric'),
    ('parts per billion', 'concentration', 'ppb', 'metric'),
    ('percentage', 'concentration', '%', 'metric'),
    ('micrometers ra', 'surface roughness', 'µm Ra', 'metric'),
    ('microinches ra', 'surface roughness', 'µin Ra', 'imperial'),

    -- Data & Discrete Tracking
    ('bytes', 'data', 'B', 'metric'),
    ('kilobytes', 'data', 'KB', 'metric'),
    ('megabytes', 'data', 'MB', 'metric'),
    ('gigabytes', 'data', 'GB', 'metric'),
    ('terabytes', 'data', 'TB', 'metric'),
    ('pieces', 'count', 'pcs', 'any'),
    ('boxes', 'count', 'box', 'any'),
    ('each', 'count', 'ea', 'any');

insert into tests (test_name, test_unit_category, documentation_references)
values
    -- Spatial & Geometry
    ('length', 'length', array []::text[]),
    ('width', 'length', array []::text[]),
    ('height', 'length', array []::text[]),
    ('thickness', 'length', array []::text[]),
    ('depth', 'length', array []::text[]),
    ('clearance', 'length', array []::text[]),
    ('diameter', 'length', array []::text[]),
    ('radius', 'length', array []::text[]),
    ('elongation', 'length', array []::text[]),
    ('travel distance', 'length', array []::text[]),
    ('deflection', 'length', array []::text[]),

    -- Volume
    ('total capacity', 'volume', array []::text[]),
    ('fill volume', 'volume', array []::text[]),
    ('headspace volume', 'volume', array []::text[]),
    ('shot size', 'volume', array []::text[]),
    ('dosage', 'volume', array []::text[]),
    ('leak volume', 'volume', array []::text[]),
    ('void volume', 'volume', array []::text[]),
    ('displacement volume', 'volume', array []::text[]),

    -- Weight (Force-based / Scale measurements)
    ('gross weight', 'weight', array []::text[]),
    ('net weight', 'weight', array []::text[]),
    ('tare weight', 'weight', array []::text[]),

    -- Mass (Intrinsic material quantity)
    ('material mass', 'mass', array []::text[]),
    ('mass variance', 'mass', array []::text[]),
    ('coating mass', 'mass', array []::text[]),

    -- Force (Push/Pull dynamics)
    ('tensile strength', 'force', array []::text[]),
    ('compressive strength', 'force', array []::text[]),
    ('insertion force', 'force', array []::text[]),
    ('extraction force', 'force', array []::text[]),
    ('actuation force', 'force', array []::text[]),

    -- Pressure (Force per unit area)
    ('burst pressure', 'pressure', array []::text[]),
    ('operating pressure', 'pressure', array []::text[]),
    ('vacuum level', 'pressure', array []::text[]),
    ('pressure drop', 'pressure', array []::text[]),

    -- Torque (Rotational / Fastener testing)
    ('fastener torque', 'torque', array []::text[]),
    ('breakaway torque', 'torque', array []::text[]),
    ('stripping torque', 'torque', array []::text[]),
    ('removal torque', 'torque', array []::text[]),

    -- Temperature (Thermal states & processing)
    ('ambient temperature', 'temperature', array []::text[]),
    ('melt temperature', 'temperature', array []::text[]),
    ('curing temperature', 'temperature', array []::text[]),
    ('operating temperature', 'temperature', array []::text[]),

    -- Density (Material purity & composition)
    ('bulk density', 'density', array []::text[]),
    ('material density', 'density', array []::text[]),
    ('specific gravity', 'density', array []::text[]),

    -- Flow Rate (Fluid dynamics & dispensing speed)
    ('volumetric flow rate', 'flow rate', array []::text[]),
    ('melt flow index', 'flow rate', array []::text[]),
    ('leak rate fluid', 'flow rate', array []::text[]),
    ('cycle time', 'duration', array []::text[]),
    ('dwell time', 'duration', array []::text[]),
    ('response time', 'duration', array []::text[]),
    ('exposure time', 'duration', array []::text[]),

    -- Frequency (Vibration, rotation, & cycles)
    ('rotational speed', 'frequency', array []::text[]),
    ('resonant frequency', 'frequency', array []::text[]),
    ('sampling rate', 'frequency', array []::text[]),

    -- Velocity (Linear & angular speeds)
    ('linear velocity', 'velocity', array []::text[]),
    ('impact velocity', 'velocity', array []::text[]),

    -- Power (Instantaneous energy consumption or output)
    ('power consumption', 'power', array []::text[]),
    ('peak power output', 'power', array []::text[]),
    ('standby power', 'power', array []::text[]),
    ('thermal power dissipation', 'power', array []::text[]),

    -- Energy (Total work done or storage capacity over time)
    ('battery capacity energy', 'energy', array []::text[]),
    ('energy consumption per cycle', 'energy', array []::text[]),
    ('impact energy', 'energy', array []::text[]),
    ('caloric heat content', 'energy', array []::text[]),

    -- Concentration (Chemical mixes & purity)
    ('chemical purity', 'concentration', array []::text[]),
    ('solution strength', 'concentration', array []::text[]),
    ('humidity level', 'concentration', array []::text[]),
    ('contamination level', 'concentration', array []::text[]),

    -- Surface Roughness (Micro-texture & finish)
    ('surface roughness ra', 'surface roughness', array []::text[]),
    ('peak to valley rz', 'surface roughness', array []::text[]),
    ('max profile height rt', 'surface roughness', array []::text[]),

    -- Data (Digital testing & file logging)
    ('log file size', 'data', array []::text[]),
    ('firmware image size', 'data', array []::text[]),

    -- Count (Discrete quantities & counting scales)
    ('batch quantity', 'count', array []::text[]),
    ('defect count', 'count', array []::text[]),
    ('sample size', 'count', array []::text[])
;

insert into modifiers (modifier)
values ('top'),
       ('bottom'),
       ('left'),
       ('right'),
       ('center'),
       ('diagonal'),
       ('inside'),
       ('outside'),
       ('front'),
       ('back'),
       ('inner'),
       ('outer'),
       ('overflow'),
       ('nominal'),
       ('headspace'),
       ('intake'),
       ('discharge'),
       ('initial'),
       ('residual'),
       ('upper-chamber'),
       ('lower-chamber'),
       ('dry'),
       ('wet'),
       ('tare'),
       ('push'),
       ('pull'),
       ('inlet'),
       ('outlet'),
       ('static'),
       ('dynamic'),
       ('internal'),
       ('surface'),
       ('input'),
       ('output'),
       ('pre-test'),
       ('post-test'),
       ('warp'),
       ('fill');


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