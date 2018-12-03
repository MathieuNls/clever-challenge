#![allow(non_snake_case)]
use std::time::Instant;
use std::path::Path;
use std::collections::{HashMap, HashSet};

mod result;
use result::Result;

// timeTrack has been modified in the conversion from go to Rust.
// Unfortunetly, the time library has many direct time functions
// still marked as nightly.
fn timeTrack(start: Instant, string: &'static str) {
    let elapsed = start.elapsed();
    print!("{} took {} seconds and {} nanoseconds\n",
           string,
           elapsed.as_secs(),
           elapsed.subsec_nanos());
}

// main is the entry point of our Rust program.
//
// The go line: defer timeTrack(time.Now(), "compute diff")
// translates as our first and last lines.  Rust has no defer
// as that is memory that is held for the duration of the function
// secretly and against the visibility that rust strives for.
//
// It also calls the compute and Display method of the returned struct
// to stdout.
fn main() {
    let now = Instant::now();

    println!("{}", compute().unwrap());

    timeTrack(now, "compute diff");
}

// compute parses the git diffs in ./diffs and returns
// a result struct that contains all the relevant informations
// about these diffs
//  list of files in the diffs
//  number of regions
//  number of line added
//  number of line deleted
//  list of function calls seen in the diffs and their number of calls
fn compute() -> std::result::Result<Result, std::io::Error> {
    let mut retVal = Result::empty();
    let data_folder = Path::new("./diffs");

    for entry in data_folder.read_dir()? {
        let entry = entry?;
        if entry.path().is_file() {
            let mut files = HashSet::new();
            files.insert(entry.file_name().into_string().unwrap());
            retVal += Result::new(files, 0, 0, 0, HashMap::new());
        }
    }

    Ok(retVal)
}
