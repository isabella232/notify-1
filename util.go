package notify

import (
	"os"
	"path/filepath"
	"strings"
)

const sep = string(os.PathSeparator)

func nonil(err ...error) error {
	for _, err := range err {
		if err != nil {
			return err
		}
	}
	return nil
}

// canonical resolves any symlink in the given path and returns it in a clean form.
// It expects the path to be absolute. It fails to resolve circular symlinks by
// maintaining a simple iteration limit.
//
// TODO(rjeczalik): replace with realpath?
func canonical(p string) (string, error) {
	for i, depth := 1, 1; i < len(p); i, depth = i+1, depth+1 {
		if depth > 128 {
			return "", &os.PathError{Op: "canonical", Path: p, Err: errDepth}
		}
		if j := strings.IndexRune(p[i:], '/'); j == -1 {
			i = len(p)
		} else {
			i = i + j
		}
		fi, err := os.Lstat(p[:i])
		if err != nil {
			return "", err
		}
		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			s, err := os.Readlink(p[:i])
			if err != nil {
				return "", err
			}
			p = "/" + s + p[i:]
			i = 1 // no guarantee s is canonical, start all over
		}
	}
	return filepath.Clean(p), nil
}

// Joinevents TODO
func joinevents(events []Event) (e Event) {
	if len(events) == 0 {
		e = All
	} else {
		for _, event := range events {
			e |= event
		}
	}
	return
}

func Split(s string) (string, string) {
	if i := LastIndexSep(s); i != -1 {
		return s[:i], s[i+1:]
	}
	return "", s
}

func Base(s string) string {
	if i := LastIndexSep(s); i != -1 {
		return s[i+1:]
	}
	return s
}

// IndexSep TODO
func IndexSep(s string) int {
	for i := 0; i < len(s); i++ {
		if s[i] == os.PathSeparator {
			return i
		}
	}
	return -1
}

// LastIndexSep TODO
func LastIndexSep(s string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == os.PathSeparator {
			return i
		}
	}
	return -1
}
