package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const locDisclaimer = "\nHey, it turns out, counting LOCs is harder than previously thought! In short, use flag -locf, where f stands for force, to get more meaingful output. Because of the discovered complexity, I will instead learn how to embed Rust code into my projects. In this case, it means that -locf will call Tokei library. I will bundle it built here for convenience. This should in theory run on any Linux distro, but shouldn't run on anything else."

// getArgsAfterFlag returns the slice of raw args that occur after the first
// occurrence of flagName in os.Args (flagName should include the leading dash,
// e.g. "-locf"). It returns nil if flagName isn't present.
func getArgsAfterFlag(flagName string) []string {
	// os.Args[0] is the program name; scan os.Args[1:]
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == flagName {
			// return everything after the flag
			if i+1 >= len(os.Args) {
				return []string{}
			}
			// Return raw args after the flag (don't interpret them with flag package)
			return os.Args[i+1:]
		}
	}
	return nil
}

func main() {
	helpFlag := flag.Bool("help", false, "Manual for JHU.")
	locFlag := flag.Bool("loc", false, "Count lines of code in current dir and nested dirs")
	locfFlag := flag.Bool("locf", false, "Use Tokei to count lines of code in current dir and nested dirs")
	oneLinerFlag := flag.Bool("ol", false, "Copy all your project files into clipboard")
	oneLinerSpecificFlag := flag.Bool("ols", false, "Copy specific files from config into clipboard (from ~.conf/jhu.conf")
	flag.Parse()

	switch {
	case *locFlag:
		// AI thought passing rootDir string was a good idea, but that's kinda implicit and obvious, so I don't see the point.
		totalLines, err := CountLOC()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Total LOC: %d\n", totalLines)
		fmt.Println(locDisclaimer)

	case *locfFlag:
		// Important: collect the raw args after the -locf occurrence and forward them
		raw := getArgsAfterFlag("-locf")
		// If user used "--locf" style (not in your code but just in case), also check that:
		if raw == nil {
			raw = getArgsAfterFlag("--locf")
		}
		// If we didn't find the raw slice (flag not present in os.Args somehow),
		// fall back to flag.Args() (parsed remaining), but normally raw should be used.
		var forward []string
		if raw != nil {
			// sanitize: if the user accidentally passed other jhu flags after -locf,
			// we still forward everything literally. That's intended for passthrough.
			forward = raw
		} else {
			forward = flag.Args()
		}

		// If the user didn't pass any args to tokei, default to current dir
		if len(forward) == 0 {
			forward = []string{"."}
		}

		// If they supplied something like "-help" or "--help" as part of forward,
		// it will be passed verbatim to tokei.
		fmt.Println("Running Tokei via embedded binary...")
		if err := runEmbeddedTokei(forward); err != nil {
			log.Fatal(err)
		}
	case *oneLinerFlag:
		err := CopyIntoClipboard() //maybe I want to pass a path? Maybe I can just get a path from the function anyway.
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Project files copied to clipboard.")

	case *oneLinerSpecificFlag:
		err := CopySpecificIntoClipboard()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Specific project files copied to clipboard.")

	case *helpFlag:
		fmt.Println("TODO: is this professional development?")

	default:
		fmt.Println("No flag provided. Use -help if not sure how to use JHU.")
	}
}
