package main

import (
	"fmt"
	"log"
	"runtime"
)

// Method 1: Using runtime.Caller
func b() {
	// Get caller information
	pc, _, _, ok := runtime.Caller(1) // 1 means one step up in the call stack
	if ok {
		// Get function name
		functionName := runtime.FuncForPC(pc).Name()
		log.Printf("b() was called by function: %s\n", functionName)
	}
}

// Method 2: More detailed information about the caller
func b2() {
	// Get caller information
	pc, file, line, ok := runtime.Caller(1)
	if ok {
		// Get function name
		fn := runtime.FuncForPC(pc).Name()
		log.Printf("b2() was called by:\nFunction: %s\nFile: %s\nLine: %d\n", fn, file, line)
	}
}

// Method 3: Print entire call stack
func b3() {
	// Create a buffer of program counters
	const depth = 5
	var pcs [depth]uintptr

	// Skip 1 frame (the frame of b3 itself)
	n := runtime.Callers(1, pcs[:])

	// Get caller frames
	frames := runtime.CallersFrames(pcs[:n])

	log.Println("Call stack:")
	for {
		frame, more := frames.Next()
		log.Printf("Function: %s\nFile: %s\nLine: %d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
}

// Method 4: Custom function with formatted output
func getCallerInfo(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "Unknown caller"
	}

	fn := runtime.FuncForPC(pc).Name()
	return fmt.Sprintf("Called by %s [%s:%d]", fn, file, line)
}

func b4() {
	log.Println(getCallerInfo(1))
}

// Test functions
func a() {
	log.Println("a")
	b()
	b2()
	b3()
	b4()
}

func main() {
	a()
}
