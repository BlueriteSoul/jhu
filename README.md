# Jerry’s Handy Utility

A personal/educational utility for developers — starting small, growing organically. Currently includes **three main features**:

1. **Counting LOC (Lines of Code)** in a repository
2. **Preparing a "one-liner"** by merging all files in a project into the clipboard — perfect for quickly sharing your project with AI or other tools.
3. **Copying specific project files** into the clipboard based on a configuration file.

---

## Features

### 1. Count Lines of Code (`-loc` / `-locf`)

``` jhu -loc ```

Counts the total number of lines in the current directory and nested directories, ignoring hidden files, `.git/` folders, and non-text files. This is very simple and naive, because it turns out, counting lines of code is not trivial. More specifically, defining a line of code is not trivial. That's why I attempted to use Tokei library.

**Advanced (`-locf`)**:
Lines Of Code Force. Runs Tokei via embedded Rust library. Currently experimental, supports passing additional flags like `--help`, `--sort`, or `--output json`. Static library (`libtokei.a`) is bundled; rebuild instructions below.

---

### 2. Prepare a One-Liner (`-ol`)

``` jhu -ol ```

Merges all project files into a single string and copies it to the clipboard. Ideal for sharing projects with AI or quickly transporting code snippets.

**Important (Wayland users only):**

* Requires [`wl-clipboard`](https://github.com/bugaevc/wl-clipboard) (`wl-copy` / `wl-paste`) installed.
* Only tested on **Wayland** on my machine — no support for X11 or other systems.

**Notes:**

* Useful only for relatively small projects, as AI chatbots can only take in so much input. If you're project is too big and you're working on a feature spanning few files, try jhu -ols
* Hidden and non-text files are automatically ignored.

---

### 3. Copy Specific Files Into Clipboard (`-ols`)

``` jhu -ols ```

One Liner Specific. Uses a manually created configuration file to copy only certain files from a project into the clipboard.

**Configuration File**:
`~/.config/jhu.conf`

**Example:**

```
#if 1
PROJECT_PATH="/home/averagearchuser/git_repos/tokei/"
FILES="Cargo.toml,src/main.rs,src/cli_runner.rs,src/lib.rs,src/cli.rs"
#endif
```

* Only active blocks (`#if 1` ... `#endif`) are used, so you can at least manually switch between projects.
* `PROJECT_PATH` is the project root; `FILES` is a comma-separated list of relative paths.

---

## Installation

1. Clone the repository:

``` git clone [https://github.com/BlueriteSoul/jerrys_handy_utility.git](https://github.com/BlueriteSoul/jerrys_handy_utility.git) ```
``` cd jerrys_handy_utility ```

2. Build the project (requires Go):

``` go install . ```

**Note:** ```GOPATH/bin is where go install places the binary. Ensure this directory is in your ```PATH to run jhu from anywhere:

``` export PATH=```PATH:```(go env GOPATH)/bin ```

3. Run the utility:

``` jhu -loc ```
``` jhu -locf ```
``` jhu -ol ```
``` jhu -ols ```

---

## Rebuilding the Rust Library (libtokei.a)

If you need to rebuild the Tokei library yourself:

1. Navigate to the modified Tokei folder:

``` cd libs/tokei_reference ```

2. Build with Cargo:

``` cargo build --release ```

3. Copy the built libraries to the `libs/` folder:

``` cp target/release/libtokei.* ../ ```

---

## Third-Party Components

* **Tokei** ([GitHub](https://github.com/XAMPPRocky/tokei))

  * Copyright (c) 2016 Erin Power
  * Licensed under **MIT and/or Apache 2.0**
  * Used for advanced LOC counting (`-locf` flag)

---

## Disclaimer

* Personal/educational project; experimental features included.
* Only guaranteed to work on **my Linux machine with Wayland + wl-clipboard installed**.
* Binary blobs like `libtokei.a` are included for convenience.
* **No warranty is provided** — use at your own risk.

---

## License

Released into the **public domain (CC0 1.0 Universal)**. See `LICENSE` and `NOTICE` files for details.
