package git

/*
#include <git2.h>

int _go_git_opts_get_search_path(int level, git_buf *buf)
{
    return git_libgit2_opts(GIT_OPT_GET_SEARCH_PATH, level, buf);
}

int _go_git_opts_set_search_path(int level, const char *path)
{
    return git_libgit2_opts(GIT_OPT_SET_SEARCH_PATH, level, path);
}

int _go_git_opts_set_size_t(int opt, size_t val)
{
    return git_libgit2_opts(opt, val);
}

int _go_git_opts_get_size_t(int opt, size_t *val)
{
    return git_libgit2_opts(opt, val);
}

int _go_git_opts_set_uint(int opt, unsigned int val)
{
    return git_libgit2_opts(opt, val);
}

int _go_git_opts_get_uint(int opt, unsigned int *val)
{
    return git_libgit2_opts(opt, val);
}

int _go_git_opts_set_int(int opt, int val)
{
    return git_libgit2_opts(opt, val);
}

int _go_git_opts_get_int(int opt, int *val)
{
    return git_libgit2_opts(opt, val);
}
*/
import "C"
import (
	"runtime"
	"unsafe"
)

func SearchPath(level ConfigLevel) (string, error) {
	var buf C.git_buf
	defer C.git_buf_free(&buf)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	err := C._go_git_opts_get_search_path(C.int(level), &buf)
	if err < 0 {
		return "", MakeGitError(err)
	}

	return C.GoString(buf.ptr), nil
}

func SetSearchPath(level ConfigLevel, path string) error {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	err := C._go_git_opts_set_search_path(C.int(level), cpath)
	if err < 0 {
		return MakeGitError(err)
	}

	return nil
}

func MwindowSize() (int, error) {
	return getSizet(C.GIT_OPT_GET_MWINDOW_SIZE)
}

func SetMwindowSize(size int) error {
	return setSizet(C.GIT_OPT_SET_MWINDOW_SIZE, size)
}

func MwindowMappedLimit() (int, error) {
	return getSizet(C.GIT_OPT_GET_MWINDOW_MAPPED_LIMIT)
}

func SetMwindowMappedLimit(size int) error {
	return setSizet(C.GIT_OPT_SET_MWINDOW_MAPPED_LIMIT, size)
}

func MwindowOpenLimit() (int, error) {
	return getSizet(C.GIT_OPT_GET_MWINDOW_FILE_LIMIT)
}

func SetMwindowOpenLimit(size int) error {
	return setSizet(C.GIT_OPT_SET_MWINDOW_FILE_LIMIT, size)
}

func EnableHTTPExpectContinue() (bool, error) {
	i, err := getInt(C.GIT_OPT_ENABLE_HTTP_EXPECT_CONTINUE)
	if err != nil {
		return false, err
	}
	return i != 0, nil
}

func SetEnableHTTPExpectContinue(enable bool) error {
	i := 0
	if enable {
		i = 1
	}
	return setInt(C.GIT_OPT_ENABLE_HTTP_EXPECT_CONTINUE, i)
}

func getSizet(opt C.int) (int, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var val C.size_t
	err := C._go_git_opts_get_size_t(opt, &val)
	if err < 0 {
		return 0, MakeGitError(err)
	}

	return int(val), nil
}

func setSizet(opt C.int, val int) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	cval := C.size_t(val)
	err := C._go_git_opts_set_size_t(opt, cval)
	if err < 0 {
		return MakeGitError(err)
	}

	return nil
}

func getUint(opt C.int) (uint, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var val C.uint
	err := C._go_git_opts_get_uint(opt, &val)
	if err < 0 {
		return 0, MakeGitError(err)
	}
	return uint(val), nil
}

func setUint(opt C.int, val uint) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	cval := C.uint(val)
	err := C._go_git_opts_set_uint(opt, cval)
	if err < 0 {
		return MakeGitError(err)
	}
	return nil
}

func getInt(opt C.int) (int, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var val C.int
	err := C._go_git_opts_get_int(opt, &val)
	if err < 0 {
		return 0, MakeGitError(err)
	}
	return int(val), nil
}

func setInt(opt C.int, val int) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	cval := C.int(val)
	err := C._go_git_opts_set_int(opt, cval)
	if err < 0 {
		return MakeGitError(err)
	}
	return nil
}
