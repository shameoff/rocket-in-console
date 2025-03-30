# Rocket in Console

**Rocket in Console** is a cross-platform rocket flight simulator game in the terminal, written in Go using the [termbox-go](https://github.com/nsf/termbox-go) library. The project demonstrates a dynamic physics model, ASCII art for visual effects, and interactive controls (using arrow keys).

**This project was created as an artifact for the "Intellectual Property" topic in the "Legal Studies" course.**

## Features

- **Flight Physics:** Real integration of acceleration, control over thrust in four directions, and power limitations for engines (main engine – up to 100, auxiliary engines – up to ±10).
- **Dynamic Landscape:** Generation of random stars, clouds, and trees to create the feeling of moving through a cosmic space.
- **Visual Effects:** Animation of engine exhaust with different colors for the main (red) and auxiliary (blue) engines, as well as an explosion effect on a crash landing.
- **Cross-Platform:** Ability to build for Windows, Linux, macOS, as well as for x86 (amd64) and ARM (arm64) architectures.

## Installation

### Prerequisites

- Go 1.18 or later.
- Git (if you clone the repository).

### Cloning the Repository

```bash
git clone https://github.com/shameoff/rocket-in-console.git
cd rocket-in-console
```

## Running the Game

You can run the game immediately with:

```bash
go run cmd/myrocketgame/main.go
```

## Building

To build an executable for your current platform, run:

```bash
go build -o rocket-in-console ./cmd/myrocketgame
```

## Cross-Platform Building with GitHub Actions

The repository includes [GitHub Workflows](.github/workflows/build.yml) for building on Windows, Linux, and macOS, as well as for amd64 and arm64 architectures. After pushing or creating a PR to the `master` or `main` branch, the build process will be triggered and artifacts will be available under the _Actions_ tab.

## Project Structure

```
rocket-in-console/
├── cmd/
│   └── myrocketgame/
│       └── main.go         # Main entry point and game loop
├── pkg/
│   ├── input/              # Keyboard input handling
│   ├── objects/            # Definitions of game objects (rocket, stars, clouds, trees, etc.)
│   ├── physics/            # Physics model and logic for updating object states
│   └── render/             # Terminal rendering functions (ASCII art, UI)
├── .github/
│   └── workflows/
│       └── build.yml       # GitHub Actions for cross-platform builds
└── go.mod                  # Go module file
```

## Contribution and Development

If you have suggestions or want to contribute changes, please open an [Issue](https://github.com/shameoff/rocket-in-console/issues) or a [Pull Request](https://github.com/shameoff/rocket-in-console/pulls).

## License

This project is licensed under the MIT License.
