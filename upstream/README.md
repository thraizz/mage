## Upstream XMage subrepository

This repo vendors the upstream XMage project (`magefree/mage`) as a git submodule at:

- `upstream/mage-upstream`

We use it to generate the upstream `cards.h2` database and migrate it into the Go server's
Postgres `cards` table.

### Initialize (first time)

```bash
git submodule update --init --recursive
```

### Update to latest upstream (optional)

```bash
git -C upstream/mage-upstream fetch origin
git -C upstream/mage-upstream checkout origin/master
git submodule update --init --recursive
```

### Update the Go server cards DB from upstream

```bash
./mage-server-go/scripts/update_cards_db_from_upstream.sh --docker
```

