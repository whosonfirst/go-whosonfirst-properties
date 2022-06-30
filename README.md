# go-whosonfirst-properties

Go package for working with Who's On First properties

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/whosonfirst/go-whosonfirst-properties.svg)](https://pkg.go.dev/github.com/whosonfirst/go-whosonfirst-properties)

## Tools

### index-properties

Crawl a series of Who's On First documents and ensure that all their properties have a corresponding property file in your `whosonfirst-properties/properties` directory.

```
$> ./bin/index-properties -h
Usage of ./bin/index-properties:
  -alternate value
    	One or more paths to alternate properties directories that will be crawled to check for existing properties (that will not be duplicated).
  -debug
    	Go through all the motions but don't write any new files.
  -exclude value
    	One or more valid regular expressions to use for excluding property names you don't want to index
  -iterator-uri string
    	A valid go-whosonfirst-iterate/v2 URI. (default "repo://")
  -properties string
    	The path to your whosonfirst-properties/properties directory
```

For example:

```
$> ./bin/index-properties \
	-mode sqlite \
	-properties ../whosonfirst-properties/properties \
	/usr/local/data/whosonfirst-data-constituency-us-latest.db
```

Or:

```
$> ./bin/index-properties \
	-exclude 'misc\:.*' \
	-alternate /usr/local/whosonfirst/whosonfirst-properties/properties \
	-properties /usr/local/sfomuseum/sfomuseum-properties \
	/usr/local/data/sfomuseum-data-*
```

Or iterating over all the repositories matching a pattern (`sfomuseum-data-flights-`) in a given organization (`sfomuseum-data`):

```
$> ./bin/index \
	-iterator-uri org:///tmp \
	-properties /usr/local/sfomuseum/sfomuseum-properties/properties \
	-alternate /usr/local/whosonfirst/whosonfirst-properties/properties \
	'sfomuseum-data://?prefix=sfomuseum-data-flights-&exclude=sfomuseum-data-flights-YYYY-MM'
```

### report

Generate a CSV report for a list of `whosonfirst-properties` properties.

```
> ./bin/report -h
Usage of ./bin/report:
  -properties string
    	      The path to your whosonfirst-properties/properties directory
  -report string
    	  The path to write your whosonfirst-properties report. Default is STDOUT.
```

For example:

```
> ./bin/report -properties ../whosonfirst-properties/properties
id,prefix,name,description
1158804491,edtf,cessation,"Indicates when a place stopped being a going concern. The semantics for something ceasing may vary from placetype to placetype. For example, a venue may cease operations or a country may split in to multiple countries."
1158844675,abrv,{lang}_x_colloquial,"The colloquial, informal abbreviation for a place."
1158808009,addr,city,
1158804493,geom,area,"The geometric area of a feature, in WGS84 (unprojected lat/lng)."
1158844669,abrv,{lang}_x_historical,The historical abbreviation for a place.
1158804489,edtf,deprecated,Indicates the date when a place was determined to be invalid (was never a going concern).
1158808003,addr,conscriptionnumber,
1158804497,geom,area_square_m,"The geometric area of a feature in square meters, in the EPSG:3410 projection."
... and so on
```

## See also

* https://github.com/whosonfirst/whosonfirst-properties
