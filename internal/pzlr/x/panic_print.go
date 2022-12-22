package x

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jpillora/ansi"
	"github.com/maruel/panicparse/v2/stack"
)

func RecoverPanicPrint() {
	r := recover()
	if r == nil {
		return
	}
	fmt.Print(ansi.Red.String("panic:\n"))
	fmt.Printf(ansi.Bright.String("\t%v\n"), r)
	PanicPrint(r)
}

func PanicPrint(r any) {
	const mb = 1 << 20
	buf := make([]byte, 1*mb)
	for i := 0; ; i++ {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			buf = buf[:n]
			break
		}
		l := len(buf) * 2
		buf = make([]byte, l)
	}
	b := bytes.Buffer{}
	b.Write(buf)
	s, _, err := stack.ScanSnapshot(&b, os.Stdout, stack.DefaultOpts())
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	// Find out similar goroutine traces and group them into buckets.
	buckets := s.Aggregate(stack.AnyValue).Buckets
	// Only show stack frames in cwd, and lengthen paths
	for _, bucket := range buckets {
		filtered := []stack.Call{}
		for _, line := range bucket.Signature.Stack.Calls {
			_, err := os.Stat(line.SrcName)
			if err == nil {
				abs, err := filepath.Abs(line.SrcName)
				if err == nil {
					parent := filepath.Dir(filepath.Dir(filepath.Dir(abs))) + string(filepath.Separator)
					line.SrcName = strings.TrimPrefix(abs, parent)
					filtered = append(filtered, line)
				}
			}
		}
		bucket.Signature.Stack.Calls = filtered
	}
	// Calculate alignment.
	srcLen := 0
	pkgLen := 0
	for _, bucket := range buckets {
		for _, line := range bucket.Signature.Stack.Calls {
			if l := len(ansi.Blue.String(fmt.Sprintf("%s:%d", line.SrcName, line.Line))); l > srcLen {
				srcLen = l
			}
			if l := len(filepath.Base(line.Func.ImportPath)); l > pkgLen {
				pkgLen = l
			}
		}
	}
	// Print the goroutine buckets
	for _, bucket := range buckets {
		// Print the goroutine header.
		extra := ""
		if s := bucket.SleepString(); s != "" {
			extra += " [" + s + "]"
		}
		if bucket.Locked {
			extra += " [locked]"
		}
		if len(bucket.CreatedBy.Calls) != 0 {
			extra += fmt.Sprintf(" [Created by %s.%s @ %s:%d]", bucket.CreatedBy.Calls[0].Func.DirName, bucket.CreatedBy.Calls[0].Func.Name, bucket.CreatedBy.Calls[0].SrcName, bucket.CreatedBy.Calls[0].Line)
		}
		fmt.Printf(ansi.Blue.String("%d: %s%s\n"), len(bucket.IDs), bucket.State, extra)
		// Print the stack lines.
		for _, line := range bucket.Stack.Calls {
			fmt.Printf(
				"\t%-*s %-*s %s(%s)\n",
				pkgLen, ansi.Green.String(line.Func.DirName), srcLen,
				ansi.Blue.String(fmt.Sprintf("%s:%d", line.SrcName, line.Line)),
				line.Func.Name, &line.Args)
		}
		if bucket.Stack.Elided {
			io.WriteString(os.Stdout, "    (...)\n")
		}
	}
}
