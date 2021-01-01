// Index on fileuniqueid
db.posts.createIndex( { "media.fileuniqueid": 1 }, { unique: true } )

// Multikey index on histogramaverage and histogramsum
db.posts.createIndex( { "media.histogramaverage": 1, "media.histogramsum": 1 } )

// Index on messageid
db.posts.createIndex( { "messageid": 1 } )
