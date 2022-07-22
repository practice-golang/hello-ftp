package config

import (
	"os"
	"strings"
)

var ReplacerSlash = strings.NewReplacer("\\", string(os.PathSeparator), "/", string(os.PathSeparator))
