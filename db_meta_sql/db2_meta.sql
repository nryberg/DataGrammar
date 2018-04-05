 SELECT     'abc' as server,
        'xyz' as database_name,
        'DB2' as Database_Type,
        TBCreator as Schema_Name,
        TBNAME as Table_Name,
        NAME as Column_Name,
        COLNO as Column_Ordinal,
        COLTYPE as Column_Type,
        LONGLENGTH as Column_Length,
        Length as Column_Precision,
        Scale as column_scale,
        CURRENT_DATE as Run_Date
FROM SYSIBM.SYSCOLUMNS
ORDER BY TBCreator, TBName, Name

