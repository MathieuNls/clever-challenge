use result::Result; // Result directly is our result.
use std::collections::{HashMap, HashSet};
use std::fs::read;
use std::path::Path;

/// Given a file, return statistical information about the file as if it were a diff file.
/// Returns Some Result if the file could be read and None if the reading failed.
///
/// # Arguments
///
/// `file` - The path to the file to be analized
///
/// # Advice
///
/// This will be the main function per file.
/// It is adviced to create a handler around this as this function is not taking exclusive
/// rights to the path.  The path should be secured by the handler for the sake of returning
/// an error message if the path was invalid and no data could be read.
pub fn diffStats(file: &Path) -> Option<Result> {
    if let Ok(_file_data) = read(file) {
        let mut retVal = Result::new(HashSet::new(), 0, 0, 0, HashMap::new());

        // TODO Process file

        Some(retVal)
    } else {
        None // During file error, we simply return nothing to indicate that the file has no contents instead of valid contents with nothing in it.
    }
}


#[cfg(test)]
mod tests {
    use std::path::Path;

    use result::Result;
    use super::diffStats;

    #[test]
    fn invaldFile() {
        let path = Path::new("/INVALID");
        assert!(match diffStats(&path) {
                    None => true,
                    _ => false,
                });
    }

    #[test]
    fn emptyFile() {
        use std::fs::{write, remove_file};
        let filename = "test.tmp"; // We create an empty file just to be sure it exists.
        write(&filename, "");
        let path = Path::new(&filename);
        assert!(match diffStats(&path) {
                    Some(_) => true,
                    _ => false,
                });
        // TODO implementing PartialEq on result to improve this test to the next line.
        // assert!(match diffStats(&path){ Some(data) => data == Result::empty(), _ => false});
        remove_file(&filename); // Cleanup that file so we don't keep creating different files.

        // TODO crate tempfile should be introduced in test setup to improve this test by
        // removing the issue with creating and removing a file.
    }
}
