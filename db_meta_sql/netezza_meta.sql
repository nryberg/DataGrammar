 SELECT     'xyz' as Database,
        'Netezza' as Database_Type,
        DBNAME as Schema,
        TABLENAME as Table,
        COLUMN_NAME as Column,
        ORDINAL_POSITION as Ordinal,
        TYPE_NAME as Type,
        COLUMN_SIZE as Length,
        COLUMN_SIZE as Precision,
        DECIMAL_DIGITS AS Scale,
        CURRENT_DATE as Run_Date
FROM TABLE_COLS 
ORDER BY DBNAME, TABLENAME, COLUMN_NAME