// +build !dynamic

package git

/*
#cgo CFLAGS: -I${SRCDIR}/libgit2/include
#cgo LDFLAGS: -L${SRCDIR}/libgit2/build/ -lgit2
#cgo windows LDFLAGS: -lwinhttp
#cgo !windows pkg-config: --static ${SRCDIR}/libgit2/build/libgit2.pc
#include <git2.h>

#if LIBGIT2_VER_MAJOR != 0 || LIBGIT2_VER_MINOR != 27
# error "Invalid libgit2 version; this git2go supports libgit2 v0.27"
#endif

*/
import "C"
