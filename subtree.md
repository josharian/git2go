This fork of git2go uses subtrees to pull in libgit2.
This is because the code must all be present in the tree
in order for Go modules to work correctly with it.
See https://github.com/golang/go/issues/34867.

To update to a new version of libgit2, run:

$ git subtree pull --prefix libgit2 --squash libgit2 DESIRED_LIBGIT2_COMMIT

TODO: figure out whether we can use https://github.com/apenwarr/git-subtrac
