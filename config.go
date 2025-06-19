package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//go:embed settings.json
var defaultSettings embed.FS

type Config struct {
	TargetFile  string `json:"target_file"`
	BackupDir   string `json:"backup_dir"`
	HotkeyCombo string `json:"hotkey_combo"`
	AutoBackup  bool   `json:"auto_backup"`
	MaxBackups  int    `json:"max_backups"`
}

var (
	config     *Config
	configPath string
)

func initializeConfig() error {
	// 설정 파일 경로 설정
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("실행 파일 경로 가져오기 실패: %v", err)
	}
	configPath = filepath.Join(filepath.Dir(exePath), "settings.json")

	// 설정 파일이 존재하는지 확인
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 설정 파일이 없으면 기본 설정으로 생성
		if err := createDefaultConfig(); err != nil {
			return fmt.Errorf("기본 설정 파일 생성 실패: %v", err)
		}
	}

	// 설정 파일 로드
	return loadConfig()
}

func createDefaultConfig() error {
	// embed된 기본 설정 읽기
	defaultData, err := defaultSettings.ReadFile("settings.json")
	if err != nil {
		return fmt.Errorf("기본 설정 읽기 실패: %v", err)
	}

	// 기본 설정을 구조체로 파싱
	var defaultConfig Config
	if err := json.Unmarshal(defaultData, &defaultConfig); err != nil {
		return fmt.Errorf("기본 설정 파싱 실패: %v", err)
	}

	// 경로 변수 치환
	defaultConfig.TargetFile = expandPath(defaultConfig.TargetFile)
	defaultConfig.BackupDir = expandPath(defaultConfig.BackupDir)

	// Steam ID 자동 감지 및 경로 수정
	if strings.Contains(defaultConfig.TargetFile, "your_steam_id") {
		steamID, err := detectSteamID()
		if err != nil {
			// Steam ID를 찾지 못하면 기본 경로 사용
			defaultConfig.TargetFile = strings.Replace(defaultConfig.TargetFile, "\\your_steam_id", "", -1)
		} else {
			defaultConfig.TargetFile = strings.Replace(defaultConfig.TargetFile, "your_steam_id", steamID, -1)
		}
	}

	// 설정 파일로 저장
	return saveConfig(&defaultConfig)
}

func loadConfig() error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("설정 파일 읽기 실패: %v", err)
	}

	config = &Config{}
	if err := json.Unmarshal(data, config); err != nil {
		return fmt.Errorf("설정 파일 파싱 실패: %v", err)
	}

	// 경로 변수 치환
	config.TargetFile = expandPath(config.TargetFile)
	config.BackupDir = expandPath(config.BackupDir)

	return nil
}

func saveConfig(cfg *Config) error {
	// 백업 디렉토리 생성
	if err := os.MkdirAll(cfg.BackupDir, 0755); err != nil {
		return fmt.Errorf("백업 디렉토리 생성 실패: %v", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("설정 JSON 생성 실패: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("설정 파일 저장 실패: %v", err)
	}

	config = cfg
	return nil
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "%localappdata%") {
		localAppData := os.Getenv("LOCALAPPDATA")
		return strings.Replace(path, "%localappdata%", localAppData, 1)
	}
	return path
}

func detectSteamID() (string, error) {
	// Steam 설치 경로에서 사용자 ID 찾기
	localAppData := os.Getenv("LOCALAPPDATA")
	sbPath := filepath.Join(localAppData, "SB", "Saved", "SaveGames")

	entries, err := os.ReadDir(sbPath)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "." && entry.Name() != ".." {
			// 숫자로만 구성된 디렉토리 찾기 (Steam ID)
			if isNumeric(entry.Name()) {
				savePath := filepath.Join(sbPath, entry.Name(), "StellarBladeSave00.sav")
				if _, err := os.Stat(savePath); err == nil {
					return entry.Name(), nil
				}
			}
		}
	}

	return "", fmt.Errorf("Steam ID를 찾을 수 없습니다")
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}

func GetConfig() *Config {
	return config
}
