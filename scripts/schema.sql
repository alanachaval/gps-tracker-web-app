CREATE TABLE gpsframes.`user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user` varchar(50) NOT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE gpsframes.`gpsTrack` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `userId` int(11) unsigned NOT NULL,
  `trackNumber` int(11) unsigned NOT NULL,
  `time` varchar(50) NOT NULL,
  `longitude` varchar(10) NOT NULL,
  `latitude` varchar(10) NOT NULL,
  `status` varchar(10) NOT NULL,
  `latitudeHemisphere` varchar(10) NOT NULL,
  `longitudeHemisphere` varchar(10) NOT NULL,
  `earthVelocity` varchar(10) NOT NULL,
  `track` varchar(10) NOT NULL,
  `date` varchar(50) NOT NULL,
  `magneticVariation` varchar(10) NOT NULL,
  `directionVariation` varchar(10) NOT NULL,
  `systemPosition` varchar(10) NOT NULL,
  `checksum` varchar(50) NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (userId) REFERENCES gpsframes.user(id),
  UNIQUE KEY (userId, trackNumber)
);
