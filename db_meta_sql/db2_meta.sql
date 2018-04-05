 SELECT     'xyz' as Database,
        'DB2' as Database_Type,
        TBCreator as Schema,
        TBNAME as Table,
        NAME as Column,
        COLNO as Ordinal,
        COLTYPE as Type,
        LONGLENGTH as Length,
        Length as Precision,
        Scale,
        CURRENT_DATE as Run_Date
FROM SYSIBM.SYSCOLUMNS
ORDER BY TBCreator, TBName, Name