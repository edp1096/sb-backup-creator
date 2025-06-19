package main

import (
	"log"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	kernel32        = windows.NewLazySystemDLL("kernel32.dll")
	procCreateMutex = kernel32.NewProc("CreateMutexW")
	procCloseHandle = kernel32.NewProc("CloseHandle")
	mutexHandle     windows.Handle
)

const (
	ERROR_ALREADY_EXISTS = 183
	MUTEX_NAME           = "Global\\SBBackCreatorMutex"
)

// ensureSingleInstance 단일 인스턴스 보장
func ensureSingleInstance() bool {
	mutexName, err := windows.UTF16PtrFromString(MUTEX_NAME)
	if err != nil {
		log.Printf("UTF16 변환 실패: %v", err)
		return false
	}

	// Named Mutex 생성
	ret, _, err := procCreateMutex.Call(
		0,                                  // lpMutexAttributes
		0,                                  // bInitialOwner
		uintptr(unsafe.Pointer(mutexName)), // lpName
	)

	mutexHandle = windows.Handle(ret)

	if mutexHandle == 0 {
		log.Printf("Mutex 생성 실패: %v", err)
		return false
	}

	// 이미 존재하는지 확인
	if err.(syscall.Errno) == ERROR_ALREADY_EXISTS {
		log.Println("이미 실행 중인 인스턴스가 있습니다.")
		showAlreadyRunningMessage()
		return false
	}

	log.Println("단일 인스턴스 확인 완료")
	return true
}

// releaseSingleInstance 단일 인스턴스 해제
func releaseSingleInstance() {
	if mutexHandle != 0 {
		procCloseHandle.Call(uintptr(mutexHandle))
		mutexHandle = 0
	}
}

// showAlreadyRunningMessage 이미 실행 중임을 알리는 메시지
func showAlreadyRunningMessage() {
	// Windows MessageBox 표시
	user32 := windows.NewLazySystemDLL("user32.dll")
	procMessageBox := user32.NewProc("MessageBoxW")

	title, _ := windows.UTF16PtrFromString("SB Backup Creator")
	message, _ := windows.UTF16PtrFromString("이미 실행 중인 SB Backup Creator가 있습니다.\n시스템 트레이를 확인해주세요.")

	procMessageBox.Call(
		0, // hWnd
		uintptr(unsafe.Pointer(message)),
		uintptr(unsafe.Pointer(title)),
		0x30, // MB_ICONWARNING | MB_OK
	)
}
