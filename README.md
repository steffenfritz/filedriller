## filedriller

[![Build status](https://ci.appveyor.com/api/projects/status/vffor64yaxd2bc3q?svg=true)](https://ci.appveyor.com/project/steffenfritz/friller)
[![Go Report Card](https://goreportcard.com/badge/github.com/dla-fritz/filedriller)](https://goreportcard.com/report/github.com/dla-fritz/filedriller)

filedriller walks a directory tree and identifies all regular files by type with [siegfried](https://www.itforarchivists.com/siegfried/). Furthermore it creates UUIDv4s, hash sums (md5, sha1, sha256, sha512 or blake2b-512) and filedriller can check if the file is in the [NSRL](https://www.nist.gov/itl/ssd/software-quality-group/national-software-reference-library-nsrl).

The NSRL check expects a Redis server that serves NSRL SHA-1 hashes. You can use [my docker image](https://hub.docker.com/r/ampoffcom/nslredis)

## Status 

v1.0-BETA

For issues see the issue tab.

## Installation

1. Binary release
    
    Download the file for your platform and execute it. The executables are named friller.
    
    _Note: If the build badge above is green and says passing, it is a good idea to install from source._
    
or

2. From source

        go get codeberg.org/steffenfritz/filedriller/cmd/friller

then

3. Download signature file

       friller -download


4. _Optional NSRL_:

        - docker pull ampoffcom/nslredis:032020

        - docker images

        - docker run -p 6379:6379 $IMAGEID        

    When you pass the -redisserv flag, friller sends a SHA-1 hash to the specified server.



## Usage Examples
0. Fetch the pronom.sig file

        friller -download

1. Without Redis / NSRL

        friller -in SOMEDIRECTORY

2. With Redis / NSRL

        friller -in SOMEDIRECTORY -redisserv localhost

3. With alternate output file

        friller -in SOMEDIRECTORY -output foo.csv

## Output

The output is written to a CSV file. Schema of the file:

    Filename, SizeInByte, Registry, PUID, Name, Version, MIME, ByteMatch, IdentificationNote, HashSum, UUID, inNSRL, Entropy

## Flags

Usage of ./friller:
  
  -download
  
    	Download siegfried's signature file
  
  -entropy

    	Calculate the entropy of files. Limited to file sizes up to 1GB
  
  -hash string
  
    	The hash algorithm to use: md5, sha1, sha256, sha512, blake2b-512 (default "sha256")
  
  -in string
  
    	Root directory to work on
  
  -output string
  
    	Output file (default "info.csv")
  
  -redisport string
  
    	Redis port number for a NSRL database (default "6379")
  
  -redisserv string
  
    	Redis server address for a NSRL database
 
  -version

    	Print version and build info
