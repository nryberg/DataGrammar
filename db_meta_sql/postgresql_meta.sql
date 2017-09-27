SELECT 'localhost' as server, -- replace as needed
      table_catalog as database,
      'Postgresql' as type,
      'System Name' as system, -- typically the source system name
      TABLE_SCHEMA as schema,
      TABLE_NAME as table,
      column_name as column,
      ordinal_position as ordinal,
      data_type as type,
      character_maximum_length as length,
      numeric_precision as precision,
      numeric_scale as scale

FROM  information_schema.columns
WHERE table_schema = 'public'
ORDER BY TABLE_NAME, column_name ASC
