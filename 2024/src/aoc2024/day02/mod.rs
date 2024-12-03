use core::str;

// import the helpers module that sits one folder up on the aoc2024 module
use crate::aoc2024::helpers::parse_content;

// Function to embed and return the appropriate file content
fn get_file_content(test_flag: bool) -> &'static str {
    if test_flag {
        include_str!("testInput")
    } else {
        include_str!("input")
    }
}

fn safe_line_p1(l: &str) -> bool {
    let parts: Vec<&str> = l.split_whitespace().collect();
    // convert the list to a vector of integers safely
    let numbers: Vec<i32> = parts.iter().map(|s| s.parse().unwrap_or(0)).collect();

    let delta: Vec<i32> = numbers.iter().zip(
            numbers.iter().skip(1)
        ).map(|(a, b)| a - b).collect();
    
    for i in 0..delta.len() {
        if delta[i].abs() > 3 {
            return false;
        }
        if delta[i] == 0 {
            return false;
        }
    }
    if delta[0] < 0 {
        return delta.iter().filter( |&x| *x < 0).count() == delta.len()
    } 
    return delta.iter().filter( |&x| *x > 0).count() == delta.len()
    
}
fn safe_line_p2(l: &str) -> bool {
    let parts: Vec<&str> = l.split_whitespace().collect();
    // convert the list to a vector of integers safely
    let numbers: Vec<i32> = parts.iter().map(|s| s.parse().unwrap_or(0)).collect();

    let safeLine: bool = safe_line_p1(l);
    if safeLine {
        return true;
    }
    // Generate a list removing a single element from the list
    // Check if the line is safe without that element
    for i in 0..numbers.len() {
        let mut newNumbers = numbers.clone();
        newNumbers.remove(i);
        let newLine = newNumbers.iter().map(|x| x.to_string()).collect::<Vec<String>>().join(" ");
        if safe_line_p1(&newLine) {
            return true;
        }
    }
    return false;

    
}
pub fn run(test: bool) {
    let lines = parse_content(get_file_content(test));
    println!("Safe lines P1: {}", lines.iter().filter(|l| safe_line_p1(l)).count());
    println!("Safe lines P2: {}", lines.iter().filter(|l| safe_line_p2(l)).count());
}
