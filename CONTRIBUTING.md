## Taking over an external contributor's PR

For security reasons, maintainers do **not** authorize GitHub Actions to run on external contributors' PRs. Instead, review the contribution, cherry-pick its commits onto a branch in `langchain-ai/langsmith-go`, and open a maintainer-owned PR so CI can run there.

1. Inspect the original PR's complete diff and commit list without approving its workflows. Review code, tests, dependency changes, and any scripts or workflow changes that CI would execute. Moving code to an internal branch does not make it safe; finish this review before pushing.
2. Start from the latest target branch in a clean checkout. Fetch the original PR's head and cherry-pick only the reviewed commits, oldest first:

   ```bash
   git fetch origin
   git switch -c external-pr-<number> origin/<target-branch>
   git fetch origin pull/<number>/head
   git cherry-pick -x <reviewed-commit-sha> [<next-reviewed-commit-sha> ...]
   ```

   Use explicit SHAs from the original PR, not a moving branch head. `-x` records the source commit, and cherry-picking preserves the contributor's authorship. If a cherry-pick conflicts, resolve and review the resulting diff before continuing, or abort with `git cherry-pick --abort`.
3. Check the final diff against the target branch, run the relevant local checks, then push the branch to `origin` (not the contributor's fork) and open a new PR against the original target branch. Link the original PR, credit the contributor, and describe any changes made during the takeover. Agents should use their PR-creation tool when available.
4. Run CI and obtain the usual review on the replacement PR. Do not approve workflows on the original PR or weaken CI security settings. Link the replacement from the original PR; close the original as superseded once the replacement is merged.

If the contributor adds commits later, review and cherry-pick those explicitly as well; do not automatically sync unreviewed updates.

## Setting up the environment

To set up the repository, run:

```sh
$ ./scripts/bootstrap
$ ./scripts/build
```

This will install all the required dependencies and build the SDK.

You can also [install go 1.22+ manually](https://go.dev/doc/install).

## Modifying/Adding code

Most of the SDK is generated code. Modifications to code will be persisted between generations, but may
result in merge conflicts between manual patches and changes from the generator. The generator will never
modify the contents of the `lib/` and `examples/` directories.

## Adding and running examples

All files in the `examples/` directory are not modified by the generator and can be freely edited or added to.

```go
# add an example to examples/<your-example>/main.go

package main

func main() {
  // ...
}
```

```sh
$ go run ./examples/<your-example>
```

## Using the repository from source

To use a local version of this library from source in another project, edit the `go.mod` with a replace
directive. This can be done through the CLI with the following:

```sh
$ go mod edit -replace github.com/langchain-ai/langsmith-go=/path/to/langsmith-go
```

## Running tests

```sh
$ ./scripts/test
```

## Formatting

This library uses the standard gofmt code formatter:

```sh
$ ./scripts/format
```
