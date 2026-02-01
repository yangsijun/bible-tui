# Bible TUI

터미널에서 성경을 읽고, 검색하고, 책갈피하는 프로그램입니다.

Go + [Bubbletea](https://github.com/charmbracelet/bubbletea) 기반의 크로스플랫폼 TUI/CLI 앱으로, 대한성서공회(bskorea.or.kr)에서 크롤링한 성경 데이터를 SQLite에 저장하여 오프라인으로 사용합니다.

## 주요 기능

- 성경 전체 66권 크롤링 및 오프라인 읽기
- 인터랙티브 TUI 모드 (책 목록 / 장 선택 / 읽기 뷰)
- 전문 검색 (FTS5) + 성경 구절 참조 검색 (`창세기 1`, `창 3:3`)
- 책갈피 & 하이라이트
- 읽기 계획 (통독 / 매쿠인 1년 완독)
- 테마 (Dark / Light / Solarized / Nord)
- 글자 크기 3단계 조절
- CLI 명령어 (read, search, random)
- 크롤링 중단/재개 (체크포인트)
- 단일 바이너리, CGO 불필요

## 설치

### GitHub Releases (추천)

[Releases 페이지](https://github.com/yangsijun/bible-tui/releases)에서 OS에 맞는 바이너리를 다운로드합니다.

Releases 페이지에서 OS/아키텍처에 맞는 파일을 다운로드한 뒤:

```bash
# macOS / Linux
tar xzf bible_*_*.tar.gz
sudo mv bible /usr/local/bin/

# Windows
# zip 파일 압축 해제 후 PATH에 추가
```

### Go Install

```bash
go install github.com/sijun-dong/bible-tui@latest
```

### 소스 빌드

```bash
git clone https://github.com/yangsijun/bible-tui.git
cd bible-tui
go build -o bible .
```

## 시작하기

### 1. 성경 데이터 크롤링

최초 1회 크롤링이 필요합니다. 바이너리에 성경 데이터가 포함되어 있지 않습니다.

```bash
bible crawl
```

특정 책만 크롤링:

```bash
bible crawl --book gen    # 창세기만
bible crawl --book mat    # 마태복음만
```

크롤링 중 중단해도 다시 실행하면 이어서 진행됩니다.

### 2. TUI 모드

```bash
bible tui
```

### 3. CLI 명령어

```bash
bible read 창세기 1       # 창세기 1장
bible read 창 1:3-5      # 창세기 1장 3~5절
bible search 사랑         # "사랑" 검색
bible random              # 랜덤 구절
bible bookmark list       # 책갈피 목록
bible highlight list      # 하이라이트 목록
```

## TUI 키바인딩

### 전역

| 키 | 기능 |
|---|---|
| `q`, `Ctrl+C` | 종료 |
| `?` | 도움말 |
| `b` | 책 목록 |
| `/` | 검색 |
| `m` | 책갈피/하이라이트 |
| `s` | 설정 |
| `p` | 읽기 계획 |
| `Esc` | 이전 화면 |

### 책 목록

| 키 | 기능 |
|---|---|
| `j`, `k`, `↑`, `↓` | 위/아래 이동 |
| `Enter` | 책 선택 |
| `/` | 필터 검색 |

### 읽기 화면

| 키 | 기능 |
|---|---|
| `j`, `k` | 위/아래 스크롤 |
| `h`, `←` | 이전 장 |
| `l`, `→` | 다음 장 |
| `g` | 맨 위 |
| `G` | 맨 아래 |

### 검색

| 키 | 기능 |
|---|---|
| `Enter` | 검색 실행 / 결과 선택 |
| `j`, `k` | 결과 탐색 |

검색창에서 구절 참조도 입력 가능합니다:

```
창세기        → 창세기 1장으로 이동
창 3:3       → 창세기 3장 3절로 이동
요한복음 3:16 → 요한복음 3장 16절로 이동
```

### 책갈피/하이라이트

| 키 | 기능 |
|---|---|
| `Tab` | 책갈피/하이라이트 탭 전환 |
| `j`, `k` | 위/아래 이동 |
| `d` | 삭제 |
| `Enter` | 해당 구절로 이동 |

### 읽기 계획

| 키 | 기능 |
|---|---|
| `n` | 새 계획 생성 |
| `Enter` | 오늘 읽기 분량 보기 |
| `Space` | 완료 체크 |
| `d` | 계획 삭제 |

## 크롤링 옵션

```bash
bible crawl                    # 전체 크롤링
bible crawl --book gen         # 특정 책만
bible crawl --dry-run          # DB 스키마만 생성
bible crawl --reset            # 데이터 삭제 후 재크롤링
bible crawl --reset --book gen # 특정 책만 재크롤링
```

## 테마

설정 화면(`s`)에서 테마를 변경할 수 있습니다:

- **Dark** (기본) — Tokyo Night 스타일 어두운 테마
- **Light** — 밝은 배경
- **Solarized** — Solarized Dark
- **Nord** — Nord 팔레트

## 데이터 저장 위치

- macOS: `~/Library/Application Support/bible-tui/bible.db`
- Linux: `~/.config/bible-tui/bible.db`
- Windows: `%AppData%\bible-tui\bible.db`

SQLite 단일 파일로 성경 데이터, 책갈피, 하이라이트, 읽기 계획, 설정이 모두 저장됩니다.

## 개발

```bash
# 테스트
go test ./... -race -count=1

# 빌드
go build -o bible .

# 정적 분석
go vet ./...
```

## 라이선스

MIT
