package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// executePS executes a PowerShell script with the provided arguments without using temporary files
func executePS(script string, args string) error {
	// Parse arguments into a slice
	var argsList []string
	if args != "" {
		argsList = strings.Split(args, " ")
	}

	// Create a PowerShell command that:
	// 1. Sets the execution policy to bypass
	// 2. Uses a ScriptBlock to execute the script content directly
	// 3. Passes the arguments to the script

	// Build argument string to pass to the script
	argsString := ""
	if len(argsList) > 0 {
		argsString = " " + strings.Join(argsList, " ")
	}

	// Create the command that executes the script with arguments
	cmdArgs := []string{
		"-ExecutionPolicy", "Bypass",
		"-NoProfile",
		"-Command",
		// We use a scriptblock to execute the script content directly
		// & ([scriptblock]::Create($scriptContent)) arg1 arg2 ...
		fmt.Sprintf("& ([scriptblock]::Create(@'\n%s\n'@))%s", script, argsString),
	}

	// Create command to execute PowerShell
	cmd := exec.Command("powershell", cmdArgs...)

	// Set up pipes for stdout and stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute PowerShell script: %w", err)
	}

	return nil
}

func main() {
	// Create a simple hello world PowerShell script that uses two arguments
	helloScript := `
param(
    [Parameter(Position=0)]
    [string]$Name = "World",
    
    [Parameter(Position=1)]
    [string]$Greeting = "Hello"
)

Write-Output "$Greeting, $Name from PowerShell!"
`

	// Execute the script with arguments
	fmt.Println("Running PowerShell script with arguments:")
	if err := executePS(helloScript, "Gopher Amazing"); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Execute with different arguments
	fmt.Println("\nRunning again with different arguments:")
	if err := executePS(helloScript, "Developer Greetings"); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
