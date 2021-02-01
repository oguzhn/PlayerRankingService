#!/bin/bash

RANKING_MONGO_ADDR=localhost:27017 \
RANKING_MONGO_DB=RankingDb \
RANKING_MONGO_COL=Players \
RANKING_ADDR=:8080 \
go run main.go