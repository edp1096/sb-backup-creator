package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	watcher     *fsnotify.Watcher
	watcherDone chan bool
	lastBackup  time.Time
)

func startFileWatcher() {
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Printf("파일 감시자 생성 실패: %v", err)
		return
	}

	watcherDone = make(chan bool)

	// 감시할 파일의 디렉토리 추가
	targetFile := GetConfig().TargetFile
	targetDir := filepath.Dir(targetFile)

	// 디렉토리가 존재하는지 확인
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		log.Printf("대상 디렉토리가 존재하지 않습니다: %s", targetDir)
		// 디렉토리가 생성될 때까지 주기적으로 확인
		go waitForDirectory(targetDir)
		return
	}

	err = watcher.Add(targetDir)
	if err != nil {
		log.Printf("디렉토리 감시 추가 실패: %v", err)
		return
	}

	log.Printf("파일 감시 시작: %s", targetFile)

	go func() {
		defer watcher.Close()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// 대상 파일이 변경된 경우만 처리
				if event.Name == targetFile {
					if event.Op&fsnotify.Write == fsnotify.Write {
						log.Printf("파일 변경 감지: %s", event.Name)

						// 너무 자주 백업하는 것을 방지하기 위한 디바운싱
						if time.Since(lastBackup) > 5*time.Second {
							go performAutoBackup()
							lastBackup = time.Now()
						}
					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("파일 감시 오류: %v", err)

			case <-watcherDone:
				return
			}
		}
	}()
}

func waitForDirectory(targetDir string) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if _, err := os.Stat(targetDir); err == nil {
				log.Printf("대상 디렉토리 생성됨. 파일 감시 시작: %s", targetDir)
				startFileWatcher()
				return
			}
		case <-watcherDone:
			return
		}
	}
}

func stopFileWatcher() {
	if watcherDone != nil {
		close(watcherDone)
	}
	if watcher != nil {
		watcher.Close()
	}
}

func restartFileWatcher() {
	stopFileWatcher()
	time.Sleep(1 * time.Second) // 잠시 대기
	startFileWatcher()
}
