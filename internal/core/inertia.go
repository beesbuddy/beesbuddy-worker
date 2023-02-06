package core

import (
	"github.com/beesbuddy/beesbuddy-worker/static"
	"github.com/beesbuddy/beesbuddy-worker/ui/views"
	"github.com/petaki/inertia-go"
	"github.com/petaki/support-go/mix"
)

func NewMixAndInertiaManager(debug bool, url, appName string) (*mix.Mix, *inertia.Inertia, error) {
	mixManager := mix.New(url, "./static", "")

	var version string
	var err error

	if debug {
		version, err = mixManager.Hash("")
		if err != nil {
			return nil, nil, err
		}
	} else {
		version, err = mixManager.HashFromFS("", static.Files)
		if err != nil {
			return nil, nil, err
		}
	}

	inertiaManager := inertia.NewWithFS(
		url,
		"web.gohtml",
		version,
		views.Templates,
	)

	inertiaManager.Share("title", appName)
	inertiaManager.ShareFunc("asset", func(path string) (string, error) {
		if debug {
			return mixManager.Mix(path, "")
		}

		return "/" + path, nil
	})

	return mixManager, inertiaManager, nil
}
