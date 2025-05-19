#!/bin/bash

# Configuration
SERVER_PATH="servers/$USER_ID/$SERVER_ID"  # Pass $USER_ID and $SERVER_ID as env vars
BACKUP_DIR="$SERVER_PATH/backups"
MAX_BACKUPS=7  # Keep last 7 backups

# Create backup directory if missing
mkdir -p "$BACKUP_DIR"

# Create timestamped backup
TIMESTAMP=$(date +"%Y%m%d-%H%M%S")
tar -czf "$BACKUP_DIR/world-$TIMESTAMP.tar.gz" -C "$SERVER_PATH" world/

# Remove old backups (keep last $MAX_BACKUPS)
ls -t "$BACKUP_DIR"/world-*.tar.gz | tail -n +$(($MAX_BACKUPS + 1)) | xargs rm -f