#!/bin/sh

# Note: This does not support alternate properties directories yet

PROPERTIES_ORG=whosonfirst
PROPERTIES_REPO=whosonfirst-properties
PROPERTIES_LOCAL=/usr/local/data/whosonfirst-properties

GITHUB_USER="whosonfirst"
GITHUB_TOKEN_URI=constant://?val=s33kret

ITERATOR_SOURCE=whosonfirst-data:\/\/?prefix=sfomuseum-data-&exclude=sfomuseum-data-venue-

GIT=`which git`
RUNTIMEVAR=`which runtimevar`
INDEX_PROPERTIES=`which index-properties`

INDEX_PROPERTIES=../bin/index-properties

USAGE=""

# Get CLI options

while getopts "o:r:s:t:h" opt; do
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
	s)
	    ITERATOR_SOURCE=$OPTARG
	    ;;		
	t)
	    GITHUB_TOKEN=$OPTARG
	    ;;
	u)
	    GITHUB_USER=$OPTARG
	    ;;
	:   )
	    echo "Unrecognized flag"
	    ;;
    esac
done

if [ "${USAGE}" = "1" ]
then
    echo "usage: ./index.sh -options"
    echo "options:"
    echo "-h Print this message."
    echo "-o The GitHub organization for the properties repo. (Default is ${PROPERTIES_ORG}.)"
    echo "-r The name of the properties repo. (Default is ${PROPERTIES_REPO}.)"
    echo "-s A whosonfirst/go-whosonfirst-iterate-organization URI source to defines repositories to index. (Default is ${ITERATOR_SOURCE}.)"    
    echo "-t A gocloud.dev/runtimevar URI referencing the GitHub API access token to use for updating {PROPERTIES_REPO}. (Default is ${GITHUB_TOKEN}.)"
    echo "-u The GitHub user associated with the GitHub API access token to use for updating {PROPERTIES_REPO}. (Default is ${GITHUB_USER}.)"
    exit 0
fi

# Start 

# First get an access token for writing changes

GITHUB_TOKEN=`${RUNTIMEVAR} ${GITHUB_TOKEN_URI}`

# Clone indexing properties

PROPERTIES_GIT="https://${GITHUB_USER}:${GITHUB_TOKEN}@github.com/${PROPERTIES_ORG}/${PROPERTIES_REPO}"

${GIT} config --global user.email "whosonfirst@localhost"
${GIT} config --global user.name "whosonfirst"

${GIT} clone --depth 1 ${PROPERTIES_GIT} ${PROPERTIES_LOCAL}

# Index properties from repos

${INDEX_PROPERTIES} -iterator-uri org:///tmp -properties ${PROPERTIES_LOCAL} \"${ITERATOR_SOURCE}\"

# Commit changes

cd ${PROPERTIES_LOCAL}

${GIT} add properties
${GIT} commit -m "update properties" properties
${GIT} push origin main
