mod cli;

fn main() {
    let matches = cli::root_command().get_matches();

    match matches.subcommand() {
        None => cli::root_handler(matches),
        _ => unreachable!(),
    }
}
