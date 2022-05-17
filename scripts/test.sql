.read database/migrations/1_init_schema.up.sql

attach database './cmd/spa.db' as l2;
attach database './cmd/eng.db' as l1;
attach database './cmd/translations.db' as translation;

.timer on
