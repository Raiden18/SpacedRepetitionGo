Telegram Bot that implements Spaced Repetition System to learn languages.

## How it fits together

- **Notion** holds the cards and decides when each is due (the `Show` formula).
- **updater** copies due cards from Notion into MySQL. It first sends any
  queued answers to Notion, and never overwrites cards whose answers are still
  queued.
- **notifier** posts or edits the "N cards to revise / memorize" messages.
- **image-replacer** downloads card images so Telegram can send them as files.
- **bot** is the webhook server (`:8080/webhook`) that handles button taps.
  Each answer is saved to the `notion_outbox` table first, then sent to Notion
  in the background with retries, so a Notion error can't lose it.
- **initializer** creates and migrates the tables. Every other binary does
  this too on start, so running it is optional.

## Database

All cards are in one `cards` table. `queue` is `revise` or `memorize`, and
`position` orders a box's cards within a queue. When a binary first starts
with an empty `cards` table, it copies cards from the old
`flash_cards_to_revise` and `flash_cards_to_memorize` tables, keeping their
order. The old tables, including `flash_cards_to_memorize_in_process`, are no
longer used. They're kept so an older release still works after a rollback;
once you're happy with the new version you can drop them:

```sql
DROP TABLE flash_cards_to_revise, flash_cards_to_memorize, flash_cards_to_memorize_in_process;
```

## Configuration (`secrets.yaml`, encrypted with git-crypt)

```yaml
notion:
  api_key: ...
  observed_databases: ...
  english_vocabulary_id: ...
  greek_vocabulary_id: ...
  greek_passive_vocabulary_id: ...
telegram:
  api_key: ...
  chat_id: ...          # only this chat can use the bot
  webhook_secret: ...   # optional, see below
```

### Webhook secret

With `telegram.webhook_secret` set, the bot rejects webhook calls that don't
carry Telegram's matching secret header. Set the same value on Telegram's side
(1–256 characters: `A-Z`, `a-z`, `0-9`, `_`, `-`):

```sh
curl "https://api.telegram.org/bot<TOKEN>/setWebhook" \
  -d url=https://<your-host>/webhook \
  -d secret_token=<the same secret>
```

Deploy the bot with the secret in `secrets.yaml` and call `setWebhook` at about
the same time; updates sent in between are rejected and redelivered by Telegram.

## Server requirements

- MySQL on the local socket, database `spaced_repetition`
- ImageMagick (`convert`) and librsvg (`rsvg-convert`) for image conversion
- `/root/repetition/images/` for downloaded images
- `/root/repetition/update-cache-notifier.sh`, run by the `/update` command

## Deploy and rollback

Pushes to `main` build every binary and upload them, with the decrypted
`secrets.yaml`, to `~/repetition/releases/<release>/` on the VPS. The files are
then copied into `~/repetition/go` (each one renamed into place in one step,
so a running job never finds a binary missing) and `ktor.service` restarts.
The 5 most recent releases are kept.

To roll back, copy an earlier release into `~/repetition/go`:

```sh
cd ~/repetition
ls -1t releases/
RELEASE=<previous release>
for file in updater notifier initializer bot image-replacer secrets.yaml; do
  cp -p "releases/$RELEASE/$file" "go/.$file.new" && mv -f "go/.$file.new" "go/$file"
done
sudo systemctl restart ktor.service
```
