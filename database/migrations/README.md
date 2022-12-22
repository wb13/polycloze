# Note to migration script author

When creating or dropping tables/views/indexes/triggers, it's important to
include the `IF NOT EXISTS` or `IF EXISTS` clause.
The migration scripts are ran every time the server connects to a database
(every request for requests that access the review or user databases).
Sometimes this causes `goose` to upgrade to the same version twice...

## Creating tables

```sql
CREATE TABLE IF NOT EXISTS ...
```

## Dropping tables

```sql
DROP TABLE IF EXISTS ...
```
