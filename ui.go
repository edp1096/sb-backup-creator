package main

import (
	"log"
	"os/exec"
	"runtime"

	"github.com/sqweek/dialog"
)

func showSettingsDialog() {
	// 간단한 설정 대화상자 표시
	// 실제 구현에서는 더 정교한 GUI를 사용할 수 있습니다

	config := GetConfig()

	// 대상 파일 경로 변경
	if result := dialog.Message("대상 파일 경로를 변경하시겠습니까?\n현재: %s", config.TargetFile).YesNo(); result {
		if newPath, err := dialog.File().Filter("Save Files", "sav").Load(); err == nil {
			config.TargetFile = newPath
		}
	}

	// 백업 디렉토리 변경
	if result := dialog.Message("백업 디렉토리를 변경하시겠습니까?\n현재: %s", config.BackupDir).YesNo(); result {
		if newDir, err := dialog.Directory().Browse(); err == nil {
			config.BackupDir = newDir
		}
	}

	// 자동 백업 토글
	autoBackupText := "비활성화"
	if config.AutoBackup {
		autoBackupText = "활성화"
	}
	if result := dialog.Message("자동 백업을 토글하시겠습니까?\n현재: %s", autoBackupText).YesNo(); result {
		config.AutoBackup = !config.AutoBackup
	}

	// 설정 저장
	if err := saveConfig(config); err != nil {
		dialog.Message("설정 저장에 실패했습니다: %v", err).Error()
		return
	}

	dialog.Message("설정이 저장되었습니다.").Info()

	// 파일 감시자 재시작
	restartFileWatcher()

	// 단축키 업데이트
	updateHotkeys()
}

func openBackupFolder() {
	backupDir := GetConfig().BackupDir

	switch runtime.GOOS {
	case "windows":
		exec.Command("explorer", backupDir).Start()
	case "darwin":
		exec.Command("open", backupDir).Start()
	case "linux":
		exec.Command("xdg-open", backupDir).Start()
	default:
		log.Printf("백업 폴더 열기를 지원하지 않는 운영체제입니다: %s", runtime.GOOS)
	}
}

func openConfigFile() {
	switch runtime.GOOS {
	case "windows":
		exec.Command("notepad", configPath).Start()
	case "darwin":
		exec.Command("open", "-t", configPath).Start()
	case "linux":
		exec.Command("xdg-open", configPath).Start()
	default:
		log.Printf("설정 파일 열기를 지원하지 않는 운영체제입니다: %s", runtime.GOOS)
	}
}

func showAboutDialog() {
	dialog.Message(`SB Backup Creator v0.0.3

Stellar Blade 세이브 파일 자동 백업 도구

기능:
- 파일 변경 시 자동 백업
- 단축키를 통한 수동 백업
- 백업 파일 개수 제한
- 시스템 트레이에서 실행

개발자: 사용자 요청으로 제작`).Title("프로그램 정보").Info()
}
