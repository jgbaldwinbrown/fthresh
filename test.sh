#!/bin/bash
set -e

bash install.sh

(cd testdata && (
	bash testgood.sh
	bash testgood_together.sh
	#bash testplot.sh
))
