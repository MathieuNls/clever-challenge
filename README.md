# Clever Initiative Challenge Implementation

Within this repo, you will find an implementation of the Clever Initiative Challenge written in rust.  The focus will be to keep the goals for extendability, maintainability and efficiency during the process of solving the problem.  

## The Challenge

The challenge for trainee applicant consists in parsing a few diffs--in the most efficient way possible (speed, maintainability, evolvability, ...)--and compute the following statistics:

- list of files in the diffs
- number of regions
- number of lines added
- number of lines deleted
- list of function calls seen in the diffs and their number of calls

All these stats are to be computed globally (i.e. for all the diffs combined).

## Permanent Positions

I am applying for a permanent positions.  I understand that this is meant for less permanent positions.  However, as I was directly sent here, I will be completing this challenge.  

## Why rust?

I wish I could do this in Go.  I don't believe in presenting something in a programing language while I am learning it.  As I would be new to Go, I would feel it is not my best work.  

While I am new to rust, my major Rust side project [FeGaBo](github.com/tormyst/FeGaBo) has given me enough experience with the language.  

Using rust has interesting benefits.  While possible to make efficient solutions in other languages, the advantages that rust offers pushes a lot of work onto the compiler.  Ensuring that a task can multithreaded is simple in rust as if the operation is not memory safe, Rust will not compile.  Casts and mutability are brought forward ensuring only what needs to be done is.  

## Building this project

This is a `cargo` project running on nightly rust.  

Once everything installed, you can use `cargo run` to run the solution
