package git

/*
#include <git2.h>
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type Patch struct {
	ptr *C.git_patch
}

func newPatchFromC(ptr *C.git_patch) *Patch {
	if ptr == nil {
		return nil
	}

	patch := &Patch{
		ptr: ptr,
	}

	runtime.SetFinalizer(patch, (*Patch).Free)
	return patch
}

func (patch *Patch) Free() error {
	if patch.ptr == nil {
		return ErrInvalid
	}
	runtime.SetFinalizer(patch, nil)
	C.git_patch_free(patch.ptr)
	patch.ptr = nil
	return nil
}

func (patch *Patch) String() (string, error) {
	if patch.ptr == nil {
		return "", ErrInvalid
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var buf C.git_buf

	ecode := C.git_patch_to_buf(&buf, patch.ptr)
	runtime.KeepAlive(patch)
	if ecode < 0 {
		return "", MakeGitError(ecode)
	}
	defer C.git_buf_free(&buf)

	return C.GoString(buf.ptr), nil
}

// LineStats reports the number of context lines, added lines, and deleted lines in the patch.
// It is useful for generating diff --numstat type of output.
// See https://libgit2.org/libgit2/#HEAD/group/patch/git_patch_line_stats.
func (patch *Patch) LineStats() (ctxt, additions, deletions int, err error) {
	if patch.ptr == nil {
		return 0, 0, 0, ErrInvalid
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var c, a, d C.size_t
	ret := C.git_patch_line_stats(&c, &a, &d, patch.ptr)
	runtime.KeepAlive(patch)
	if ret < 0 {
		return 0, 0, 0, MakeGitError(ret)
	}
	return int(c), int(a), int(d), nil
}

func toPointer(data []byte) (ptr unsafe.Pointer) {
	if len(data) > 0 {
		ptr = unsafe.Pointer(&data[0])
	}
	return
}

func (v *Repository) PatchFromBuffers(oldPath, newPath string, oldBuf, newBuf []byte, opts *DiffOptions) (*Patch, error) {
	var patchPtr *C.git_patch

	oldPtr := toPointer(oldBuf)
	newPtr := toPointer(newBuf)

	cOldPath := C.CString(oldPath)
	defer C.free(unsafe.Pointer(cOldPath))

	cNewPath := C.CString(newPath)
	defer C.free(unsafe.Pointer(cNewPath))

	copts, _ := diffOptionsToC(opts)
	defer freeDiffOptions(copts)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ecode := C.git_patch_from_buffers(&patchPtr, oldPtr, C.size_t(len(oldBuf)), cOldPath, newPtr, C.size_t(len(newBuf)), cNewPath, copts)
	runtime.KeepAlive(oldBuf)
	runtime.KeepAlive(newBuf)
	if ecode < 0 {
		return nil, MakeGitError(ecode)
	}
	return newPatchFromC(patchPtr), nil
}
