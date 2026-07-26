package update

import "watsap/utils/config"

func WatsapUpdate() {
	if config.UPDATE_URL == "" {
		return
	}
	compare()
}
