# Majipoor - a tool to synchronize mysql into a postgresql datalake

This tool is a simple program to synchronize a mysql database to a postgresql database
for datalake purposes. 

It does not try to match the schema perfectly, but instead tries to be loose enough so that 
the laxity of mysql doesn't cause problem with a more strongly enforced schema in Postgresql.

It is inspired by [pg_chameleon](https://pgchameleon.org/), which is better suited if you 
are looking for schema enforcement.

I am building this tool to synchronize a woocommerce system into a postgresql datalake, 
so that I can reliably run dbt to transform WP/WC's schema into something usable.

## Mysql Schema Introspection