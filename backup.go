package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func performAutoBackup() {
	if !GetConfig().AutoBackup {
		return
	}

	sourceFile := GetConfig().TargetFile
	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		log.Printf("백업할 파일이 존재하지 않습니다: %s", sourceFile)
		return
	}

	backupDir := GetConfig().BackupDir
	autoBackup0 := filepath.Join(backupDir, "StellarBladeSave00_auto_0.sav")
	autoBackup1 := filepath.Join(backupDir, "StellarBladeSave00_auto_1.sav")

	// 자동 백업 파일 순환 관리
	if err := rotateAutoBackups(autoBackup0, autoBackup1); err != nil {
		log.Printf("자동 백업 순환 실패: %v", err)
		return
	}

	// 새로운 백업을 _auto_0으로 생성
	if err := copyFile(sourceFile, autoBackup0); err != nil {
		log.Printf("자동 백업 실패: %v", err)
		return
	}

	log.Printf("자동 백업 완료: %s", autoBackup0)

	// 자동 백업은 cleanupOldBackups 호출하지 않음 (항상 2개만 유지)
}

// rotateAutoBackups 자동 백업 파일 순환 관리
func rotateAutoBackups(autoBackup0, autoBackup1 string) error {
	// auto_0과 auto_1 파일 존재 확인
	_, err0 := os.Stat(autoBackup0)
	_, err1 := os.Stat(autoBackup1)

	exists0 := err0 == nil
	exists1 := err1 == nil

	if exists0 && exists1 {
		// 둘 다 있을 때: auto_1 삭제, auto_0을 auto_1로 이동
		if err := os.Remove(autoBackup1); err != nil {
			return fmt.Errorf("기존 auto_1 삭제 실패: %v", err)
		}
		if err := os.Rename(autoBackup0, autoBackup1); err != nil {
			return fmt.Errorf("auto_0을 auto_1로 이동 실패: %v", err)
		}
		log.Printf("자동 백업 순환: auto_0 → auto_1")
	} else if exists0 && !exists1 {
		// auto_0만 있을 때: auto_0을 auto_1로 이동
		if err := os.Rename(autoBackup0, autoBackup1); err != nil {
			return fmt.Errorf("auto_0을 auto_1로 이동 실패: %v", err)
		}
		log.Printf("자동 백업 순환: auto_0 → auto_1")
	}
	// exists1만 있거나 둘 다 없는 경우는 그대로 진행 (새로운 auto_0 생성)

	return nil
}

func performManualBackup() {
	sourceFile := GetConfig().TargetFile
	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		log.Printf("백업할 파일이 존재하지 않습니다: %s", sourceFile)
		return
	}

	backupDir := GetConfig().BackupDir
	now := time.Now()
	backupFileName := fmt.Sprintf("StellarBladeSave00_%s.sav", now.Format("20060102_150405"))
	backupPath := filepath.Join(backupDir, backupFileName)

	if err := copyFile(sourceFile, backupPath); err != nil {
		log.Printf("수동 백업 실패: %v", err)
		return
	}

	log.Printf("수동 백업 완료: %s", backupPath)

	// 오래된 백업 파일 정리
	cleanupOldBackups()
}

func copyFile(src, dst string) error {
	// 백업 디렉토리 생성
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("백업 디렉토리 생성 실패: %v", err)
	}

	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("원본 파일 열기 실패: %v", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("백업 파일 생성 실패: %v", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("파일 복사 실패: %v", err)
	}

	return nil
}

func cleanupOldBackups() {
	maxBackups := GetConfig().MaxBackups
	if maxBackups <= 0 {
		return
	}

	backupDir := GetConfig().BackupDir
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		log.Printf("백업 디렉토리 읽기 실패: %v", err)
		return
	}

	// 백업 파일만 필터링 (자동 백업 파일 제외)
	var backupFiles []os.DirEntry
	for _, entry := range entries {
		if !entry.IsDir() &&
			strings.HasPrefix(entry.Name(), "StellarBladeSave00_") &&
			strings.HasSuffix(entry.Name(), ".sav") &&
			!strings.Contains(entry.Name(), "_auto_") { // 자동 백업 파일 제외
			backupFiles = append(backupFiles, entry)
		}
	}

	// 파일 개수가 최대값을 초과하지 않으면 정리하지 않음
	if len(backupFiles) <= maxBackups {
		return
	}

	// 파일명으로 정렬 (날짜 기준)
	sort.Slice(backupFiles, func(i, j int) bool {
		return backupFiles[i].Name() < backupFiles[j].Name()
	})

	// 오래된 파일 삭제
	filesToDelete := len(backupFiles) - maxBackups
	for i := 0; i < filesToDelete; i++ {
		filePath := filepath.Join(backupDir, backupFiles[i].Name())
		if err := os.Remove(filePath); err != nil {
			log.Printf("오래된 백업 파일 삭제 실패: %v", err)
		} else {
			log.Printf("오래된 백업 파일 삭제: %s", backupFiles[i].Name())
		}
	}
}
