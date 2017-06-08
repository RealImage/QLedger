#!/usr/bin/env bash

DB=`echo $DATABASE_URL|cut -d/ -f4|cut -d? -f1`
pg_dump -s -x -O $DB | grep -v -e "^--" -e "^$" > schema.sql