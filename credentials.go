package git

/*
#include <git2.h>
// TODO: restore these? See comments at (*Cred).Type.
// #include <git2/sys/cred.h>
// git_credtype_t _go_git_cred_credtype(git_cred *cred);
*/
import "C"
import "unsafe"

type CredType uint

const (
	CredTypeUserpassPlaintext CredType = C.GIT_CREDTYPE_USERPASS_PLAINTEXT
	CredTypeSshKey            CredType = C.GIT_CREDTYPE_SSH_KEY
	CredTypeSshCustom         CredType = C.GIT_CREDTYPE_SSH_CUSTOM
	CredTypeDefault           CredType = C.GIT_CREDTYPE_DEFAULT
)

type Cred struct {
	ptr *C.git_cred
}

func (o *Cred) HasUsername() bool {
	if C.git_cred_has_username(o.ptr) == 1 {
		return true
	}
	return false
}

// TODO: restore. This was broken by v28+1,
// I don't see clearly how to fix it,
// and we don't use it.
//
// Relevant changelog entry:
// * The "private" implementation details of the `git_cred` structure have been
//   moved to a dedicated `git2/sys/cred.h` header, to clarify that the underlying
//   structures are only provided for custom transport implementers.
//   The breaking change is that the `username` member of the underlying struct
//   is now hidden, and a new `git_cred_get_username` function has been provided.
//
// Purportedly, this is fixed by these two commits:
//   https://github.com/libgit2/git2go/commit/97e6392d3ab67bbf3e3e59b86a0bc9ebf7430e98
//   https://github.com/libgit2/git2go/commit/c5159e624e55cb14c56a3e5f36200be409fba9d6
// which I have integrated here and in the cgo preamble at the top of this file,
// but it still doesn't work for me. No idea why.
//
// func (o *Cred) Type() CredType {
// 	return (CredType)(C._go_git_cred_credtype(o.ptr))
// }

func credFromC(ptr *C.git_cred) *Cred {
	return &Cred{ptr}
}

func NewCredUserpassPlaintext(username string, password string) (int, Cred) {
	cred := Cred{}
	cusername := C.CString(username)
	defer C.free(unsafe.Pointer(cusername))
	cpassword := C.CString(password)
	defer C.free(unsafe.Pointer(cpassword))
	ret := C.git_cred_userpass_plaintext_new(&cred.ptr, cusername, cpassword)
	return int(ret), cred
}

// NewCredSshKey creates new ssh credentials reading the public and private keys
// from the file system.
func NewCredSshKey(username string, publicKeyPath string, privateKeyPath string, passphrase string) (int, Cred) {
	cred := Cred{}
	cusername := C.CString(username)
	defer C.free(unsafe.Pointer(cusername))
	cpublickey := C.CString(publicKeyPath)
	defer C.free(unsafe.Pointer(cpublickey))
	cprivatekey := C.CString(privateKeyPath)
	defer C.free(unsafe.Pointer(cprivatekey))
	cpassphrase := C.CString(passphrase)
	defer C.free(unsafe.Pointer(cpassphrase))
	ret := C.git_cred_ssh_key_new(&cred.ptr, cusername, cpublickey, cprivatekey, cpassphrase)
	return int(ret), cred
}

// NewCredSshKeyFromMemory creates new ssh credentials using the publicKey and privateKey
// arguments as the values for the public and private keys.
func NewCredSshKeyFromMemory(username string, publicKey string, privateKey string, passphrase string) (int, Cred) {
	cred := Cred{}
	cusername := C.CString(username)
	defer C.free(unsafe.Pointer(cusername))
	cpublickey := C.CString(publicKey)
	defer C.free(unsafe.Pointer(cpublickey))
	cprivatekey := C.CString(privateKey)
	defer C.free(unsafe.Pointer(cprivatekey))
	cpassphrase := C.CString(passphrase)
	defer C.free(unsafe.Pointer(cpassphrase))
	ret := C.git_cred_ssh_key_memory_new(&cred.ptr, cusername, cpublickey, cprivatekey, cpassphrase)
	return int(ret), cred
}

func NewCredSshKeyFromAgent(username string) (int, Cred) {
	cred := Cred{}
	cusername := C.CString(username)
	defer C.free(unsafe.Pointer(cusername))
	ret := C.git_cred_ssh_key_from_agent(&cred.ptr, cusername)
	return int(ret), cred
}

func NewCredDefault() (int, Cred) {
	cred := Cred{}
	ret := C.git_cred_default_new(&cred.ptr)
	return int(ret), cred
}
