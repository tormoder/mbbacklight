package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

const (
	kdbFileBrightness    = "/sys/class/leds/spi::kbd_backlight/brightness"
	kdbFileMaxBrightness = "/sys/class/leds/spi::kbd_backlight/max_brightness"

	screenFileBrightness    = "/sys/class/backlight/gmux_backlight/brightness"
	screenFileMaxBrightness = "/sys/class/backlight/gmux_backlight/max_brightness"

	defaultKdbStep    = 25
	defaultScreenStep = 25
)

func main() {
	fstep := flag.Int("step", 0, "step value for command")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: mbbacklight [flags] [system] [operation]\n\n")
		fmt.Fprintf(os.Stderr, "Systems:\n")
		fmt.Fprintf(os.Stderr, "  -kbd\n  -screen\n\n")
		fmt.Fprintf(os.Stderr, "Operations:\n")
		fmt.Fprintf(os.Stderr, "  -get\n  -up\n  -down\n  -max\n  -set [value]\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	var (
		system    string
		operation string
		value     string
		step      int
	)

	switch flag.NArg() {
	case 3:
		value = flag.Arg(2)
		fallthrough
	case 2:
		system = flag.Arg(0)
		operation = flag.Arg(1)
		step = *fstep
	default:
		exitWithReasonAndUsage("%s", "invalid number of arguments")
	}

	switch system {
	case "kbd":
		if step == 0 {
			step = defaultKdbStep
		}
		handleBrightnessCmd(operation, value, step, kdbFileBrightness, kdbFileMaxBrightness)
	case "screen":
		if step == 0 {
			step = defaultScreenStep
		}
		handleBrightnessCmd(operation, value, step, screenFileBrightness, screenFileMaxBrightness)
	default:
		exitWithReasonAndUsage("unknown system: %q", system)
	}
}

func handleBrightnessCmd(op, value string, step int, briPath, briMaxPath string) {
	switch op {
	case "get", "set", "up", "down", "max":
	default:
		exitWithReasonAndUsage("unknown command: %q", op)
	}

	currentAsBytes, err := ioutil.ReadFile(briPath)
	if err != nil {
		exitWithReason("error getting brightness: %v", err)
	}
	currentAsBytes = bytes.TrimSpace(currentAsBytes)
	current, err := strconv.Atoi(string(currentAsBytes))
	if err != nil {
		exitWithReason("error parsing current brightness value: %v", err)
	}

	maxAsBytes, err := ioutil.ReadFile(briMaxPath)
	if err != nil {
		exitWithReason("error getting max brightness: %v", err)
	}
	maxAsBytes = bytes.TrimSpace(maxAsBytes)
	max, err := strconv.Atoi(string(maxAsBytes))
	if err != nil {
		exitWithReason("error parsing max brightness value: %v", err)
	}

	switch op {
	case "get":
		fmt.Fprintf(os.Stdout, "%s\n", currentAsBytes)
	case "max":
		fmt.Fprintf(os.Stdout, "%s\n", maxAsBytes)
	case "set":
		val, err := strconv.Atoi(value)
		if err != nil {
			exitWithReason("error parsing brightness value (%q): %v", val, err)
		}
		if val < 0 {
			val = 0
		} else if val > max {
			val = max
		}
		writeBrightnessValue(val, briPath)
	case "up":
		newVal := current + step
		if newVal > max {
			newVal = max
		}
		writeBrightnessValue(newVal, briPath)
	case "down":
		newVal := current - step
		if newVal < 0 {
			newVal = 0
		}
		writeBrightnessValue(newVal, briPath)
	default:
		exitWithReasonAndUsage("unknown command: %v", op)
	}
}

func writeBrightnessValue(value int, path string) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 644)
	if err != nil {
		exitWithReason("error opening brightness file: %v", err)
	}
	defer file.Close()

	val := strconv.Itoa(value) + "\n"
	_, err = file.WriteString(val)
	if err != nil {
		exitWithReason("error writing brightness value to file: %v", err)
	}
}

func exitWithReason(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(1)
}

func exitWithReasonAndUsage(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	flag.Usage()
	os.Exit(1)
}
