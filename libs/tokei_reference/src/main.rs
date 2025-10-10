use std::io::{self, Write};
use tokei::cli_runner;

fn main() {
    println!("Halleluja, Rust is in da house. (main.rs)");
    io::stdout().flush().unwrap();

    if let Err(e) = cli_runner::run_cli() {
        eprintln!("Error: {}", e);
        std::process::exit(1);
    }
}
