package embeds

import (
	"fmt"
	"path"
	"strings"
)

type EmbeddedFile struct {
	Path *EmbeddedFilePath
	Data []byte
}

type EmbeddedFilePath struct {
	path   string
	parent string
	base   string
}

func (e *EmbeddedFilePath) Path() string {
	return e.path
}

func (e *EmbeddedFilePath) Parent() string {
	return e.parent
}

func (e *EmbeddedFilePath) Base() string {
	return e.base
}

func ParsePath(str string) (*EmbeddedFilePath, error) {
	if str == "" || !strings.HasPrefix(str, "/") {
		return nil, fmt.Errorf("path must be absolute, got %q", str)
	}
	rel := strings.TrimPrefix(str, "/")
	if rel == "" || rel == "." || strings.Contains(rel, "..") {
		return nil, fmt.Errorf("invalid path %q", str)
	}

	return &EmbeddedFilePath{
		path:   str,
		parent: path.Dir(rel),  // e.g. "etc/postgresql-custom"
		base:   path.Base(rel), // e.g. "postgresql.custom.conf"
	}, nil
}

func MustParsePath(str string) *EmbeddedFilePath {
	if efp, err := ParsePath(str); err != nil {
		panic(err)
	} else {
		return efp
	}
}
