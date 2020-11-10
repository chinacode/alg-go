#!/bin/sh

cd /data/bloom_server;

nohup ./bloom_server -u edm -p 3306  -h rm-j6ct209xj16ef37jb90150.mysql.rds.aliyuncs.com -P EDM@20201110 -d edm_crawl &
