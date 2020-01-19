#!/bin/sh
set -euxo pipefail

export MAILME_APP_TEMPLATE_PATH=${MAILME_APP_TEMPLATE_PATH}
export MAILME_APP_TW_PROXY_HOST=${MAILME_APP_TW_PROXY_HOST}
export MAILME_APP_DSN=${MAILME_APP_DSN}
export MAILME_APP_PEM_FILE=${MAILME_APP_PEM_FILE}
export MAILME_APP_KEY_FILE=${MAILME_APP_KEY_FILE}

CRON_CONFIRM_SEND='*/3 * * * *'
echo "$CRON_CONFIRM_SEND /app/mailmeapp send-confirmation" >> /var/spool/cron/crontabs/root

CRON_SCHEDULE_CHECK='*/5 * * * *'
echo "$CRON_SCHEDULE_CHECK /app/mailmeapp check-new-subscriptions" >> /var/spool/cron/crontabs/root

CRON_SCHEDULE_PREPARE='5 18 * * *'
echo "$CRON_SCHEDULE_PREPARE /app/mailmeapp prepare-subscriptions" >> /var/spool/cron/crontabs/root


CRON_SCHEDULE_SEND='20 18 * * *'
echo "$CRON_SCHEDULE_SEND /app/mailmeapp send-subscriptions" >> /var/spool/cron/crontabs/root

crond -f
