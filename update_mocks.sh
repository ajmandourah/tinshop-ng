#!/bin/bash

mkdir -p mock_repository
mockgen github.com/DblK/tinshop/repository Config > mock_repository/mock_config.go 
mockgen github.com/DblK/tinshop/repository Source > mock_repository/mock_source.go 