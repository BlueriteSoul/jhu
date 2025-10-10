//! # Tokei: Count your code quickly.
//!
//! A simple, efficient library for counting code in directories. This
//! functionality is also provided as a
//! [CLI utility](//github.com/XAMPPRocky/tokei). Tokei uses a small state
//! machine rather than regular expressions found in other code counters. Tokei
//! can accurately count a lot more edge cases such as nested comments, or
//! comment syntax inside string literals.
//!
//! # Examples
//!
//! Gets the total lines of code from all rust files in current directory,
//! and all subdirectories.
//!
//! ```no_run
//! use std::collections::BTreeMap;
//! use std::fs::File;
//! use std::io::Read;
//!
//! use crate::{Config, Languages, LanguageType};
//!
//! // The paths to search. Accepts absolute, relative, and glob paths.
//! let paths = &["src", "tests"];
//! // Exclude any path that contains any of these strings.
//! let excluded = &["target"];
//! // `Config` allows you to configure what is searched and counted.
//! let config = Config::default();
//!
//! let mut languages = Languages::new();
//! languages.get_statistics(paths, excluded, &config);
//! let rust = &languages[&LanguageType::Rust];
//!
//! println!("Lines of code: {}", rust.code);
//! ```

#![deny(
    trivial_casts,
    trivial_numeric_casts,
    unused_variables,
    unstable_features,
    unused_import_braces,
    missing_docs
)]

#[macro_use]
extern crate log;
#[macro_use]
extern crate serde;

#[macro_use]
mod utils;
mod config;
mod consts;
mod language;
mod sort;
mod stats;

pub use self::{
    config::Config,
    consts::*,
    language::{Language, LanguageType, Languages},
    sort::Sort,
    stats::{find_char_boundary, CodeStats, Report},
};
// above is all of the lib.rs unmodified, I was supposed to only append the following:
//
//
//
//
// soooo, are you saying I don't need any of this (below)?
/*use std::ffi::{CStr, CString};
use std::os::raw::c_char;

/// # Safety
/// The `dir` must be a valid null-terminated C string.
#[no_mangle]
pub extern "C" fn count_loc(dir: *const c_char) -> u64 {
    if dir.is_null() {
        return 0;
    }

    // Convert C string to Rust string
    let c_str = unsafe { CStr::from_ptr(dir) };
    let path = match c_str.to_str() {
        Ok(s) => s,
        Err(_) => return 0,
    };

    // Prepare Tokei
    let paths = &[path];
    let excluded: &[&str] = &["target", ".git"];
    let config = Config::default();

    let mut languages = Languages::new();
    languages.get_statistics(paths, excluded, &config);

    // Sum lines of code across all languages
    let total: u64 = languages.iter().map(|(_, lang)| lang.code as u64).sum();
    total
}*/
// I'm backing up the below
/*
use std::os::raw::c_int;

mod cli;
/// CLI runner for Tokei.
pub mod cli_runner;
mod cli_utils;
mod input;

/// Run the Tokei CLI. Returns 0 on success, 1 on error.
#[no_mangle]
pub extern "C" fn run_tokei_cli() -> c_int {
    match cli_runner::run_cli() {
        Ok(_) => 0,
        Err(_) => 1,
    }
}
*/
use std::os::raw::{c_char, c_int};

mod cli;
/// IDK, man, just work pls
pub mod cli_runner;
mod cli_utils;
mod input;

/// Can I just say anything here?
#[no_mangle]
pub extern "C" fn run_tokei_with_args(argc: c_int, argv: *const *const c_char) -> c_int {
    use std::ffi::CStr;

    let args = unsafe {
        std::slice::from_raw_parts(argv, argc as usize)
            .iter()
            .map(|&c| CStr::from_ptr(c).to_string_lossy().into_owned())
            .collect::<Vec<String>>()
    };

    match cli_runner::run_cli_with_args(args) {
        Ok(_) => 0,
        Err(_) => 1,
    }
}
