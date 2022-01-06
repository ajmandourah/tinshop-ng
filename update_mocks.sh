#!/bin/bash

mkdir -p mock_repository
mockgen github.com/DblK/tinshop/repository Config > mock_repository/mock_config.go 
mockgen github.com/DblK/tinshop/repository Source > mock_repository/mock_source.go 
mockgen github.com/DblK/tinshop/repository Collection > mock_repository/mock_collection.go 
mockgen github.com/DblK/tinshop/repository Sources > mock_repository/mock_sources.go 
mockgen github.com/DblK/tinshop/repository Stats > mock_repository/mock_stats.go 
mockgen github.com/DblK/tinshop/repository API > mock_repository/mock_api.go 