package errorsx

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func printSourceContext(
	path string,
	line int,
	col int,
	msg string,
	phase Phase,
) {
	redBold := color.New(color.FgRed, color.Bold).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	gray := color.New(color.FgHiBlack).SprintFunc()

	// If no path provided, fall back to header-only error
	if path == "" {
		fmt.Printf(
			"%s %s\n",
			redBold(fmt.Sprintf("[%s Error]:", phase)),
			gray(msg),
		)
		return
	}

	f, err := os.Open(path)
	if err != nil {
		// Graceful fallback instead of exit
		fmt.Printf(
			"%s:%d:%d\n%s %s\n",
			path,
			line,
			col,
			redBold(fmt.Sprintf("[%s Error]:", phase)),
			gray(msg),
		)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	current := 1
	var srcLine string

	for scanner.Scan() {
		if current == line {
			srcLine = scanner.Text()
			break
		}
		current++
	}

	// If line not found, fall back gracefully
	if srcLine == "" {
		fmt.Printf(
			"%s:%d:%d\n%s %s\n",
			path,
			line,
			col,
			redBold(fmt.Sprintf("[%s Error]:", phase)),
			gray(msg),
		)
		return
	}

	// Header
	fmt.Printf("%s:%d:%d\n", path, line, col)

	// Source line
	fmt.Printf("  %s\n", srcLine)

	// Caret underline
	if col < 1 {
		col = 1
	}
	prefix := strings.Repeat(" ", col-1)
	fmt.Printf("  %s%s\n", prefix, red("^"))

	// Error message
	fmt.Printf(
		"%s %s\n",
		redBold(fmt.Sprintf("[%s Error]:", phase)),
		gray(msg),
	)
}
