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
	backupFileName := "StellarBladeSave00_auto.sav"
	backupPath := filepath.Join(backupDir, backupFileName)

	if err := copyFile(sourceFile, backupPath); err != nil {
		log.Printf("자동 백업 실패: %v", err)
		return
	}

	log.Printf("자동 백업 완료: %s", backupPath)

	// 오래된 백업 파일 정리
	cleanupOldBackups()
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

	// 백업 파일만 필터링
	var backupFiles []os.DirEntry
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), "StellarBladeSave00_") && strings.HasSuffix(entry.Name(), ".sav") {
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
