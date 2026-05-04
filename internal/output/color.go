package output

import "fmt"

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

func red(s string) string    { return colorRed + s + colorReset }
func green(s string) string  { return colorGreen + s + colorReset }
func yellow(s string) string { return colorYellow + s + colorReset }
func cyan(s string) string   { return colorCyan + s + colorReset }
func bold(s string) string   { return colorBold + s + colorReset }

func printMissing(env, key string) {
	fmt.Printf("  %s  %-30s  %s\n", red("MISSING"), cyan(key), yellow(env))
}

func printMismatch(key string, envVals map[string]string) {
	fmt.Printf("  %s  %s\n", yellow("MISMATCH"), cyan(key))
	for env, val := range envVals {
		fmt.Printf("      %-20s = %s\n", green(env), val)
	}
}

func printMatch(key string) {
	fmt.Printf("  %s  %s\n", green("OK     "), cyan(key))
}

func printHeader(envNames []string) {
	fmt.Println(bold("\nenvdiff — environment comparison"))
	fmt.Printf("Comparing: %s\n\n", cyan(fmt.Sprintf("%v", envNames)))
}
