# go-whosonfirst-properties

Go package for working with Who's On First properties

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/whosonfirst/go-whosonfirst-properties.svg)](https://pkg.go.dev/github.com/whosonfirst/go-whosonfirst-properties)

## Tools

### index

Crawl a series of Who's On First documents and ensure that all their properties have a corresponding property file in your `whosonfirst-properties/properties` directory.

```
> ./bin/index -h
Usage of ./bin/index:
  -debug
	Go through all the motions but don't write any new files
  -mode string
    	The mode to use importing data. Valid modes are: directory,feature,feature-collection,files,geojson-ls,meta,path,repo,sqlite. (default "repo")
  -properties string
    	      The path to your whosonfirst-properties/properties directory
```

For example:

```
./bin/index -mode sqlite -properties ../whosonfirst-properties/properties /usr/local/data/whosonfirst-data-constituency-us-latest.db
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
