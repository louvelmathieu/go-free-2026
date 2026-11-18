package fake

import "flag"

var leak = flag.Bool("leak", false, "exécute aussi le test qui fuit volontairement")
