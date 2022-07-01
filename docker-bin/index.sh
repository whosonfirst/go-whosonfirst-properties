#!/bin/bash

# Note: bash is necessary for the `FOO+=" ${OPTARG}"` stuff below which does not
# work using plain-old sh under alpine.

PROPERTIES_ORG=whosonfirst
PROPERTIES_REPO=whosonfirst-properties

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
	    PROPERTIES_ALTERNATES+=" ${OPTARG}"
	    ;;
	c)
	    CHECKOUT_BRANCH=$OPTARG
	    ;;
	e)
	    PROPERTIES_EXCLUSIONS+=" ${OPTARG}"
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
	    PUSH_BRANCH=$OPTARG
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
    echo "-a Zero or more Git URLs for alternate properties repositories to clone."
    echo "-c An optional branch to checkout when performing updates. If not empty then this value will be used to set the -u (update branch) flag. (Default is ${CHECKOUT_BRANCH})."
    echo "-e Zero or more regular expressions to specify properties that should not be indexed."
    echo "-o The GitHub organization for the properties repo. (Default is ${PROPERTIES_ORG}.)"
    echo "-r The name of the properties repo. (Default is ${PROPERTIES_REPO}.)"
    echo "-s A whosonfirst/go-whosonfirst-iterate-organization URI source to defines repositories to index. (Default is ${ITERATOR_SOURCE}.)"    
    echo "-t A gocloud.dev/runtimevar URI referencing the GitHub API access token to use for updating {PROPERTIES_REPO}. (Default is ${GITHUB_TOKEN_URI}.)"
    echo "-u The branch name where updates should be pushed. (Default is ${PUSH_BRANCH})."
    exit 0
fi

if [ "${CHECKOUT_BRANCH}" != "" ]
then
    PUSH_BRANCH=$CHECKOUT_BRANCH
fi

# First get an access token for writing changes
# See also: https://github.com/sfomuseum/runtimevar

GITHUB_TOKEN=`${RUNTIMEVAR} ${GITHUB_TOKEN_URI}`

# Git housekeeping

${GIT} config --global user.email "whosonfirst@localhost"
${GIT} config --global user.name "whosonfirst"

# Clone indexing properties

PROPERTIES_GIT="https://${GITHUB_TOKEN}@github.com/${PROPERTIES_ORG}/${PROPERTIES_REPO}"

PROPERTIES_LOCAL=/usr/local/data/${PROPERTIES_REPO}

${GIT} clone --depth 1 ${PROPERTIES_GIT} ${PROPERTIES_LOCAL}

# Clone the alternates

for ALTERNATE_GIT in "${PROPERTIES_ALTERNATES}"
do
    ALTERNATE_FNAME=`basename ${ALTERNATE_GIT}`
    ALTERNATE_LOCAL="/usr/local/data/${ALTERNATE_FNAME}"
    ${GIT} clone --depth 1 ${ALTERNATE_GIT} ${ALTERNATE_LOCAL}
done

# Build indexing command from flags

INDEXING_CMD="${INDEX_PROPERTIES} -iterator-uri org:///tmp -properties ${PROPERTIES_LOCAL}/properties"

for ALTERNATE_GIT in "${PROPERTIES_ALTERNATES}"
do
    ALTERNATE_FNAME=`basename ${ALTERNATE_GIT}`
    ALTERNATE_LOCAL="/usr/local/data/${ALTERNATE_FNAME}"
    INDEXING_CMD="${INDEXING_CMD} -alternate ${ALTERNATE_LOCAL}/properties"
done

for EXCLUSION in "${PROPERTIES_EXCLUSIONS}"
do
    INDEXING_CMD="${INDEXING_CMD} -exclude ${EXCLUSION}"
done

# Actually do the indexing

echo ${INDEXING_CMD} ${ITERATOR_SOURCE}
${INDEXING_CMD} ${ITERATOR_SOURCE}

# Commit changes

cd ${PROPERTIES_LOCAL}

if [ "${CHECKOUT_BRANCH}" != "" ]
then
    ${GIT} checkout -b ${CHECKOUT_BRANCH}
fi

${GIT} add properties
${GIT} commit -m "update properties" properties
${GIT} push origin ${PUSH_BRANCH}
