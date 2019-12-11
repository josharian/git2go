package git

import (
	"runtime"
	"sort"
	"testing"
	"testing/quick"
	"time"
)

func TestRefModification(t *testing.T) {
	t.Parallel()
	repo := createTestRepo(t)
	defer cleanupTestRepo(t, repo)

	commitId, treeId := seedTestRepo(t, repo)

	_, err := repo.References.Create("refs/tags/tree", treeId, true, "testTreeTag")
	checkFatal(t, err)

	tag, err := repo.References.Lookup("refs/tags/tree")
	checkFatal(t, err)
	checkRefType(t, tag, ReferenceOid)

	ref, err := repo.References.Lookup("HEAD")
	checkFatal(t, err)
	checkRefType(t, ref, ReferenceSymbolic)

	if target := ref.Target(); target != nil {
		t.Fatalf("Expected nil *Oid, got %v", target)
	}

	ref, err = ref.Resolve()
	checkFatal(t, err)
	checkRefType(t, ref, ReferenceOid)

	if target := ref.Target(); target == nil {
		t.Fatalf("Expected valid target got nil")
	}

	if target := ref.SymbolicTarget(); target != "" {
		t.Fatalf("Expected empty string, got %v", target)
	}

	if commitId.String() != ref.Target().String() {
		t.Fatalf("Wrong ref target")
	}

	_, err = tag.Rename("refs/tags/renamed", false, "")
	checkFatal(t, err)
	tag, err = repo.References.Lookup("refs/tags/renamed")
	checkFatal(t, err)
	checkRefType(t, ref, ReferenceOid)

}

func TestReferenceIterator(t *testing.T) {
	t.Parallel()
	repo := createTestRepo(t)
	defer cleanupTestRepo(t, repo)

	loc, err := time.LoadLocation("Europe/Berlin")
	checkFatal(t, err)
	sig := &Signature{
		Name:  "Rand Om Hacker",
		Email: "random@hacker.com",
		When:  time.Date(2013, 03, 06, 14, 30, 0, 0, loc),
	}

	idx, err := repo.Index()
	checkFatal(t, err)
	err = idx.AddByPath("README")
	checkFatal(t, err)
	treeId, err := idx.WriteTree()
	checkFatal(t, err)

	message := "This is a commit\n"
	tree, err := repo.LookupTree(treeId)
	checkFatal(t, err)
	commitId, err := repo.CreateCommit("HEAD", sig, sig, message, tree)
	checkFatal(t, err)

	_, err = repo.References.Create("refs/heads/one", commitId, true, "headOne")
	checkFatal(t, err)

	_, err = repo.References.Create("refs/heads/two", commitId, true, "headTwo")
	checkFatal(t, err)

	_, err = repo.References.Create("refs/heads/three", commitId, true, "headThree")
	checkFatal(t, err)

	iter, err := repo.NewReferenceIterator()
	checkFatal(t, err)

	var list []string
	expected := []string{
		"refs/heads/master",
		"refs/heads/one",
		"refs/heads/three",
		"refs/heads/two",
	}

	// test some manual iteration
	nameIter := iter.Names()
	name, err := nameIter.Next()
	for err == nil {
		list = append(list, name)
		name, err = nameIter.Next()
	}
	if !IsErrorCode(err, ErrIterOver) {
		t.Fatal("Iteration not over")
	}

	sort.Strings(list)
	compareStringList(t, expected, list)

	// test the iterator for full refs, rather than just names
	iter, err = repo.NewReferenceIterator()
	checkFatal(t, err)
	count := 0
	_, err = iter.Next()
	for err == nil {
		count++
		_, err = iter.Next()
	}
	if !IsErrorCode(err, ErrIterOver) {
		t.Fatal("Iteration not over")
	}

	if count != 4 {
		t.Fatalf("Wrong number of references returned %v", count)
	}

}

func TestReferenceOwner(t *testing.T) {
	t.Parallel()
	repo := createTestRepo(t)
	defer cleanupTestRepo(t, repo)

	commitId, _ := seedTestRepo(t, repo)

	ref, err := repo.References.Create("refs/heads/foo", commitId, true, "")
	checkFatal(t, err)

	owner := ref.Owner()
	if owner == nil {
		t.Fatal("nil owner")
	}

	if owner.ptr != repo.ptr {
		t.Fatalf("bad ptr, expected %v have %v\n", repo.ptr, owner.ptr)
	}
}

func TestUtil(t *testing.T) {
	t.Parallel()
	repo := createTestRepo(t)
	defer cleanupTestRepo(t, repo)

	commitId, _ := seedTestRepo(t, repo)

	ref, err := repo.References.Create("refs/heads/foo", commitId, true, "")
	checkFatal(t, err)

	ref2, err := repo.References.Dwim("foo")
	checkFatal(t, err)

	if ref.Cmp(ref2) != 0 {
		t.Fatalf("foo didn't dwim to the right thing")
	}

	if ref.Shorthand() != "foo" {
		t.Fatalf("refs/heads/foo has no foo shorthand")
	}

	hasLog, err := repo.References.HasLog("refs/heads/foo")
	checkFatal(t, err)
	if !hasLog {
		t.Fatalf("branches have logs by default")
	}
}

func TestIsNote(t *testing.T) {
	t.Parallel()
	repo := createTestRepo(t)
	defer cleanupTestRepo(t, repo)

	commitID, _ := seedTestRepo(t, repo)

	sig := &Signature{
		Name:  "Rand Om Hacker",
		Email: "random@hacker.com",
		When:  time.Now(),
	}

	refname, err := repo.Notes.DefaultRef()
	checkFatal(t, err)

	_, err = repo.Notes.Create(refname, sig, sig, commitID, "This is a note", false)
	checkFatal(t, err)

	ref, err := repo.References.Lookup(refname)
	checkFatal(t, err)

	if !ref.IsNote() {
		t.Fatalf("%s should be a note", ref.Name())
	}

	ref, err = repo.References.Create("refs/heads/foo", commitID, true, "")
	checkFatal(t, err)

	if ref.IsNote() {
		t.Fatalf("%s should not be a note", ref.Name())
	}
}

func TestReferenceIsValidName(t *testing.T) {
	t.Parallel()
	if !ReferenceIsValidName("HEAD") {
		t.Errorf("HEAD should be a valid reference name")
	}
	if ReferenceIsValidName("HEAD1") {
		t.Errorf("HEAD1 should not be a valid reference name")
	}
}

func compareStringList(t *testing.T, expected, actual []string) {
	for i, v := range expected {
		if actual[i] != v {
			t.Fatalf("Bad list")
		}
	}
}

func checkRefType(t *testing.T, ref *Reference, kind ReferenceType) {
	if ref.Type() == kind {
		return
	}

	// The failure happens at wherever we were called, not here
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		t.Fatalf("Unable to get caller")
	}
	t.Fatalf("Wrong ref type at %v:%v; have %v, expected %v", file, line, ref.Type(), kind)
}

func TestSplitNamespace(t *testing.T) {
	tests := []struct{ Name, Namespace, Rest string }{
		// Normal cases.
		{"b", "", "b"},
		{"b/c/d", "", "b/c/d"},
		{"refs/namespaces/ns/b", "ns", "b"},
		{"refs/namespaces/ns/b/c/d", "ns", "b/c/d"},
		{"refs/namespaces/ns/refs/head", "ns", "refs/head"},
		{"refs/namespaces/ns/refs/namespaces/nested/b", "ns/nested", "b"},
		{"refs/namespaces/ns/refs/namespaces/nested/refs/namespaces/three/b", "ns/nested/three", "b"},
		// Invalid/corner cases.
		{"refs/namespaces/ns/", "ns", ""},
		{"refs/namespaces/ns", "", "refs/namespaces/ns"},
	}
	for _, test := range tests {
		ns, rest := SplitNamespace(test.Name)
		if ns != test.Namespace || rest != test.Rest {
			t.Errorf("SplitNamespace(%q)=(%q, %q), want (%q, %q)", test.Name, ns, rest, test.Namespace, test.Rest)
		}
	}
}

func TestNamespacePrefix(t *testing.T) {
	tests := []struct{ Namespace, Prefix string }{
		// Normal cases.
		{"ns", "refs/namespaces/ns/"},
		{"1/2", "refs/namespaces/1/refs/namespaces/2/"},
		{"", ""},
		// TODO: test some unusual cases?
		// Ideally we'd match how command line git handles things like "1/" and "/2".
		// But I don't want to bother finding out now. Maybe later.
	}
	for _, test := range tests {
		prefix := NamespacePrefix(test.Namespace)
		if prefix != test.Prefix {
			t.Errorf("NamespacePrefix(%q)=%q, want %q", test.Namespace, prefix, test.Prefix)
		}
	}
}

func TestNamespaceRoundTrip(t *testing.T) {
	fromName := func(name string) bool {
		ns, rest := SplitNamespace(name)
		return name == NamespacePrefix(ns)+rest
	}
	fromNS := func(ns, suffix string) bool {
		prefix := NamespacePrefix(ns)
		parsed, rest := SplitNamespace(prefix + suffix)
		return parsed == ns && rest == suffix
	}
	if err := quick.Check(fromName, nil); err != nil {
		t.Error(err)
	}
	if err := quick.Check(fromNS, nil); err != nil {
		t.Error(err)
	}
}
