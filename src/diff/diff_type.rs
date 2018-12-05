/// Represents a kind of diff line.  
/// The diff is split into 2 different parts:
///     Header information
///     File information
/// The problem is that this information is not split clearly.  The header is not constant and may
/// lack any part of the header or body.  
///
/// In this enum, all posible types of lines are sorted.  This is important even if the line is not
/// needed now, because if data is needed from a line later, we have the infromation labeled.
#[derive(Debug, Clone)]
pub enum DiffType {
    /// The begining of a diff of the form: 'diff --git a/path/to/file b/path/tofile'
    Header,
    /// Index information about a file: 'index 0123456..789abcd 100644'
    Index,
    /// The path to the original file, the first file in the header: '--- a/path/to/file'
    OriginalFile,
    /// The path to the new file, the second file in the header: '+++ b/path/to/file'
    NewFile,
    /// Optional header to indicate that a files mode differs changed: 'new file mode 100644'
    NewMode,
    /// The start of a file region: '@@ -888,12 +1002,33 @@ part of the file here'
    NewRegion,
    /// The indication that the file was deleted: 'deleted file mode 100644'
    FileDeleted,
    /// A line of a file which is the same between the given files: ' always with a space'
    FileLine,
    /// A line that exists in the new file, but not the original: '+the line of the file'
    Addition,
    /// A line that exists in the original file, but not the new: '-the line of the file'
    Subtraction,
}


impl DiffType {
    /// Returns if a type is part of the file body, or has file body parts to it.  
    /// This is useful as the slowest part of the process would be parcing the entire file for
    /// functions.
    pub fn is_file_body(&self) -> bool {
        match self {
            DiffType::FileLine | DiffType::Addition | DiffType::Subtraction |
            DiffType::NewRegion => true,
            _ => false,
        }
    }
}

/// Returns the diff status of a line.  There are a few states, including the different header
/// lines, normal file lines, additions, and subtractions.
///
/// This will do the minimum work required to figure out what type a line is because it assumes the
/// string slice it is given starts at the begining of the line and that the file is a valid diff
/// file.  Any unknown first letter of the line is treated as FileLine, as are empty lines.
///
/// # Known Issues
///
/// Currently, OriginalFile and NewFile are treeted as Addition and Subtraction lines.
/// This is an issue of only using the first letter.  While it could be improved, there are edge
/// cases that cannot be resolved such as an additon line adding a line with 2 or more '+'
/// characters.  To avoid this, this system uses a blind approch that is not attempting to solve
/// this issue.
///
/// # Arguments
///
/// `line` - A string slice of a line of a diff file.
///
/// # Example
/// ```
/// assert!(diff_type("diff --git ...") == DiffType::Header)
fn diff_type(line: &str) -> DiffType {
    match line.chars().next().unwrap_or(' ') {
        'd' => DiffType::Header,
        'i' => DiffType::Index,
        'n' => DiffType::NewMode,
        '@' => DiffType::NewRegion,
        '+' => DiffType::Addition,
        '-' => DiffType::Subtraction,
        ' ' => DiffType::FileLine,
        _ => DiffType::FileLine,
    }
}

/// Stateful Diff reading to avoid confusion over header and body.
/// Place all stateful information here.
///
/// TODO: If you would want to know how many functions came from a given file,
/// putting a current file calculation here would be how to do it.
pub struct DiffFormatTyper {
    last: Option<DiffType>,
}

impl DiffFormatTyper {
    /// A new reader that should be started at the top of the file, however, it can start at any
    /// header line.
    pub fn new() -> Self {
        DiffFormatTyper { last: None }
    }

    /// A state machine that uses the guess of diff_type to create the actual type of a line and to
    /// fail when it is not sure.  
    ///
    /// Many just check that this line is expected and return it.  Some however take the guess, do
    /// a string compare to the minimum degree required and return the updated value.  
    ///
    /// # Returns
    ///
    /// A value of Some(data) is considered a correct responce from the function.  
    /// None is the universal error.  We also reset the state machine on error, and will continue
    /// to output correct responces at the next header. 
    ///
    /// A line to STDERR is sent indicating the error based on the string given and the current
    /// state of the machine.  
    ///
    /// # Note
    ///
    /// This section is messy and should be replaced with macros.
    pub fn type_line(&mut self, line: &str) -> Option<DiffType> {
        let diff_type = diff_type(line);
        let diff_type_fixed = match self.last.clone() {
            None => {
                match diff_type {
                    DiffType::Header => Some(DiffType::Header),
                    _ => None,
                }
            }
            Some(last_type) => {
                match last_type {
                    DiffType::Header => {
                        match diff_type {
                            this @ DiffType::Index |
                            this @ DiffType::NewMode => Some(this),
                            DiffType::Header => {
                                if line.starts_with("deleted") {
                                    Some(DiffType::FileDeleted)
                                } else {
                                    None
                                }
                            }
                            _ => None,
                        }
                    }
                    DiffType::Index => {
                        match diff_type {
                            this @ DiffType::Header => Some(this),
                            DiffType::Subtraction => {
                                if line.starts_with("---") {
                                    Some(DiffType::OriginalFile)
                                } else {
                                    None
                                }
                            }
                            _ => None,
                        }
                    }
                    DiffType::OriginalFile => {
                        match diff_type {
                            DiffType::Addition => {
                                if line.starts_with("+++") {
                                    Some(DiffType::NewFile)
                                } else {
                                    None
                                }
                            }
                            _ => None,
                        }
                    }
                    DiffType::NewFile => {
                        match diff_type {
                            this @ DiffType::NewRegion => Some(this),
                            _ => None,
                        }
                    }
                    DiffType::NewMode => {
                        match diff_type {
                            this @ DiffType::Index => Some(this),
                            _ => None,
                        }
                    }
                    DiffType::FileDeleted => {
                        match diff_type {
                            this @ DiffType::Index => Some(this),
                            _ => None,
                        }
                    }
                    DiffType::NewRegion | DiffType::FileLine | DiffType::Addition |
                    DiffType::Subtraction => {
                        match diff_type {
                            this @ DiffType::FileLine |
                            this @ DiffType::Addition |
                            this @ DiffType::Subtraction |
                            this @ DiffType::NewRegion |
                            this @ DiffType::Header => Some(this),
                            _ => None,
                        }
                    }
                }
            }
        };
        if let None = diff_type_fixed {
            eprintln!("Error: Diff file is not formated correctly.  Line \"{}\" after {:?} was unexpected.",
                      line,
                      self.last);
        }
        self.last = diff_type_fixed;
        self.last.clone()
    }
}
