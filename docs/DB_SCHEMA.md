# Orthocal Database Schema

Orthocal supports database schema version 4.

The CLI is read-only for normal calendar commands. The `update` command validates a replacement database before installing it.

## Required Tables

| Table                | Purpose                         |
| -------------------- | ------------------------------- |
| `app_metadata`       | Generator and schema metadata   |
| `calendar_days`      | One row per Gregorian day       |
| `saints`             | Saints and commemorations       |
| `service_ranks`      | Service rank lookup data        |
| `scripture_readings` | Daily scripture readings        |
| `hymns`              | Hymn text and hymn metadata     |

## Event Tables

Schema version 4 adds calculated calendar-event support.

| Table                 | Purpose                                   |
| --------------------- | ----------------------------------------- |
| `calendar_events`     | One row per calculated event or range     |
| `calendar_day_events` | Links active calculated events to a date  |

Older databases that lack event tables can still be read. Event views and event sections will be empty.

## app_metadata

Important keys:

| Key                        | Meaning                            |
| -------------------------- | ---------------------------------- |
| `schema_version`           | Expected value is `4`              |
| `calendar_year`            | Calendar year contained in the DB  |
| `generated_at`             | Generator timestamp                |
| `generator`                | Generator name                     |
| `includes_calendar_events` | `true` when event tables are built |
| `source_calendar_url`      | Calendar source URL                |
| `source_script_root_url`   | Script source root URL             |
| `source_scripts`           | Source scripts used by generator   |

## calendar_days

Important columns:

| Column                      | Meaning                                      |
| --------------------------- | -------------------------------------------- |
| `id`                        | Primary key                                  |
| `gregorian_date`            | Gregorian date in `YYYY-MM-DD` format        |
| `gregorian_weekday`         | Gregorian weekday name                       |
| `julian_date`               | Julian date in `YYYY-MM-DD` format           |
| `dataheader`                | Full source date header                      |
| `headerheader`              | Liturgical header                            |
| `fasting_rule`              | Source fasting rule text                     |
| `feasts`                    | Pipe-separated calculated feast titles       |
| `fasts`                     | Pipe-separated fasting season titles         |
| `remembrances`              | Pipe-separated remembrance titles            |
| `fast_free_periods`         | Pipe-separated fast-free period titles       |
| `fasting_level_code`        | Source-style fasting code                    |
| `fasting_level_name`        | Human fasting level name                     |
| `fasting_level_description` | Longer fasting level description             |
| `is_holiday`                | `1` for calculated major holiday or feast    |
| `is_lent_day`               | `1` during calculated lenten or fast periods |

## saints

Important columns:

| Column              | Meaning                         |
| ------------------- | ------------------------------- |
| `day_id`            | References `calendar_days.id`   |
| `saint_order`       | Source order for the day        |
| `name`              | Saint or commemoration text     |
| `icon_file`         | Source icon filename            |
| `is_primary`        | `1` for primary saints          |
| `is_western`        | `1` for western saints          |
| `service_rank_code` | Service rank code               |
| `service_rank_name` | Service rank name               |

Orthocal sorts day saint output as primary saints, western saints, ranked saints from 6 down to 0, then ordinary `o` ranks.

## scripture_readings

Important columns:

| Column            | Meaning                       |
| ----------------- | ----------------------------- |
| `day_id`          | References `calendar_days.id` |
| `reading_order`   | Reading order for the day     |
| `verse_reference` | Bible reference               |
| `description`     | Optional service description  |
| `reading_url`     | Source reading URL            |
| `display_text`    | Source display text           |

## hymns

Important columns:

| Column          | Meaning                       |
| --------------- | ----------------------------- |
| `day_id`        | References `calendar_days.id` |
| `section_order` | Section order                 |
| `hymn_order`    | Hymn order within section     |
| `hymn_type`     | Hymn type                     |
| `tone`          | Tone text                     |
| `title`         | Hymn title                    |
| `text`          | Hymn text                     |

## calendar_events

Important columns:

| Column            | Meaning                              |
| ----------------- | ------------------------------------ |
| `id`              | Primary key                          |
| `event_key`       | Unique event key                     |
| `category`        | Event category                       |
| `title`           | Event title                          |
| `start_date`      | Start date in `YYYY-MM-DD` format    |
| `end_date`        | End date in `YYYY-MM-DD` format      |
| `is_range`        | `1` when the event spans a range     |
| `source_script`   | Source script that produced the row  |
| `source_root_url` | Source root URL                      |
| `notes`           | Generator notes                      |
| `sort_order`      | Sort order                           |

Known categories include `fixed_feast`, `movable_feast`, `fasting_season`, `fast_free_week`, and `remembrance`.

## calendar_day_events

Important columns:

| Column       | Meaning                       |
| ------------ | ----------------------------- |
| `day_id`     | References `calendar_days.id` |
| `event_id`   | References `calendar_events.id` |
| `event_date` | Active date in `YYYY-MM-DD` format |
| `category`   | Event category                |
| `title`      | Event title                   |
| `sort_order` | Sort order                    |

## Backward Compatibility

`orthocal info` reports the detected schema version and compatibility status.

Older databases can still serve basic date, saint, reading, hymn, and search commands. Missing schema version 4 event columns or event tables produce empty event data rather than crashing date views.

Newer databases show a warning in `orthocal info`. Read-only commands continue when the columns they need are still compatible. `orthocal update` rejects newer schema versions unless `--force` is passed.
