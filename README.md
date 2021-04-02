# minikube-image-benchmark

## Purpose
The purpose of this project is to create a simple to run application that benchmarks different methods of building & pushing an image to minikube.
Each benchmark is run multiple times and the average run time for the runs is calculated and output to a csv file to review the results.

## Requirements
* Docker needs to be installed
* Currently only supported on Linux (only tested on Debian)

## Methods
The three current methods the benchmarks tests is using minikube docker-env, minikube image load, and minikube registry addon, with more being added in the future.

## How to Run Benchmarks
```
make
```
```
./out/benchmark # defaults to 100 runs per method
```
or
```
./out/benchmark --runs 20 # will run 20 runs per method
```
```
cat ./out/results.csv # where the output is stored
```

## Non-Iterative vs Iterative Flow
In the iterative flow the images/cache is cleared after every image build, making it so each build is on a brand new Docker.

In the non-iterative flow the images/cache is cleared at the end of a set of benchmarks. So if 20 runs per benchmark, no cache is cleared until all 20 runs have completed, just the last layer of the image is changed between runs.
