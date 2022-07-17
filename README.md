# Port information intestion batch service

The service accepts a path to a file containing the information about ports. The struct definition is 

# Building
The binary can be ran natively, of it one wants to use a docker image, the makefile has `build` and `docker` targets to build
a linux binary and a docker image. 

When using docker the input file dir must be mounted as a volume inside the container for the binary to be able to access it.

# Input file structure
IMVHO it would have been easier if the input file contained an array of elements, and wasn't a one large object. That could be done by running `jq` with some magic 
commands I guess. In the first place, there might be duplicate entries there, but what's the primary key? If it's the json key, then no duplicates are allowed, no?

Also, there's no prescribed order to keys in a JSON object, but the reuqirements state that the database should contain the "latest version found in the JSON".
Another point for JSON array ?

# Input validation 
Does not exist. There's no schema for these, so I have very limited options here. Applying common sense. Assuming that `unlocks` is the same as the object key
and the object key is not relevant to the excercise (as it's present in the unlocks array whatever that mean). Not enough data to confirm or deny this assumption. 

# Mocks 
Generally,  I am using counterfeiter for mocks. In this ecxcercise the `Port` type is defined in the service package, which means an inport cycle however I slice the problem.
Hence I've written dead simple mocks myself. Happy to discuss alternatives (I really like counterfeiter).

# Git history rewrite(s)
Did not do that, there are some small fix commits here and there. 

# Testing
Integration testing is not done. It should be. To do it properly I'd need to spin up at least sqlite, (dockerised postgres would be better), set up some migrations to create the table (goose), load up the provided example file and then check what's in the database. To me it seems like ~2 hours of work? Give or take. Not enough time. 

# Things that should be done 
* Support flags and `--help` output. One of the things that end up on the chopping block :(
* Verify if `regions` and `alias` are arrays of strings or something else. I could not see any examples of non empty arrays, and if I have 2 hours
to finish, then I don't have time to figure out how to access these fields with `jq`.I bet it's super easy, but I have more experience with working with arrays
(think logs) than a single gigantic object.
* Better makefile
* Multi-stage dockerfile where we use one image to build the binary and then copy the result from that image into a clean base where the service runs
* Input validation
* Could come up with a lot more, really.
