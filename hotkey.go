package main

import (
	"log"
	"strings"

	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

var (
	hotkeyRegistered bool
	currentHotkey    *hotkey.Hotkey
)

func registerHotkeys() {
	config := GetConfig()
	if config.HotkeyCombo == "" {
		log.Println("단축키가 설정되지 않았습니다")
		return
	}

	log.Printf("단축키 등록: %s", config.HotkeyCombo)

	// 단축키 조합 파싱
	modifiers, key := parseHotkeyCombo(config.HotkeyCombo)
	if key == hotkey.Key0 {
		log.Println("잘못된 단축키 형식입니다")
		return
	}

	// mainthread에서 실행
	mainthread.Call(func() {
		// 기존 단축키 해제
		if currentHotkey != nil {
			currentHotkey.Unregister()
		}

		// 새 단축키 등록
		hk := hotkey.New(modifiers, key)
		err := hk.Register()
		if err != nil {
			log.Printf("단축키 등록 실패: %v", err)
			return
		}

		currentHotkey = hk
		hotkeyRegistered = true
		log.Printf("단축키 등록 성공: %v", hk)

		// 고루틴에서 키 이벤트 대기
		go func() {
			for {
				<-hk.Keydown()
				log.Println("백업 단축키 감지")
				go performManualBackup() // 수동 백업과 동일한 로직 사용
			}
		}()
	})
}

func unregisterHotkeys() {
	if hotkeyRegistered && currentHotkey != nil {
		mainthread.Call(func() {
			currentHotkey.Unregister()
			currentHotkey = nil
			hotkeyRegistered = false
		})
	}
}

func parseHotkeyCombo(combo string) ([]hotkey.Modifier, hotkey.Key) {
	// "ctrl+shift+b" -> modifiers: [ModCtrl, ModShift], key: KeyB
	parts := strings.Split(strings.ToLower(combo), "+")
	var modifiers []hotkey.Modifier
	var key hotkey.Key = hotkey.Key0

	for _, part := range parts {
		part = strings.TrimSpace(part)
		switch part {
		case "ctrl", "control":
			modifiers = append(modifiers, hotkey.ModCtrl)
		case "shift":
			modifiers = append(modifiers, hotkey.ModShift)
		case "alt":
			modifiers = append(modifiers, hotkey.ModAlt)
		case "win", "windows", "cmd":
			modifiers = append(modifiers, hotkey.ModWin)
		default:
			// 일반 키 매핑
			key = mapStringToKey(part)
		}
	}

	return modifiers, key
}

func mapStringToKey(s string) hotkey.Key {
	// 알파벳
	switch s {
	case "a":
		return hotkey.KeyA
	case "b":
		return hotkey.KeyB
	case "c":
		return hotkey.KeyC
	case "d":
		return hotkey.KeyD
	case "e":
		return hotkey.KeyE
	case "f":
		return hotkey.KeyF
	case "g":
		return hotkey.KeyG
	case "h":
		return hotkey.KeyH
	case "i":
		return hotkey.KeyI
	case "j":
		return hotkey.KeyJ
	case "k":
		return hotkey.KeyK
	case "l":
		return hotkey.KeyL
	case "m":
		return hotkey.KeyM
	case "n":
		return hotkey.KeyN
	case "o":
		return hotkey.KeyO
	case "p":
		return hotkey.KeyP
	case "q":
		return hotkey.KeyQ
	case "r":
		return hotkey.KeyR
	case "s":
		return hotkey.KeyS
	case "t":
		return hotkey.KeyT
	case "u":
		return hotkey.KeyU
	case "v":
		return hotkey.KeyV
	case "w":
		return hotkey.KeyW
	case "x":
		return hotkey.KeyX
	case "y":
		return hotkey.KeyY
	case "z":
		return hotkey.KeyZ
	// 숫자
	case "0":
		return hotkey.Key0
	case "1":
		return hotkey.Key1
	case "2":
		return hotkey.Key2
	case "3":
		return hotkey.Key3
	case "4":
		return hotkey.Key4
	case "5":
		return hotkey.Key5
	case "6":
		return hotkey.Key6
	case "7":
		return hotkey.Key7
	case "8":
		return hotkey.Key8
	case "9":
		return hotkey.Key9
	// 펑션키
	case "f1":
		return hotkey.KeyF1
	case "f2":
		return hotkey.KeyF2
	case "f3":
		return hotkey.KeyF3
	case "f4":
		return hotkey.KeyF4
	case "f5":
		return hotkey.KeyF5
	case "f6":
		return hotkey.KeyF6
	case "f7":
		return hotkey.KeyF7
	case "f8":
		return hotkey.KeyF8
	case "f9":
		return hotkey.KeyF9
	case "f10":
		return hotkey.KeyF10
	case "f11":
		return hotkey.KeyF11
	case "f12":
		return hotkey.KeyF12
	// 기타
	case "space":
		return hotkey.KeySpace
	case "enter":
		return hotkey.KeyReturn
	case "tab":
		return hotkey.KeyTab
	case "esc", "escape":
		return hotkey.KeyEscape
	default:
		return hotkey.Key0 // 잘못된 키
	}
}

func updateHotkeys() {
	// 기존 단축키 해제
	unregisterHotkeys()

	// 새로운 단축키 등록
	registerHotkeys()
}
