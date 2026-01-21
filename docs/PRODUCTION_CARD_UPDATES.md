# Production Card Data Updates

This document explains how to update card data on the production server using Scryfall bulk data.

## Quick Start

### Update Production Cards (Recommended)

```bash
# Update with latest Scryfall data
./update-cards-prod.sh --download

# Update with existing local data
./update-cards-prod.sh

# Dry run (see what would happen without executing)
./update-cards-prod.sh --download --dry-run
```

### Rollback (if needed)

```bash
./update-cards-prod.sh --rollback
```

---

## What the Script Does

The `update-cards-prod.sh` script safely updates card data on production:

1. **Checks connectivity** to production server
2. **Downloads** latest Scryfall data (optional)
3. **Creates backup** of current card data
4. **Copies data** to production server
5. **Runs migration** to create Scryfall tables
6. **Imports cards** using the streaming importer
7. **Creates compatibility view** for backward compatibility
8. **Verifies** the update was successful

### Safety Features

- ✅ **Automatic backups** before any changes
- ✅ **Dry-run mode** to preview changes
- ✅ **One-command rollback** if issues arise
- ✅ **Zero downtime** - server stays running during update
- ✅ **Compatibility layer** - existing code continues to work

---

## Usage Examples

### Standard Update (Using Existing Data)

```bash
./update-cards-prod.sh
```

This uses the most recent Scryfall data file found locally (e.g., `all-cards-20260111103023.json`).

### Update with Fresh Download

```bash
./update-cards-prod.sh --download
```

This downloads the latest Scryfall bulk data (~2.3GB) before updating production.

### Preview Changes (Dry Run)

```bash
./update-cards-prod.sh --download --dry-run
```

Shows what would be executed without making any changes.

### Rollback to Previous Data

```bash
./update-cards-prod.sh --rollback
```

Restores the previous card data from the XMage backup.

---

## Configuration

The script uses the same configuration as `deploy.sh`:

```bash
# Set via environment variables (optional)
export REMOTE_USER="hkdebiandocker"
export REMOTE_HOST="192.168.178.24"
export REMOTE_PATH="gomage"

# Then run
./update-cards-prod.sh --download
```

Default values match your production setup.

---

## Update Frequency

### Recommended Schedule

- **Weekly**: Update cards every Sunday to catch new releases
- **On demand**: After major set releases
- **Emergency**: If cards are missing or incorrect

### Automated Updates (Optional)

Set up a cron job to automatically update weekly:

```bash
# Add to your local crontab
# Every Sunday at 3 AM
0 3 * * 0 cd /Users/aron/dev/opensource/mage && ./update-cards-prod.sh --download >> logs/card-updates.log 2>&1
```

Or use GitHub Actions (see `.github/workflows/mtg-fetch-cards.yml` for inspiration).

---

## Verification

After updating, verify the production server:

### 1. Check Card Counts

```bash
ssh hkdebiandocker@192.168.178.24 'cd gomage && docker compose exec postgres psql -U mage -d mage -c "SELECT COUNT(*) FROM scryfall_cards WHERE lang='\''en'\'';"'
```

Expected: ~101,000+ cards

### 2. Test Card Lookup

```bash
ssh hkdebiandocker@192.168.178.24 'cd gomage && docker compose exec postgres psql -U mage -d mage -c "SELECT name, card_type FROM cards WHERE name LIKE '\''Lightning Bolt'\'' LIMIT 1;"'
```

Should return: Lightning Bolt | Instant

### 3. Check Server Status

```bash
curl http://192.168.178.24:17171/status
```

Should return JSON with server stats.

### 4. Test Frontend

Visit: http://mage.aronschueler.de

Try creating a deck and searching for cards.

---

## Backups

### Automatic Backups

Every update creates:
- **SQL backup**: `cards_backup_YYYYMMDD_HHMMSS.sql`
- **Info file**: `backup_YYYYMMDD_HHMMSS_info.txt`

Location: `${REMOTE_HOST}:${REMOTE_PATH}/mage-server-go/data/`

### Manual Backup

```bash
# Before manual changes
ssh hkdebiandocker@192.168.178.24 'cd gomage && docker compose exec postgres pg_dump -U mage -d mage -t cards -t scryfall_cards > backup_manual.sql'
```

### Restore from Backup

```bash
# Option 1: Use rollback script
./update-cards-prod.sh --rollback

# Option 2: Manual restore
ssh hkdebiandocker@192.168.178.24 'cd gomage && docker compose exec -T postgres psql -U mage -d mage < data/cards_backup_YYYYMMDD_HHMMSS.sql'
```

---

## Troubleshooting

### "Cannot connect to remote server"

**Check:**
1. VPN/network connection to 192.168.178.24
2. SSH key is configured
3. Server is running

**Test:**
```bash
ssh hkdebiandocker@192.168.178.24 "echo 'Connection OK'"
```

### "No Scryfall data file found"

**Solution:**
```bash
# Download data first
cd mage-server-go
./scripts/download_scryfall_bulk.sh

# Then run update
cd ..
./update-cards-prod.sh
```

### "Import failed with errors"

**Check:**
1. Database has enough disk space
2. PostgreSQL container is healthy
3. Review logs:
```bash
ssh hkdebiandocker@192.168.178.24 'cd gomage && docker compose logs postgres'
```

**Rollback if needed:**
```bash
./update-cards-prod.sh --rollback
```

### "Cards not found after update"

**Likely cause:** Modal DFC cards (double-faced) need special handling.

**Check:**
```bash
ssh hkdebiandocker@192.168.178.24 'cd gomage && docker compose exec postgres psql -U mage -d mage -c "SELECT name FROM cards WHERE name LIKE '\''%Birgi%'\'';"'
```

Should show: `Birgi, God of Storytelling // Harnfel, Horn of Bounty`

The code now handles these automatically (see `SCRYFALL_DFC_CARDS.md`).

---

## Performance Impact

### During Update

- **Duration**: ~2-3 minutes
- **Server impact**: Minimal - import runs in background
- **Downtime**: None - server continues running
- **Database load**: Moderate during import, then normal

### After Update

- **Query performance**: Similar or better (optimized indexes)
- **Disk space**: +500MB for Scryfall data
- **Memory**: No significant change

---

## Migration from XMage

The production update script automatically:

1. **Backs up** the old XMage `cards` table as `cards_xmage_backup`
2. **Creates** the new `scryfall_cards` table with full data
3. **Creates** a `cards` view for backward compatibility
4. **Preserves** all existing deck data and user data

Your existing decks and code continue to work unchanged.

---

## Advanced Options

### Custom Remote Server

```bash
REMOTE_HOST="different-server.com" \
REMOTE_USER="myuser" \
REMOTE_PATH="custom/path" \
./update-cards-prod.sh --download
```

### Skip Download (Use Specific File)

```bash
# Ensure the file is in mage-server-go/data/ or root directory
ls -lh all-cards-20260111103023.json

# Then run without --download
./update-cards-prod.sh
```

### Update Only Migrations (No Data)

```bash
# Copy just the migration files
rsync -avz mage-server-go/migrations/ hkdebiandocker@192.168.178.24:gomage/mage-server-go/migrations/

# Apply manually
ssh hkdebiandocker@192.168.178.24 'cd gomage && docker compose exec postgres psql -U mage -d mage < mage-server-go/migrations/009_create_scryfall_tables.up.sql'
```

---

## Related Documentation

- **Migration Guide**: `SCRYFALL_MIGRATION_GUIDE.md` - User guide for Scryfall data
- **Migration Plan**: `DATA_MIGRATION.md` - Technical implementation details
- **DFC Cards**: `SCRYFALL_DFC_CARDS.md` - Handling double-faced cards
- **Deployment**: `deploy.sh` - Full application deployment
- **Migration Summary**: `MIGRATION_COMPLETE.md` - Initial migration results

---

## Support

### Check Production Status

```bash
# Server health
curl http://192.168.178.24:8080/health

# API status
curl http://192.168.178.24:17171/status

# Database status
ssh hkdebiandocker@192.168.178.24 'cd gomage && docker compose exec postgres psql -U mage -d mage -c "SELECT version();"'
```

### View Logs

```bash
# All services
ssh hkdebiandocker@192.168.178.24 'cd gomage && docker compose logs -f'

# Just database
ssh hkdebiandocker@192.168.178.24 'cd gomage && docker compose logs -f postgres'

# Just backend
ssh hkdebiandocker@192.168.178.24 'cd gomage && docker compose logs -f mage-server'
```

### Emergency Recovery

If something goes wrong:

```bash
# 1. Rollback cards
./update-cards-prod.sh --rollback

# 2. Restart services
ssh hkdebiandocker@192.168.178.24 'cd gomage && docker compose restart'

# 3. Check status
curl http://192.168.178.24:17171/status

# 4. If still broken, full redeploy
./deploy.sh
```
