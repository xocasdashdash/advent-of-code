mod aoc2024;

use std::env;

fn main() {
    let args: Vec<String> = env::args().collect();
    if args.len() < 2 {
        eprintln!("Usage: cargo run -- <day> [test]");
        return;
    }

    let day = &args[1];
    let test_flag = args.get(2).map(|s| s == "test").unwrap_or(false);

    match day.as_str() {
        "01" => aoc2024::day01::run(test_flag),
        "02" => aoc2024::day02::run(test_flag),
        _ => eprintln!("Day {} not implemented", day),
    }
}
