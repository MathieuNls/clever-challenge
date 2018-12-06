use result::Result;
use regex::Regex;

/// Returns the two filenames present in the diff header.  
/// We are assuming --git, and this will ensure it as any non diff header will crash.  
/// This is important as otherwise, we would be storing each filename twice: once for a path and
/// once for b path.
///
/// This implementation assumes a path is any collection of characters that is not a space.
pub fn header_filenames(header: &str) -> (String, String) {
    lazy_static! {
        static ref HEADER_GIT_FILENAME: Regex = Regex::new(r"^diff --git a/([^ ]+) b/([^ ]+)$").unwrap();
    }
    let groups = HEADER_GIT_FILENAME.captures(header).unwrap();
    (groups.get(1).unwrap().as_str().to_string(), groups.get(2).unwrap().as_str().to_string())
}

/// Finds and adds all functions to a given result structure.  
/// This is the slowest part of the diff stats process as we have to iterate through a utf8 string
/// doing constant comparisons.  
/// For the purposes of this chalange, I am defining a function given the regex bellow.  A number
/// of word characters followed imidiatly by an open bracket.  The issue with this structure is
/// that it is unlikely to be correct.  Languages like C can have as much whitespace as wanted
/// between the name of the function and the open bracket.  Likewise, the string 'for(' should not
/// be a function identifier, however this regex would.  
///
/// All in all, this is the best that can be done in a short time as the alternetive would be to
/// actually understand whatever language we are sifting though and then find funtions as defined
/// by the language itself.
///
/// # Why are we passing in a result structure?
///
/// This function has the beautiful disadvantage of being a variable size responce.  
/// We don't know how many bytes any return of find_functions will take as it may return 0 or 5.  
/// We however could return a Results structure: Initializing say 10_000 result structures, with
/// 10_000 empty sets for filenames is fairly useless.  
/// Returning HashSet: We could return just the part of the results we wanted, the hashset and
/// create a function to add them together, or to open the underlying structure of the result.  
/// I went with passing the result structure in directly.  This is something I feel needs to be
/// fixed, however many other options seem unoptimal.
pub fn find_functions(string: &str, result: &mut Result) {
    lazy_static! {
        static ref FUNCTION_REGEX: Regex = Regex::new(r"\w+\(").unwrap();
    }
    FUNCTION_REGEX
        .find_iter(string)
        .map(|function| format!("{})", function.as_str()))
        .for_each(|string| result.add_function_call(string))
}
