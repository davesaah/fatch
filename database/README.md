# The Sweet Database

## Install Postgres 17 Options

- Use database-compose.yml file.
- Install postgres17, postgresql17-contrib, postgresql17-server on your distro.
  - `sudo systemctl enable postgresql`
  - `sudo postgresql-setup --initdb`.
  - `sudo systemctl start postgresql`
  - `sudo nano /var/lib/pgsql/data/pg_hba.conf   # RHEL/Fedora`. Change all
    the options for peer into md5.

## Order of Loading Schema

1. Extensions
2. Types
3. Tables
4. Indexes
5. Constraints
6. Functions
7. Procedures
8. Triggers
