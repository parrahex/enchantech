# enchan.tech

My personal page. Go + Gin on the backend, and the sky behind the page follows
the real weather in my city.

The server picks a scene from the current forecast and puts it into the page as
a data attribute. CSS has a palette per scene, a small canvas script draws and
animates the sky. The template also gets the temperature, a short weather
description, the moon phase, the season and the local time.

## Scenes

Anything falling from the sky wins over the time of day:

| Weather | Scene |
| --- | --- |
| thunderstorm, WMO code 95+ | `thunder` |
| snow: 71, 73, 75, 77, 85, 86 | `snow` |
| rain, drizzle or showers: 51-67, 80-82 | `rain` |
| sun is down | `night` |
| overcast: code 3 or 85%+ cloud cover | `overcast` |
| before 12:00 | `morning` |
| 12:00-17:59 | `day` |
| 18:00 and later | `evening` |

Night uses the forecast's `is_day` flag, so it follows the real sunset and
sunrise. Morning, day and evening are just clock hours.

Seasons go by month at the location: winter from December, spring from March,
summer from June, autumn from October. South of the equator everything is
shifted by six months.

Cloud cover sets how many clouds there are, wind sets how fast they move. At
85% cover or in fog (codes 45, 48) the sun and moon are hidden, and fog adds
drifting bands of mist.

The forecast is refreshed every 15 minutes. If no location is set, the server
doesn't call anything and just shows the default `day` sky.

## Location

Weather comes from [Open-Meteo](https://open-meteo.com), no account or key
needed. Set all three:

```sh
APP_CITY_LATITUDE=48.85
APP_CITY_LONGITUDE=2.35
APP_CITY_TIMEZONE=Europe/Paris
```

The time zone is an IANA name. The tz database is compiled into the binary, so
the image doesn't need `tzdata`. `docker compose` reads these from `.env`.

## Front end

Templates, styles, fonts and images aren't in this repo. Without them the
server still starts and serves a plain text placeholder.

Keep yours in `./content` or point `APP_CONTENT_DIR` at it:

```
content/
├── templates/
│   └── home.html
└── assets/
    ├── site.css
    └── ...
```

`home.html` is a regular Go template and gets:

| Field | Type | Meaning |
| --- | --- | --- |
| `.Profile.Name` | string | display name |
| `.Profile.Username` | string | handle without `@` |
| `.Profile.Bio` | string | one-liner |
| `.Profile.Avatar` | string | path to the avatar in `assets/` |
| `.Profile.AvatarAlt` | string | alt text for the avatar |
| `.Profile.Links` | list | `.ID`, `.Label`, `.URL` |
| `.Scene` | string | scene name from the table above |
| `.Known` | bool | false if the forecast couldn't be loaded |
| `.Night` | bool | the sun is down |
| `.Fog` | bool | there is fog |
| `.Cover` | int | cloud cover, 0-100 |
| `.Wind` | float | wind at 10 m, km/h |
| `.Condition` | string | `light rain`, `overcast` and so on; empty if unknown |
| `.Temperature` | int | °C, rounded |
| `.Moon` | float | moon phase: 0 new, 0.5 full, 1 new again; mirrored in the southern hemisphere |
| `.Season` | string | `winter`, `spring`, `summer` or `autumn` |
| `.Clock` | bool | a location is set, so its local time is known |
| `.Offset` | int | UTC offset of the location in minutes |

Everything in `assets/` is served under `/assets/`. CSS can key off
`:root[data-scene='rain']` and so on. The canvas script is optional, a page
that ignores `.Scene` works too.

Name, bio and links live in `internal/app/profile.go`.

## Running

```sh
docker compose up --build
```

Without Docker:

```sh
make run     # starts the server with .env loaded
make sky     # compiles content/sky.ts into content/assets/sky.js
make check   # gofmt, go vet, go build
```

## Deploy

`content/` gets baked into the image, so build it on a machine that has it:

```sh
docker buildx build --platform linux/amd64 -t enchantech:latest --load .
docker save enchantech:latest | gzip | ssh user@server 'gunzip | docker load'
ssh user@server 'cd enchantech && docker compose up -d'
```

On the server the compose file uses `image: enchantech:latest` instead of
`build: .`, with `Caddyfile` and `.env` in the same folder. Caddy gets the
certificates on its own. For Graviton build with `linux/arm64`.

## Configuration

| Variable | Default | Purpose |
| --- | --- | --- |
| `APP_ADDRESS` | `:8080` | listen address |
| `APP_SHUTDOWN_TIMEOUT` | `5s` | graceful shutdown limit |
| `APP_CONTENT_DIR` | `content` | where the front end lives |
| `APP_CITY_LATITUDE` | unset | forecast latitude |
| `APP_CITY_LONGITUDE` | unset | forecast longitude |
| `APP_CITY_TIMEZONE` | unset | IANA time zone for the location |
