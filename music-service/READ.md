### Сейчас фронт использует такие эндпоинты:

Auth

GET /auth/me
POST /auth/login
POST /auth/signup
POST /auth/logout
Tracks

GET /tracks
GET /tracks/:id
POST /tracks (multipart/form-data)
POST /tracks/:id/like
DELETE /tracks/:id/like
Artists

GET /artists
GET /artists/:id
Playlists

GET /playlists
GET /playlists/:id
POST /playlists
DELETE /playlists/:id
DELETE /playlists/:id/tracks/:trackId 

### Project Proposal
Music Streaming Service (Spotify-like)
Project Title

SoundWave – Music Streaming Service

1. Project Overview & Relevance

Music streaming services have become the primary way people listen to music today. Platforms such as Spotify and Apple Music offer large music libraries but are often overloaded with complex features, advertisements, and recommendation systems that may confuse users or distract them from the core purpose — listening to music.

The goal of this project is to design a simple and lightweight music streaming service that focuses on core functionality such as browsing music, managing playlists, and playing tracks. The project is intended as an educational software system that demonstrates backend architecture, data modeling, and modular design using the Go programming language.

This project is relevant because it reflects a real-world system while remaining manageable in scope for an academic environment.

2. Target Users

The target users of the system are:

University students

Young users who prefer minimalistic applications

Users who want a simple music service without complex recommendations

The system is not intended to replace commercial platforms but to demonstrate how such systems are designed and structured.

3. Competitor Analysis

Several popular music streaming services already exist:

Spotify

Pros: Large music library, advanced recommendations

Cons: Advertisements in free version, complex interface

Apple Music

Pros: High audio quality, integration with Apple ecosystem

Cons: Subscription required, limited cross-platform flexibility

YouTube Music

Pros: Free access, large content base

Cons: Less focus on structured playlists and albums

Our solution differs by focusing on:

Simplicity and clarity

Core music functionality only

Clean backend architecture suitable for educational purposes

4. Planned Features

At this stage, only planned features are defined. UI design and detailed workflows are intentionally excluded.

User Features

User registration and authentication

Browse music by artist and album

Search for tracks

Play music tracks

Create and manage playlists

Add and remove tracks from playlists

Mark tracks as favorites

Admin Features

Add, update, and delete artists

Add, update, and delete albums

Add, update, and delete tracks

5. Project Scope

The project will be implemented as a monolithic backend application written in Go.
It will focus on:

Clear separation of responsibilities (handlers, services, repositories)

Proper data modeling

Scalability for future extensions

Frontend implementation and advanced features such as recommendation algorithms are outside the scope of this milestone.

6. Conclusion

The proposed music streaming service is a realistic and relevant project that demonstrates understanding of software architecture, backend development, and system design. The project is suitable for incremental development and fulfills the requirements of the Advanced Programming 1 course milestone.


## План (MVP → “полноценный бэк” по твоему списку)

Ниже — **по шагам** и **список файлов**, которые нужно создать/изменить. Код пока не пишу, только архитектура и что где будет лежать (как просил).

---

# Шаг 0. Договоримся о ролях и доступе

### Роли
- `user` — слушает, лайкает, плейлисты, история, поиск, популярное.
- `artist` — всё как user + может **загружать треки/обложки**, управлять своими релизами.

### Правила доступа (вкратце)
- Любые действия “сохранить/лайк/плейлист/история” → **только залогиненный**.
- Загрузка музыки/обложек, создание треков/альбомов → **только role=artist**.

---

# Шаг 1. Миграции и схема БД (Postgres)

## Что добавить в БД
Тебе нужны таблицы:

1) `users` (уже есть, но роль должна быть `user|artist`)
2) `artists` (профиль артиста, привязка к user или отдельная сущность)
3) `genres`
4) `tracks` (метаданные + пути к локальным файлам)
5) `albums` (обложка локально)
6) `album_tracks` (если нужно упорядочивание треков в альбоме)
7) `playlists`
8) `playlist_tracks`
9) `likes` (лайки треков)
10) `listens` / `play_history` (история прослушиваний, статистика)
11) индексы для поиска и популярного (ILIKE/tsvector — позже можно улучшить)

## Файлы миграции (рекомендую перейти на “SQL миграции по файлам”)
Создать папку:
- `music-service/migrations/`

Файлы (примерный список):
- `music-service/migrations/0001_init_users.sql`
- `music-service/migrations/0002_artists.sql`
- `music-service/migrations/0003_genres.sql`
- `music-service/migrations/0004_albums.sql`
- `music-service/migrations/0005_tracks.sql`
- `music-service/migrations/0006_album_tracks.sql`
- `music-service/migrations/0007_playlists.sql`
- `music-service/migrations/0008_playlist_tracks.sql`
- `music-service/migrations/0009_likes.sql`
- `music-service/migrations/0010_listens_history.sql`
- `music-service/migrations/0011_indexes_search_popular.sql`

И обновить мигратор (см. шаг 2).

---

# Шаг 2. Новый нормальный мигратор (вместо “одной строки SQL”)

Создать/изменить:

## Новые файлы
- `music-service/repositories/migrations.go`  
  (читает файлы из `migrations/` и применяет по порядку, хранит применённые в таблице `schema_migrations`)

- `music-service/cmd/migrate/main.go`  
  (у тебя есть, но он будет вызывать новый мигратор)

## Изменить/удалить
- `music-service/router/migrate.go`  
  Лучше миграции держать не в `router`, а в `cmd/migrate` + `repositories`.

---

# Шаг 3. Локальные аудио-файлы и обложки

## Структура хранения на диске (пример)
- `music-service/storage/`
    - `covers/` (обложки)
    - `audio/` (mp3/flac/etc)
    - `tmp/` (временные)

## В config добавить
- `StorageAudioDir`
- `StorageCoversDir`
- `PublicBaseURL` (если нужно строить публичные ссылки)

### Файлы
- `music-service/services/storage_service.go`  
  (сохраняет файл на диск, генерирует уникальное имя, возвращает путь/URL)
- `music-service/handlers/upload_handler.go`  
  (эндпоинты загрузки: multipart/form-data)
- `music-service/router/static.go` (или внутри setup)  
  (раздача статических файлов через Gin: `/static/audio/...`, `/static/covers/...`)

---

# Шаг 4. Авторизация по ролям (user/artist)

## Что нужно
- Middleware `RequireAuth` (у тебя уже `CheckAuth`)
- Middleware `RequireRole("artist")`

### Файлы
- `music-service/router/role_middleware.go`  
  (`RequireRole(role string)` читает `currentUser` и проверяет `user.Role`)
- (опционально) `music-service/services/auth_service.go`  
  (JWT, генерация/проверка токенов — позже можно туда вынести)

---

# Шаг 5. История прослушиваний + популярное

## Логика
- endpoint “пользователь начал/закончил слушать трек” → записать событие в `listens`
- “популярное” = top по количеству прослушиваний за период/всего

### Файлы
**Models**
- `music-service/models/listen.go`

**Repositories**
- `music-service/repositories/listen_repository.go`
- (обновить) `music-service/repositories/interface.go` (добавить `ListenRepository`)

**Services**
- `music-service/services/listen_service.go`
    - `RecordListen(userID, trackID)`
    - `TopTracks(limit, period)`
    - `UserHistory(userID, limit, offset)`

**Handlers**
- `music-service/handlers/listen_handler.go`
    - `POST /tracks/:id/listen` (auth)
    - `GET /users/me/history` (auth)
    - `GET /tracks/top` (public)

---

# Шаг 6. Лайки (favorites)

## Логика
- `POST /tracks/:id/like` (auth)
- `DELETE /tracks/:id/like` (auth)
- `GET /users/me/likes` (auth)

### Файлы
**Models**
- `music-service/models/like.go`

**Repositories**
- `music-service/repositories/like_repository.go`
- `music-service/repositories/interface.go` (добавить `LikeRepository`)

**Services**
- `music-service/services/like_service.go`

**Handlers**
- `music-service/handlers/like_handler.go`

---

# Шаг 7. Плейлисты (у тебя уже почти есть)
Довести до полноценного:
- приватность (public/private) — опционально
- получить плейлист с треками
- сортировка внутри плейлиста

### Файлы (добавить/расширить)
- `music-service/repositories/playlist_repository.go` (добавить метод `ListTracks(playlistID)` и/или `GetWithTracks`)
- `music-service/services/playlist_service.go` (методы чтения с треками)
- `music-service/handlers/playlist_handler.go` (эндпоинт `GET /playlists/:id/tracks`)

---

# Шаг 8. Поиск (название трека, имя артиста, жанр, альбом)

## Логика
- `GET /search?q=...&type=all|tracks|artists|albums`
- `GET /tracks/search?q=...&genre=...&artist=...` (фильтры)

### Файлы
**Repositories**
- `music-service/repositories/search_repository.go`  
  или расширить `music_repository.go` методами:
    - `SearchTracks(query, genreID/genreName, artistName)`
    - `SearchArtists(query)`
    - `SearchAlbums(query)`
- `music-service/repositories/interface.go` (добавить `SearchRepository` или расширить `MusicRepository`)

**Services**
- `music-service/services/search_service.go`

**Handlers**
- `music-service/handlers/search_handler.go`

---

# Шаг 9. Эндпоинты для артиста (загрузка трека/обложки)

### Пример набора API
- `POST /artist/tracks` (role=artist, multipart: audio + metadata)
- `POST /artist/albums` (role=artist)
- `POST /artist/albums/:id/cover` (role=artist, multipart)
- `GET /artists/:id/tracks` (public)
- `GET /tracks/:id` (public)

### Файлы
**Handlers**
- `music-service/handlers/artist_handler.go`
- `music-service/handlers/track_handler.go`
- `music-service/handlers/album_handler.go`

**Services**
- `music-service/services/artist_service.go`
- `music-service/services/track_service.go`
- `music-service/services/album_service.go`
- `music-service/services/storage_service.go` (см. шаг 3)

**Repositories**
- `music-service/repositories/artist_repository.go` (если хочешь отдельный репозиторий)
- расширить `music_service`/`music_repository` или выделить:
    - `album_repository.go`
    - `track_repository.go`

---

# Шаг 10. Оставить место под OAuth2.0

Ты просил “создай файлы для OAuth2.0 если возможно” — да, можно сделать **заглушки/каркас**, без конкретного провайдера.

### Файлы (скелет)
- `music-service/services/oauth_service.go`  
  Интерфейс/каркас:
    - `GetAuthURL(provider)`
    - `HandleCallback(provider, code, state)` → получает профиль → привязывает к user

- `music-service/handlers/oauth_handler.go`
    - `GET /oauth/:provider/login`
    - `GET /oauth/:provider/callback`

- `music-service/models/oauth_account.go`  
  (связь user ↔ provider account)

- `music-service/repositories/oauth_repository.go`  
  (сохранить/найти привязку)

- `music-service/repositories/interface.go`  
  (добавить `OAuthRepository`)

- `music-service/router/oauth_routes.go`  
  (подключение oauth-роутов)

И миграция:
- `music-service/migrations/0012_oauth_accounts.sql`  
  Таблица типа `oauth_accounts(user_id, provider, provider_user_id, access_token?, refresh_token?)` (токены можно не хранить или хранить шифрованно — это отдельная тема).

---

# Итог: список файлов, которые нужно СОЗДАТЬ (сводно)

## migrations/
- `music-service/migrations/0001_init_users.sql`
- `music-service/migrations/0002_artists.sql`
- `music-service/migrations/0003_genres.sql`
- `music-service/migrations/0004_albums.sql`
- `music-service/migrations/0005_tracks.sql`
- `music-service/migrations/0006_album_tracks.sql`
- `music-service/migrations/0007_playlists.sql`
- `music-service/migrations/0008_playlist_tracks.sql`
- `music-service/migrations/0009_likes.sql`
- `music-service/migrations/0010_listens_history.sql`
- `music-service/migrations/0011_indexes_search_popular.sql`
- `music-service/migrations/0012_oauth_accounts.sql`

## models/
- `music-service/models/genre.go`
- `music-service/models/listen.go`
- `music-service/models/like.go`
- `music-service/models/oauth_account.go`

## repositories/
- `music-service/repositories/migrations.go`
- `music-service/repositories/genre_repository.go`
- `music-service/repositories/track_repository.go` *(или расширить music_repository.go)*
- `music-service/repositories/album_repository.go` *(или расширить)*
- `music-service/repositories/artist_repository.go` *(если выделять)*
- `music-service/repositories/listen_repository.go`
- `music-service/repositories/like_repository.go`
- `music-service/repositories/search_repository.go`
- `music-service/repositories/oauth_repository.go`

## services/
- `music-service/services/storage_service.go`
- `music-service/services/track_service.go`
- `music-service/services/album_service.go`
- `music-service/services/artist_service.go`
- `music-service/services/listen_service.go`
- `music-service/services/like_service.go`
- `music-service/services/search_service.go`
- `music-service/services/oauth_service.go`
- *(опционально)* `music-service/services/errors.go` (общие ошибки типа `ErrNotFound/ErrForbidden/ErrConflict`)

## handlers/
- `music-service/handlers/track_handler.go`
- `music-service/handlers/album_handler.go`
- `music-service/handlers/artist_handler.go`
- `music-service/handlers/listen_handler.go`
- `music-service/handlers/like_handler.go`
- `music-service/handlers/search_handler.go`
- `music-service/handlers/upload_handler.go`
- `music-service/handlers/oauth_handler.go`

## router/
- `music-service/router/role_middleware.go`
- `music-service/router/oauth_routes.go`
- *(опционально)* `music-service/router/static.go` (раздача `/static`)

---

## Чтобы не расползлось: предлагаю порядок внедрения
1) **Миграции** (users + tracks + playlists + likes + listens)
2) **storage_service + upload endpoints** (локальная загрузка)
3) **role middleware** (artist-only загрузка)
4) **история + популярное**
5) **поиск по track/artist/genre**
6) **OAuth каркас**

---

Если ок — скажи, с чего начинаем:  
**(1) миграции и схема БД** или **(2) загрузка аудио/обложек локально**?





--------------------------------------------------------------------1
Ниже — короткая “дорожная карта”, что **ещё осталось**, **сколько шагов**, и **что именно писать/куда** (по файлам). Это чтобы ты мог продолжать даже без меня.

---

## Что уже закрыто (по сути готово)
- **DB миграции** (users/artists/genres/albums/tracks/playlists/likes/listens + индексы).
- **JWT auth** (signup/login + middleware).
- **Роли** `RequireRole("artist")`.
- **Локальное хранилище** (storage + upload audio/covers + static).
- **Плейлисты** (CRUD + add/remove tracks).
- **Лайки** (like/unlike + мои лайки).
- **История** (listens + play_count++ + топ треков).
- **Поиск базовый** (ILIKE через pg_trgm индексы).
- **Профиль артиста** (создание профиля).

---

## Что осталось до “полного MVP под сайт музыки”
Я бы оценил, что осталось **6–8 шагов** (в зависимости от того, насколько “по-взрослому” делать).

### Шаг 1 (обязательный): Нормальные эндпоинты каталога (публичные списки)
**Цель:** фронту нужно “главная/каталог”.

**Нужно сделать:**
- `GET /tracks` (пагинация, сортировка: new/popular)
- `GET /artists` (у тебя уже есть, но можно нормализовать под новую схему)
- `GET /albums` или хотя бы `GET /artists/:id/albums`
- `GET /genres`
- `GET /tracks/:id`

**Куда писать:**
- `repositories/track_repository.go`: `List(...)`, `ListByArtistID(...)`, `ListByAlbumID(...)`
- `services/track_service.go`: обёртки + валидация
- `handlers/track_handler.go`: методы `List`, `ListByArtist`, `ListByAlbum`
- `router/router_setup.go`: подключить роуты

---

### Шаг 2 (обязательный): Популярное “по периоду” + топы
Сейчас топ у тебя по `tracks.play_count` “за всё время”. Обычно нужно:
- `top/week`, `top/month`, `top/all`

**Куда писать:**
- `repositories/listen_repository.go`: запросы агрегации по `listens` за период
- `services/listen_service.go`: `TopTracks(period)`
- `handlers/listen_handler.go` или `track_handler.go`: `GET /tracks/top?period=week`

---

### Шаг 3 (обязательный): Возврат плейлиста **с треками**
Сейчас плейлист есть, но фронту надо показать содержимое.

**Куда писать:**
- `repositories/playlist_repository.go`: `ListTracks(playlistID)` (join playlist_tracks + tracks)
- `services/playlist_service.go`: `GetWithTracks(...)`
- `handlers/playlist_handler.go`: `GET /playlists/:id/tracks` или `GET /playlists/:id` возвращает вместе

---

### Шаг 4 (важный): Поиск “всё-в-одном” (track + artist + genre) и фильтры
Сейчас поиск раздельный. Обычно удобно:
- `GET /search?q=...` возвращает одновременно tracks+artists+albums
- фильтры: `genre`, `artist_id`

**Куда писать:**
- `repositories/search_repository.go`: join’ы на `artists/genres`, фильтры
- `services/search_service.go`
- `handlers/search_handler.go`

---

### Шаг 5 (важный): Artist контент менеджмент
Артист должен:
- создать альбом
- загрузить cover
- загрузить трек
- привязать трек к альбому/жанру
- видеть “мои треки/альбомы” (owner view)

**Куда писать:**
- `handlers/artist_handler.go`: `GET /artist/me/tracks`, `GET /artist/me/albums`
- `repositories/track_repository.go`: `ListByArtistID`
- `repositories/album_repository.go`: `ListByArtistID`

---

### Шаг 6 (важный): Валидация и безопасность
Чтобы “не развалилось на проде”:
- rate limit на login
- ограничение размера upload (Gin)
- проверка типов файлов (у тебя уже по расширениям — норм)
- нормальный формат ошибок

**Куда писать:**
- `router/middleware_rate_limit.go` (простая реализация)
- `handlers/*`: единый формат ошибок (можно helper)
- `services/storage_service.go`: лимиты

---

### Шаг 7 (опционально, но очень желательно): Refresh tokens / sessions
Чтобы не перелогиниваться часто.
**Куда:** `services/auth_service.go`, `repositories/session_repository.go`, миграция `sessions`.

---

### Шаг 8 (последний, как ты сказал): OAuth2 skeleton
Оставить место и сделать заглушки:
- `GET /oauth/:provider/login`
- `GET /oauth/:provider/callback`

**Куда:** `handlers/oauth_handler.go`, `services/oauth_service.go`, `repositories/oauth_repository.go`, миграция `oauth_accounts`.

---

## Сколько шагов “точно” до твоего текущего ТЗ
Твоё ТЗ: локальные файлы, роли, история, лайки, плейлисты, популярное, поиск, артист загружает, юзер сохраняет после логина, место под OAuth2.

Если строго по пунктам, то **осталось примерно 4 шага**:
1) Публичные списки и нормальные каталожные endpoints (tracks/artists/albums/genres) — сейчас частично.
2) Плейлист “с треками” (полноценный).
3) Популярное “по периодам” (не только общий счётчик).
4) Поиск с join’ами и фильтрами (более “музыкальный”).

OAuth2 — потом.

---

## Что делать прямо сейчас (следующий шаг, конкретно)
Я бы выбрал **Шаг 1: Track list endpoints**, потому что это основа для фронта.

**Тебе нужно написать:**

### `repositories/track_repository.go`
- `List(ctx, limit, offset, sort)`
- `ListByArtistID(ctx, artistID, limit, offset)`
- `ListByAlbumID(ctx, albumID, limit, offset)`

### `services/track_service.go`
- методы-обёртки + дефолтные limit/offset + валидировать sort

### `handlers/track_handler.go`
- `List`
- `ListByArtist`
- `ListByAlbum`

### `router/router_setup.go`
- добавить:
  - `GET /tracks`
  - `GET /artists/:id/tracks`
  - `GET /albums/:id/tracks`

---

Если ты скажешь: **“делаем шаг 1”**, я в следующем ответе дам готовые куски кода **с названиями файлов**, как ты просил, без лишней болтовни.