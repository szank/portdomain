# Port information intestion batch service

The service accepts a path to a file containing the information about ports. The struct definition is 

# Building
The binary can be ran natively, of it one wants to use a docker image, the makefile has `build` and `docker` targets to build
a linux binary and a docker image. 

When using docker the input file dir must be mounted as a volume inside the container for the binary to be able to access it.

# Things that should be done 
* Better makefile
* Multi-stage dockerfile where we use one image to build the binary and then copy the result from that image into a clean base where the service runs