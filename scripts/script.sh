#!/usr/bin/env bash
count=0
for i in {0..4}
do
  count=$((i+5))
  curl -d "figure 0.5 0.5
  move 0.${count} 0.${count}
  update" http://localhost:17000
 sleep 1
done
