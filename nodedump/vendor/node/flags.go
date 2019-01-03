package node

import "flag"
import "log"

var specFlag bool
var condFlag bool
var labelFlag bool
var annotationFlag bool
var infoFlag bool
var allFlag bool

func InitFlags() {
	flag.BoolVar(&specFlag, "spec", false, "Node Spec")
	flag.BoolVar(&condFlag, "cond", false, "Node Conditions")
	flag.BoolVar(&labelFlag, "label", false, "Node Labels")
	flag.BoolVar(&annotationFlag, "annotation", false, "Node Annotations")
	flag.BoolVar(&infoFlag, "info", false, "Node Info")
	flag.BoolVar(&allFlag, "all", false, "Enable all Flags")

	flag.Parse()
	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		log.Fatalf("%s\n", "No flags provided. Refer to the usage above.")
	}
}
