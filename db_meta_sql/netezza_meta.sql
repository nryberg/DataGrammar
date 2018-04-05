 SELECT   'abc' as server, 
        'Netezza' as Database_Type,
        Current_DB as Database_Name,
        DBNAME as Schema_Name,
        TABLENAME as Table_Name,
        COLUMN_NAME as Column_Name,
        ORDINAL_POSITION as Column_Ordinal,
        TYPE_NAME as Column_Type,
        COLUMN_SIZE as Column_Length,
        COLUMN_SIZE as Column_Precision,
        DECIMAL_DIGITS AS Column_Scale,
        CURRENT_DATE as Run_Date
FROM TABLE_COLS 
ORDER BY DBNAME, TABLENAME, COLUMN_NAME

