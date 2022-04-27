#!/usr/bin/env bash

# Include colors.sh
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/colors.sh"

echoHeader "Scraping documentation from godoc webserver"
wget -r -nv -np -N -E -p -k -e robots=off --include-directories="/pkg,/lib" --exclude-directories="*" http://${DOCS_HOST}/pkg/${DOCS_PATH}/

echoHeader "Rewriting docs folder with fresh docs"
mkdir -p ${DOCS_DIR}
rm -rf ${DOCS_DIR}/{pkg,lib,index.html}
cp -a ${DOCS_HOST}/lib ${DOCS_DIR}
cp -a ${DOCS_HOST}/pkg/${DOCS_PATH}/* ${DOCS_DIR}
rm -rf ${DOCS_HOST}

echoHeader "Rebasing links to our docs root directory"
echo "s,http://${DOCS_HOST}/src/${DOCS_PATH}/,,g" > .pattern
grep godoc/style.css docs/index.html | sed 's/.*href="/s,/; s/lib.*/,,g/; s/\./\\./g' >> .pattern
find docs -name "*.html" | xargs sed -i.bak -f .pattern
find docs -name "*.bak" | xargs rm
rm .pattern

echoTitle "Done! Go check your new docs"
