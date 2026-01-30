#!/bin/bash
# Telegram notification script
# Usage: ./notify.sh "Your message"

MESSAGE="$1"

# Configure these variables for Telegram notifications
# TELEGRAM_BOT_TOKEN=""
# TELEGRAM_CHAT_ID=""

if [ -n "$TELEGRAM_BOT_TOKEN" ] && [ -n "$TELEGRAM_CHAT_ID" ]; then
    curl -s -X POST "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendMessage" \
        -d chat_id="${TELEGRAM_CHAT_ID}" \
        -d text="${MESSAGE}" \
        -d parse_mode="HTML"
else
    echo "[notify] $MESSAGE"
fi
