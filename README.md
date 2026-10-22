# SpacedRepetitionGo

A personal Telegram bot for learning vocabulary (English and Greek) with
spaced repetition. You keep your words in Notion; the bot sends you the cards
that are due in a Telegram chat, and your answers go back to Notion, which
decides when you see each card again.

## How it works

### In Notion

Notion is the single source of truth: it holds the cards and the review
schedule.

- **Boxes.** Each vocabulary list (for example *English Vocabulary* or *Greek
  Vocabulary*) is a Notion database, called a box. A separate "observed
  databases" table lists the boxes the bot should read. Each of its rows has a
  `Name`, the box's database `Id`, and an `Observable` checkbox that turns the
  box on or off.
- **Cards.** Each page in a box is a card:
  - `Name`: the word or phrase
  - `Example`: an example sentence (optional)
  - `Explanation` and `Answers`: the meaning or translation (optional), shown
    hidden behind a spoiler
  - the page cover: an image for the card (optional)
  - `Know Level 1` … `Know Level 14`: checkboxes recording how well you know
    the card
  - `Show`: a formula that is true when the card is due. The intervals behind
    it (`Interval 1` … `Interval 14`) live in Notion too; the bot only reads
    `Show`.

### In Telegram

The bot works in one chat (`telegram.chat_id`) and ignores everyone else.

1. **Notifications.** It posts a message like *"You have 12 flashcards to
   revise 🧠"*, with one button per box showing how many of its cards are due.
   There is a second message for cards to memorize.
2. **Cards.** Tapping a box sends its first card: the word in bold, the
   example in italics, and the explanation hidden as a spoiler. Every card has
   the same button layout:

   | Row | Buttons |
   |---|---|
   | 1 | ⬅️ Previous · Next ➡️ |
   | 2 | To the beginning ⏮️ · To the end ⏭️ |
   | 3 | Card to revise: Forgot ❌ · Recalled ✅. Card to memorize: Memorized ✅ |
   | 4 | Dictionary link: Cambridge Dictionary for English, Forvo for Greek |
   | 5 | Finish 🏁 (closes the card) |

3. **`/update`** runs `update-cache-notifier.sh` on the server, which refreshes
   the cards and notifications right away.

### Two queues

- **To memorize:** new cards, where `Know Level 1`–`4` are all unchecked and
  `Show` is true.
- **To revise:** cards you've learned (`Know Level 1` checked) that are due
  (`Show` is true).

What each answer does:

| Button | Effect in Notion | Where the card goes |
|---|---|---|
| Recalled ✅ | Checks the next unchecked `Know Level`, so the next review is further away | Leaves the queue until Notion shows it again |
| Forgot ❌ | Unchecks every `Know Level` | Goes to the end of the box's memorize queue |
| Memorized ✅ | Checks only `Know Level 1` | Leaves the queue; Notion later shows it for revision |

For English cards, every answer also sets `Interval 14` to 560. This
gives old cards with all 13 levels checked a 14th interval.

## Architecture

Five Go programs (in `cmd/`) run on one VPS with a local MySQL database that
caches the due cards:

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

### Code layout

| Package | What it does |
|---|---|
| `cmd/*` | Entry points for the five programs |
| `jobs` | The updater, notifier, image and initializer jobs |
| `bot` | Commands and button handlers |
| `flashcard` | Cards: parsing from Notion, the `cards` table, Telegram messages and buttons |
| `box` | Boxes (Notion databases) and the observed-databases list |
| `notification` | The "N cards to revise / memorize" messages |
| `outbox` | Reliable write-back of answers to Notion |
| `notion`, `telegram` | Thin wrappers around the Notion and Telegram APIs |
| `image` | Downloading, converting and caching card images |
| `schema` | Creates and migrates all tables on start |
| `config` | Reads `secrets.yaml` |

## Database

All cards are in one `cards` table. `queue` is `revise` or `memorize`, and
`position` orders a box's cards within a queue. The first time a binary
starts after the upgrade, it copies cards from the old
`flash_cards_to_revise` and `flash_cards_to_memorize` tables, keeping their
order. This happens only once, when the `cards` table is first created. The
old tables, including `flash_cards_to_memorize_in_process`, are no
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
