# Archive Feature Documentation

## Overview

The archive feature automatically stores up to 10 historical versions of each configuration. Every time a config is updated, the previous version is automatically archived before the update is applied.

## Features

### Automatic Archiving
- **Trigger**: Every config update automatically archives the previous version
- **Version Numbering**: Archives are numbered sequentially (1, 2, 3, etc.)
- **Complete Snapshots**: Each archive contains the complete config state at the time of archiving

### Archive Management
- **Limit**: Maximum 10 archives per config
- **Automatic Cleanup**: When the limit is exceeded, the oldest archive is automatically removed
- **Storage**: Archives are stored in a separate `config_archives` collection
- **Deletion**: When a config is deleted, all its archives are automatically deleted

### Archive Data
Each archive contains:
- Complete config snapshot
- Version number
- Archive timestamp
- User who triggered the update
- All original config fields

## API Endpoints

### Get Config Archives
```
GET /config/archives?id={config_id}
```

**Response:**
```json
{
  "archives": [
    {
      "id": "archive_id",
      "config_id": "config_id",
      "name": "config_name",
      "type": "config_type",
      "subtype": "config_subtype",
      "tags": ["tag1", "tag2"],
      "tenant_id": "tenant_id",
      "created_by": "user_id",
      "last_updated_by": "user_id",
      "metadata": {...},
      "version": 1,
      "archived_at": "2024-01-01T00:00:00Z",
      "archived_by": "user_id",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 1
}
```

## Database Schema

### Config Archives Collection
```json
{
  "_id": "ObjectId",
  "config_id": "ObjectId",
  "name": "string",
  "type": "string",
  "subtype": "string",
  "tags": ["string"],
  "tenant_id": "string",
  "created_by": "string",
  "last_updated_by": "string",
  "metadata": {},
  "version": 1,
  "archived_at": "timestamp",
  "archived_by": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

## Usage Examples

### Creating and Updating Configs
```bash
# Create a config
curl -X POST http://localhost:8080/config \
  -H "Content-Type: application/json" \
  -d '{
    "name": "database-config",
    "type": "database",
    "subtype": "postgres",
    "tags": ["production"],
    "metadata": {"host": "localhost", "port": 5432}
  }'

# Update the config (creates archive automatically)
curl -X PUT http://localhost:8080/config?id=<config_id> \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": {"host": "localhost", "port": 5432, "ssl": true}
  }'
```

### Viewing Archives
```bash
# Get version history
curl -X GET http://localhost:8080/config/archives?id=<config_id>
```

## Testing

Use the provided test script to verify archive functionality:
```bash
chmod +x test_archive_feature.sh
./test_archive_feature.sh
```

## Implementation Details

### Archive Creation Process
1. When a config update is requested, the current version is retrieved
2. **Transaction begins** - All operations are wrapped in a MongoDB transaction
3. A new archive entry is created with the current config state
4. Version number is calculated (current archive count + 1)
5. Archive is saved to the `config_archives` collection
6. If archive count exceeds 10, the oldest archive is removed
7. The config is updated with the new values
8. **Transaction commits** - If any step fails, all changes are rolled back

### Archive Cleanup
- **Automatic cleanup** happens during update operations
- Only the oldest archive (lowest version number) is removed
- Cleanup is triggered when archive count exceeds `MaxArchiveHistory` (10)
- **Complete cleanup** happens when config is deleted - all archives are removed
- **Transaction safety** - All cleanup operations are wrapped in transactions for data consistency

### Error Handling
- Archive creation failures are handled gracefully
- If archiving fails, the update operation is aborted
- Archive retrieval errors return appropriate HTTP status codes
- **Transaction rollback** - If any operation in the transaction fails, all changes are automatically rolled back
- **Data consistency** - Ensures archive and config operations are always in sync

## Transaction Safety

### Data Consistency
The archive feature uses MongoDB transactions to ensure data consistency:

- **Update Operations**: Archive creation and config updates are atomic
- **Delete Operations**: Archive deletion and config deletion are atomic
- **Rollback Protection**: If any operation fails, all changes are automatically rolled back
- **Session Management**: Uses MongoDB sessions for transaction isolation

### Transaction Flow
```
Update Config:
1. Start Transaction
2. Create Archive
3. Update Config
4. Commit Transaction (or Rollback on failure)

Delete Config:
1. Start Transaction
2. Delete All Archives
3. Delete Config
4. Commit Transaction (or Rollback on failure)
```

## Performance Considerations

### Indexes
Consider adding these indexes for better performance:
```javascript
// Archive collection
db.config_archives.createIndex({"config_id": 1, "version": -1})
db.config_archives.createIndex({"tenant_id": 1})
```

### Storage
- Monitor archive collection size and growth
- Consider implementing archive retention policies if needed
- Archives are included in backup strategies

## Migration Notes

### For New Deployments
- Archive collection will be created automatically
- No existing data migration required
- Archive functionality is enabled by default

### For Existing Deployments
- Archive collection will be created on first update
- Existing configs will not have archives until updated
- No impact on existing functionality
