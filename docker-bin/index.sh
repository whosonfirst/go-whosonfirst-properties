#!bin/sh

PROPERTIES_ORG=whosonfirst
PROPERTIES_REPO=whosonfirst-properties
PROPERTIES_LOCAL=/usr/local/data/whosonfirst-properties

GITHUB_TOKEN_URI=constant://?val=s33kret
ITERATOR_URI=whosonfirst-data://?prefix=sfomuseum-data-&exclude=sfomuseum-data-venue-

GIT=`which git`
RUNTIMEVAR=`which runtimevar`
INDEX=`which index-properties`

USAGE=""

# Get CLI options

while getopts "o:r:t:i:h" opt; do
    case "$opt" in
        h) 
	    USAGE=1
	    ;;
	o)
	    PROPERTIES_ORG=$OPTARG
	    ;;
	r)
	    PROPERTIES_REPO=$OPTARG
	      ;;	
	t)
	    GITHUB_TOKEN=$OPTARG
	    ;;
	i)
	    ITERATOR_URI=$OPTARG
	    ;;
	:   )
	    echo "Unrecognized flag"
	    ;;
    esac
done

if [ "${USAGE}" = "1" ]
then
    echo "usage: update.sh"
    echo "options:"
    echo "...please write me"
    exit 0
fi

# Start 

# First get an access token for writing changes

GITHUB_TOKEN=`${RUNTIMEVAR} ${GITHUB_TOKEN_URI}`

# Clone indexing properties

PROPERTIES_GIT="https://${GITHUB_USER}:${GITHUB_TOKEN}@github.com/${PROPERTIES_ORG}/${PROPERTIES_REPO}"

${GIT} clone --depth 1 ${PROPERTIES_GIT} ${PROPERTIES_LOCAL}

# Index properties from repos

${INDEX} -iterator-uri org:///tmp -properties ${PROPERTIES_LOCAL} "${ITERATOR_URI}"

# Commit changes

cd ${PROPERTIES_LOCAL}

${GIT} add properties
${GIT} commit -m "update properties" properties
${GIT} push origin main
