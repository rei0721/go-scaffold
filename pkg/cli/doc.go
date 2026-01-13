/*
Package cli provides a general-purpose CLI framework for building
command-line tools with standardized structure, type-safe parameter
binding, and testability-first design.

# Features

  - Interface-oriented design for DI and mocking
  - Type-safe flag parsing with automatic type conversion
  - Testable with custom I/O (RunWithIO)
  - Standard error codes following Unix conventions
  - Automatic help generation
  - Environment variable fallback for flags

# Usage

Basic example:

	app := cli.NewApp("mytool")
	app.SetVersion("1.0.0")
	app.SetDescription("My CLI tool")

	app.AddCommand(&MyCommand{})

	if err := app.Run(os.Args[1:]); err != nil {
	    os.Exit(cli.GetExitCode(err))
	}

Defining a command:

	type GenerateCommand struct{}

	func (c *GenerateCommand) Name() string {
	    return "generate"
	}

	func (c *GenerateCommand) Description() string {
	    return "Generate code from templates"
	}

	func (c *GenerateCommand) Usage() string {
	    return "generate --model=<name> [--output=<dir>]"
	}

	func (c *GenerateCommand) Flags() []cli.Flag {
	    return []cli.Flag{
	        {
	            Name:        "model",
	            Type:        cli.FlagTypeString,
	            Required:    true,
	            Description: "Model name to generate",
	        },
	        {
	            Name:        "output",
	            Type:        cli.FlagTypeString,
	            Default:     "./models",
	            Description: "Output directory",
	            EnvVar:      "OUTPUT_DIR",
	        },
	        {
	            Name:        "force",
	            ShortName:   "f",
	            Type:        cli.FlagTypeBool,
	            Default:     false,
	            Description: "Overwrite existing files",
	        },
	    }
	}

	func (c *GenerateCommand) Execute(ctx *cli.Context) error {
	    model := ctx.GetString("model")
	    output := ctx.GetString("output")
	    force := ctx.GetBool("force")

	    fmt.Fprintf(ctx.Stdout, "Generating %s to %s (force=%v)\n", model, output, force)
	    return nil
	}

# Testing

Commands are testable using custom I/O:

	func TestGenerateCommand(t *testing.T) {
	    app := cli.NewApp("test")
	    app.AddCommand(&GenerateCommand{})

	    var stdout bytes.Buffer
	    err := app.RunWithIO(
	        []string{"generate", "--model", "User"},
	        nil,
	        &stdout,
	        io.Discard,
	    )

	    if err != nil {
	        t.Fatalf("unexpected error: %v", err)
	    }

	    output := stdout.String()
	    if !strings.Contains(output, "User") {
	        t.Errorf("expected output to contain 'User', got: %s", output)
	    }
	}

# Error Handling

The package defines standard error types with exit codes:

  - UsageError: Exit code 2 (parameter errors)
  - CommandError: Exit code 1 (execution errors)
  - CancelledError: Exit code 130 (user interruption)

Extract exit code from error:

	if err := app.Run(os.Args[1:]); err != nil {
	    fmt.Fprintln(os.Stderr, err)
	    os.Exit(cli.GetExitCode(err))
	}
*/
package cli
