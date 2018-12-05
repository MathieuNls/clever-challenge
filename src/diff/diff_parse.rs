use regex::Regex;

pub fn header_filenames(header: &str) -> (String, String){
    lazy_static! {
        static ref HEADER_GIT_FILENAME: Regex = Regex::new(r"^diff --git a/([^ ]+) b/([^ ]+)$").unwrap();
    }
    let groups = HEADER_GIT_FILENAME.captures(header).unwrap();
    (groups.get(1).unwrap().as_str().to_string(), groups.get(2).unwrap().as_str().to_string())
}
