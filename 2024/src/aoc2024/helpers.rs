

/// # Panics
///
#[inline]
pub fn parse_content(content: &str) -> Vec<String> {
    content.lines().map(|line| line.to_string()).collect()
}