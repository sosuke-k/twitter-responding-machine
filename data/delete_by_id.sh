#!/bin/bash
while IFS='' read -r line || [[ -n "$line" ]]; do
    echo "Text read from file: $line"
    mysql -u root trm -e "delete from tweets where id = $line"
done < "$1"
