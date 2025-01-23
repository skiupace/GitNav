use syntect::{
    easy::HighlightLines,
    highlighting::ThemeSet,
    parsing::SyntaxSet,
    util::LinesWithEndings,
};

use ratatui::{
    text::{Line, Text, Span},
    style::{Color, Style},
};

use std::{ fs, path::{Path, PathBuf} };
use crate::common::Result;
use ignore::WalkBuilder;


pub fn is_git_repo(path: &Path) -> bool {
    let git_dir = path.join(".git");
    git_dir.exists() && git_dir.is_dir()
}

pub fn get_name(path: &Path) -> String {
    path.file_name()
        .and_then(|name| name.to_str())
        .unwrap_or("")
        .to_string()
}

pub fn get_icon(path: &Path) -> &'static str {
    if path.is_dir() {
        "📁"
    } else {
        "📄"
    }
}

pub fn list_files_and_folders(current_path: &Path) -> Vec<PathBuf> {
    let mut items = Vec::new();

    // Use the `ignore` crate to walk the directory and exclude .gitignore content + .git
    for result in WalkBuilder::new(current_path)
        .max_depth(Some(1))
        .hidden(false)
        .git_ignore(true)
        .filter_entry(|entry| {
            !entry.path().to_string_lossy().contains(".git")
        })
        .build()
    {
        match result {
            Ok(entry) => {
                // Skip the root directory itself
                if entry.path() != current_path {
                    items.push(entry.into_path());
                }
            } Err(err) => {
                eprintln!("Error walking directory: {}", err);
            }
        }
    }
    
    items
}

pub fn highlight_file_content(file_path: &str) -> Result<Text> {
    // Load syntax definitions and theme
    let ps = SyntaxSet::load_defaults_newlines();
    let ts = ThemeSet::load_defaults();
    let theme = &ts.themes["base16-ocean.dark"]; // Choose a theme

    // Find the syntax for the file
    let syntax = ps.find_syntax_for_file(file_path)?
        .unwrap_or_else(|| ps.find_syntax_plain_text());

    // Create a highlighter
    let mut h = HighlightLines::new(syntax, theme);

    // Read the file content
    let content = fs::read_to_string(file_path)?;

    // Highlight the content (with a line limit)
    let mut highlighted_lines = Vec::new();
    let max_lines = 100; // Limit the number of lines rendered
    for (line_number, line) in LinesWithEndings::from(&content).enumerate() {
        if line_number >= max_lines {
            break;
        }

        let ranges = h.highlight_line(line, &ps)?;
        let spans = ranges
            .into_iter()
            .map(|(style, text)| {
                Span::styled(
                    text.to_string(),
                    Style::default()
                        .fg(Color::Rgb(style.foreground.r, style.foreground.g, style.foreground.b))
                        .bg(Color::Reset),
                )
            })
            .collect::<Vec<_>>();
        highlighted_lines.push(Line::from(spans));
    }

    Ok(Text::from(highlighted_lines))
}

