#!/usr/bin/env bash

BASEDIR=$2

echo $BASEDIR

inkscape -z -e $BASEDIR/icon-72x72.png -w 72 -h 72 $1
inkscape -z -e $BASEDIR/icon-96x96.png -w 96 -h 96 $1
inkscape -z -e $BASEDIR/icon-128x128.png -w 128 -h 128 $1
inkscape -z -e $BASEDIR/icon-144x144.png -w 144 -h 144 $1
inkscape -z -e $BASEDIR/icon-152x152.png -w 152 -h 152 $1
inkscape -z -e $BASEDIR/icon-192x192.png -w 192 -h 192 $1
inkscape -z -e $BASEDIR/icon-384x384.png -w 384 -h 384 $1
inkscape -z -e $BASEDIR/icon-512x512.png -w 512 -h 512 $1
