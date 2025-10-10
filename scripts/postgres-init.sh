#!/bin/bash
# PostgreSQL initialization script

set -e

# Create additional databases if needed
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Create extensions
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
    CREATE EXTENSION IF NOT EXISTS "pg_stat_statements";
    
    -- Create additional schemas if needed
    CREATE SCHEMA IF NOT EXISTS garmin;
    CREATE SCHEMA IF NOT EXISTS analytics;
    
    -- Grant permissions
    GRANT ALL PRIVILEGES ON DATABASE $POSTGRES_DB TO $POSTGRES_USER;
    GRANT ALL PRIVILEGES ON SCHEMA garmin TO $POSTGRES_USER;
    GRANT ALL PRIVILEGES ON SCHEMA analytics TO $POSTGRES_USER;
    
    -- Create indexes for performance
    -- (Add your specific indexes here)
    
EOSQL

echo "PostgreSQL initialization completed successfully"