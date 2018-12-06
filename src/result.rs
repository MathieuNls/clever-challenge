use std::fmt;
use std::collections::{HashSet, HashMap};
use std::ops::AddAssign;

// result contains an analysis of a set of commits
pub struct Result {
    // The name of the files seen
    files: HashSet<String>,
    // How many regions we have (i.e. seperated by @@)
    regions: usize,
    // How many lines were added total
    lineAdded: usize,
    // How many lines were deleted total
    lineDeleted: usize,
    // How many times the function seen in the code are called
    functionCalls: HashMap<String, usize>,
}

impl Result {
    // Creates an empty Result struct.
    pub fn empty() -> Self {
        Self::new(HashSet::new(), 0, 0, 0, HashMap::new())
    }

    // Constructor for a Result struct that expects everything to be handed to it.
    pub fn new(files: HashSet<String>,
               regions: usize,
               lineAdded: usize,
               lineDeleted: usize,
               functionCalls: HashMap<String, usize>)
               -> Self {
        Result {
            files,
            regions,
            lineAdded,
            lineDeleted,
            functionCalls,
        }
    }

    pub fn add_filename(&mut self, filename: String) {
        self.files.insert(filename);
    }

    pub fn add_function_call(&mut self, function: String) {
        let current = *self.functionCalls.get(&function).unwrap_or(&0);
        self.functionCalls.insert(function,current + 1);
    }

    pub fn add_region(&mut self) {
        self.regions += 1
    }

    pub fn count_added_line(&mut self) {
        self.lineAdded += 1
    }

    pub fn count_removed_line(&mut self) {
        self.lineDeleted += 1
    }
}

// Implementation of the formating system in rust ensuring that the structure may be
// printed in a way close to that of the print statment provided in the go structure.
impl fmt::Display for Result {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "Files: \n")?;
        for file in &self.files {
            write!(f, "    -{}\n", file)?;
        }
        write!(f, "Regions: {}\n", self.regions)?;
        write!(f, "LA: {}\n", self.lineAdded)?;
        write!(f, "LD: {}\n", self.lineDeleted)?;

        write!(f, "Functions calls: \n")?;
        for (key, value) in &self.functionCalls {
            write!(f, "   {}: {}\n", key, value)?;
        }
        Ok(())
    }
}

// Iplementation of the += opperator which is useful syntacticaly to clearly demonstrait
// how results are combined.
impl AddAssign for Result {
    fn add_assign(&mut self, mut other: Result) {
        other
            .files
            .drain()
            .for_each(|file| { self.files.insert(file); });
        self.regions += other.regions;
        self.lineAdded += other.lineAdded;
        self.lineDeleted += other.lineDeleted;
        other
            .functionCalls
            .drain()
            .for_each(|(key, value)| {
                          let to_add = *self.functionCalls.get(&key).unwrap_or(&0);
                          self.functionCalls.insert(key, value + to_add);
                      });
    }
}
