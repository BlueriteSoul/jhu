# Jerry’s Handy Utility

A personal/educational utility for developers — starting small, growing organically. Currently includes **two main features**:

1. **Counting LOC (Lines of Code)** in a repository
2. **Preparing a "one-liner"** by merging all files in a project into the clipboard — perfect for quickly sharing your project with AI or other tools.

---

## Features

### 1. Count Lines of Code (`-loc` / `-locf`)

``` jhu -loc ```

Counts the total number of lines in the current directory and nested directories.

**Disclaimer (Easter egg)**:

```
Hey, it turns out, counting LOCs is harder than previously thought!
Use flag `-locf`, where `f` stands for force, to get more meaningful output.
This version embeds the Rust-based Tokei library.
It should run on my Linux machine; other systems are not supported.
```

---

### 2. Prepare a One-Liner (`-ol`)

``` jhu -ol ```

Merges all project files into a single string and copies it to the clipboard. Ideal for sharing projects with AI or quickly transporting code snippets.

**Important (Wayland users only):**

* Requires [`wl-clipboard`](https://github.com/bugaevc/wl-clipboard) (`wl-copy` / `wl-paste`) installed.
* Only tested on **Wayland** on my machine — no support for X11 or other systems.

**Notes:**

* Very large projects may exceed clipboard limits. For extremely large codebases, consider writing the output to a file instead.

---

## Installation

1. Clone the repository:

``` git clone [https://github.com/BlueriteSoul/jerrys_handy_utility.git](https://github.com/BlueriteSoul/jerrys_handy_utility.git) ```
``` cd jerrys_handy_utility ```

2. Build the project (requires Go):

``` go install . ```

3. Run the utility:

``` jhu -loc ```
``` jhu -ol ```

---

## Third-Party Components

* **Tokei** ([GitHub](https://github.com/XAMPPRocky/tokei))

  * Copyright (c) 2016 Erin Power
  * Licensed under **MIT and/or Apache 2.0**
  * Used here for advanced LOC counting (`-locf` flag)

---

## Disclaimer

* This is a **personal/educational project**.
* Only guaranteed to work on **my Linux machine with Wayland + wl-clipboard installed**.
* Binary blobs (like the bundled Tokei library) are included for convenience.
* **No warranty is provided** — use at your own risk.

---

## License

Released into the **public domain (CC0 1.0 Universal)**. See `LICENSE` and `NOTICE` files for details.
