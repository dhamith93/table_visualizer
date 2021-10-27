# Table visualizer

A tool to visually represent database table connections. Supports MySQL/MariaDB/Postgres.

## Usage

```bash
./table_viz -u <db-user> -p <db-password> -h <ip:port> -d <db-name> -t <table-name> -s mysql|mariadb|postgres -o <output-file-name.jpg>
```

* Example for postgress for all tables
```bash
./table_viz -u postgres -p 1234 -h 172.17.0.2:5432 -d shakespeare -s postgres -o out_p.jpg
```

* Example for postgress for single table. This shows given table and all depended tables recursively
```bash
./table_viz -u postgres -p 1234 -h 172.17.0.2:5432 -d shakespeare -t character_work -s postgres -o out_p_char_work.jpg
```

* Example for mariadb for all tables
```bash
./table_viz -u root -p 1234 -h 172.17.0.3:3306 -d test_db -s mariadb -o out_m.jpg
```

* Example for mariadb for single table. This shows given table and all depended tables recursively
```bash
./table_viz -u root -p 1234 -h 172.17.0.3:3306 -d test_db -t user_logs -s mariadb -o out_m_user_logs.jpg
```

## Examples

* Output when ran for single table `character`. Shows `character` table and all dependent tables.
![sample 1](/sample_out/out_p_char_work.jpg)

* Output whe ran without specifing a table. Shows all tables and connections.
![sample 2](/sample_out/out_p.jpg)