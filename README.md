```
(\ 
\'\ 
 \'\     __________  		___  ___          _             
 / '|   ()_________)		|  \/  |         | |            
 \ '/    \ ~~~~~~~~ \		| .  . | __ _ ___| |_ _ __ ___  
   \       \ ~~~~~~   \		| |\/| |/ _  / __| __| __/  _  \ 
   ==).      \__________\	| |  | | (_| \__ \ |_| | | (_) |
  (__)       ()__________)	\_|  |_/\__,_|___/\__|_|  \___/ 
```

# 
Data and Feature Catalogue in Go

#

## Feature Store

A feature store is a service to store and version features.

### FeatureSets and FeatureStates

A Feature can either be computed on a dataset or a data stream, respectively using a batch or a stream processing pipeline.
This is due to the different life cycle and performance requirements for collecting and serving those data to end applications.

### Example
Have a look at the conf folder for an example configuration, using either ElasticSearch or Mongo as backend.

```
./mastro --configpath conf/featurestore/elastic/example_elastic.cfg
```


### Data Catalogue
Data providers can describe and publish data using a shared definition format.
Consequently, Data definitions can be crawled from networked and distributed file systems, as well as directly published to a common endpoint.

