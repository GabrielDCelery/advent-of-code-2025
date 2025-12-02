# Advent of Code 2025

## What is this repository for?

This is a repo to solve the puzzles published on [Advent of Code 2025](https://adventofcode.com/2025).

## What do you need to run this?

Have [mise](https://mise.jdx.dev/) installed.

## Running the app

The puzzles' inputs are stored in the `inputs` folder. As the data is specific to each competitor, if you want this repository to solve your puzzle, you must first replace the files with your own.

1. First bootstrap the project:

```sh
mise run bootstrap
```

2. Build the binaries

```sh
mise run build
```

3. Run the binaries to get the soltuon.

```sh
./bin/day01 -file $(realpath ./inputs/day01.txt) -passwordMethod click -logLevel debug
# OR
./bin/day02 -file $(realpath ./inputs/day02.txt) -validator anyrepeat -logLevel debug
```

> [!TIP]
> To get the flags for each binary run them with the help flag, example: ./bin/day01 -help
