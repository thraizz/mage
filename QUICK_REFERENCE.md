# MAGE Quick Reference - Production Updates

## 🚀 Deploy Full Application

```bash
./deploy.sh
```

Deploys entire application (backend + frontend + database) to production.

---

## 🎴 Update Cards Only

### Standard Update
```bash
./update-cards-prod.sh --download
```

Updates card data with latest Scryfall data (~2-3 minutes).

### Preview Changes
```bash
./update-cards-prod.sh --download --dry-run
```

Shows what would happen without executing.

### Rollback
```bash
./update-cards-prod.sh --rollback
```

Restores previous card data from backup.

---

## 📅 Automated Weekly Updates

### Setup Cron Job
```bash
# Edit crontab
crontab -e

# Add this line (runs every Sunday at 3 AM):
0 3 * * 0 cd /Users/aron/dev/opensource/mage && ./update-cards-weekly.sh
```

### Manual Weekly Run
```bash
./update-cards-weekly.sh
```

Logs to `logs/card-update-YYYYMMDD_HHMMSS.log`

---

## 🔍 Verify Production

### Check Status
```bash
curl http://192.168.178.24:17171/status
```

### Check Card Counts
```bash
ssh hkdebiandocker@192.168.178.24 \
  'cd gomage && docker compose exec postgres psql -U mage -d mage -tAc "SELECT COUNT(*) FROM cards;"'
```

### Test Card Lookup
```bash
ssh hkdebiandocker@192.168.178.24 \
  'cd gomage && docker compose exec postgres psql -U mage -d mage -c "SELECT name FROM cards WHERE name LIKE '\''Lightning%'\'' LIMIT 3;"'
```

---

## 📋 View Logs

### Production Logs
```bash
ssh hkdebiandocker@192.168.178.24 'cd gomage && docker compose logs -f'
```

### Just Backend
```bash
ssh hkdebiandocker@192.168.178.24 'cd gomage && docker compose logs -f mage-server'
```

### Just Database
```bash
ssh hkdebiandocker@192.168.178.24 'cd gomage && docker compose logs -f postgres'
```

### Update Logs (Local)
```bash
tail -f logs/card-update-*.log
```

---

## 🆘 Troubleshooting

### Restart Production Services
```bash
ssh hkdebiandocker@192.168.178.24 'cd gomage && docker compose restart'
```

### Check Container Health
```bash
ssh hkdebiandocker@192.168.178.24 'cd gomage && docker compose ps'
```

### Emergency Full Redeploy
```bash
./deploy.sh
```

### Rollback Cards
```bash
./update-cards-prod.sh --rollback
```

---

## 📚 Full Documentation

| File | Purpose |
|------|---------|
| `PRODUCTION_CARD_UPDATES.md` | Complete card update guide |
| `SCRYFALL_MIGRATION_GUIDE.md` | Scryfall data usage guide |
| `SCRYFALL_DFC_CARDS.md` | Double-faced card handling |
| `DATA_MIGRATION.md` | Technical migration details |
| `MIGRATION_COMPLETE.md` | Initial migration results |

---

## ⚙️ Configuration

### Environment Variables
```bash
export REMOTE_USER="hkdebiandocker"
export REMOTE_HOST="192.168.178.24"
export REMOTE_PATH="gomage"
```

### Production URLs
- **Frontend**: http://mage.aronschueler.de
- **Backend API**: http://api.mage.aronschueler.de
- **Status**: http://192.168.178.24:17171/status
- **Health**: http://192.168.178.24:8080/health

---

## 🎯 Common Tasks

### Update Everything
```bash
# Full deployment
./deploy.sh

# Then update cards
./update-cards-prod.sh --download
```

### Update Just Cards
```bash
./update-cards-prod.sh --download
```

### Check if Update Needed
```bash
# Check Scryfall for latest data date
curl -s https://api.scryfall.com/bulk-data/all-cards | jq -r '.updated_at'

# Check your local file date
ls -lh all-cards-*.json
```

### Download Data Locally (No Update)
```bash
cd mage-server-go
./scripts/download_scryfall_bulk.sh
```

---

## 🔒 Safety Checklist

Before updating production:

- [ ] Verify local changes work: `make test`
- [ ] Check production is healthy: `curl http://192.168.178.24:17171/status`
- [ ] Have rollback plan ready: `./update-cards-prod.sh --rollback`
- [ ] Update during low-traffic time (recommended)
- [ ] Monitor logs during update
- [ ] Verify after update: Test card search and deck loading

---

## 📞 Quick Commands

```bash
# Update prod cards
./update-cards-prod.sh --download

# Preview changes
./update-cards-prod.sh --download --dry-run

# Rollback
./update-cards-prod.sh --rollback

# Deploy all
./deploy.sh

# Check status
curl http://192.168.178.24:17171/status

# View logs
ssh hkdebiandocker@192.168.178.24 'cd gomage && docker compose logs -f'
```
