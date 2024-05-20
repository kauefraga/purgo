use clap::{arg, command, ArgMatches, Command};

pub fn root_command() -> Command {
    command!()
        .arg(arg!(-d --depth <DEPTH> "Define how deep the search should be").default_value("3"))
}

pub fn root_handler(matches: ArgMatches) {
    let depth = matches.get_one::<String>("depth").unwrap();
    let depth = depth.parse().unwrap_or(3);

    println!("Depth: {}", depth);

    todo!();
}
