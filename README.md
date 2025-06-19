# SB save file backup tool

Stellar Blade 세이브 파일 백업 도구

## 기능

- **자동 백업**: 세이브 파일이 변경될 때마다 자동으로 백업
- **수동 백업**: 단축키 또는 트레이 메뉴를 통한 즉시 백업
- **날짜별 백업**: 날짜 형식으로 백업 파일 저장 (`StellarBladeSave00_yyyymmdd.sav`)
- **백업 개수 제한**: 오래된 백업 파일 자동 정리
- **시스템 트레이**: 백그라운드에서 조용히 실행
- **설정 관리**: JSON 파일을 통한 유연한 설정
- **단일 인스턴스**: 중복 실행 방지, 하나의 인스턴스만 실행됨

## 사용법

### 첫 실행
1. `sb-backup-creator.exe` 실행
2. 시스템 트레이에 아이콘이 나타남
3. 첫 실행 시 자동으로 설정 파일 생성
4. Steam ID가 자동으로 감지되어 경로 설정

### 트레이 메뉴
- **지금 백업**: 즉시 수동 백업 실행
- **백업 폴더 열기**: 백업 파일들이 저장된 폴더 열기
- **설정**: GUI를 통한 설정 변경
- **설정 파일 편집**: `settings.json` 파일 직접 편집
- **정보**: 프로그램 정보 표시
- **종료**: 프로그램 종료

### 단축키
기본 단축키: `Ctrl + Shift + F9`
- 언제든지 이 단축키를 눌러 날짜 형식 백업 실행
- `settings.json`에서 단축키 변경 가능

## 설정 파일 (settings.json)

* 트레이 메뉴에서 `설정 편집` 클릭
```json
{
  "target_file": "C:\\Users\\USERNAME\\AppData\\Local\\SB\\Saved\\SaveGames\\STEAM_ID\\StellarBladeSave00.sav",
  "backup_dir": "C:\\Users\\USERNAME\\AppData\\Local\\SB\\Backups",
  "hotkey_combo": "ctrl+shift+f9",
  "auto_backup": true,
  "max_backups": 50
}
```

* 항목 설명
    - `target_file`: 백업할 세이브 파일 경로
    - `backup_dir`: 백업 파일들이 저장될 디렉토리
    - `hotkey_combo`: 수동 백업 단축키 (ctrl+shift+b 형식)
    - `auto_backup`: 자동 백업 활성화 여부
    - `max_backups`: 최대 백업 파일 개수 (0은 무제한)

### 단축키 설정 예시
- `ctrl+shift+b`
- `alt+f1`
- `ctrl+alt+s`

## 백업 파일 형식

- **자동 백업**: `StellarBladeSave00_auto.sav` (덮어쓰기, 1개만 유지)
- **수동 백업**: `StellarBladeSave00_20240619_143022.sav` (누적)
- **단축키 백업**: `StellarBladeSave00_20240619_143022.sav` (누적, 수동 백업과 동일)

## 문제 해결

### 중복 실행 시도 시
1. 이미 실행 중이라는 메시지 박스 표시
2. 시스템 트레이에서 기존 인스턴스 확인 가능
3. 새로운 인스턴스는 자동으로 종료됨

### Steam ID 자동 감지 실패
1. 수동으로 Steam ID 확인:
   - `%localappdata%\\SB\\Saved\\SaveGames\\` 폴더 열기
   - 숫자로 된 폴더명이 Steam ID
2. `settings.json`에서 `target_file` 경로 직접 수정

### 백업이 실행되지 않음
1. 대상 파일이 존재하는지 확인
2. 백업 디렉토리 쓰기 권한 확인
3. 프로그램을 관리자 권한으로 실행

### 단축키가 작동하지 않음
1. 다른 프로그램과 단축키 충돌 확인
2. `settings.json`에서 다른 조합으로 변경
3. 프로그램 재시작


## 설치 및 빌드

### 필요 조건
- Go 1.21 이상
- Windows 10/11 (cross-platform 지원하지만 주로 Windows 환경에서 테스트됨)

### 소스 빌드

* Golang 1.24 이상 필요
```
make setup
make
```

## 주의사항

- 바이러스 백신이 오탐지할 수 있습니다
- 이 프로그램은 Stellar Blade 게임과 관련이 없습니다
- 백업은 정기적으로 확인하시기 바랍니다
- 중요한 세이브 파일은 추가로 수동 백업을 권장합니다