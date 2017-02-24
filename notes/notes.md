# Notes.md

## 2/22/2017

Solve what's in front of you - one thing at a time

A tree of database information is unique for every descriptor : database, system, schema, etc...

The system piece is driven by the local needs.  In a lot of cases, the schema and system are always the same.  So how to compress this effectively.

Database, Schema are usually programmatic.  In a data warehouse situation, you can have many systems feeding in data.  That's important, but more from a tagging perspective.  

One server can have multiple databases.  At the end of the data, many of the top level descriptors (server, database, schema) can be repetitive to a set of tables.  So maybe we go back to the basics.  The two required columns are Table and Column.  Everything above that is important if there are multiple databases being tracked (highly likely).  So even if in the circumstance of a single set of tables situation, you have to remember that the meta server might end up containing multiple unrelated databases.

Todo
[x] Figure out function to collapse repetitive values for server, database and schema

## 2/23/2017

There are three core pieces to the mapping:

* Database
* Table
* Column

Everything else is nuance.

Databases have server names, but usually just one.
Tables have schemas, but usually it's not relevant (not always true)
Columns have a bunch of detail attached to them - that's where the fun stuff really lies.

Search has to cover everything.

That was a lot of work getting the pieces together.  For now, databases and tables are functional.  Next to do the column page. and then think about the overall theme.  Then search again.
