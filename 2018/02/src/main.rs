use std::env;
use std::fs::File;
use std::io::prelude::*;
use std::iter::Iterator;
use std::collections::HashMap;

fn count_chars(mut hash_map: HashMap<char, u32>, word: char) -> HashMap<char, u32> {
    {
        let c = hash_map.entry(word).or_insert(0);
        *c += 1;
    }

    hash_map
}
fn unique_lines(mut hash_map: HashMap<String, String>, l: String) -> HashMap<String, String> {
    hash_map.insert(l.clone(),l);
    hash_map
}
fn part_1(input_string: String) -> u32 {
    let mut twos = 0;
    let mut threes = 0;
    for line in input_string.lines() {
        let s: HashMap<char, u32> = line.
                chars().
                fold(HashMap::new(), count_chars).
                into_iter().
                filter(|&(_,v)| v == 2 || v==3).
                collect();
        let list: Vec<u32> = s.iter().map(|(_, v)| v.clone()).collect();
        
        if list.contains(&2) {
            twos +=1;
        }
        if list.contains(&3) {
            threes +=1;
        }
    }
    return twos*threes;
}
fn part_2(input_string: String) -> Vec<String> {
    let u = input_string.lines();
    u.map(|x| x.unwrap()).collect()
}
fn main() {
    // --snip--
    let args: Vec<String> = env::args().collect();

    let config = parse_config(&args);
    let file = File::open(config.filename).expect("file not found").unwrap();
    let mut input_string = String::new();
    
    expect("something went wrong reading the file");
    //println!("Checksum {}",part_1(input_string.clone()));
    part_2(input_string);
}

struct Config {
    filename: String,
}

fn parse_config(args: &[String]) -> Config {
    let filename;
    if args.len() <1 {
        filename = "input".to_string()
    }else{
        filename = args[1].clone();
    }

    Config {  filename}
}