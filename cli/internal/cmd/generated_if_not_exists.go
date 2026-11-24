package cmd

import (
	"slices"
	"strings"

	"github.com/urfave/cli/v3"
)

// ifNotExistsLookup names the operation that reads the resource a generated
// create operation conflicted with, and the flag that carries its id.
type ifNotExistsLookup struct {
	// op reads one resource, and lookupFlag is its id flag.
	op, lookupFlag string
	// idFlag is the create operation's flag holding the same id.
	idFlag string
}

// generatedIfNotExists lists generated create operations, keyed by
// "<resource> <operation>", that take --if-not-exists. When the API answers
// such a create with 409, the CLI prints the existing resource instead.
var generatedIfNotExists = map[string]ifNotExistsLookup{
	"agents create": {op: "retrieve", lookupFlag: "--agent-id", idFlag: "--id"},
}

const ifNotExistsFlagUsage = "Return the existing resource with this id instead of failing with a conflict"

// addGeneratedIfNotExistsFlags documents --if-not-exists on the operations in
// generatedIfNotExists. runGenerated removes it before the generated tree
// parses the arguments.
func addGeneratedIfNotExistsFlags(resources []*cli.Command) {
	for _, r := range resources {
		for _, op := range r.Commands {
			if _, ok := generatedIfNotExists[r.Name+" "+op.Name]; !ok {
				continue
			}
			// The generated tree is global, and each root command adds to it.
			if slices.ContainsFunc(op.Flags, func(f cli.Flag) bool { return f.Names()[0] == "if-not-exists" }) {
				continue
			}
			op.Flags = append(op.Flags, &cli.BoolFlag{Name: "if-not-exists", Usage: ifNotExistsFlagUsage})
		}
	}
}

// extractGeneratedIfNotExists removes --if-not-exists from args when the
// operation takes it. Elsewhere the flag stays, so the generated tree rejects
// it as unknown.
func extractGeneratedIfNotExists(resource *cli.Command, args []string) ([]string, bool) {
	op := operationName(resource, args, generatedValueFlags())
	if _, ok := generatedIfNotExists[resource.Name+" "+op]; !ok {
		return args, false
	}
	return extractFlag(args, "--if-not-exists")
}

// ifNotExistsLookupArgs turns the arguments of a create operation in
// generatedIfNotExists into those of its lookup: the operation name is
// swapped, the create operation's own flags are dropped, the generated
// root's flags (such as --format) are kept, and the id is passed to the
// lookup's id flag. It reports false when args name no literal id.
func ifNotExistsLookupArgs(resource *cli.Command, args []string) ([]string, bool) {
	valueFlags := generatedValueFlags()
	opName := operationName(resource, args, valueFlags)
	spec, ok := generatedIfNotExists[resource.Name+" "+opName]
	if !ok {
		return nil, false
	}
	idx := slices.IndexFunc(resource.Commands, func(c *cli.Command) bool { return c.Name == opName })
	if idx < 0 {
		return nil, false
	}
	// An id read from a file (@path) or from a piped body is resolved only
	// inside the generated tree, so the lookup can't name it.
	_, id, err := extractValueFlag(args, spec.idFlag)
	if err != nil || id == "" || strings.HasPrefix(id, "@") {
		return nil, false
	}
	opFlags := map[string]bool{} // spelling -> takes a separate value
	for _, f := range resource.Commands[idx].Flags {
		takesValue := true
		if v, ok := f.(interface{ TakesValue() bool }); ok {
			takesValue = v.TakesValue()
		}
		for _, n := range f.Names() {
			if len(n) == 1 {
				opFlags["-"+n] = takesValue
			} else {
				opFlags["--"+n] = takesValue
			}
		}
	}
	var lookup []string
	swapped := false
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			break
		}
		name, _, hasValue := strings.Cut(arg, "=")
		if takesValue, isOpFlag := opFlags[name]; isOpFlag && !valueFlags[name] {
			if takesValue && !hasValue {
				i++
			}
			continue
		}
		if strings.HasPrefix(arg, "-") {
			lookup = append(lookup, arg)
			if valueFlags[arg] && i+1 < len(args) {
				lookup = append(lookup, args[i+1])
				i++
			}
			continue
		}
		if !swapped && arg == opName {
			arg, swapped = spec.op, true
		}
		lookup = append(lookup, arg)
	}
	return append(lookup, spec.lookupFlag+"="+id), true
}
