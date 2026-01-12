#!/bin/bash
# Automated Weekly Card Update Script
# Suitable for cron jobs or scheduled tasks
#
# This script:
# - Downloads latest Scryfall data
# - Updates production cards
# - Sends notification on completion/failure
# - Logs everything for review
#
# Usage in crontab:
#   0 3 * * 0 /path/to/update-cards-weekly.sh >> /path/to/logs/card-updates.log 2>&1

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LOG_FILE="${SCRIPT_DIR}/logs/card-update-$(date +%Y%m%d_%H%M%S).log"

# Create logs directory if it doesn't exist
mkdir -p "${SCRIPT_DIR}/logs"

# Function to log with timestamp
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*" | tee -a "$LOG_FILE"
}

# Function to send notification (customize as needed)
notify() {
    local status="$1"
    local message="$2"
    
    # Log to file
    log "$status: $message"
    
    # Optional: Send email, Slack, Discord, etc.
    # Uncomment and configure as needed:
    
    # Email example:
    # echo "$message" | mail -s "MAGE Card Update: $status" admin@example.com
    
    # Slack example:
    # curl -X POST -H 'Content-type: application/json' \
    #   --data "{\"text\":\"MAGE Card Update: $status\\n$message\"}" \
    #   "$SLACK_WEBHOOK_URL"
    
    # Discord example:
    # curl -X POST -H 'Content-type: application/json' \
    #   --data "{\"content\":\"**MAGE Card Update: $status**\\n$message\"}" \
    #   "$DISCORD_WEBHOOK_URL"
}

log "=== Starting Weekly Card Update ==="
log "Script directory: $SCRIPT_DIR"
log "Log file: $LOG_FILE"

# Change to script directory
cd "$SCRIPT_DIR"

# Check if update script exists
if [[ ! -f "update-cards-prod.sh" ]]; then
    notify "ERROR" "update-cards-prod.sh not found in $SCRIPT_DIR"
    exit 1
fi

# Run the update
log "Running update-cards-prod.sh --download"
START_TIME=$(date +%s)

if ./update-cards-prod.sh --download >> "$LOG_FILE" 2>&1; then
    END_TIME=$(date +%s)
    DURATION=$((END_TIME - START_TIME))
    MINUTES=$((DURATION / 60))
    SECONDS=$((DURATION % 60))
    
    log "Update completed successfully in ${MINUTES}m ${SECONDS}s"
    notify "SUCCESS" "Production cards updated successfully in ${MINUTES}m ${SECONDS}s"
    
    # Optional: Keep only last 30 days of logs
    find "${SCRIPT_DIR}/logs" -name "card-update-*.log" -mtime +30 -delete
    
    exit 0
else
    END_TIME=$(date +%s)
    DURATION=$((END_TIME - START_TIME))
    MINUTES=$((DURATION / 60))
    SECONDS=$((DURATION % 60))
    
    log "Update FAILED after ${MINUTES}m ${SECONDS}s"
    notify "FAILURE" "Production card update failed after ${MINUTES}m ${SECONDS}s. Check logs: $LOG_FILE"
    
    exit 1
fi
