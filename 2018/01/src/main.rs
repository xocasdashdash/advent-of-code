use std::env;
use std::fs::File;
use std::io::prelude::*;
use std::collections::HashSet;

fn main() {
    // --snip--
    let args: Vec<String> = env::args().collect();

    let config = parse_config(&args);
    let mut f = File::open(config.filename).expect("file not found");
    let mut input_string = String::new();
    f.read_to_string(&mut input_string).
    expect("something went wrong reading the file");
    let mut seen_frequencies: HashSet<i32> = HashSet::new();
    let mut frequency = 0;
    let mut i = 0;

    let suma: i32 = input_string
        .lines()
        .map(|w| w.parse::<i32>().unwrap())
        .sum();

    let input: Vec<i32> = input_string
        .lines()
        .map(|w| w.parse::<i32>().unwrap())
        .collect();

    println!("SUM: {}", suma);
    let mut hits=0;
    while seen_frequencies.insert(frequency) {
        hits +=1;
        frequency += input[i];
        i = (i + 1) % input.len();
    }
    println!("FREQ: {}, hits = {}. Loops {}", frequency, hits, hits/input_string.lines().count() );
}

struct Config {
    filename: String,
}

fn parse_config(args: &[String]) -> Config {
    let filename;
    if args.len() <1 {
        filename = "input"
    }else{
        filename = &args[1].clone().to_string();
    }

    Config {  filename.to_s }
}