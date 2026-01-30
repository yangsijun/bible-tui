# Bible TUI/CLI Application — Work Plan

## TL;DR

> **Quick Summary**: 대한성서공회(bskorea.or.kr)에서 개역개정 성경을 크롤링하여 SQLite에 저장하고, Go + Bubbletea 기반의 크로스플랫폼 TUI/CLI 성경 읽기 프로그램을 구축한다. 다중 역본 확장 가능한 아키텍처로 설계하며, 테마/폰트 설정을 포함한다.
>
> **Deliverables**:
> - `bible crawl` — 성경 데이터 크롤러 (사용자 실행, 바이너리에 데이터 미포함)
> - `bible read <책> <장>` — CLI로 성경 본문 출력
> - `bible search <검색어>` — CLI 검색
> - `bible random` — 랜덤 구절 출력
> - `bible tui` — 인터랙티브 TUI 모드 (읽기/검색/책갈피/하이라이트/읽기계획/테마/폰트설정)
> - 크로스플랫폼 바이너리 (Mac/Linux/Windows)
>
> **Estimated Effort**: Large (4 phases)
> **Parallel Execution**: YES — 2-3 waves per phase
> **Critical Path**: Project Setup → DB Layer → Crawler → CLI Read → TUI Shell → Search → Bookmarks → Settings → Reading Plans

---

## Context

### Original Request
Mac/Linux/Windows용 TUI와 CLI를 제공하는 성경 프로그램. 성경 데이터는 대한성서공회 성경읽기 홈페이지에서 개역개정 크롤링. 다른 역본 추가 가능성 고려. 테마/폰트 설정 기능 포함.

### Interview Summary
**Key Discussions**:
- **Language**: Go + Bubbletea — 한글 렌더링 안정성, 단일 바이너리 배포
- **Data**: 사전 크롤링 + SQLite (오프라인 사용, FTS5 검색)
- **TUI Features**: 읽기, 검색, 책갈피/하이라이트, 읽기계획(통독+M'Cheyne+사용자정의), 테마, 폰트
- **CLI**: read, search, tui, random
- **Test**: TDD (Go 내장 testing)
- **Build**: 단계별 (4 phases)
- **아키텍처**: 다중 역본 대비 (DB에 version 컬럼, 크롤러에 version 파라미터)

**Research Findings**:
- URL 패턴: `korbibReadpage.php?version=GAE&book={code}&chap={num}`
- 66권 전체 영문코드+한글명+장 수 확보 (`bible.list.js`의 `szHANBook` 배열)
- HTML 구조: `<div id="tdBible1" class="bible_read">` > `<span class="number">` + 텍스트
- 소제목: `<font class="smallTitle">`, 각주: class `D2` div
- `modernc.org/sqlite`에서 커스텀 FTS5 토크나이저 불가 → LIKE 검색 + FTS5 보조 전략
- 역본 코드: GAE(개역개정), HAN(개역한글), SAE(표준새번역), SAENEW(새번역), COG(공동번역), COGNEW(공동번역개정판), CEV

### Metis Review
**Identified Gaps** (addressed):
- **저작권**: 크롤러를 별도 `bible crawl` 명령으로 분리, 바이너리에 데이터 미포함 → 사용자가 직접 실행
- **한국어 검색**: 커스텀 FTS5 토크나이저 불가 → LIKE 검색 기본, FTS5 보조
- **크롤링 안정성**: 체크포인트(crawl_status 테이블)로 중단/재개 지원
- **인코딩**: UTF-8 확인됨, 방어적 charset 감지 추가
- **Windows**: Windows Terminal 최소 요구, cmd.exe 레거시 미지원

---

## Work Objectives

### Core Objective
대한성서공회 개역개정 성경을 크롤링하여 오프라인 SQLite DB에 저장하고, Go+Bubbletea 기반 크로스플랫폼 TUI/CLI 성경 읽기 앱을 구축한다. 다중 역본 확장 가능한 아키텍처, 테마/폰트 커스터마이징을 포함한다.

### Concrete Deliverables
- Go 모듈 프로젝트 (`bible-tui`)
- SQLite 데이터베이스 스키마 (다중 역본 대비)
- 웹 크롤러 (`bible crawl` 명령)
- CLI 명령어 (read, search, random, tui)
- Bubbletea TUI (읽기, 검색, 책갈피, 하이라이트, 읽기계획, 설정)
- 테마/폰트 설정 시스템
- 크로스플랫폼 빌드 (darwin/linux/windows)

### Definition of Done
- [ ] `go test ./... -race` — ALL PASS
- [ ] `go vet ./...` — clean
- [ ] 3개 플랫폼 빌드 성공: `GOOS=darwin/linux/windows go build`
- [ ] `bible crawl` → 66권 전체 크롤링 완료, 무결성 검증 통과
- [ ] `bible read 창세기 1` → 창세기 1장 전체 출력
- [ ] `bible search 사랑` → 관련 구절 검색 결과 출력
- [ ] `bible random` → 랜덤 구절 1개 출력
- [ ] `bible tui` → TUI 모드 진입, 정상 동작

### Must Have
- 크로스플랫폼 단일 바이너리 (Mac/Linux/Windows)
- 오프라인 동작 (크롤링 후)
- 한글 텍스트 정상 렌더링 (TUI에서 CJK 폭 처리)
- 다중 역본 대비 DB 스키마 (version 컬럼)
- 크롤러 중단/재개 (체크포인트)
- 크롤링 후 무결성 검증 (66권, 장수 일치, 빈 절 없음)
- 테마 설정 (다크/라이트 + 커스텀)
- 폰트 크기 설정

### Must NOT Have (Guardrails)
- 바이너리에 성경 데이터 임베딩 금지 (저작권)
- `mattn/go-sqlite3` 사용 금지 (CGO → 크로스컴파일 불가)
- `len()` 으로 표시 너비 계산 금지 → `runewidth.StringWidth()` 사용
- 2개 이상 구현체 없이 인터페이스 생성 금지
- 캐싱 레이어 금지 (SQLite가 31K 행 충분히 처리)
- cobra에서 `Run` 사용 금지 → 항상 `RunE`
- 형태소 분석 약속 금지 (LIKE 검색이 기본)
- 오디오 재생, 웹/모바일 UI, 히브리어/그리스어 원문, 클라우드 동기화 금지
- 글로벌 state 금지 → DB는 constructor injection으로 전달

---

## Verification Strategy (MANDATORY)

### Test Decision
- **Infrastructure exists**: NO (새 프로젝트)
- **User wants tests**: TDD
- **Framework**: Go 내장 `testing` 패키지 + `github.com/charmbracelet/x/exp/teatest`

### TDD Workflow
각 TODO는 RED-GREEN-REFACTOR:
1. **RED**: 실패하는 테스트 작성
2. **GREEN**: 테스트를 통과하는 최소 코드 구현
3. **REFACTOR**: 그린 유지하며 정리

### QA Commands (Agent-Executable)
```bash
# 전체 테스트
go test ./... -v -count=1
go test ./... -race

# 빌드 검증 (3 플랫폼)
GOOS=darwin GOARCH=arm64 go build -o /dev/null .
GOOS=linux GOARCH=amd64 go build -o /dev/null .
GOOS=windows GOARCH=amd64 go build -o /dev/null .

# 정적 분석
go vet ./...

# CLI 동작 확인
go run . --help
go run . crawl --help
go run . read --help
go run . search --help
go run . random --help
go run . tui --help
```

---

## Execution Strategy

### Phase Overview

| Phase | Focus | Tasks |
|-------|-------|-------|
| Phase 1 | 기반: 프로젝트+DB+크롤러+기본CLI | Tasks 1-7 |
| Phase 2 | 검색+책갈피+하이라이트 | Tasks 8-11 |
| Phase 3 | TUI 핵심 | Tasks 12-16 |
| Phase 4 | 설정(테마/폰트)+읽기계획+마무리 | Tasks 17-22 |

### Parallel Execution Waves

```
=== Phase 1: Foundation ===
Wave 1 (Start Immediately):
├── Task 1: Project scaffold + Go module + Cobra root
└── (sequential only — everything depends on this)

Wave 2 (After Task 1):
├── Task 2: SQLite DB layer + schema (multi-version ready)
├── Task 3: Bible reference parser (한글 책이름 → 코드)
└── (parallel: 2 and 3 are independent)

Wave 3 (After Task 2):
├── Task 4: HTML parser (크롤링 대상 HTML → 구조화된 데이터)
└── (depends on DB schema for output format)

Wave 4 (After Tasks 2, 4):
├── Task 5: Crawler engine (HTTP + rate limiting + checkpoint)
└── (depends on DB + parser)

Wave 5 (After Tasks 2, 3):
├── Task 6: CLI read + random commands
└── Task 7: CLI crawl command (wires crawler engine)
    (6 and 7 parallel — both depend on DB + reference parser)

=== Phase 2: Search & Bookmarks ===
Wave 6 (After Phase 1):
├── Task 8: FTS5 + LIKE search layer
├── Task 9: Bookmark/Highlight data layer
└── (parallel: independent DB features)

Wave 7 (After Tasks 8, 9):
├── Task 10: CLI search command
└── Task 11: CLI bookmark commands
    (parallel)

=== Phase 3: TUI Core ===
Wave 8 (After Phase 2):
├── Task 12: TUI shell + navigation framework
└── (foundation for all TUI features)

Wave 9 (After Task 12):
├── Task 13: TUI reading view (book→chapter→verse)
├── Task 14: TUI search view
└── (parallel views)

Wave 10 (After Tasks 13, 14):
├── Task 15: TUI bookmark/highlight view
└── Task 16: TUI help/keybindings view
    (parallel)

=== Phase 4: Settings & Reading Plans ===
Wave 11 (After Phase 3):
├── Task 17: Settings data layer (테마/폰트 persistence)
├── Task 18: Reading plan data layer
└── (parallel: independent data layers)

Wave 12 (After Tasks 17, 18):
├── Task 19: Theme engine (다크/라이트/커스텀)
├── Task 20: Font size engine
└── (parallel)

Wave 13 (After Tasks 17-20):
├── Task 21: TUI settings view
├── Task 22: TUI reading plan view
└── (parallel views)

Critical Path: 1 → 2 → 4 → 5 → 7 → 8 → 10 → 12 → 13 → 19 → 21
```

### Dependency Matrix

| Task | Depends On | Blocks | Parallel With |
|------|-----------|--------|---------------|
| 1 | None | ALL | None |
| 2 | 1 | 4,5,6,7,8,9 | 3 |
| 3 | 1 | 6,7 | 2 |
| 4 | 2 | 5 | 3 |
| 5 | 2,4 | 7 | None |
| 6 | 2,3 | 12 | 7 |
| 7 | 5 | Phase 2 | 6 |
| 8 | Phase 1 | 10 | 9 |
| 9 | Phase 1 | 11 | 8 |
| 10 | 8 | 14 | 11 |
| 11 | 9 | 15 | 10 |
| 12 | Phase 2 | 13-16 | None |
| 13 | 12 | 19 | 14 |
| 14 | 12 | None | 13 |
| 15 | 13,14 | None | 16 |
| 16 | 12 | None | 15 |
| 17 | Phase 3 | 19,20,21 | 18 |
| 18 | Phase 3 | 22 | 17 |
| 19 | 17 | 21 | 20 |
| 20 | 17 | 21 | 19 |
| 21 | 19,20 | None | 22 |
| 22 | 18 | None | 21 |

---

## TODOs

---

### Phase 1: Foundation (프로젝트 기반 + DB + 크롤러 + 기본 CLI)

---

- [ ] 1. Project Scaffold & Cobra Root Command

  **What to do**:
  - `go mod init github.com/user/bible-tui` (실제 모듈 경로에 맞게)
  - 디렉토리 구조 생성:
    ```
    bible-tui/
    ├── main.go                  # cobra root 실행
    ├── cmd/
    │   └── root.go              # cobra root command
    ├── internal/
    │   ├── db/                  # SQLite 데이터 레이어
    │   ├── crawler/             # 웹 크롤러
    │   ├── parser/              # HTML 파서
    │   ├── bible/               # 성경 참조 파싱 (책이름→코드)
    │   ├── tui/                 # Bubbletea TUI
    │   │   ├── app.go           # root TUI model
    │   │   └── styles/          # lipgloss 스타일 + 테마
    │   └── config/              # 설정 관리 (테마/폰트/경로)
    ├── testdata/                # 테스트 fixtures (HTML 샘플 등)
    └── .gitattributes           # *.golden -text
    ```
  - `go get` 핵심 의존성: `bubbletea`, `lipgloss`, `bubbles`, `cobra`, `goquery`, `modernc.org/sqlite`, `go-runewidth`, `golang.org/x/time/rate`, `golang.org/x/net/html/charset`
  - `cmd/root.go`: cobra root command (`bible`)에 `PersistentPreRunE`로 DB 초기화 훅 설정
  - `main.go`: `cmd.Execute()` 호출
  - 테스트: `go build ./...` 성공, `go run . --help` 출력 확인
  - `.gitattributes` 생성: `*.golden -text`

  **Must NOT do**:
  - 아직 서브커맨드 구현 안 함 (root만)
  - 비즈니스 로직 작성 안 함

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: 프로젝트 초기 스캐폴딩, 파일 생성 + 의존성 설치 위주
  - **Skills**: [`git-master`]
    - `git-master`: 첫 커밋 생성 + .gitattributes 관리

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 1 (solo)
  - **Blocks**: Tasks 2, 3
  - **Blocked By**: None

  **References**:
  - **Pattern**: cobra root command 패턴 — `PersistentPreRunE`로 DB 초기화, `RunE`로 에러 전파
  - **Pattern**: `charmbracelet/mods` 프로젝트의 cobra+bubbletea 통합 패턴 참고
  - **Library**: `github.com/spf13/cobra` — subcommand routing
  - **Library**: `github.com/charmbracelet/bubbletea` — TUI framework
  - **Library**: `modernc.org/sqlite` — pure Go SQLite (CGO-free)
  - **Library**: `github.com/PuerkitoBio/goquery` — HTML parsing
  - **Library**: `github.com/mattn/go-runewidth` — CJK 문자 너비
  - **Library**: `golang.org/x/time/rate` — rate limiter
  - **Library**: `golang.org/x/net/html/charset` — 인코딩 감지

  **Acceptance Criteria**:
  - [ ] `go build ./...` → 성공 (exit 0)
  - [ ] `go run . --help` → "bible" 또는 프로그램 설명 포함 출력
  - [ ] `go vet ./...` → clean
  - [ ] 디렉토리 구조 확인: `ls internal/db internal/crawler internal/parser internal/bible internal/tui internal/config` 모두 존재
  - [ ] `.gitattributes` 존재, `*.golden -text` 포함

  **Commit**: YES
  - Message: `feat: project scaffold with cobra root command and directory structure`
  - Files: `go.mod`, `go.sum`, `main.go`, `cmd/root.go`, `internal/*/`, `.gitattributes`

---

- [ ] 2. SQLite DB Layer + Schema (Multi-Version Ready)

  **What to do**:
  - `internal/db/db.go`: DB 연결 관리
    - `Open(path string) (*DB, error)` — SQLite 열기, WAL mode, foreign keys on
    - `Close() error`
    - `Migrate() error` — 스키마 생성/마이그레이션
    - DB 경로: `os.UserConfigDir()` + `/bible-tui/bible.db`
  - 스키마 (다중 역본 대비):
    ```sql
    CREATE TABLE IF NOT EXISTS versions (
      id INTEGER PRIMARY KEY,
      code TEXT UNIQUE NOT NULL,     -- 'GAE', 'HAN', 'SAE', etc.
      name TEXT NOT NULL,            -- '개역개정', '개역한글', etc.
      lang TEXT NOT NULL DEFAULT 'ko'
    );

    CREATE TABLE IF NOT EXISTS books (
      id INTEGER PRIMARY KEY,
      version_id INTEGER NOT NULL REFERENCES versions(id),
      code TEXT NOT NULL,            -- 'gen', 'exo', etc.
      name_ko TEXT NOT NULL,         -- '창세기', '출애굽기', etc.
      abbrev_ko TEXT NOT NULL,       -- '창', '출', etc.
      testament TEXT NOT NULL,       -- 'old', 'new'
      chapter_count INTEGER NOT NULL,
      sort_order INTEGER NOT NULL,
      UNIQUE(version_id, code)
    );

    CREATE TABLE IF NOT EXISTS verses (
      id INTEGER PRIMARY KEY,
      book_id INTEGER NOT NULL REFERENCES books(id),
      chapter INTEGER NOT NULL,
      verse_num INTEGER NOT NULL,
      text TEXT NOT NULL,
      section_title TEXT,            -- 소제목 (nullable)
      has_footnote BOOLEAN NOT NULL DEFAULT 0,
      UNIQUE(book_id, chapter, verse_num)
    );

    CREATE TABLE IF NOT EXISTS footnotes (
      id INTEGER PRIMARY KEY,
      verse_id INTEGER NOT NULL REFERENCES verses(id),
      marker TEXT,                   -- '1)', '2)' etc.
      content TEXT NOT NULL
    );

    -- FTS5 for search
    CREATE VIRTUAL TABLE IF NOT EXISTS verses_fts USING fts5(
      text,
      content=verses,
      content_rowid=id,
      tokenize='unicode61'
    );

    -- FTS sync triggers
    CREATE TRIGGER IF NOT EXISTS verses_ai AFTER INSERT ON verses BEGIN
      INSERT INTO verses_fts(rowid, text) VALUES (new.id, new.text);
    END;
    CREATE TRIGGER IF NOT EXISTS verses_ad AFTER DELETE ON verses BEGIN
      INSERT INTO verses_fts(verses_fts, rowid, text) VALUES('delete', old.id, old.text);
    END;
    CREATE TRIGGER IF NOT EXISTS verses_au AFTER UPDATE ON verses BEGIN
      INSERT INTO verses_fts(verses_fts, rowid, text) VALUES('delete', old.id, old.text);
      INSERT INTO verses_fts(rowid, text) VALUES (new.id, new.text);
    END;

    -- User data
    CREATE TABLE IF NOT EXISTS bookmarks (
      id INTEGER PRIMARY KEY,
      verse_id INTEGER NOT NULL REFERENCES verses(id),
      note TEXT,
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS highlights (
      id INTEGER PRIMARY KEY,
      verse_id INTEGER NOT NULL REFERENCES verses(id),
      color TEXT NOT NULL DEFAULT 'yellow',
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      UNIQUE(verse_id)
    );

    -- Reading plans
    CREATE TABLE IF NOT EXISTS reading_plans (
      id INTEGER PRIMARY KEY,
      name TEXT NOT NULL,
      plan_type TEXT NOT NULL,       -- 'sequential', 'mcheyne', 'custom'
      version_id INTEGER NOT NULL REFERENCES versions(id),
      total_days INTEGER NOT NULL,
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS reading_plan_entries (
      id INTEGER PRIMARY KEY,
      plan_id INTEGER NOT NULL REFERENCES reading_plans(id),
      day_number INTEGER NOT NULL,
      book_code TEXT NOT NULL,
      chapter_start INTEGER NOT NULL,
      chapter_end INTEGER NOT NULL,
      completed BOOLEAN NOT NULL DEFAULT 0,
      completed_at DATETIME,
      UNIQUE(plan_id, day_number, book_code, chapter_start)
    );

    -- Settings (theme, font, etc.)
    CREATE TABLE IF NOT EXISTS settings (
      key TEXT PRIMARY KEY,
      value TEXT NOT NULL
    );

    -- Crawler state
    CREATE TABLE IF NOT EXISTS crawl_status (
      version_code TEXT NOT NULL,
      book_code TEXT NOT NULL,
      chapter INTEGER NOT NULL,
      status TEXT NOT NULL DEFAULT 'pending',  -- 'pending', 'done', 'error'
      verse_count INTEGER,
      crawled_at DATETIME,
      error_msg TEXT,
      PRIMARY KEY(version_code, book_code, chapter)
    );
    ```
  - CRUD 함수들:
    - `InsertVersion(code, name, lang) error`
    - `InsertBook(versionID, code, nameKo, abbrevKo, testament, chapterCount, sortOrder) (int64, error)`
    - `InsertVerse(bookID, chapter, verseNum, text, sectionTitle, hasFootnote) (int64, error)`
    - `InsertFootnote(verseID, marker, content) error`
    - `GetVerses(versionCode, bookCode string, chapter int) ([]Verse, error)`
    - `GetRandomVerse(versionCode string) (*Verse, error)`
    - `GetSetting(key string) (string, error)` / `SetSetting(key, value string) error`
  - 테스트 (TDD):
    - `internal/db/db_test.go`
    - TestMigrate: 스키마 생성 확인
    - TestInsertAndGetVerses: 삽입→조회 round trip
    - TestGetRandomVerse: 랜덤 구절 반환 확인
    - TestSettings: key-value 저장/조회
    - 모든 테스트는 `:memory:` SQLite 사용

  **Must NOT do**:
  - 검색 함수 (Task 8에서)
  - 북마크/하이라이트 CRUD (Task 9에서)
  - 읽기계획 CRUD (Task 18에서)

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: DB 스키마 설계 + Go 구현 + TDD — 중간-높은 복잡도
  - **Skills**: []
    - 특수 스킬 불필요 (Go 표준 패턴)

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Task 3)
  - **Blocks**: Tasks 4, 5, 6, 7, 8, 9
  - **Blocked By**: Task 1

  **References**:
  - **Library**: `modernc.org/sqlite` — pure Go SQLite, FTS5 지원
  - **Pattern**: `zk-org/zk` — FTS5 content sync 패턴 (content table + virtual table + triggers)
  - **Pattern**: WAL mode + foreign keys: `PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON;`
  - **API**: `os.UserConfigDir()` — OS별 표준 설정 경로 (`~/Library/Application Support` / `~/.config` / `%AppData%`)
  - **Schema Note**: `versions` 테이블로 다중 역본 지원. `books.version_id`로 역본별 책 분리. 향후 HAN/SAE 등 추가 시 같은 스키마 사용.

  **Acceptance Criteria**:
  - [ ] `go test ./internal/db/... -v -count=1` → ALL PASS
  - [ ] `go test ./internal/db/... -race` → no race conditions
  - [ ] TestMigrate: 모든 테이블 생성 확인 (`SELECT name FROM sqlite_master WHERE type='table'`)
  - [ ] TestInsertAndGetVerses: 3개 절 삽입 → 조회 → text 일치 확인
  - [ ] TestGetRandomVerse: nil이 아닌 Verse 반환 확인
  - [ ] TestSettings: SetSetting → GetSetting → 값 일치

  **Commit**: YES
  - Message: `feat(db): SQLite schema with multi-version support and CRUD operations`
  - Files: `internal/db/db.go`, `internal/db/db_test.go`, `internal/db/models.go`

---

- [ ] 3. Bible Reference Parser (한글 책이름 → 코드 매핑)

  **What to do**:
  - `internal/bible/reference.go`:
    - 66권 매핑 테이블: 한글 전체명, 한글 약어, 영문 코드
      ```
      {"창세기", "창", "gen", "old", 50},
      {"출애굽기", "출", "exo", "old", 40},
      ...
      {"요한계시록", "계", "rev", "new", 22},
      ```
    - `ParseReference(input string) (*Reference, error)` — 입력 파싱
      - `"창세기 1"` → `{BookCode: "gen", Chapter: 1, Verse: 0}`
      - `"창 1"` → `{BookCode: "gen", Chapter: 1, Verse: 0}`
      - `"gen 1"` → `{BookCode: "gen", Chapter: 1, Verse: 0}`
      - `"창세기 1:3"` → `{BookCode: "gen", Chapter: 1, Verse: 3}`
      - `"창 1:3-5"` → `{BookCode: "gen", Chapter: 1, VerseStart: 3, VerseEnd: 5}`
    - `GetBookName(code string) string` — 코드 → 한글명
    - `GetBookAbbrev(code string) string` — 코드 → 한글 약어
    - `AllBooks() []BookInfo` — 전체 목록 (TUI 목록 표시용)
    - `BookInfo` struct: `Code, NameKo, AbbrevKo, Testament, ChapterCount`
  - 테스트 (TDD):
    - `internal/bible/reference_test.go`
    - TestParseReference_FullName: "창세기 1" → gen, 1
    - TestParseReference_Abbrev: "창 1" → gen, 1
    - TestParseReference_English: "gen 1" → gen, 1
    - TestParseReference_WithVerse: "창 1:3" → gen, 1, 3
    - TestParseReference_Range: "창 1:3-5" → gen, 1, 3, 5
    - TestParseReference_Invalid: "없는책 1" → error
    - TestAllBooks: 66권 반환 확인

  **Must NOT do**:
  - DB 연동 (이것은 순수 로직)
  - 다국어 매핑 (한글+영문만)

  **Recommended Agent Profile**:
  - **Category**: `unspecified-low`
    - Reason: 순수 데이터 매핑 + 파서, 복잡도 낮음
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Task 2)
  - **Blocks**: Tasks 6, 7
  - **Blocked By**: Task 1

  **References**:
  - **Data Source**: `bible.list.js`의 `szHANBook` 배열 — 66권 영문코드 + 한글명 + 장수
    ```
    szHANBook[0] = new Array("창세기", "gen", "1"..."50");  // 50장
    szHANBook[1] = new Array("출애굽기", "exo", "1"..."40"); // 40장
    ...
    szHANBook[65] = new Array("요한계시록", "rev", "1"..."22"); // 22장
    ```
  - **한글 약어 참고**: 대한성서공회 약어 표준 (창/출/레/민/신/수/삿/룻/삼상/삼하/왕상/왕하...)
  - **Pattern**: Go regex for 파싱 — `regexp.MustCompile(`^([가-힣a-zA-Z0-9]+)\s+(\d+)(?::(\d+)(?:-(\d+))?)?$`)`

  **Acceptance Criteria**:
  - [ ] `go test ./internal/bible/... -v -count=1` → ALL PASS
  - [ ] TestParseReference_FullName: "창세기 1" → BookCode=="gen", Chapter==1
  - [ ] TestParseReference_Abbrev: "창 1" → BookCode=="gen", Chapter==1
  - [ ] TestAllBooks: len == 66
  - [ ] `go vet ./internal/bible/...` → clean

  **Commit**: YES
  - Message: `feat(bible): reference parser with Korean name/abbreviation mapping`
  - Files: `internal/bible/reference.go`, `internal/bible/reference_test.go`, `internal/bible/books.go`

---

- [ ] 4. HTML Parser (크롤링 HTML → 구조화된 데이터)

  **What to do**:
  - `internal/parser/parser.go`:
    - `ParseChapterHTML(html string) (*ChapterData, error)`
    - `ChapterData` struct:
      ```go
      type ChapterData struct {
        BookCode     string
        Chapter      int
        VersionCode  string
        SectionTitle string        // 첫 소제목
        Verses       []VerseData
      }
      type VerseData struct {
        Number       int
        Text         string       // HTML 태그 제거된 순수 텍스트
        SectionTitle string       // 이 절 앞에 나오는 소제목 (nullable)
        Footnotes    []FootnoteData
      }
      type FootnoteData struct {
        Marker  string            // "1)", "2)"
        Content string
      }
      ```
    - 파싱 로직:
      1. `goquery`로 `div#tdBible1` 선택
      2. `<span class="number">N&nbsp;&nbsp;&nbsp;</span>` → 절 번호 추출
      3. 절 번호 다음 텍스트 → 절 본문 (HTML 태그 strip)
      4. `<font class="smallTitle">` → 소제목 추출
      5. `<div class="D2">` 또는 `<a class="comment">` → 각주 추출
      6. `<font size='1'>` 내부 텍스트 → 본문에 포함 (작은 글씨 주석)
  - `testdata/genesis_1.html` — 실제 크롤링한 HTML 샘플 저장 (fixture)
  - `testdata/psalm_119.html` — 긴 장 테스트용
  - 테스트 (TDD — fixture 기반):
    - TestParseGenesis1: 31절, 소제목 "천지 창조" 포함
    - TestParseVerse1Text: "태초에 하나님이 천지를 창조하시니라" 정확히 일치
    - TestParseFootnotes: 각주 마커와 내용 일치
    - TestParseSectionTitles: 소제목 위치 정확성
    - TestParseEmptyHTML: 에러 반환

  **Must NOT do**:
  - HTTP 요청 (파서는 순수 HTML 문자열 입력)
  - DB 저장 (Task 5에서)

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: HTML 파싱 복잡도 높음 (불규칙한 HTML, 여러 요소 추출)
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO (depends on Task 2 for output struct alignment)
  - **Parallel Group**: Wave 3 (after Task 2)
  - **Blocks**: Task 5
  - **Blocked By**: Task 2

  **References**:
  - **Data Source**: 크롤링된 실제 HTML (리서치에서 확보)
    - `<div id="tdBible1" class="bible_read">` — 전체 본문 컨테이너
    - `<span class="number">1&nbsp;&nbsp;&nbsp;</span>` — 절 번호
    - `<font class="smallTitle">천지 창조</font>` — 소제목
    - `<font class="chapNum">제 1 장</font>` — 장 번호
    - `<div class=D2>또는 형체가 없는</div>` — 각주 내용
    - `<a class=comment>1)</a>` — 각주 마커
  - **Library**: `github.com/PuerkitoBio/goquery` — jQuery-like HTML selector
  - **Fixture 생성**: `testdata/genesis_1.html`은 `https://www.bskorea.or.kr/bible/korbibReadpage.php?version=GAE&book=gen&chap=1` 에서 다운로드

  **Acceptance Criteria**:
  - [ ] `go test ./internal/parser/... -v -count=1` → ALL PASS
  - [ ] TestParseGenesis1: `len(result.Verses) == 31`
  - [ ] TestParseVerse1Text: `result.Verses[0].Text` contains "태초에 하나님이 천지를 창조하시니라"
  - [ ] TestParseFootnotes: Genesis 1에서 각주 3개 이상 추출
  - [ ] `testdata/genesis_1.html` 파일 존재

  **Commit**: YES
  - Message: `feat(parser): HTML parser for bskorea.or.kr Bible pages with fixture tests`
  - Files: `internal/parser/parser.go`, `internal/parser/parser_test.go`, `testdata/genesis_1.html`

---

- [ ] 5. Crawler Engine (HTTP + Rate Limiting + Checkpoint)

  **What to do**:
  - `internal/crawler/crawler.go`:
    - `Crawler` struct:
      ```go
      type Crawler struct {
        db          *db.DB
        client      *http.Client
        limiter     *rate.Limiter      // 0.5 req/sec (2초 간격)
        userAgent   string             // "BibleTUI/1.0 (Personal; non-commercial)"
        baseURL     string
        versionCode string             // "GAE" (파라미터화 — 다중 역본 대비)
      }
      ```
    - `New(db, opts...) *Crawler` — 생성자 (functional options)
    - `CrawlAll(ctx context.Context) error` — 전체 크롤링
      1. `versions` 테이블에 현재 버전 upsert
      2. `books` 테이블에 66권 bulk insert (bible 패키지의 AllBooks() 사용)
      3. 각 book/chapter에 대해:
         - `crawl_status` 체크 → 이미 done이면 skip
         - rate limiter 대기
         - HTTP GET `baseURL?version={code}&book={bookCode}&chap={chapter}`
         - `charset.NewReader()` 로 인코딩 안전하게 처리
         - parser로 HTML 파싱
         - 결과를 DB에 insert (verse + footnotes)
         - `crawl_status` → 'done' 업데이트
         - 에러 시 `crawl_status` → 'error' + error_msg
      4. 완료 후 무결성 검증:
         - 66권 존재 확인
         - 각 권의 장 수 일치 확인
         - 빈 text 절 없음 확인
    - `CrawlBook(ctx, bookCode string) error` — 특정 책만 크롤링
    - Progress 콜백: `func(book string, chapter, total int)` — 진행률 표시용
  - 테스트 (TDD):
    - `internal/crawler/crawler_test.go`
    - TestCrawlAll_DryRun: HTTP 모킹, DB에 올바르게 저장되는지 확인
    - TestResumability: 일부 완료 후 재실행 → skip 확인
    - TestRateLimiting: limiter 동작 확인
    - TestValidation: 무결성 검증 로직 테스트
    - HTTP 모킹: `httptest.NewServer` + fixture HTML 반환

  **Must NOT do**:
  - CLI command 연결 (Task 7에서)
  - 실제 사이트 호출 (테스트에서는 모킹)

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: HTTP + rate limiting + checkpoint + 무결성 검증 — 높은 복잡도
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 4 (sequential)
  - **Blocks**: Task 7
  - **Blocked By**: Tasks 2, 4

  **References**:
  - **Library**: `golang.org/x/time/rate` — `rate.NewLimiter(0.5, 1)` (2초당 1 요청)
  - **Library**: `golang.org/x/net/html/charset` — `charset.NewReader(resp.Body, contentType)`
  - **Library**: `net/http` — `http.Client` with timeout
  - **Pattern**: `httptest.NewServer` — HTTP 모킹 for 테스트
  - **URL Pattern**: `https://www.bskorea.or.kr/bible/korbibReadpage.php?version=GAE&book=gen&chap=1`
  - **Data**: `bible.list.js` — 66권 코드 + 장수 (무결성 검증 기준)
  - **User-Agent**: `BibleTUI/1.0 (Personal; non-commercial)`

  **Acceptance Criteria**:
  - [ ] `go test ./internal/crawler/... -v -count=1` → ALL PASS
  - [ ] TestCrawlAll_DryRun: mock 서버에서 3권×2장 크롤링 → DB에 삽입 확인
  - [ ] TestResumability: 중단 후 재실행 시 이미 완료된 장 skip 확인
  - [ ] TestValidation: 장수 불일치 시 에러 반환

  **Commit**: YES
  - Message: `feat(crawler): web crawler with rate limiting, checkpointing, and validation`
  - Files: `internal/crawler/crawler.go`, `internal/crawler/crawler_test.go`

---

- [ ] 6. CLI Read & Random Commands

  **What to do**:
  - `cmd/read.go`:
    - `bible read <책이름> <장>` — 해당 장 전체 출력
    - `bible read <책이름> <장>:<절>` — 특정 절 출력
    - `bible read <책이름> <장>:<시작절>-<끝절>` — 범위 출력
    - `--version` 플래그 (기본값: "GAE", 다중 역본 대비)
    - 출력 포맷: 절 번호 + 본문, 소제목은 볼드/별도 라인
    - `lipgloss`로 터미널 컬러 출력 (절 번호 색상 등)
    - `go-runewidth`로 정렬
  - `cmd/random.go`:
    - `bible random` — 랜덤 구절 1개 출력
    - `--version` 플래그
    - 출력: "창세기 1:1 — 태초에 하나님이 천지를 창조하시니라"
  - 테스트 (TDD):
    - `cmd/read_test.go`, `cmd/random_test.go`
    - TestReadCommand: in-memory DB에 fixture 데이터 → read 실행 → stdout 확인
    - TestReadCommand_InvalidBook: 에러 메시지 확인
    - TestRandomCommand: 출력 비어있지 않음 확인

  **Must NOT do**:
  - TUI 연동
  - search command (Task 10)

  **Recommended Agent Profile**:
  - **Category**: `unspecified-low`
    - Reason: DB 조회 + 포맷팅 출력, 중간 복잡도
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 5 (with Task 7)
  - **Blocks**: Task 12
  - **Blocked By**: Tasks 2, 3

  **References**:
  - **Pattern**: cobra `RunE` + `cmd.Context()` for context propagation
  - **Library**: `github.com/charmbracelet/lipgloss` — 터미널 스타일링
  - **Library**: `github.com/mattn/go-runewidth` — `runewidth.StringWidth()` for 한글 너비
  - **Internal**: `internal/bible/reference.go` — `ParseReference()` 입력 파싱
  - **Internal**: `internal/db/db.go` — `GetVerses()`, `GetRandomVerse()` 조회

  **Acceptance Criteria**:
  - [ ] `go test ./cmd/... -v -run TestRead` → PASS
  - [ ] `go test ./cmd/... -v -run TestRandom` → PASS
  - [ ] `go run . read --help` → "book" 포함 도움말
  - [ ] `go run . random --help` → 도움말 출력
  ```bash
  # Integration test (after crawl):
  go run . read 창세기 1 | head -5
  # Assert: "태초에" 포함
  go run . random
  # Assert: output non-empty, contains ":"
  ```

  **Commit**: YES
  - Message: `feat(cli): read and random commands with Korean reference parsing`
  - Files: `cmd/read.go`, `cmd/read_test.go`, `cmd/random.go`, `cmd/random_test.go`

---

- [ ] 7. CLI Crawl Command

  **What to do**:
  - `cmd/crawl.go`:
    - `bible crawl` — 전체 크롤링 실행
    - `--version` 플래그 (기본값: "GAE")
    - `--book` 플래그 (특정 책만 크롤링)
    - `--dry-run` 플래그 (DB 스키마만 생성, 크롤링 안 함)
    - 진행률 표시: `[32/1189] 창세기 32장 크롤링 중...`
    - 완료 후 무결성 검증 결과 출력
    - 에러 시 재시도 안내 메시지
  - 테스트 (TDD):
    - `cmd/crawl_test.go`
    - TestCrawlCommand_DryRun: `--dry-run` → 스키마 생성만 확인
    - TestCrawlCommand_Help: 도움말에 version/book/dry-run 플래그 포함

  **Must NOT do**:
  - 크롤러 로직 수정 (Task 5에서 완성)

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: 크롤러를 cobra command에 연결하는 간단 작업
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 5 (with Task 6)
  - **Blocks**: Phase 2
  - **Blocked By**: Task 5

  **References**:
  - **Internal**: `internal/crawler/crawler.go` — `New()`, `CrawlAll()`
  - **Pattern**: cobra `RunE` + progress callback
  - **Library**: `log/slog` — structured logging for 크롤링 진행 상황

  **Acceptance Criteria**:
  - [ ] `go test ./cmd/... -v -run TestCrawl` → PASS
  - [ ] `go run . crawl --help` → version, book, dry-run 플래그 표시
  - [ ] `go run . crawl --dry-run` → exit 0, 스키마 생성 확인

  **Commit**: YES
  - Message: `feat(cli): crawl command with progress display and dry-run mode`
  - Files: `cmd/crawl.go`, `cmd/crawl_test.go`

---

### Phase 2: Search & Bookmarks

---

- [ ] 8. Search Layer (FTS5 + LIKE Fallback)

  **What to do**:
  - `internal/db/search.go`:
    - `SearchVerses(versionCode, query string, limit int) ([]SearchResult, error)`
      - 전략: FTS5 MATCH 먼저 시도 → 결과 < 5개면 LIKE fallback 추가
      - FTS5: `SELECT ... FROM verses_fts WHERE verses_fts MATCH ? ...`
      - LIKE: `SELECT ... FROM verses WHERE text LIKE '%' || ? || '%' ...`
      - `SearchResult`: Verse + BookName + Highlight snippet
    - 결과에 하이라이트: 검색어 위치 표시 (lipgloss 스타일용 offset 정보)
  - 테스트 (TDD):
    - TestSearchVerses_FTS: "하나님" 검색 → 결과 반환
    - TestSearchVerses_LIKE: LIKE 모드로 부분 문자열 검색
    - TestSearchVerses_NoResult: 없는 단어 → 빈 배열
    - TestSearchVerses_Limit: limit 적용 확인

  **Must NOT do**:
  - 형태소 분석
  - CLI command (Task 10)

  **Recommended Agent Profile**:
  - **Category**: `unspecified-low`
    - Reason: SQL 쿼리 위주, 복잡도 중간
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 6 (with Task 9)
  - **Blocks**: Task 10
  - **Blocked By**: Phase 1 complete

  **References**:
  - **SQLite FTS5**: `SELECT rowid, snippet(verses_fts, 0, '**', '**', '...', 32) FROM verses_fts WHERE verses_fts MATCH ?`
  - **Pattern**: FTS5 + LIKE 병행 전략 — `zk-org/zk` 참고
  - **Note**: `modernc.org/sqlite`는 커스텀 토크나이저 미지원 → `unicode61` 기본 토크나이저 사용

  **Acceptance Criteria**:
  - [ ] `go test ./internal/db/... -v -run TestSearch` → ALL PASS
  - [ ] TestSearchVerses_FTS: "하나님" → 결과 1개 이상
  - [ ] TestSearchVerses_Limit: limit=3 → 최대 3개 반환

  **Commit**: YES
  - Message: `feat(db): full-text search with FTS5 and LIKE fallback`
  - Files: `internal/db/search.go`, `internal/db/search_test.go`

---

- [ ] 9. Bookmark & Highlight Data Layer

  **What to do**:
  - `internal/db/bookmark.go`:
    - `AddBookmark(verseID int64, note string) error`
    - `RemoveBookmark(id int64) error`
    - `ListBookmarks(limit, offset int) ([]Bookmark, error)` — 최신순
    - `IsBookmarked(verseID int64) (bool, error)`
  - `internal/db/highlight.go`:
    - `AddHighlight(verseID int64, color string) error`
    - `RemoveHighlight(verseID int64) error`
    - `ListHighlights(limit, offset int) ([]Highlight, error)`
    - `GetHighlightColor(verseID int64) (string, error)` — 없으면 ""
    - 지원 색상: "yellow", "green", "blue", "pink", "purple"
  - 테스트 (TDD):
    - TestAddAndListBookmarks: 추가 → 목록에 표시
    - TestRemoveBookmark: 삭제 확인
    - TestAddAndGetHighlight: 추가 → 색상 조회
    - TestHighlightUpsert: 같은 절에 다른 색상 → 업데이트

  **Must NOT do**:
  - CLI/TUI 연동
  - 내보내기/가져오기

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: 단순 CRUD 작업
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 6 (with Task 8)
  - **Blocks**: Task 11
  - **Blocked By**: Phase 1 complete

  **References**:
  - **Internal**: `internal/db/db.go` — DB 구조, 모델 패턴
  - **Schema**: `bookmarks`, `highlights` 테이블 (Task 2에서 생성)

  **Acceptance Criteria**:
  - [ ] `go test ./internal/db/... -v -run TestBookmark` → ALL PASS
  - [ ] `go test ./internal/db/... -v -run TestHighlight` → ALL PASS

  **Commit**: YES
  - Message: `feat(db): bookmark and highlight CRUD operations`
  - Files: `internal/db/bookmark.go`, `internal/db/bookmark_test.go`, `internal/db/highlight.go`, `internal/db/highlight_test.go`

---

- [ ] 10. CLI Search Command

  **What to do**:
  - `cmd/search.go`:
    - `bible search <검색어>` — 검색 결과 출력
    - `--version` 플래그 (기본값: "GAE")
    - `--limit` 플래그 (기본값: 20)
    - 출력 포맷:
      ```
      Found 15 results for "사랑":

      [1] 요한복음 3:16
          하나님이 세상을 이처럼 **사랑**하사 독생자를 주셨으니...

      [2] 고린도전서 13:4
          **사랑**은 오래 참고 **사랑**은 온유하며...
      ```
    - 검색어 하이라이트 (lipgloss bold/color)
  - 테스트 (TDD):
    - TestSearchCommand: fixture DB → search → stdout에 결과 포함
    - TestSearchCommand_NoResults: 없는 단어 → "No results" 메시지

  **Must NOT do**:
  - 검색 로직 수정 (Task 8에서)

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: search layer를 CLI에 연결하는 간단 작업
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 7 (with Task 11)
  - **Blocks**: Task 14
  - **Blocked By**: Task 8

  **References**:
  - **Internal**: `internal/db/search.go` — `SearchVerses()`
  - **Library**: `lipgloss` — 검색어 하이라이트 스타일링

  **Acceptance Criteria**:
  - [ ] `go test ./cmd/... -v -run TestSearch` → PASS
  - [ ] `go run . search --help` → limit, version 플래그 표시

  **Commit**: YES
  - Message: `feat(cli): search command with highlighted results`
  - Files: `cmd/search.go`, `cmd/search_test.go`

---

- [ ] 11. CLI Bookmark Commands

  **What to do**:
  - `cmd/bookmark.go`:
    - `bible bookmark add <참조> [--note "메모"]` — 책갈피 추가
    - `bible bookmark list` — 책갈피 목록
    - `bible bookmark remove <id>` — 책갈피 삭제
    - `bible highlight add <참조> [--color yellow]` — 하이라이트
    - `bible highlight list` — 하이라이트 목록
    - `bible highlight remove <참조>` — 하이라이트 삭제
  - 테스트 (TDD):
    - TestBookmarkAdd: 추가 → list에 표시
    - TestHighlightAdd: 추가 → list에 표시

  **Must NOT do**:
  - TUI 연동

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: DB CRUD를 CLI에 연결
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 7 (with Task 10)
  - **Blocks**: Task 15
  - **Blocked By**: Task 9

  **References**:
  - **Internal**: `internal/db/bookmark.go`, `internal/db/highlight.go`
  - **Internal**: `internal/bible/reference.go` — `ParseReference()`

  **Acceptance Criteria**:
  - [ ] `go test ./cmd/... -v -run TestBookmark` → PASS
  - [ ] `go run . bookmark --help` → add, list, remove 서브커맨드 표시

  **Commit**: YES
  - Message: `feat(cli): bookmark and highlight management commands`
  - Files: `cmd/bookmark.go`, `cmd/bookmark_test.go`

---

### Phase 3: TUI Core

---

- [ ] 12. TUI Shell & Navigation Framework

  **What to do**:
  - `internal/tui/app.go`:
    - Root Model:
      ```go
      type AppModel struct {
        state       AppState        // current view
        db          *db.DB
        width       int
        height      int
        theme       *Theme          // 현재 테마
        fontSize    int             // 현재 폰트 크기 설정값
        // sub-models (lazy init)
        bookList    *BookListModel
        reading     *ReadingModel
        search      *SearchModel
        bookmarks   *BookmarkModel
        settings    *SettingsModel
        plans       *PlanModel
        help        *HelpModel
      }
      type AppState int
      const (
        StateBookList AppState = iota
        StateChapterList
        StateReading
        StateSearch
        StateBookmarks
        StateSettings
        StatePlans
        StateHelp
      )
      ```
    - 키바인딩:
      - `q` / `Ctrl+C` — 종료
      - `?` — 도움말
      - `b` — 책 목록
      - `/` — 검색
      - `m` — 책갈피
      - `p` — 읽기 계획
      - `s` — 설정
      - `Esc` — 이전 화면
    - `tea.WindowSizeMsg` 핸들링 for responsive layout
    - 상태 전환: state machine 패턴
    - 하단 상태바: 현재 위치 + 키 힌트
  - `cmd/tui.go`:
    - `bible tui` command → `bubbletea.NewProgram(app.New(db))` 실행
  - `internal/tui/styles/theme.go`:
    - `Theme` struct 기본 정의 (상세 구현은 Task 19)
    - `DefaultDarkTheme()`, `DefaultLightTheme()` 기본 골격만
  - 테스트 (TDD):
    - `internal/tui/app_test.go` — `teatest` 사용
    - TestAppInit: 초기 상태 == StateBookList
    - TestAppKeyQ: `q` → 종료 메시지
    - TestAppKeyHelp: `?` → StateHelp 전환
    - `lipgloss.SetColorProfile(termenv.Ascii)` in test setup

  **Must NOT do**:
  - 각 view의 상세 구현 (다음 tasks에서)
  - 테마/폰트 상세 로직 (Task 19, 20)

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
    - Reason: TUI 레이아웃, 상태 머신, 키바인딩 — UI 중심
  - **Skills**: [`frontend-ui-ux`]
    - `frontend-ui-ux`: TUI도 UI/UX 설계가 중요 — 네비게이션 흐름, 레이아웃

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 8 (solo — TUI foundation)
  - **Blocks**: Tasks 13, 14, 15, 16
  - **Blocked By**: Phase 2 complete

  **References**:
  - **Library**: `github.com/charmbracelet/bubbletea` — Model-Update-View 패턴
  - **Library**: `github.com/charmbracelet/lipgloss` — 터미널 스타일링
  - **Library**: `github.com/charmbracelet/bubbles` — 재사용 컴포넌트
  - **Library**: `github.com/charmbracelet/x/exp/teatest` — TUI 테스트
  - **Pattern**: `charmbracelet/mods` — cobra에서 bubbletea 실행하는 패턴
  - **Pattern**: 상태 머신: switch state in Update → delegate to sub-model
  - **Pattern**: `lipgloss.SetColorProfile(termenv.Ascii)` in tests

  **Acceptance Criteria**:
  - [ ] `go test ./internal/tui/... -v -count=1` → ALL PASS
  - [ ] `go run . tui --help` → 도움말 표시
  - [ ] TestAppInit: initial state == StateBookList
  - [ ] TestAppKeyQ: sends tea.Quit

  **Commit**: YES
  - Message: `feat(tui): app shell with state machine navigation and keybindings`
  - Files: `internal/tui/app.go`, `internal/tui/app_test.go`, `internal/tui/styles/theme.go`, `cmd/tui.go`

---

- [ ] 13. TUI Reading View (Book → Chapter → Verse)

  **What to do**:
  - `internal/tui/booklist.go`:
    - 구약 39권 + 신약 27권 목록 표시
    - `bubbles/list` 사용
    - Enter → 장 목록으로 이동
  - `internal/tui/chapterlist.go`:
    - 선택된 책의 장 목록 (그리드 레이아웃: 5열)
    - Enter → 해당 장 읽기 화면
  - `internal/tui/reading.go`:
    - `bubbles/viewport` 사용 — 스크롤 가능한 성경 본문
    - 절 번호 + 본문 (lipgloss 스타일)
    - 소제목 표시 (별도 스타일)
    - `h` / `←` → 이전 장, `l` / `→` → 다음 장
    - `j` / `k` / 스크롤 — 위아래 이동
    - 상단: "창세기 1장" 타이틀
    - 하단: 키 힌트 바
    - 하이라이트된 절은 배경색 표시
    - 북마크된 절은 아이콘 표시
    - DB에서 async 로딩: `tea.Cmd` → custom message
  - 테스트 (TDD):
    - TestBookListModel: 66권 항목 표시
    - TestReadingModel: 절 데이터 → viewport 내용에 텍스트 포함
    - TestReadingNavigation: h/l → 장 변경 메시지

  **Must NOT do**:
  - 검색 뷰 (Task 14)
  - 테마 적용 상세 (Task 19에서 통합)

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
    - Reason: TUI 핵심 읽기 화면, 레이아웃 디자인 중심
  - **Skills**: [`frontend-ui-ux`]
    - `frontend-ui-ux`: 읽기 경험의 UX가 프로그램 핵심

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 9 (with Task 14)
  - **Blocks**: Tasks 15, 19
  - **Blocked By**: Task 12

  **References**:
  - **Library**: `bubbles/viewport` — 스크롤 가능한 텍스트 뷰
  - **Library**: `bubbles/list` — 아이템 목록 (책 목록)
  - **Library**: `lipgloss` — 절 번호, 소제목, 본문 스타일링
  - **Library**: `go-runewidth` — 한글 텍스트 너비 계산
  - **Pattern**: async DB query → `tea.Cmd` returning `VersesLoadedMsg`
  - **Internal**: `internal/db/db.go` — `GetVerses()`
  - **Internal**: `internal/db/highlight.go` — `GetHighlightColor()`
  - **Internal**: `internal/db/bookmark.go` — `IsBookmarked()`

  **Acceptance Criteria**:
  - [ ] `go test ./internal/tui/... -v -run TestBookList` → PASS
  - [ ] `go test ./internal/tui/... -v -run TestReading` → PASS
  - [ ] TestBookListModel: 66개 항목 + 구약/신약 구분
  - [ ] TestReadingModel: "태초에" 텍스트 viewport에 포함

  **Commit**: YES
  - Message: `feat(tui): reading view with book/chapter navigation and viewport`
  - Files: `internal/tui/booklist.go`, `internal/tui/chapterlist.go`, `internal/tui/reading.go`, + tests

---

- [ ] 14. TUI Search View

  **What to do**:
  - `internal/tui/search.go`:
    - `bubbles/textinput` — 검색어 입력
    - 검색 결과 `bubbles/list`로 표시
    - 검색어 하이라이트 (lipgloss bold/color)
    - Enter on result → 해당 구절의 읽기 화면으로 이동
    - 실시간 검색 (debounce 300ms) 또는 Enter로 검색
    - 결과 포맷: `[책이름 장:절] 본문 (검색어 강조)`
  - 테스트 (TDD):
    - TestSearchModel_Input: 텍스트 입력 → 검색 실행
    - TestSearchModel_Results: 결과 목록 표시

  **Must NOT do**:
  - 검색 로직 변경 (Task 8)
  - 정규식 검색

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
    - Reason: 검색 UI/UX, 하이라이트, debounce 등 UI 중심
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 9 (with Task 13)
  - **Blocks**: None directly
  - **Blocked By**: Task 12

  **References**:
  - **Library**: `bubbles/textinput` — 텍스트 입력 위젯
  - **Library**: `bubbles/list` — 검색 결과 목록
  - **Internal**: `internal/db/search.go` — `SearchVerses()`
  - **Pattern**: debounce: `time.AfterFunc` + `tea.Cmd`

  **Acceptance Criteria**:
  - [ ] `go test ./internal/tui/... -v -run TestSearch` → PASS
  - [ ] TestSearchModel_Input: 입력 → 검색 메시지 발생
  - [ ] TestSearchModel_Results: 결과 → 리스트 아이템으로 변환

  **Commit**: YES
  - Message: `feat(tui): search view with text input and highlighted results`
  - Files: `internal/tui/search.go`, `internal/tui/search_test.go`

---

- [ ] 15. TUI Bookmark & Highlight View

  **What to do**:
  - `internal/tui/bookmarks.go`:
    - 책갈피 목록 (`bubbles/list`)
    - 하이라이트 목록 (탭 전환 또는 필터)
    - Enter → 해당 구절 읽기 화면
    - `d` → 삭제
    - 각 항목: `[책이름 장:절] 본문 미리보기... (메모)`
    - 하이라이트: 색상 표시
  - 읽기 화면에서 책갈피/하이라이트 토글:
    - `internal/tui/reading.go` 수정:
    - `Ctrl+B` → 현재 절 책갈피 토글
    - `Ctrl+H` → 현재 절 하이라이트 (색상 선택 팝업)
  - 테스트 (TDD):
    - TestBookmarkView: 목록 표시 + 삭제
    - TestReadingBookmarkToggle: Ctrl+B → 북마크 추가/제거

  **Must NOT do**:
  - 내보내기/공유 기능

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
    - Reason: UI 중심 작업
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 10 (with Task 16)
  - **Blocks**: None
  - **Blocked By**: Tasks 13, 14

  **References**:
  - **Internal**: `internal/db/bookmark.go`, `internal/db/highlight.go`
  - **Library**: `bubbles/list` — 목록 위젯
  - **Internal**: `internal/tui/reading.go` — 읽기 화면에 키바인딩 추가

  **Acceptance Criteria**:
  - [ ] `go test ./internal/tui/... -v -run TestBookmark` → PASS
  - [ ] TestBookmarkView: 목록에 bookmark 항목 표시
  - [ ] TestReadingBookmarkToggle: Ctrl+B → DB에 반영

  **Commit**: YES
  - Message: `feat(tui): bookmark and highlight views with toggle keybindings`
  - Files: `internal/tui/bookmarks.go`, `internal/tui/bookmarks_test.go`, 수정: `internal/tui/reading.go`

---

- [ ] 16. TUI Help & Keybindings View

  **What to do**:
  - `internal/tui/help.go`:
    - 모든 키바인딩 목록 표시
    - 각 화면별 사용 가능한 키 설명
    - 스크롤 가능
    - `?` 또는 `Esc`로 닫기
  - 키바인딩 테이블:
    ```
    Global:
      q, Ctrl+C  종료
      ?          도움말
      b          책 목록
      /          검색
      m          책갈피
      p          읽기 계획
      s          설정
      Esc        이전 화면

    읽기 화면:
      h, ←       이전 장
      l, →       다음 장
      j, k       위/아래 스크롤
      Ctrl+B     책갈피 토글
      Ctrl+H     하이라이트
      g          맨 위
      G          맨 아래
    ```

  **Must NOT do**:
  - 키바인딩 커스터마이징

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: 단순 정적 텍스트 뷰
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 10 (with Task 15)
  - **Blocks**: None
  - **Blocked By**: Task 12

  **References**:
  - **Library**: `bubbles/viewport` — 스크롤 가능한 도움말 텍스트
  - **Internal**: `internal/tui/app.go` — 키바인딩 정의

  **Acceptance Criteria**:
  - [ ] `go test ./internal/tui/... -v -run TestHelp` → PASS
  - [ ] TestHelpView: "Ctrl+B" 등 키바인딩 텍스트 포함

  **Commit**: YES
  - Message: `feat(tui): help view with keybinding reference`
  - Files: `internal/tui/help.go`, `internal/tui/help_test.go`

---

### Phase 4: Settings (Theme/Font) + Reading Plans + Polish

---

- [ ] 17. Settings Data Layer (Theme/Font Persistence)

  **What to do**:
  - `internal/config/config.go`:
    - `Config` struct:
      ```go
      type Config struct {
        ThemeName   string  // "dark", "light", "solarized", "nord", "custom"
        FontSize    int     // 1 (small), 2 (medium/default), 3 (large)
        VersionCode string  // 기본 역본 코드 "GAE"
      }
      ```
    - `LoadConfig(db *db.DB) (*Config, error)` — settings 테이블에서 로드, 기본값 반환
    - `SaveConfig(db *db.DB, cfg *Config) error` — settings 테이블에 저장
    - 기본값: dark theme, font size 2, version "GAE"
    - key 매핑: `theme_name`, `font_size`, `default_version`
  - 테스트 (TDD):
    - TestLoadConfig_Defaults: 빈 DB → 기본값 반환
    - TestSaveAndLoadConfig: 저장 → 로드 → 일치
    - TestConfigPersistence: 앱 재시작 시뮬레이션 → 설정 유지

  **Must NOT do**:
  - 테마 렌더링 (Task 19)
  - TUI 뷰 (Task 21)
  - YAML/TOML 설정 파일 (DB settings 테이블 사용)

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: 단순 key-value 저장/로드
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 11 (with Task 18)
  - **Blocks**: Tasks 19, 20, 21
  - **Blocked By**: Phase 3 complete

  **References**:
  - **Internal**: `internal/db/db.go` — `GetSetting()`, `SetSetting()`
  - **Schema**: `settings` 테이블 (key TEXT PK, value TEXT)

  **Acceptance Criteria**:
  - [ ] `go test ./internal/config/... -v -count=1` → ALL PASS
  - [ ] TestLoadConfig_Defaults: ThemeName=="dark", FontSize==2

  **Commit**: YES
  - Message: `feat(config): settings persistence layer for theme, font, and version`
  - Files: `internal/config/config.go`, `internal/config/config_test.go`

---

- [ ] 18. Reading Plan Data Layer

  **What to do**:
  - `internal/db/plans.go`:
    - `CreateSequentialPlan(versionID int64, name string) (int64, error)` — 통독 계획 생성
      - 창세기→요한계시록 순서, 1일 3장 기본 (총 ~394일)
    - `CreateMcCheynePlan(versionID int64, name string) (int64, error)` — 매쿠인 계획
      - 1일 4구간 (구약1 + 구약2 + 신약 + 시편/복음) — 1년 완독
      - 매쿠인 데이터: 하드코딩 or internal data file
    - `CreateCustomPlan(versionID int64, name string, entries []PlanEntry) (int64, error)` — 사용자 정의
    - `GetActivePlan(versionID int64) (*ReadingPlan, error)` — 활성 계획
    - `GetTodayEntries(planID int64) ([]PlanEntry, error)` — 오늘 읽을 분량
    - `MarkEntryCompleted(entryID int64) error`
    - `GetPlanProgress(planID int64) (completed, total int, error)` — 진행률
    - `ListPlans() ([]ReadingPlan, error)` — 전체 계획 목록
    - `DeletePlan(planID int64) error`
  - `internal/bible/mcheyne.go`:
    - 매쿠인 읽기표 365일 데이터
    - `McCheyneSchedule() [][]McCheyneEntry` — [day][sections]
  - 테스트 (TDD):
    - TestCreateSequentialPlan: 생성 → entries 확인 (1189장 ÷ 3 = ~396일)
    - TestMarkCompleted: 완료 표시 → progress 증가
    - TestGetTodayEntries: 오늘 날짜 기준 해당 항목 반환
    - TestMcCheyneSchedule: 365일 분량, 각 날 4구간

  **Must NOT do**:
  - TUI 뷰 (Task 22)
  - 알림/리마인더

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: 읽기 계획 생성 로직 (특히 매쿠인 365일 데이터) 복잡
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 11 (with Task 17)
  - **Blocks**: Task 22
  - **Blocked By**: Phase 3 complete

  **References**:
  - **Data**: Robert Murray M'Cheyne 성경 읽기표 — 365일 × 4구간 (공개 데이터)
  - **Internal**: `internal/bible/books.go` — `AllBooks()` 순서대로 통독 계획 생성
  - **Schema**: `reading_plans` + `reading_plan_entries` 테이블

  **Acceptance Criteria**:
  - [ ] `go test ./internal/db/... -v -run TestPlan` → ALL PASS
  - [ ] `go test ./internal/bible/... -v -run TestMcCheyne` → PASS
  - [ ] TestCreateSequentialPlan: entries > 300
  - [ ] TestMcCheyneSchedule: len == 365, 각 날 len >= 2

  **Commit**: YES
  - Message: `feat(db): reading plan engine with sequential, M'Cheyne, and custom plans`
  - Files: `internal/db/plans.go`, `internal/db/plans_test.go`, `internal/bible/mcheyne.go`, `internal/bible/mcheyne_test.go`

---

- [ ] 19. Theme Engine (Dark/Light/Custom)

  **What to do**:
  - `internal/tui/styles/theme.go` (확장):
    - `Theme` struct:
      ```go
      type Theme struct {
        Name           string
        // Base colors
        Background     lipgloss.Color
        Foreground     lipgloss.Color
        // Text elements
        VerseNumber    lipgloss.Color
        SectionTitle   lipgloss.Color
        FootnoteMarker lipgloss.Color
        // UI elements
        StatusBar      lipgloss.Style
        TitleBar       lipgloss.Style
        SelectedItem   lipgloss.Style
        SearchHighlight lipgloss.Style
        // Highlight colors (for user highlights)
        HighlightColors map[string]lipgloss.Color
        // Bookmarked verse indicator
        BookmarkIndicator lipgloss.Style
      }
      ```
    - 프리셋 테마:
      - `DarkTheme()` — 어두운 배경, 밝은 글자 (기본)
      - `LightTheme()` — 밝은 배경, 어두운 글자
      - `SolarizedTheme()` — Solarized Dark
      - `NordTheme()` — Nord 팔레트
    - `GetTheme(name string) *Theme`
    - `AllThemeNames() []string`
    - TUI 전체에 테마 적용: `AppModel.theme` 변경 → 모든 뷰 스타일 업데이트
  - 기존 TUI 뷰들의 하드코딩된 스타일을 `theme.XxxStyle` 참조로 변경
  - 테스트 (TDD):
    - TestGetTheme: "dark" → non-nil Theme
    - TestAllThemeNames: 4개 테마
    - TestThemeColorsNotZero: 모든 필드에 값 존재

  **Must NOT do**:
  - 사용자 커스텀 컬러 에디터 (프리셋만 제공)
  - 테마 import/export

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
    - Reason: 색상 팔레트 디자인, 시각적 일관성
  - **Skills**: [`frontend-ui-ux`]
    - `frontend-ui-ux`: 테마 디자인은 UI/UX 핵심

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 12 (with Task 20)
  - **Blocks**: Task 21
  - **Blocked By**: Task 17

  **References**:
  - **Library**: `lipgloss` — `lipgloss.Color()`, `lipgloss.AdaptiveColor{}`
  - **Library**: `github.com/muesli/termenv` — `termenv.HasDarkBackground()` for auto-detect
  - **Design**: Nord palette — `#2E3440`, `#3B4252`, `#D8DEE9`, `#88C0D0`
  - **Design**: Solarized — `#002b36`, `#586e75`, `#93a1a1`, `#fdf6e3`
  - **Internal**: `internal/tui/app.go` — `AppModel.theme` 필드
  - **Pattern**: `lipgloss.AdaptiveColor{Light: "#333", Dark: "#EEE"}` for auto theme

  **Acceptance Criteria**:
  - [ ] `go test ./internal/tui/styles/... -v -count=1` → ALL PASS
  - [ ] TestGetTheme("dark"): non-nil, Background != Foreground
  - [ ] TestAllThemeNames: >= 4 테마

  **Commit**: YES
  - Message: `feat(tui): theme engine with dark, light, solarized, and nord presets`
  - Files: `internal/tui/styles/theme.go`, `internal/tui/styles/theme_test.go`, 수정: 기존 TUI 뷰 파일들

---

- [ ] 20. Font Size Engine

  **What to do**:
  - `internal/tui/styles/font.go`:
    - 폰트 "크기" 개념 (터미널에서는 실제 폰트 크기 변경 불가):
      - **Small (1)**: compact layout, 줄간격 최소
      - **Medium (2)**: 기본, 절 사이 1줄 간격
      - **Large (3)**: 넓은 간격, 절 사이 빈 줄, 절 번호 강조 크게
    - `FontSizeConfig` struct:
      ```go
      type FontSizeConfig struct {
        Level         int    // 1, 2, 3
        VersePadding  int    // 절 사이 빈 줄 수
        VerseIndent   int    // 본문 들여쓰기
        NumberWidth   int    // 절 번호 표시 너비
        SectionGap    int    // 소제목 전후 간격
      }
      ```
    - `GetFontSizeConfig(level int) *FontSizeConfig`
    - 읽기 뷰에서 렌더링 시 `FontSizeConfig` 반영
    - `+` / `-` 키로 실시간 크기 조절 (읽기 화면에서)
  - 테스트 (TDD):
    - TestFontSizeConfig: 각 레벨별 설정값 확인
    - TestFontSizeBounds: 0, 4 등 범위 밖 → clamp

  **Must NOT do**:
  - 실제 터미널 폰트 변경 (불가능)
  - 4단계 이상 크기

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: 단순 레이아웃 수치 매핑
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 12 (with Task 19)
  - **Blocks**: Task 21
  - **Blocked By**: Task 17

  **References**:
  - **Internal**: `internal/tui/reading.go` — 읽기 뷰 렌더링에서 config 적용
  - **Internal**: `internal/config/config.go` — FontSize 설정 저장

  **Acceptance Criteria**:
  - [ ] `go test ./internal/tui/styles/... -v -run TestFontSize` → PASS
  - [ ] TestFontSizeConfig: level 1 → VersePadding 0, level 3 → VersePadding 2

  **Commit**: YES
  - Message: `feat(tui): font size engine with 3-level layout density`
  - Files: `internal/tui/styles/font.go`, `internal/tui/styles/font_test.go`

---

- [ ] 21. TUI Settings View

  **What to do**:
  - `internal/tui/settings.go`:
    - 설정 화면 레이아웃:
      ```
      ╔══════════════════════════╗
      ║         설정             ║
      ╠══════════════════════════╣
      ║ 테마:    [◀ Dark    ▶]  ║
      ║ 글자크기: [◀ 보통   ▶]  ║
      ║ 기본역본: [◀ 개역개정▶]  ║
      ╠══════════════════════════╣
      ║ [저장] [취소]            ║
      ╚══════════════════════════╝
      ```
    - 좌우 화살표로 옵션 변경
    - 테마 변경 시 실시간 미리보기 (배경색 즉시 변경)
    - 글자크기 변경 시 미리보기 텍스트 표시
    - Enter/S → 저장 (DB에 persist)
    - Esc → 취소 (변경사항 버림)
    - 역본 선택: 크롤링된 역본 목록에서 선택 (현재 GAE만)
  - 테스트 (TDD):
    - TestSettingsView_Init: 현재 설정값 표시
    - TestSettingsView_ChangeTheme: 좌우 → 테마명 변경
    - TestSettingsView_Save: Enter → DB에 반영

  **Must NOT do**:
  - 키바인딩 커스터마이징
  - 언어 설정

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
    - Reason: 설정 UI 레이아웃, 미리보기, 인터랙션
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 13 (with Task 22)
  - **Blocks**: None
  - **Blocked By**: Tasks 19, 20

  **References**:
  - **Internal**: `internal/config/config.go` — LoadConfig, SaveConfig
  - **Internal**: `internal/tui/styles/theme.go` — GetTheme, AllThemeNames
  - **Internal**: `internal/tui/styles/font.go` — GetFontSizeConfig
  - **Pattern**: 좌우 선택 위젯: 커스텀 Model (option list + index)

  **Acceptance Criteria**:
  - [ ] `go test ./internal/tui/... -v -run TestSettings` → PASS
  - [ ] TestSettingsView_Init: 테마 "dark" 표시
  - [ ] TestSettingsView_Save: 저장 후 LoadConfig → 변경값 일치

  **Commit**: YES
  - Message: `feat(tui): settings view with theme, font size, and version selection`
  - Files: `internal/tui/settings.go`, `internal/tui/settings_test.go`

---

- [ ] 22. TUI Reading Plan View

  **What to do**:
  - `internal/tui/plans.go`:
    - 읽기 계획 관리 화면:
      - 활성 계획 목록
      - 진행률 바 표시 (e.g., `[████████░░] 78%`)
      - 새 계획 생성 (통독/매쿠인/사용자정의 선택)
    - 오늘의 읽기 분량 화면:
      - 오늘 읽을 구간 목록
      - Enter → 해당 구간 읽기 화면으로 이동
      - Space → 완료 체크 토글
      - 진행률 표시
    - 사용자 정의 계획 생성:
      - 시작 책/장, 끝 책/장 선택
      - 기간 설정 (일수)
      - 자동 분배
  - 테스트 (TDD):
    - TestPlanView_List: 계획 목록 표시
    - TestPlanView_Today: 오늘 항목 표시
    - TestPlanView_Complete: Space → 완료 토글

  **Must NOT do**:
  - 알림/푸시
  - 소셜 공유

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
    - Reason: 계획 UI, 진행률 바, 인터랙션 복잡
  - **Skills**: [`frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 13 (with Task 21)
  - **Blocks**: None
  - **Blocked By**: Task 18

  **References**:
  - **Internal**: `internal/db/plans.go` — 계획 CRUD + 진행률
  - **Library**: `bubbles/progress` — 진행률 바 위젯
  - **Library**: `bubbles/list` — 계획 목록, 오늘 읽기 목록
  - **Internal**: `internal/tui/app.go` — StateReading으로 전환

  **Acceptance Criteria**:
  - [ ] `go test ./internal/tui/... -v -run TestPlan` → PASS
  - [ ] TestPlanView_List: 계획 항목 표시
  - [ ] TestPlanView_Complete: Space → completed_at 설정됨

  **Commit**: YES
  - Message: `feat(tui): reading plan view with progress tracking and daily readings`
  - Files: `internal/tui/plans.go`, `internal/tui/plans_test.go`

---

## Commit Strategy

| After Task | Message | Key Files | Pre-commit Verification |
|------------|---------|-----------|-------------------------|
| 1 | `feat: project scaffold with cobra root command and directory structure` | go.mod, main.go, cmd/root.go | `go build ./...` |
| 2 | `feat(db): SQLite schema with multi-version support and CRUD operations` | internal/db/*.go | `go test ./internal/db/...` |
| 3 | `feat(bible): reference parser with Korean name/abbreviation mapping` | internal/bible/*.go | `go test ./internal/bible/...` |
| 4 | `feat(parser): HTML parser for bskorea.or.kr Bible pages with fixture tests` | internal/parser/*.go, testdata/ | `go test ./internal/parser/...` |
| 5 | `feat(crawler): web crawler with rate limiting, checkpointing, and validation` | internal/crawler/*.go | `go test ./internal/crawler/...` |
| 6 | `feat(cli): read and random commands with Korean reference parsing` | cmd/read.go, cmd/random.go | `go test ./cmd/...` |
| 7 | `feat(cli): crawl command with progress display and dry-run mode` | cmd/crawl.go | `go test ./cmd/...` |
| 8 | `feat(db): full-text search with FTS5 and LIKE fallback` | internal/db/search.go | `go test ./internal/db/...` |
| 9 | `feat(db): bookmark and highlight CRUD operations` | internal/db/bookmark.go, highlight.go | `go test ./internal/db/...` |
| 10 | `feat(cli): search command with highlighted results` | cmd/search.go | `go test ./cmd/...` |
| 11 | `feat(cli): bookmark and highlight management commands` | cmd/bookmark.go | `go test ./cmd/...` |
| 12 | `feat(tui): app shell with state machine navigation and keybindings` | internal/tui/app.go, cmd/tui.go | `go test ./internal/tui/...` |
| 13 | `feat(tui): reading view with book/chapter navigation and viewport` | internal/tui/reading.go, booklist.go | `go test ./internal/tui/...` |
| 14 | `feat(tui): search view with text input and highlighted results` | internal/tui/search.go | `go test ./internal/tui/...` |
| 15 | `feat(tui): bookmark and highlight views with toggle keybindings` | internal/tui/bookmarks.go | `go test ./internal/tui/...` |
| 16 | `feat(tui): help view with keybinding reference` | internal/tui/help.go | `go test ./internal/tui/...` |
| 17 | `feat(config): settings persistence layer for theme, font, and version` | internal/config/*.go | `go test ./internal/config/...` |
| 18 | `feat(db): reading plan engine with sequential, M'Cheyne, and custom plans` | internal/db/plans.go, bible/mcheyne.go | `go test ./internal/db/... ./internal/bible/...` |
| 19 | `feat(tui): theme engine with dark, light, solarized, and nord presets` | internal/tui/styles/theme.go | `go test ./internal/tui/styles/...` |
| 20 | `feat(tui): font size engine with 3-level layout density` | internal/tui/styles/font.go | `go test ./internal/tui/styles/...` |
| 21 | `feat(tui): settings view with theme, font size, and version selection` | internal/tui/settings.go | `go test ./internal/tui/...` |
| 22 | `feat(tui): reading plan view with progress tracking and daily readings` | internal/tui/plans.go | `go test ./internal/tui/...` |

---

## Success Criteria

### Verification Commands
```bash
# All tests pass
go test ./... -v -count=1        # Expected: ALL PASS
go test ./... -race               # Expected: no races

# Cross-platform build
GOOS=darwin GOARCH=arm64 go build -o /dev/null .   # Expected: exit 0
GOOS=linux GOARCH=amd64 go build -o /dev/null .    # Expected: exit 0
GOOS=windows GOARCH=amd64 go build -o /dev/null .  # Expected: exit 0

# Static analysis
go vet ./...                      # Expected: clean

# CLI commands
go run . --help                   # Expected: subcommands listed
go run . crawl --dry-run          # Expected: schema created
go run . read 창세기 1            # Expected: 31 verses (after crawl)
go run . search 사랑              # Expected: results (after crawl)
go run . random                   # Expected: 1 verse (after crawl)
go run . tui                      # Expected: TUI launches
```

### Final Checklist
- [ ] All "Must Have" present (크로스플랫폼, 오프라인, 한글렌더링, 다중역본DB, 체크포인트, 무결성검증, 테마, 폰트)
- [ ] All "Must NOT Have" absent (바이너리에 데이터 없음, CGO 없음, len() 너비계산 없음, 불필요한 인터페이스 없음)
- [ ] All tests pass (go test ./... -race)
- [ ] 3 platforms build (darwin/linux/windows)
- [ ] README.md with usage instructions
