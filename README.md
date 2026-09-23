# enchan.tech

Personal profile page built with Go and Gin.

The server reads live weather for a configured location and picks a sky scene
from it. The scene is decided on the server and passed to the page as a data
attribute; CSS holds the palette for each scene, and a small canvas script
animates it.

## Sky scenes

Real weather comes first: if it is raining, the sky rains, whatever the hour. The
time of day is used only when nothing is falling.

| Condition | Scene |
| --- | --- |
| Thunderstorm — WMO code 95 or above | `thunder` |
| Snow — 71, 73, 75, 77, 85, 86 | `snow` |
| Drizzle, rain, or showers — 51-67, 80-82 | `rain` |
| Dark outside | `night` |
| Light, overcast — WMO code 3, or cloud cover of 85% or more | `overcast` |
| Light, before 12:00 local time | `morning` |
| Light, 12:00 to 17:59 | `day` |
| Light, 18:00 onward | `evening` |

Night follows real sunrise and sunset: it comes from the forecast's own day and
night flag for the configured coordinates, not from the clock. The last three
are clock hours, because the forecast reports whether the sun is up but not what
part of the day it is.

Cloud cover sets how many clouds are drawn, and wind speed sets how fast they
drift. The sun or moon is hidden when cover reaches 85%, or in fog (codes 45 and
48), which also lays flat bands across the sky.

With no location configured the server makes no outbound request at all and
serves the calm default `day` sky.

## Configuring a location

Weather comes from [Open-Meteo](https://open-meteo.com), which needs no account
or API key. Set all three variables together:

```sh
APP_CITY_LATITUDE=48.85
APP_CITY_LONGITUDE=2.35
APP_CITY_TIMEZONE=Europe/Paris
```

The time zone is an IANA name; its database is compiled into the binary, so no
system `tzdata` is required. `docker compose` reads these from a `.env` file in
the project directory.

## Front end

Templates, styles, fonts, and images are not part of this repository. The
backend runs without them: a clean clone builds, starts, and serves a plain
placeholder instead of failing.

Bring your own by providing a directory like this:

```
content/
├── templates/
│   └── home.html
└── assets/
    ├── site.css
    └── ...
```

`home.html` is an ordinary Go template. It receives:

| Field | Type | Meaning |
| --- | --- | --- |
| `.Profile.Name` | string | display name |
| `.Profile.Username` | string | handle, rendered without `@` |
| `.Profile.Bio` | string | one-line description |
| `.Profile.Avatar` | string | image path, served from `assets/` |
| `.Profile.AvatarAlt` | string | alt text for the avatar |
| `.Profile.Links` | list | each with `.ID`, `.Label`, `.URL` |
| `.Scene` | string | one of the scenes above, lower case |
| `.Known` | bool | false when no forecast was available |
| `.Night` | bool | the sun is down at the configured location |
| `.Fog` | bool | fog is present |
| `.Cover` | int | cloud cover, 0–100 |
| `.Wind` | float | wind speed at 10 m |
| `.Condition` | string | short weather description, such as `light rain`; empty when unknown |
| `.Temperature` | int | air temperature at 2 m, rounded, °C |
| `.Moon` | float | moon phase, 0 new → 0.5 full → 1 new, mirrored south of the equator |

Everything under `assets/` is served at `/assets/`. A scene-aware stylesheet can
key off `:root[data-scene='rain']` and friends; a canvas script is optional, and
a page that ignores `.Scene` entirely works fine.

Edit `internal/app/profile.go` to change the name, bio, and links.

## Running

```sh
docker compose up --build
```

## Configuration

| Variable | Default | Purpose |
| --- | --- | --- |
| `APP_ADDRESS` | `:8080` | listen address |
| `APP_SHUTDOWN_TIMEOUT` | `5s` | graceful shutdown limit |
| `APP_CONTENT_DIR` | `content` | where the front end lives |
| `APP_CITY_LATITUDE` | unset | forecast latitude |
| `APP_CITY_LONGITUDE` | unset | forecast longitude |
| `APP_CITY_TIMEZONE` | unset | IANA time zone for the location |
