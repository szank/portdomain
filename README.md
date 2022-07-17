# Port information intestion batch service

The service accepts a path to a file containing the information about ports. The struct definition is 

# Building
The binary can be ran natively, of it one wants to use a docker image, the makefile has `build` and `docker` targets to build
a linux binary and a docker image. 

When using docker the input file dir must be mounted as a volume inside the container for the binary to be able to access it.

# Input validation 
Does not exist. There's no schema for these, so I have very limited options here. Applying common sense. Assuming that `unlocks` is the same as the object key
and the object key is not relevant to the excercise (as it's present in the unlocks array whatever that mean). Not enough data to confirm or deny this assumption. 

# Things that should be done 
* Verify if `regions` and `alias` are arrays of strings or something else. I could not see any examples of non empty arrays, and if I have 2 hours
to finish, then I don't have time to figure out how to access these fields with `jq`.I bet it's super easy, but I have more experience with working with arrays
(think logs) than a single gigantic object.
* Better makefile
* Multi-stage dockerfile where we use one image to build the binary and then copy the result from that image into a clean base where the service runs
* Input validation