# go-whosonfirst-properties

Go package for working with Who's On First properties

## Install

You will need to have both `Go` (specifically a version of Go more recent than 1.7 so let's just assume you need [Go 1.9](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Tools

### wof-properties-report

Generate a CSV report for a list of `whosonfirst-properties` properties.

```
> ./bin/wof-properties-report -h
Usage of ./bin/wof-properties-report:
  -properties string
    	      The path to your whosonfirst-properties/properties directory
  -report string
    	  The path to write your whosonfirst-properties report. Default is STDOUT.
```

For example:

```
> ./bin/wof-properties-report -properties ../whosonfirst-properties/properties
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
