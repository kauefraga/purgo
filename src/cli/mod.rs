use std::{fs, path::PathBuf};

use clap::{arg, command, ArgMatches, Command};
use colorized::{Color, Colors};
use inquire::Confirm;
use size::{consts, Base, Size};
use walkdir::WalkDir;

pub fn root_command() -> Command {
    command!()
        .arg(arg!([directory] "Set the local to run"))
        .arg(arg!(-d --depth <DEPTH> "Define how deep the search should be").default_value("3"))
    // -i --interactive (bool), define flags via prompts
    // -f --follow-symlinks (bool), disabled by default
    // -y --yes (bool), show the files that will be deleted and ask yes/no just one time
    // -s --size define the size of a considered big file, > 50MB by default
}

pub fn root_handler(matches: ArgMatches) {
    let directory = matches.get_one::<String>("directory");
    let directory = match directory {
        Some(d) => d,
        None => ".",
    };

    if !PathBuf::from(directory).is_dir() {
        eprintln!("The given path is not a directory.");
        return;
    }

    let depth = matches.get_one::<String>("depth").unwrap();
    let depth = depth.parse().unwrap_or(3);

    let walker = WalkDir::new(directory)
        .max_depth(depth)
        .into_iter()
        .filter_map(|e| e.ok());

    for entry in walker {
        let size = entry.metadata().unwrap().len();

        // hard coded 50MB
        if size < 50 * consts::MB as u64 {
            continue;
        }

        println!(
            "{} path: {}",
            "λ".color(Colors::YellowFg),
            entry.path().display()
        );
        println!(
            "{} size: {}",
            "λ".color(Colors::YellowFg),
            Size::from_bytes(size).format().with_base(Base::Base10)
        );

        let answer = Confirm::new("Do you wanna delete this file?")
            .with_default(false)
            .prompt();

        match answer {
            Ok(true) => {
                if fs::remove_file(entry.path()).is_err() {
                    eprintln!("Failed deleting the file.");
                }
            }
            Ok(false) => continue,
            Err(e) => {
                eprintln!("{}", e);
                return;
            }
        }
    }
}
