#!/bin/bash

go build
sass static/style.scss static/style.css
export FPROOT=`pwd`
export PATH=`pwd`:$PATH

