-- Development data for testing
-- This file is loaded only in development environment

-- Insert sample data for testing
INSERT INTO garmin.activity_types (type_id, type_name, created_at, updated_at)
VALUES 
    (1, 'Running', NOW(), NOW()),
    (2, 'Cycling', NOW(), NOW()),
    (3, 'Swimming', NOW(), NOW()),
    (4, 'Walking', NOW(), NOW()),
    (5, 'Hiking', NOW(), NOW())
ON CONFLICT (type_id) DO NOTHING;

-- Add more development-specific data here
-- INSERT INTO ... 

-- Create test users or sample activities if needed
-- INSERT INTO garmin.activities ...

-- Print completion message
SELECT 'Development data loaded successfully' as status;