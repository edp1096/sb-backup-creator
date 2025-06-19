package main

import (
	"log"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"golang.design/x/hotkey/mainthread"
)

func main() {
	// 단일 인스턴스 확인
	if !ensureSingleInstance() {
		return // 이미 실행 중이면 종료
	}
	defer releaseSingleInstance()

	// 콘솔창 숨기기 (Windows)
	hideConsole()

	// mainthread에서 실행 (hotkey 라이브러리 요구사항)
	mainthread.Init(func() {
		systray.Run(onReady, onExit)
	})
}

func onReady() {
	// 트레이 아이콘 설정
	systray.SetIcon(icon.Data)
	systray.SetTitle("SB Backup Creator")
	systray.SetTooltip("Stellar Blade Save Backup Tool")

	// 설정 초기화
	if err := initializeConfig(); err != nil {
		log.Printf("설정 초기화 실패: %v", err)
		return
	}

	// 파일 감시 시작
	go startFileWatcher()

	// 단축키 등록
	go registerHotkeys()

	// 메뉴 아이템 생성
	mBackupNow := systray.AddMenuItem("지금 백업", "수동 백업 실행")
	mOpenBackup := systray.AddMenuItem("백업 폴더 열기", "백업 파일들이 저장된 폴더 열기")
	systray.AddSeparator()
	// mSettings := systray.AddMenuItem("설정", "설정 변경")
	// mConfigFile := systray.AddMenuItem("설정 파일 편집", "settings.json 파일 직접 편집")
	mConfigFile := systray.AddMenuItem("설정 편집", "settings.json 파일 직접 편집")
	systray.AddSeparator()
	// mAbout := systray.AddMenuItem("정보", "프로그램 정보")
	mExit := systray.AddMenuItem("종료", "프로그램 종료")

	// 메뉴 이벤트 처리
	go func() {
		for {
			select {
			case <-mBackupNow.ClickedCh:
				performManualBackup()
			case <-mOpenBackup.ClickedCh:
				openBackupFolder()
			// case <-mSettings.ClickedCh:
			// 	showSettingsDialog()
			case <-mConfigFile.ClickedCh:
				openConfigFile()
			// case <-mAbout.ClickedCh:
			// 	showAboutDialog()
			case <-mExit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	cleanup()
}

func hideConsole() {
	// Windows에서 콘솔창 숨기기
	// go build -ldflags "-H windowsgui" 로 빌드하면 콘솔창이 안 뜹니다
}

func cleanup() {
	stopFileWatcher()
	unregisterHotkeys()
}
