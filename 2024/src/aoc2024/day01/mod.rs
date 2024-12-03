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

pub fn run(test: bool) {
    let lines = parse_content(get_file_content(test));
    // Add your logic here
    let (left, right) = process_numbers(lines);
    // loop through the list by index
    let mut distance = 0;
    for i in 0..left.len() {
        let d = right[i] - left[i];
        distance += d.abs();
    }
    println!("Distance: {}", distance);

    // convert the right list to a map
    let mut right_map = std::collections::HashMap::new();
    
    for i in 0..right.len() {
      let key = &right[i];
      // check if the key exists in the map
      // if it does, update the value
      // if it doesn't, insert the key and value
      *right_map.entry(key).or_insert(0)+= 1;
    }

    // loop through the left list and check if the key exists in the map
    // if it does, multiply by it
    // if it doesn't, return 0
    let mut product = 0;
    for i in 0..left.len() {
        let key = &left[i];
        let p = *right_map.entry(key).or_insert(0)*left[i];
        println!("{} * {} = {}",left[i], right_map.entry(key).or_insert(0),  p);
        product += p;
    }
    println!("Product: {}", product);



}
// Function to process the file content into sorted left and right lists
fn process_numbers(lines: Vec<String>) -> (Vec<i32>, Vec<i32>) {
    let mut left = Vec::new();
    let mut right = Vec::new();

    for line in lines {
        // Split the line into numbers based on whitespace
        let numbers: Vec<&str> = line.split_whitespace().collect();

        if numbers.len() == 2 {
            // Parse the numbers and push to respective lists
            if let (Ok(left_num), Ok(right_num)) = (numbers[0].parse::<i32>(), numbers[1].parse::<i32>()) {
                left.push(left_num);
                right.push(right_num);
            }
        }
    }

    // Sort both lists
    left.sort_unstable();
    right.sort_unstable();

    (left, right)
}
