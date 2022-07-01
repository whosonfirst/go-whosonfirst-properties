#!/bin/sh

# Note: This does not support alternate properties directories yet

PROPERTIES_ORG=whosonfirst
PROPERTIES_REPO=whosonfirst-properties
PROPERTIES_LOCAL=/usr/local/data/whosonfirst-properties

PROPERTIES_ALTERNATES=""
PROPERTIES_EXCLUSIONS=""

GITHUB_TOKEN_URI=constant://?val=s33kret

ITERATOR_SOURCE="whosonfirst-data:\/\/?prefix=whosonfirst-data-&exclude=whosonfirst-data-venue-"

CHECKOUT_BRANCH=
PUSH_BRANCH=main

GIT=`which git`
RUNTIMEVAR=`which runtimevar`
INDEX_PROPERTIES=`which index-properties`

INDEX_PROPERTIES="./bin/index-properties"

USAGE=""

# Get CLI options

while getopts "a:c:e:o:r:s:t:u:h" opt; do
    case "$opt" in
        h) 
	    USAGE=1
	    ;;
	a)
	    PROPERTIES_ALTERNATES+=("$OPTARG")
	    ;;
	c)
	    CHECKOUT_BRANCH=$OPTARG
	    ;;
	e)
	    PROPERTIES_EXCLUSIONS+=("$OPTARG")
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
	    GITHUB_TOKEN_URI=$OPTARG
	    ;;
	u)
	    UPDATE_BRANCH=$OPTARG
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
    echo "-c An optional branch to checkout when performing updates. If not empty then this value will be used to set the -u (update branch) flag. (Default is ${CHECKOUT_BRANCH})."
    echo "-o The GitHub organization for the properties repo. (Default is ${PROPERTIES_ORG}.)"
    echo "-r The name of the properties repo. (Default is ${PROPERTIES_REPO}.)"
    echo "-s A whosonfirst/go-whosonfirst-iterate-organization URI source to defines repositories to index. (Default is ${ITERATOR_SOURCE}.)"    
    echo "-t A gocloud.dev/runtimevar URI referencing the GitHub API access token to use for updating {PROPERTIES_REPO}. (Default is ${GITHUB_TOKEN_URI}.)"
    echo "-u The branch name where updates should be pushed. (Default is ${UPDATE_BRANCH})."
    exit 0
fi

if [ "${CHECKOUT_BRANCH}" != "" ]
then
    $UPDATE_BRANCH=$CHECKOUT_BRANCH
fi

# First get an access token for writing changes

GITHUB_TOKEN=`${RUNTIMEVAR} ${GITHUB_TOKEN_URI}`

# Clone indexing properties

PROPERTIES_GIT="https://${GITHUB_TOKEN}@github.com/${PROPERTIES_ORG}/${PROPERTIES_REPO}"

${GIT} config --global user.email "whosonfirst@localhost"
${GIT} config --global user.name "whosonfirst"

${GIT} clone --depth 1 ${PROPERTIES_GIT} ${PROPERTIES_LOCAL}

# Index properties from repos

echo ${INDEX_PROPERTIES} -iterator-uri org:///tmp -properties ${PROPERTIES_LOCAL}/properties ${ITERATOR_SOURCE}
${INDEX_PROPERTIES} -iterator-uri org:///tmp -properties ${PROPERTIES_LOCAL}/properties ${ITERATOR_SOURCE}

# Commit changes

cd ${PROPERTIES_LOCAL}

if [ "${CHECKOUT_BRANCH}" != "" ]
then
    ${GIT} checkout -b ${CHECKOUT_BRANCH}
fi

${GIT} add properties
${GIT} commit -m "update properties" properties
${GIT} push origin ${UPDATE_BRANCH}
