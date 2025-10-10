package main

import (
	"flag"
	"fmt"
	"log"
)

const locDisclaimer = "\nHey, it turns out, counting LOCs is harder than previously thought! In short, use flag -locf, where f stands for force, to get more meaingful output. Because of the discovered complexity, I will instead learn how to embed Rust code into my projects. In this case, it means that -locf will call Tokei library. I will bundle it built here for convenience. This should in theory run on any Linux distro, but shouldn't run on anything else."

func main() {
	helpFlag := flag.Bool("help", false, "Manual for JHU.")
	locFlag := flag.Bool("loc", false, "Count lines of code in current dir and nested dirs")
	oneLinerFlag := flag.Bool("ol", false, "Copy all your project files into clipboard")
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

	case *oneLinerFlag:
		/*err := */ CopyIntoClipboard() //maybe I want to pass a path? Maybe I can just get a path from the function anyway.
		//if err != nil {
		//	log.Fatal(err)
		//}
		fmt.Println("Project files copied to clipboard.")

	case *helpFlag:
		fmt.Println("TODO: is this professional development?")

	default:
		fmt.Println("No flag provided. Use -help if not sure how to use JHU.")
	}
}
