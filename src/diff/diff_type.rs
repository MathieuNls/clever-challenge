#[derive(Debug, Clone)]
pub enum DiffType {
    Header,
    Index,
    OriginalFile,
    NewFile,
    NewMode,
    NewRegion,
    FileDeleted,
    FileLine,
    Addition,
    Subtraction,
}

impl DiffType {
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

pub struct DiffFormatTyper {
    last: Option<DiffType>,
}

impl DiffFormatTyper {
    pub fn new() -> Self {
        DiffFormatTyper { last: None }
    }

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
            eprintln!("Error: Diff file is not formated correctly.  Line {} after {:?} was unexpected.",
                      line,
                      self.last);
        }
        self.last = diff_type_fixed;
        self.last.clone()
    }
}
