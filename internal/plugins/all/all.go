package all

import (
	// Plugins need to register themselves
	_ "github.com/drewstinnett/labdoc/internal/plugins/builtin"
	_ "github.com/drewstinnett/labdoc/internal/plugins/gitlab"
)
